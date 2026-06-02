# http-request-capture

An HTTP server that records each incoming request as a single line of JSON on stdout.

## What it does

http-request-capture listens for HTTP requests on any path and logs the method, path, query parameters, headers, and body of each request as a structured JSON record using `log/slog`. JSON and form-urlencoded bodies are decoded into structured values; any other body is recorded as a string. A request that declares a JSON content type but cannot be parsed is recorded and answered with `400 Bad Request`. Because the output is one JSON object per line, it can be inspected with `jq` or forwarded to a log pipeline.

## Usage

Start the server (defaults to `:8000`):

```
$ ./http-request-capture
{"time":"2026-06-02T17:15:54.668175236+09:00","level":"INFO","msg":"HttpRequestCapture starting","addr":":8000"}
```

Listen on a different address with `-addr`:

```
$ ./http-request-capture -addr :9000
```

Send requests with any HTTP client. Each request produces one JSON line on stdout.

### GET request

```
$ curl "http://localhost:8000/test?param1=value1&param2=value2"
Hello from GET
```

```
{"time":"2026-06-02T17:16:17.152773634+09:00","level":"INFO","msg":"request captured","method":"GET","path":"/test","query":{"param1":["value1"],"param2":["value2"]},"headers":{"Accept":["*/*"],"User-Agent":["curl/8.19.0-rc3"]},"body":null,"body_parsed":true}
```

Any path is accepted and recorded in the `path` field:

```
$ curl "http://localhost:8000/test/hello?param1=value1&param2=value2"
Hello from GET
```

```
{"time":"2026-06-02T17:16:34.292500152+09:00","level":"INFO","msg":"request captured","method":"GET","path":"/test/hello","query":{"param1":["value1"],"param2":["value2"]},"headers":{"Accept":["*/*"],"User-Agent":["curl/8.19.0-rc3"]},"body":null,"body_parsed":true}
```

### POST request (application/x-www-form-urlencoded)

```
$ curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -d "username=johndoe&password=secret123" http://localhost:8000/
ok
```

```
{"time":"2026-06-02T17:17:07.667334867+09:00","level":"INFO","msg":"request captured","method":"POST","path":"/","query":{},"headers":{"Accept":["*/*"],"Content-Length":["35"],"Content-Type":["application/x-www-form-urlencoded"],"User-Agent":["curl/8.19.0-rc3"]},"body":{"password":["secret123"],"username":["johndoe"]},"body_parsed":true}
```

### POST request (application/json)

```
$ curl -X POST -H "Content-Type: application/json" -d '{"username":"johndoe","password":"secret123"}' http://localhost:8000/
ok
```

```
{"time":"2026-06-02T17:17:31.000000000+09:00","level":"INFO","msg":"request captured","method":"POST","path":"/","query":{},"headers":{"Accept":["*/*"],"Content-Length":["45"],"Content-Type":["application/json"],"User-Agent":["curl/8.19.0-rc3"]},"body":{"password":"secret123","username":"johndoe"},"body_parsed":true}
```

### Malformed JSON body

A body that declares `application/json` but does not parse is recorded with
`body_parsed` set to false and answered with `400 Bad Request`. The raw body is
kept in the `body` field.

```
$ curl -X POST -H "Content-Type: application/json" -d '{"broken":' http://localhost:8000/
Malformed request body for declared Content-Type
```

```
{"time":"2026-06-02T17:18:02.000000000+09:00","level":"INFO","msg":"request captured","method":"POST","path":"/","query":{},"headers":{"Accept":["*/*"],"Content-Length":["10"],"Content-Type":["application/json"],"User-Agent":["curl/8.19.0-rc3"]},"body":"{\"broken\":","body_parsed":false}
```

### Filtering output with jq

```
$ ./http-request-capture | jq 'select(.msg == "request captured") | {path, body}'
```

## License

This project is licensed under the [MIT License](LICENSE).
