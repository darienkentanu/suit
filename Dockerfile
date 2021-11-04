FROM golang:1.17

RUN mkdir /app

WORKDIR /app

COPY go.mod /app
COPY go.sum /app
RUN go mod download

ADD . /app

RUN go build -o main .

CMD ["/app/main"]