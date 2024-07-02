FROM golang:latest


COPY go.mod go.sum ./
RUN go mod download

COPY . ./
WORKDIR .
RUN go build -o main .
CMD ["./main"]