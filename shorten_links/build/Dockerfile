#FROM ubuntu:18.04
FROM golang:latest


COPY go.mod go.sum ./
 #RUN chmod -R 755 .
#RUN go mod tidy
RUN go mod download

COPY . ./
WORKDIR ./cmd/shortener
RUN go build -o main .
CMD ["./main"]