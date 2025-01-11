# BUILD: docker build -t tracker .
# STOP: docker stop tracker
# REMOVE: docker rm tracker
# RUN: docker run -d --name tracker -p 8000:8000 -v ./data:/app/data tracker

FROM golang:1.23 AS build
WORKDIR /app/src
COPY ./src /app/src
COPY ./public /app/public

# Add data folder and test.db
RUN mkdir /app/data
RUN touch /app/data/test.db
RUN mkdir /app/content

# Run Tests and remove test.db
RUN go test
RUN rm /app/data/test.db

# Build and run binary
RUN go build -o tracker
CMD ["/app/src/tracker"]
