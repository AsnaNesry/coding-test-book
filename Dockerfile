
FROM golang:1.16.3-alpine3.13

RUN apk add --no-cache gcc musl-dev libc-dev

WORKDIR /app

COPY . . 

RUN go get -d -v ./...

RUN go build -o coding-test .

EXPOSE 8080

# Run the executable
CMD ["./coding-test"]