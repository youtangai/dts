package server

import (
	"io"
	"os"
	"path/filepath"

	"github.com/youtangai/fts/lib/errors"
	"github.com/youtangai/fts/lib/status"
	pb "github.com/youtangai/fts/proto"
)

const (
	//DefFilePerm is
	DefFilePerm = 0755
)

//FileTransferServer is
type FileTransferServer struct {
	Dir string
}

//TransferFile is
func (s *FileTransferServer) TransferFile(stream pb.FileTransferService_FileTransferServer) error {
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
