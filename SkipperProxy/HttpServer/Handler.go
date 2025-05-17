package HttpServer

import (
	tcpserver "SkipperTunnelProxy/TcpServer"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"github.com/google/uuid"

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
	RequestID string   `json:"requestID"`
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
	RequestID string			`json:"requestID"` 
}

func (s *Server) ParseHttpRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Ejecutando ParseHttpRequest...")

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

	target := r.Host
if idx := strings.Index(target, "."); idx != -1 {
    target = target[:idx] // Elimina el puerto (todo después de ":")
}

	requestID:= uuid.New().String()

	request := HttpRequest{
		Method:    r.Method,
		Proto:     r.Proto,
		TargetUri: target,
		Path:      r.URL.RequestURI(),
		Header:    headers,
		Body:      string(bodyBytes),
		RequestID: requestID,

	}

	
	requestBytes, err := json.Marshal(request)
	if err != nil {
		http.Error(w, "Error al convertir la solicitud a JSON", http.StatusInternalServerError)
		return
	}

	ResponseChannel:= make(chan []byte)

	s.RegisterResponseChannel(requestID, ResponseChannel)


	tcpMessage := tcpserver.TcpMessage{
		Target: request.TargetUri,
		Data:   requestBytes,
	}
	fmt.Println("envio un mensajito")

	s.TcpRequestChannel <- tcpMessage

	responseBytes := <-ResponseChannel

	// Responder al cliente HTTP
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBytes)

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
