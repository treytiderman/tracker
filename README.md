# Tracker

## Run

```
go run .
```

## Build - Windows

When on windows

```
go build -o ./build/app.exe .
```

When on Linux

```
GOOS=windows GOARCH=amd64 go build .
```

## Test

```
go test
```

Verbose

```
go test -v
```

## Docker

```
sudo docker build -t tracker .
```
