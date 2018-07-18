package main

import (
	pb "github.com/youtangai/zanshin/proto"
	"io"
	"log"
	"net"
	"flag"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"os"
)

const (
	DEF_FILE_PERM = 0755
)

type fileTransferService struct {
	folderName string
	fileName string
	mode uint32
}

func (fts *fileTransferService) GetFolderInfo(ctx context.Context, info *pb.FolderInfo) (*pb.Res, error) {
	fts.folderName = info.Name
	log.Printf("recieve folder name: %s", info.Name)
	if err := os.Mkdir(info.Name, DEF_FILE_PERM); err != nil {
		return &pb.Res{Message: "failed to create dir"}, err
	}
	return &pb.Res{Message:"success"}, nil
}

func (fts *fileTransferService) GetFileInfo(ctx context.Context, info *pb.FileInfo) (*pb.Res, error) {
	fts.fileName = info.Name
	fts.mode = info.Mode
	log.Printf("recieve file info{ name:%s, mode:%d}\n", info.Name, info.Mode)
	return &pb.Res{Message:"success!!"}, nil
}

func (fts *fileTransferService) TransferFile(stream pb.FileTransferService_TransferFileServer) error {
	file, err := os.OpenFile(fts.folderName +"/"+ fts.fileName, os.O_WRONLY | os.O_CREATE, os.FileMode(fts.mode))
	if err != nil {
		log.Fatalf("cannot open file: %v\n", err)
	}
	defer file.Close()
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
		file.Write(req.Data)
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
