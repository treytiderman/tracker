# Tracker

Track your life



## Docker

Requires Docker to be installed and the source code to be downloaded

Start in the `./` directory

```bash
cd ~/tracker
```

### Docker

```bash
docker build -t tracker .
docker rm -f tracker
docker run -d --name tracker -p 8000:8000 -v ./data:/app/data tracker
```



### Docker Compose

Start

```bash
docker compose up -d
```

Stop

```bash
docker compose down
```



## Docker Compose Dev

```bash
docker compose -f docker-compose-dev.yaml up
```



## Source Code

Requires Go to be installed and the source code to be downloaded

Start in the `./src` directory

```bash
cd ~/tracker
cd ./src
```

### Source Code: Build / Compile

```bash
go build .
```

### Source Code: Build / Compile for Windows

When on windows

```bash
go build -o ./build/app.exe .
```

When on Linux

```bash
GOOS=windows GOARCH=amd64 go build .
```

### Source Code: Test

```bash
go test
```

Verbose

```bash
go test -v
```

### Source Code: Dev

Then run

```bash
go run .
```
