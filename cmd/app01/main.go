package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Name string `json:"name"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	_, _ = w.Write([]byte("Hello, World!"))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("welcome to go"))
}

func queryHandler(w http.ResponseWriter, r *http.Request) { //*http.Request is struct
	//http://localhost:5000//query?name="mahesh"
	name := r.URL.Query().Get("name")
	age := r.URL.Query().Get("age")

	if name == "" {
		name = "Guest"
	}

	_, _ = w.Write([]byte(name))
	_, _ = w.Write([]byte(age))

}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	// Correctly handle /product/<id>
	productID := r.URL.Path[len("/product/"):]
	if productID == "" {
		http.Error(w, "Product ID not provided", http.StatusBadRequest)
		return
	}
	_, _ = w.Write([]byte(fmt.Sprintf("Product ID: %s", productID)))
}

func successHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := map[string]any{
		"ok":       true,
		"message":  "JSAON encose successfull",
		"datetime": time.Now().UTC(),
	}
	_ = json.NewEncoder(w).Encode(res)
}

func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func responseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, http.StatusMethodNotAllowed, map[string]any{
			"ok":    false,
			"error": "Only post is allowed",
		})
		return
	}
	defer r.Body.Close()

	var req Request

	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&req); err != nil {
		writeJson(w, http.StatusMethodNotAllowed, map[string]any{
			"ok":    false,
			"error": "Invalid json format",
		})
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		writeJson(w, http.StatusMethodNotAllowed, map[string]any{
			"ok":    false,
			"error": "Name is required",
		})
		return
	}
	writeJson(w, http.StatusOK, map[string]any{
		"ok":      true,
		"message": fmt.Sprintf("Hello %s", req.Name),
	})

}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "not allowed", http.StatusMethodNotAllowed)
		return
	}
	_, _ = w.Write([]byte("Welcome"))
}

// Step 1: Program starts
// Go runtime starts the program
// It looks for and runs main()
// At this point:
// No server exists yet
// Nothing is listening on the network
func main() {
	// Step 2: Register handlers (routing setup)
	// 	What this does internally:
	// Registers a route
	// Associates:
	// Path: /hello
	// Handler function: helloHandler
	// Stores this mapping in Go’s default ServeMux
	// Important concepts:
	// ServeMux = request router
	//"hello"->helloHandler
	// Default one is http.DefaultServeMux
	//No network activity yet — just configuration.
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/welcome", welcomeHandler)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/query", queryHandler)
	http.HandleFunc("/product/", pathHandler)
	http.HandleFunc("/ok", successHandler)
	http.HandleFunc("/response", responseHandler)

	// Step 3: Start the HTTP server
	// What happens here:ytr
	// Creates an http.Server
	// Binds to port 8080
	// Starts listening for TCP connections
	// Blocks forever (main goroutine stops here)
	// Key detail:
	// nil means:
	// use http.DefaultServeMux as the handler
	// server := &http.Server{
	//     Addr:    ":8080",
	//     Handler: http.DefaultServeMux,
	// }
	// server.ListenAndServe()

	err := http.ListenAndServe(":5000", nil)
	// 	If you pass nil, Go uses the default HTTP request multiplexer:

	// http.DefaultServeMux

	// That means it will use any routes you registered earlier like:
	if err != nil {
		panic(err)
	}
}
