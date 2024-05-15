FROM golang:1.22.3

WORKDIR /app

COPY go.mod go.sum /
RUN go mod download

COPY cmd/client/ /cmd/client/
COPY internal/ /internal/

RUN go build -o /client /cmd/client/*.go

CMD ["/client"]
