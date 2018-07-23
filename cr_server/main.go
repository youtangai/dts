package main

import (
	"flag"
	pb "github.com/youtangai/zanshin/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	client pb.FileTransferServiceClient
)

const (
	MAX_BUFF = 16384
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
		Mode: mode,
	}
	res, err := client.GetFileInfo(ctx, fileInfo)
	if err != nil {
		return err
	}
	log.Printf("return message: %s\n", res.Message)
	return nil
}

func runTransferFile(file *os.File) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.TransferFile(ctx)
	buff := make([]byte, MAX_BUFF)
	n := 0
	for {
		log.Println("n:", n)
		n++
		count, err := file.Read(buff)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read: %v\n", err)
		}
		fileData := &pb.FileData{
			Data: buff[:count],
		}
		if err := stream.Send(fileData); err != nil {
			return err
		}
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
	dir := flag.String("dir", "", "transfer dir")
	containerId := flag.String("c", "", "container id")
	flag.Parse()

	//todo: container checkpoint process

	log.Printf("container id : %s\n", *containerId)

	log.Printf("request to %s:%s", *host, *port)

	conn, err := grpc.Dial(*host+":"+*port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client = pb.NewFileTransferServiceClient(conn)

	runGetFolderInfo(*dir)

	files, err := ioutil.ReadDir(*dir)
	if err != nil {
		log.Fatalf("failed to read dir: %+v\n", err)
	}

	err = os.Chdir(*dir)
	if err != nil {
		log.Fatalf("failed to change dir: %v\n", err)
	}

	for _, file := range files {
		runGetFileInfo(file.Name(), file.Size(), uint32(file.Mode()))
		f, err := os.Open(file.Name())
		if err != nil {
			log.Fatalf("cannot open file: %s err: %v\n", file.Name(), err)
		}
		runTransferFile(f)
	}
}
