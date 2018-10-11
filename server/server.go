package server

import (
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/youtangai/fts/lib/errors"
	"github.com/youtangai/fts/lib/logging"
	"github.com/youtangai/fts/lib/status"
	"github.com/youtangai/fts/lib/util"
	pb "github.com/youtangai/fts/proto"
	"google.golang.org/grpc"
)

//FileTransferServer is
type FileTransferServer struct {
	Dir string
	URL string
}

func NewFileTransferServer(dir, host, port string) FileTransferServer {
	return FileTransferServer{
		Dir: dir,
		URL: util.GetURL(host, port),
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
		logging.RecvFileData(fileData)

		dirPath := filepath.Join(s.Dir, fileData.Dir)
		if !util.IsExistDir(dirPath) {
			if err = os.Mkdir(dirPath, 0755); err != nil {
				return errors.MkdirError(dirPath, err)
			}
		}

		filePath := filepath.Join(dirPath, fileData.Filename)
		mode := os.FileMode(fileData.Mode)
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, mode)
		if err != nil {
			return errors.OpenFileError(filePath, err)
		}
		logging.OpenFileInfo(filePath)

		_, err = file.Write(fileData.Data)
		if err != nil {
			return errors.FileWriteError(err)
		}
	}
}

func (s FileTransferServer) Run() error {
	url := s.URL

	lis, err := net.Listen("tcp", url)
	if err != nil {
		return errors.ListenError(url, err)
	}

	server := grpc.NewServer()
	pb.RegisterFileTransferServiceServer(server, &s)

	logging.WillStartServerInfo(s.URL)
	if err := server.Serve(lis); err != nil {
		return errors.GrpcServeError(err)
	}

	return nil
}
