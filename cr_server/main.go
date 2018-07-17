package main

import (
	pb "github.com/youtangai/zanshin/proto"
	"golang.org/x/net/context"
	"time"
	"log"
	"google.golang.org/grpc"
	"flag"
)

func runGetFileInfo(client pb.FileTransferServiceClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fileInfo := &pb.FileInfo{
		Name: "name",
		Size: 128,
		Mode: 755,
	}
	res, err := client.GetFileInfo(ctx, fileInfo)
	if err != nil {
		return err
	}
	log.Printf("return message: %s\n", res.Message)
	return nil
}

func runTransferFile(client pb.FileTransferServiceClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.TransferFile(ctx)
	for {
		fileData := &pb.FileData {
			Data: []byte("hogehoge"),
		}
		if err := stream.Send(fileData);  err != nil {
			return err
		}
		break
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Printf("return message: %s\n", res.Message)
	
	return nil
} 

func main() {
	host := flag.String("host", "localhost", "server host name or ip addr")
	port := flag.String("port", "12345", "server port number")
	flag.Parse()
	log.Printf("request to %s:%s", *host, *port)
	conn, err := grpc.Dial(*host + ":" + *port, grpc.WithInsecure())
	if err != nil {
		log.Fatal("failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileTransferServiceClient(conn)

	runGetFileInfo(client)
	runTransferFile(client)
}
