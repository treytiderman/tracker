# docker build -t tracker .
# docker stop tracker
# docker rm tracker
# docker run -d --name tracker -p 8000:8000 -v ./data:/app/data tracker

FROM golang:1.23 AS build
WORKDIR /app/src
COPY ./src /app/src
COPY ./public /app/public
RUN go build -o tracker
CMD "/app/src/tracker"

# FROM scratch
# WORKDIR /app
# COPY --from=build /app/tracker /app/tracker
# CMD ["/app/tracker"]
