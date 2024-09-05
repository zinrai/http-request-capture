# http-request-capture

Lightweight Go application that serves as an HTTP server, capturing and displaying detailed information about incoming HTTP requests. It's an invaluable tool for developers and testers who need to inspect HTTP requests in real-time.

## Features

- Captures and displays full HTTP request details
- Provides a formatted dump of request headers
- Handles and displays POST request bodies for both `application/x-www-form-urlencoded` and `application/json` content types
- Easy to run and use with minimal setup

## Usage

1. Run the application:
   ```
   $ go run main.go
   ```
2. The server will start on port 8000. You should see the message:
   ```
   HttpRequestCapture: Serving on port 8000...
   ```
3. Send HTTP requests to `http://localhost:8000` using any HTTP client (browser, curl, Postman, etc.)
4. The application will capture and display the details of each incoming request in the console

## Example Output

When you send a request to the server, you'll see output similar to this:

### GET request

```bash
$ curl "http://localhost:8000/test?param1=value1&param2=value2"
Hello from GET
```

```
HttpRequestCapture: Serving on port 8000...
Full request dump:
GET /test?param1=value1&param2=value2 HTTP/1.1
Host: localhost:8000
Accept: */*
User-Agent: curl/8.9.1



path = /test
parsed: path = /test, query = {
  "param1": [
    "value1"
  ],
  "param2": [
    "value2"
  ]
}

Headers:
-----
Accept: */*
User-Agent: curl/8.9.1
-----

```

### POST request ( application/x-www-form-urlencoded )

```bash
$ curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -d "username=johndoe&password=secret123" http://localhost:8000
ok
```

```
HttpRequestCapture: Serving on port 8000...
Full request dump:
POST / HTTP/1.1
Host: localhost:8000
Accept: */*
Content-Length: 35
Content-Type: application/x-www-form-urlencoded
User-Agent: curl/8.9.1

username=johndoe&password=secret123

Request Headers:
Accept: */*
Content-Length: 35
Content-Type: application/x-www-form-urlencoded
User-Agent: curl/8.9.1

Body (application/x-www-form-urlencoded):
username=johndoe&password=secret123
```

### POST request ( application/json )

```bash
$ curl -X POST -H "Content-Type: application/json" -d '{"username":"johndoe","password":"secret123"}' http://localhost:8000
ok
```

```
HttpRequestCapture: Serving on port 8000...
Full request dump:
POST / HTTP/1.1
Host: localhost:8000
Accept: */*
Content-Length: 45
Content-Type: application/json
User-Agent: curl/8.9.1

{"username":"johndoe","password":"secret123"}

Request Headers:
Accept: */*
Content-Length: 45
Content-Type: application/json
User-Agent: curl/8.9.1

Body (application/json):
{
  "password": "secret123",
  "username": "johndoe"
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
