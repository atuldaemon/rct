FROM golang:alpine as builder
RUN go version

COPY . "/go/src/github.com/atuldaemon/rct"
WORKDIR "/go/src/github.com/atuldaemon/rct"

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /rct .
CMD ["/rct"]
EXPOSE 8080

FROM scratch
COPY --from=builder /rct .
EXPOSE 8080
CMD ["/rct"]
