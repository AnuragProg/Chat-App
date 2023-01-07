FROM golang:latest

EXPOSE 5000 4000

COPY . /app

WORKDIR /app

ENV LOADBALANCER=http://localhost:3000

RUN go build -o app

CMD ["/app/app"]