package main

import (
	"encoding/json"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// parsedBody is the decoded body plus whether decoding succeeded.
// value holds the decoded structure (or the raw string when the
// Content-Type is not one we structurally decode).
type parsedBody struct {
	value any
	ok    bool
}

// parseBody decodes the request body according to its Content-Type.
// For application/json it reports ok=false on a decode failure so the
// handler can respond 400, while still returning the raw bytes for logging.
func parseBody(contentType string, body []byte) parsedBody {
	if len(body) == 0 {
		return parsedBody{value: nil, ok: true}
	}

	switch {
	case isContentType(contentType, "application/json"):
		var data any
		if err := json.Unmarshal(body, &data); err != nil {
			return parsedBody{value: string(body), ok: false}
		}
		return parsedBody{value: data, ok: true}

	case isContentType(contentType, "application/x-www-form-urlencoded"):
		values, err := url.ParseQuery(string(body))
		if err != nil {
			return parsedBody{value: string(body), ok: false}
		}
		return parsedBody{value: map[string][]string(values), ok: true}

	default:
		return parsedBody{value: string(body), ok: true}
	}
}

// isContentType matches a media type, ignoring any parameters such as
// "; charset=utf-8".
func isContentType(header, mediaType string) bool {
	if header == mediaType {
		return true
	}
	return strings.HasPrefix(header, mediaType+";")
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read request body",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("error", err.Error()),
		)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	contentType := r.Header.Get("Content-Type")
	pb := parseBody(contentType, body)

	// Record the receipt first, regardless of how we respond.
	slog.Info("request captured",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Any("query", map[string][]string(r.URL.Query())),
		slog.Any("headers", map[string][]string(r.Header)),
		slog.Any("body", pb.value),
		slog.Bool("body_parsed", pb.ok),
	)

	// A declared content type that does not parse is a malformed request.
	if !pb.ok {
		http.Error(w, "Malformed request body for declared Content-Type", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if r.Method == http.MethodGet {
		w.Write([]byte("Hello from GET\n"))
	} else {
		w.Write([]byte("ok\n"))
	}
}

func main() {
	addr := flag.String("addr", ":8000", "Address and port to listen on")
	flag.Parse()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	http.HandleFunc("/", handler)
	slog.Info("HttpRequestCapture starting", slog.String("addr", *addr))
	if err := http.ListenAndServe(*addr, nil); err != nil {
		slog.Error("server stopped", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
