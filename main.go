package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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

func parseQuery(query string) map[string][]string {
	parsed, _ := url.ParseQuery(query)
	return parsed
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Print full request details
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Full request dump:")
	fmt.Println(string(requestDump))

	// Print path and parsed query
	fmt.Printf("\npath = %s\n", r.URL.Path)
	parsedQuery := parseQuery(r.URL.RawQuery)
	prettyQuery, _ := json.MarshalIndent(parsedQuery, "", "  ")
	fmt.Printf("parsed: path = %s, query = %s\n", r.URL.Path, string(prettyQuery))

	// Print headers
	fmt.Println("\nHeaders:")
	fmt.Println("-----")
	fmt.Println(dumpHeaders(r.Header))
	fmt.Println("-----")

	if r.Method == "POST" {
		// Read body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		contentType := r.Header.Get("Content-Type")

		if contentType == "application/x-www-form-urlencoded" {
			fmt.Printf("\nBody (application/x-www-form-urlencoded):\n%s\n", string(body))
		}

		if contentType == "application/json" {
			var data interface{}
			err := json.Unmarshal(body, &data)
			if err != nil {
				http.Error(w, "Error parsing JSON", http.StatusBadRequest)
				return
			}
			prettyJSON, _ := json.MarshalIndent(data, "", "  ")
			fmt.Printf("\nBody (application/json):\n%s\n", string(prettyJSON))
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if r.Method == "GET" {
		w.Write([]byte("Hello from GET\n"))
	} else {
		w.Write([]byte("ok\n"))
	}
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("HttpRequestCapture: Serving on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
