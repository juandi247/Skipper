package HttpServer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "welcome to juandi http")
}

type test_struct struct {
	Test string
}

type HttpRequest struct {
	Method    string            `json:"method"`  // GET, POST, PUT, etc.
	Proto   string            `json:"version"` // HTTP/1.1, HTTP/2
	TargetUri string            `json:"target"`  // todo: the subdomain that the user will select on tunnel
	Path string	  `json:"path"`  // "/api/endpoint", "/login", etc.
	Header   map[string]string `json:"headers"` // Cookies, tokens, user-agent, etc.
	Body      string            `json:"body"`    // Body
}

// func ParseHttpRequest(w http.ResponseWriter, r *http.Request) {

// 	request := HttpRequest{
// 		Method:  r.Method,
// 		Proto: r.Proto,
// 		TargetUri: r.URL.Host,
// 		Path: r.URL.RequestURI(),
// 		Header: make(map[string]string),
// 		Body: "".
// 	}


// 	headers:=make(map[string]string)

// 	for key, value := range(r.Header){
// 		headers[k] = strings.Join(v, ", ")
// 	}

// 	bodyBytes, err := io.ReadAll(r.Body)
// if err != nil {
// 	http.Error(w, "error al leer body", http.StatusInternalServerError)
// 	return
// }
// req.Body = string(bodyBytes)



// }

func ParsePost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var t test_struct
	err := decoder.Decode(&t)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, t.Test)

}

func TimeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	currentTime := time.Now().Format(time.RFC3339)
	json.NewEncoder(w).Encode(map[string]string{"time": currentTime})
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "404 - Page not found")

}
