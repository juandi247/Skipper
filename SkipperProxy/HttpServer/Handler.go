package HttpServer

import (
	"SkipperTunnelProxy/gen"
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
		fmt.Println("nos toco entrsar a localhost", host, r.URL.RequestURI())
		s.Templates.ExecuteTemplate(w, "index.html", nil)
		return
	}

	target := strings.Join(parts[:len(parts)-1], ".")
	// prod
	// target := strings.Join(parts[:len(parts)-2], ".")

	_, err := s.ConnectionManager.GetTunnelConnection(target)
	if err != nil {
		s.Templates.ExecuteTemplate(w, "error.html", nil)
		return
	}
	ResponeChannel := make(chan []byte)
	requestId := uuid.New().String()

	// headers read
	headers := make(map[string]string)
	for key, value := range r.Header {
		headers[key] = strings.Join(value, ", ")
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
		Headers:   headers,
		Body:      string(bodyBytes),
		RequestId: requestId,
	}
	fmt.Printf(" la request se va a hacer con la siguente url %v ", request.Path)

	requestBytes, err := proto.Marshal(request)
	if err != nil {
		http.Error(w, "Error al convertir la solicitud a BYTES", http.StatusInternalServerError)
		return
	}

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
	// fmt.Println("envie una request", request)
	// fmt.Println("mensaje enviado", string(requestBytes))
	s.ConnectionManager.SaveResponseChannel(requestId, ResponeChannel)

	select {
	case respBytes := <-ResponeChannel:
		var response gen.Response
		err := proto.Unmarshal(respBytes, &response)
		// fmt.Println("LLEGO MENSJAEEEEE")
		if err != nil {
			http.Error(w, "Error parsing response", http.StatusInternalServerError)
			return
		}

		// fmt.Printf("recibi una respuesta", response.Body)

		// TODO : change this becasue somethigns would be a json, so depending on the response we get w.write depending on them

		// w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// w.Write([]byte(response.Body))

		// Puedes también enviar otros headers que recibas, si quieres
		for key, value := range response.Headers {
			w.Header().Set(key, value)
		}
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
