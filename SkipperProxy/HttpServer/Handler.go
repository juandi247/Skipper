package HttpServer

import (
	"SkipperTunnelProxy/message"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
	"time"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "welcome to juandi http")
}

type HttpRequest struct {
	Method    string            `json:"method"`  // GET, POST, PUT, etc.
	Proto     string            `json:"version"` // HTTP/1.1, HTTP/2
	TargetUri string            `json:"target"`  // todo: the subdomain that the user will select on tunnel
	Path      string            `json:"path"`    // "/api/endpoint", "/login", etc.
	Header    map[string]string `json:"headers"` // Cookies, tokens. It is always a map
	Body      string            `json:"body"`    // Body
	RequestID string            `json:"requestID"`
}

type HttpResponse struct {
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
	// dont think i need those protos but for now
	ProtoMajor int               // e.g. 1
	ProtoMinor int               // e.g. 0
	Proto      string            `json:"version"` // HTTP/1.1, HTTP/2
	Header     map[string]string `json:"headers"` // Cookies, tokens. It is always a map
	Body       string            `json:"body"`    // Body
	RequestID  string            `json:"requestID"`
}

func (s *Server) HandleClientRequest(w http.ResponseWriter, r *http.Request) {

	target := r.Host
	if idx := strings.Index(target, "."); idx != -1 {
		target = target[:idx]
	}
	_, err := s.ConnectionManager.GetTunnelConnection(target)
	if err != nil {
		fmt.Fprintln(w, "404 - Page not found")
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

	request := HttpRequest{
		Method:    r.Method,
		Proto:     r.Proto,
		TargetUri: r.Host,
		Path:      r.URL.RequestURI(),
		Header:    headers,
		Body:      string(bodyBytes),
		RequestID: requestId,
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		http.Error(w, "Error al convertir la solicitud a JSON", http.StatusInternalServerError)
		return
	}
	tcpMessage := message.TcpMessage{
		Target: request.TargetUri,
		Data:   requestBytes,
	}

	err = s.ConnectionManager.SendMessageToTunnel(target, tcpMessage)
	fmt.Println("envio un mensajito", tcpMessage)
	fmt.Println("mensaje enviado", string(requestBytes))
	s.ConnectionManager.SaveResponseChannel(requestId, ResponeChannel)

	fmt.Println("gamos a esperarrrr")

	select {
	case respBytes := <-ResponeChannel:
		var response HttpResponse
		err := json.Unmarshal(respBytes, &response)
		fmt.Println("LLEGO MENSJAEEEEE")
		if err != nil {
			http.Error(w, "Error parsing response", http.StatusInternalServerError)
			return
		}

		fmt.Printf("\n \n RESPUESTA %v \n", response)

		// TODO : change this becasue somethigns would be a json, so depending on the response we get w.write depending on them

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(response.Body))

	case <-time.After(10 * time.Second):
		http.Error(w, "Timeout waiting for response", http.StatusGatewayTimeout)
	}

	// s.TcpRequestChannel <- tcpMessage
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(request)

}

// func ParseResponse(response HttpResponse){
// 	   // Configurar el tipo de contenido correctamente
//     contentType, ok := response.Header["Content-Type"]
//     if ok {
//         w.Header().Set("Content-Type", contentType)
//     } else {
//         w.Header().Set("Content-Type", "text/plain; charset=utf-8")
//     }
// }

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
