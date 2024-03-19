package main

import (
	"context"
	"io"
	"log"

	pb "github.com/quocquann/locallibrary/protos/book"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial("0.0.0.0:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer cc.Close()

	client := pb.NewBookClient(cc)
	getBook(client)
}

func getBook(client pb.BookClient) {
	stream, err := client.GetBook(context.Background(), &pb.BookRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Print("Server finish streaming")
			return
		}
		log.Printf("{\n\tTitle: %v\n\tImage: %v\n\tAuthor: %v\n\tGenre: %v\n}\n", res.Title, res.Image, res.Author, res.Genre)
	}
}
