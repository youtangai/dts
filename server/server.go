package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/youtangai/fts/lib/errors"
	"github.com/youtangai/fts/lib/status"
	pb "github.com/youtangai/fts/proto"
	"google.golang.org/grpc"
)

const (
	//DefFilePerm is
	DefFilePerm = 0755
)

//FileTransferServer is
type FileTransferServer struct {
	Dir string
}

func NewFileTransferServer(dir string) FileTransferServer {
	return FileTransferServer{
		Dir: dir,
	}
}

//TransferFile is
func (s *FileTransferServer) FileTransfer(stream pb.FileTransferService_FileTransferServer) error {
	for {
		fileData, err := stream.Recv()
		if err == io.EOF {
			err := stream.SendAndClose(&pb.Response{Status: status.NoError})
			if err != nil {
				return errors.StreamCloseError(err)
			}
			return nil
		}
		if err != nil {
			return errors.StreamRecvError(err)
		}
		filePath := filepath.Join(s.Dir, fileData.Filename)
		mode := os.FileMode(fileData.Mode)
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDONLY, mode)
		if err != nil {
			return errors.OpenFileError(filePath, err)
		}

		_, err = file.Write(fileData.Data)
		if err != nil {
			return errors.FileWriteError(err)
		}
	}
}

func (s FileTransferServer) Run(url, dir string) error {
	lis, err := net.Listen("tcp", url)
	if err != nil {
		return errors.ListenError(url, err)
	}

	server := grpc.NewServer()
	pb.RegisterFileTransferServiceServer(server, &s)

	if err := server.Serve(lis); err != nil {
		return errors.GrpcServeError(err)
	}

	return nil
}

func GetURL(host, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}
