package HttpServer

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"github.com/google/uuid"
 	"google.golang.org/protobuf/proto"
	"SkipperTunnelProxy/gen"

)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "welcome to juandi http")
}

func (s *Server) HandleClientRequest(w http.ResponseWriter, r *http.Request) {
	host := r.Host

	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	// const baseDomain = "skipper.lat"
	const baseDomain = "localhost:8080"


	parts := strings.Split(host, ".")

	if len(parts) <= 1 {
	// prod
	// if len(parts) <= 2 {
		s.Templates.ExecuteTemplate(w, "index.html", nil)
		return
	}

	target := strings.Join(parts[:len(parts)-1], ".") 
	// prod
	// target := strings.Join(parts[:len(parts)-2], ".") 

	fmt.Println("TARGETTETTT", target)

	_, err := s.ConnectionManager.GetTunnelConnection(target)
	if err != nil {
		s.Templates.ExecuteTemplate(w, "error.html", nil)
		return
	}
	ResponeChannel := make(chan []byte)
	requestId := uuid.New().String()

	fmt.Println("YA PASE VALIDACIONNNN!!!")

	// headers read
	headers := make(map[string]string)
	validHeaders := map[string]bool{
		// "User-Agent":      true,
		"Content-Type":  true,
		"Authorization": true,
		// "Accept":          true,
		"Content-Length": true,
	}
	for key, value := range r.Header {
		if validHeaders[key] {
			headers[key] = strings.Join(value, ", ")
		}
	}


	// body read
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error al leer body", http.StatusInternalServerError)
		return
	}

	request := &gen.Request{
		Method:    r.Method,
		Proto:     r.Proto,
		TargetUri: r.Host,
		Path:      r.URL.RequestURI(),
		Headers:    headers,
		Body:      string(bodyBytes),
		RequestId: requestId,
	}
	
	
	requestBytes, err := proto.Marshal(request)
	if err != nil {
		http.Error(w, "Error al convertir la solicitud a BYTES", http.StatusInternalServerError)
		return
	}

	fmt.Println("sirvieron los marshalsssss")


	// !length frame  (4 bytes) for handling the request on the tunnel
	requestLength := uint32(len(requestBytes))
	// we create a raw buffer
	requestBuffer := new(bytes.Buffer)
	/* we append to the buffer the request length
	now we append to the requestBuffer (containing the 4 byte integer with the length),
	we would only have on the request buffer the integer
	*/
	binary.Write(requestBuffer, binary.BigEndian, requestLength)

	// so we now write the request on the buffer
	requestBuffer.Write(requestBytes)

	err = s.ConnectionManager.SendMessageToTunnel(target, requestBuffer.Bytes())
	// fmt.Println("envio un mensajito")
	// fmt.Println("mensaje enviado", string(requestBytes))
	s.ConnectionManager.SaveResponseChannel(requestId, ResponeChannel)

	fmt.Println("gamos a esperarrrr")

	select {
	case respBytes := <-ResponeChannel:
		var response gen.Response
		err := proto.Unmarshal(respBytes, &response)
		// fmt.Println("LLEGO MENSJAEEEEE")
		if err != nil {
			http.Error(w, "Error parsing response", http.StatusInternalServerError)
			return
		}

		// fmt.Printf("\n \n RESPUESTA %v \n", response)

		// TODO : change this becasue somethigns would be a json, so depending on the response we get w.write depending on them

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(response.Body))

	case <-time.After(10 * time.Second):
		s.Templates.ExecuteTemplate(w, "timeout.html", nil)

	}

	// s.TcpRequestChannel <- tcpMessage
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(request)

}

// ! for testing too
type test_struct struct {
	Test string
}

func ParsePost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var t test_struct
	err := decoder.Decode(&t)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, t.Test)

}

// ! just for testing the server
func TimeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	currentTime := time.Now().Format(time.RFC3339)
	json.NewEncoder(w).Encode(map[string]string{"time": currentTime})
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "404 - Page not found")

}

// RegisterHandlers sets up the HTTP routes and static file serving.
// You should call this function when initializing your HTTP server.
func (s *Server) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", s.HandleClientRequest)
	// Serve static files from the 'assets' directory
	fs := http.FileServer(http.Dir("./templates/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Other existing handlers (adjust paths if necessary)
	mux.HandleFunc("/parse", ParsePost)
	mux.HandleFunc("/time", TimeHandler)
	mux.HandleFunc("/404", NotFoundHandler) // Assuming you might have a direct route for 404
}
