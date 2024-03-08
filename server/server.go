package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/quocquann/locallibrary/crawler"
	"github.com/quocquann/locallibrary/db"
	pb "github.com/quocquann/locallibrary/protos/book"
	"github.com/quocquann/locallibrary/queries"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedBookServer
	db *sql.DB
}

func (s *server) GetBook(req *pb.BookRequest, stream pb.Book_GetBookServer) error {
	bookQueries := queries.NewBookQueries(s.db)
	books, err := crawler.CrawlBook("https://gacxepbookstore.vn")
	if err != nil {
		return err
	}

	if err = bookQueries.AddBooks(books); err != nil {
		log.Println(err)
		return nil
	}

	for _, book := range books {
		stream.Send(&pb.BookResponse{Title: book.Title, Image: book.Image, Author: book.Author.Name, Genre: book.Genre})
	}
	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Fail to load .env file")
	}

	lis, err := net.Listen("tcp", "0.0.0.0:8080")

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	db, err := db.OpenConnection()
	if err != nil {
		log.Fatalf("Cannot connect to db: %v", err)
	}
	pb.RegisterBookServer(s, &server{db: db})

	stop := make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Println(err)
			stop <- syscall.SIGINT
		}
	}()
	<-stop
	s.GracefulStop()
	log.Println("server's shutting down")
}
