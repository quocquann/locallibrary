# Overview
A service crawl a website of book store to get information of books using [goquery](https://github.com/PuerkitoBio/goquery) and [gRPC](https://grpc.io/docs/languages/go/).

# Setup
- Download and install [Go 1.22.0]("https://go.dev/doc/install)
- Run below command to get all dependences that you need to run this service:
```
go mod tidy
```
# Commands
- Run server:
```
go run server/server.go
```

- Test the server with client:
    
```
go run client/client.go
```

