FROM golang:1.23 AS build
WORKDIR /app
COPY ./ ./
RUN go build -o tracker
CMD "/app/tracker"

# FROM scratch
# WORKDIR /app
# COPY --from=build /app/tracker /app/tracker
# CMD ["/app/tracker"]
