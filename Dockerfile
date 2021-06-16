FROM golang:1.16-alpine AS builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .
EXPOSE 8080
CMD ["/app/main"]