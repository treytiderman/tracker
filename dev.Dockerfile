FROM golang:1.23
WORKDIR /app
RUN go install github.com/air-verse/air@latest
CMD "air"
