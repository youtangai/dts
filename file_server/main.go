package main

import (
	pb "github.com/youtangai/zanshin/proto"
	"io"
	"log"
	"net"
	"flag"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
)

type fileTransferService struct {
	fileName string
	size int64
	mode uint32
}

func (fts *fileTransferService) GetFileInfo(ctx context.Context, info *pb.FileInfo) (*pb.Res, error) {
	fts.fileName = info.Name
	fts.size = info.Size
	fts.mode = info.Mode
	return &pb.Res{Message:"success!!"}, nil
}

func (fts *fileTransferService) TransferFile(stream pb.FileTransferService_TransferFileServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&pb.Res{Message:"file read done!"})
			return nil
		}
		if err != nil {
			return err
		}
		log.Println(string(req.Data))
	}
}

func main() {
	host := flag.String("host", "localhost", "hostname or ip")
	port := flag.String("port", "12345", "port number")
	flag.Parse()
	lis, err := net.Listen("tcp", *host + ":" + *port)
	if err != nil {
		log.Fatalf("failed to listen on host:%s port:%s\n", *host, *port)
	}
	log.Printf("start listening on host:%s port:%s\n", *host, *port)
	grpcServer := grpc.NewServer()
	pb.RegisterFileTransferServiceServer(grpcServer, &fileTransferService{})
	log.Printf("start grpcServer!!\n")
	grpcServer.Serve(lis)
}
