# Overview
A service crawl a website of book store to get information of books using [goquery](https://github.com/PuerkitoBio/goquery) and [gRPC](https://grpc.io/docs/languages/go/).

# Setup
Run below command to get all dependences that you need to run this service:
```
go mod tidy
```
# Scripts
- Run server:
```
go run server/server.go
```

- Test the server with client:
    
```
go run client/client.go
```

