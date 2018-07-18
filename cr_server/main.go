package main

import (
	pb "github.com/youtangai/zanshin/proto"
	"golang.org/x/net/context"
	"time"
	"log"
	"google.golang.org/grpc"
	"flag"
	"io/ioutil"
	//"path/filepath"
	"os"
	"io"
)

var (
	client pb.FileTransferServiceClient
)

func runGetFolderInfo(folderName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	folderInfo := &pb.FolderInfo{
		Name: folderName,
	}
	res, err := client.GetFolderInfo(ctx, folderInfo)
	if err != nil {
		return err
	}
	log.Printf("return message: %s\n", res.Message)
	return nil
}

func runGetFileInfo(name string, size int64, mode uint32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fileInfo := &pb.FileInfo{
		Name: name,
		Size: size,
		Mode: mode,
	}
	res, err := client.GetFileInfo(ctx, fileInfo)
	if err != nil {
		return err
	}
	log.Printf("return message: %s\n", res.Message)
	return nil
}

func runTransferFile(data []byte) error {
	log.Println(string(data))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.TransferFile(ctx)
	for {
		fileData := &pb.FileData {
			Data: data,
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
	dir := flag.String("dir", ".", "transfer dir")
	flag.Parse()
	log.Printf("request to %s:%s", *host, *port)

	conn, err := grpc.Dial(*host + ":" + *port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client = pb.NewFileTransferServiceClient(conn)

	runGetFolderInfo(*dir)
	
	files, err := ioutil.ReadDir(*dir)
	if err != nil {
		log.Fatalf("failed to read dir: %+V\n", err)
	}
	for _, file := range files {
		runGetFileInfo(file.Name(), file.Size(), uint32(file.Mode()))
		err = os.Chdir(*dir)
		if err != nil {
			log.Fatalf("failed to change dir: %v\n", err)
		}
		f, err := os.Open(file.Name())
		if err != nil {
			log.Fatalf("cannot open file: %s err: %v\n", file.Name(), err)
		}
		buff := make([]byte, file.Size())
		for {
			count, err := f.Read(buff)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("failed to read: %v\n", err)
			}
			runTransferFile(buff[:count])
		}
	}
}
