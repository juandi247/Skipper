package HttpUserClient

import (
	"SkipperTunnel/TcpUserClient"
	"SkipperTunnel/gen"
	"bytes"
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// todo: GRCP communication easier and fasterrr
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

type HttpClient struct {
	Addr   string
	Client *http.Client
}

func NewHttpCliennt(addr string, timeout time.Duration) *HttpClient {
	client := &HttpClient{
		Addr: addr,
		Client: &http.Client{
			Timeout: timeout,
		},
	}
	return client
}

func ReceiveRequest(url string, workerId int, requestChannel chan []byte, client *HttpClient, tcpConn net.Conn, ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			// fmt.Println("turning off the handle responses goroutine")
			wg.Done()
			return
		case request := <-requestChannel:
			// send the request.
			// fmt.Printf("Worker %d, executing the request\n", workerId)
			var httpRequest gen.Request
			err := proto.Unmarshal((request), &httpRequest)
			if err != nil {
				fmt.Println("Error al decodificar la request:", err)
				continue
			}

			requestID := httpRequest.RequestId

			response, _ := ConvertToHttpRequest(url, &httpRequest, client, requestID)
			TcpUserClient.HandleSendToTCP(response, tcpConn)
		}
	}
}

func ConvertToHttpRequest(url string, hr *gen.Request, client *HttpClient, requestID string) ([]byte, error) {

	body := bytes.NewBufferString(hr.Body)

	req, err := http.NewRequest(hr.Method, url+hr.Path, body)
	if err != nil {
		return nil, err
	}

	// Agregar los headers
	for k, v := range hr.Headers {
		req.Header.Set(k, v)
	}

	req.Proto = hr.Proto
	parts := strings.Split(hr.Proto, "/")
	if len(parts) == 2 {
		req.ProtoMajor = 1
		req.ProtoMinor = 1
	}

	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al enviar la solicitud: %v", err)
	}
	defer resp.Body.Close()

	// Procesar la respuesta
	response, err := ParseHttpResponse(resp, requestID)
	if err != nil {
		return nil, fmt.Errorf("error al parsear la respuesta: %v", err)
	}
	return response, nil

}

func ParseHttpResponse(r *http.Response, requestID string) ([]byte, error) {
	// fmt.Println("Ejecutando ParseHttpRequest...")

	// headers read
	headers := make(map[string]string)
	validHeaders := map[string]bool{
		"User-Agent":     true,
		"Content-Type":   true,
		"Authorization":  true,
		"Accept":         true,
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
		return nil, fmt.Errorf("error reading body")
	}

	response := &gen.Response{
		Status:    r.Status,
		Proto:     r.Proto,
		Headers:   headers,
		Body:      string(bodyBytes),
		RequestId: requestID,
	}
	requestBytes, err := proto.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("error converting to json")
	}
	// fmt.Println("respuesta del propio server", string(requestBytes))

	return requestBytes, nil
}
