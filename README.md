# Overview
A service crawl a website of book store to get information of books using [gofiber](https://gofiber.io/) [goquery](https://github.com/PuerkitoBio/goquery) and [gRPC](https://grpc.io/docs/languages/go/).

<<<<<<< HEAD
# Setup
- Download and install [Go 1.22.0]("https://go.dev/doc/install)
- Run below command to get all dependences that you need to run this service:
=======
Fiber server:
Crawl [website](https://gacxepbookstore.vn) to get book information and add to database

# Commands
- Run Fiber server:
>>>>>>> e1a2eddb1e41603b6a9aa14f08428cb7306a29ff
```
go run main.go
```
<<<<<<< HEAD
# Commands
- Run server:
=======
- Run gRPC server:
>>>>>>> e1a2eddb1e41603b6a9aa14f08428cb7306a29ff
```
go run server/server.go
```

- Test gRPC server with client:
    
```
go run client/client.go
```

