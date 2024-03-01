package main

import (
	"log"
	"net"

	"github.com/quocquann/locallibrary/crawler"
	pb "github.com/quocquann/locallibrary/protos/book"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedBookServer
}

func (*server) GetBook(req *pb.BookRequest, stream pb.Book_GetBookServer) error {
	books, err := crawler.CrawlBook("https://gacxepbookstore.vn")
	if err != nil {
		return err
	}
	for _, book := range books {
		stream.Send(&pb.BookResponse{Title: book.Title, Image: book.Image, Author: book.Author, Genre: book.Genre})
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:8000")

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	pb.RegisterBookServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
