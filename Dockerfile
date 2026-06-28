FROM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /http-request-capture

FROM scratch

COPY --from=build /http-request-capture /http-request-capture
EXPOSE 8000
USER 65534:65534
ENTRYPOINT ["/http-request-capture"]
CMD ["-addr", "0.0.0.0:8000"]
