FROM golang:1.9.2

COPY ./ /go/src/github.com/jeremyroberts0/itunes-to-spotify
WORKDIR /go/src/github.com/jeremyroberts0/itunes-to-spotify

RUN go get ./

# Statically link so we can run the binary anywhere
RUN CGO_ENABLED=0 GOOS=linux go build -a -o app .

# Multi-stage build, 2nd container, running app is just alpine; much small, such secure, wow.
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/jeremyroberts0/itunes-to-spotify/app .

EXPOSE 8081:8081

CMD ["./app"]
