# Overview
A service crawl a website of book store to get information of books using [gofiber](https://gofiber.io/) [goquery](https://github.com/PuerkitoBio/goquery) and [gRPC](https://grpc.io/docs/languages/go/).

Fiber server:
Crawl [website](https://gacxepbookstore.vn) to get book information and add to database

# Commands
- Run Fiber server:
```
go run main.go
```
- Run gRPC server:
```
go run server/server.go
```

- Test gRPC server with client:
    
```
go run client/client.go
```

