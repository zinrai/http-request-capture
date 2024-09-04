package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"sort"
	"strings"
)

func dumpHeaders(headers http.Header) string {
	var headerLines []string
	for name, values := range headers {
		for _, value := range values {
			headerLines = append(headerLines, fmt.Sprintf("%s: %s", name, value))
		}
	}
	sort.Strings(headerLines)
	return strings.Join(headerLines, "\n")
}

func handler(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Full request dump:")
	fmt.Println(string(requestDump))

	fmt.Println("\nRequest Headers:")
	fmt.Println(dumpHeaders(r.Header))

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	contentType := r.Header.Get("Content-Type")

	if r.Method == "POST" && contentType == "application/x-www-form-urlencoded" {
		fmt.Printf("\nBody (application/x-www-form-urlencoded):\n%s\n", string(body))
	}

	if r.Method == "POST" && contentType == "application/json" {
		var data interface{}
		err := json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}
		prettyJSON, _ := json.MarshalIndent(data, "", "  ")
		fmt.Printf("\nBody (application/json):\n%s\n", string(prettyJSON))
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok\n"))
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("HttpRequestCapture: Serving on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
