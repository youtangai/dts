package client

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/youtangai/dts/lib/errors"
	"github.com/youtangai/dts/lib/util"
	pb "github.com/youtangai/dts/proto"
	"google.golang.org/grpc"
)

//Client is file transfer service's client
type Client struct {
	Client pb.FileTransferServiceClient //grpc client
	Dir    string                       //directory where file are atored
}

const (
	//MaxBuff is max buffer file transfer
	MaxBuff = 4096
)

func NewClient(dir, host, port string) (Client, error) {
	url := util.GetURL(host, port)
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return Client{}, errors.GrpcDialError(err)
	}

	return Client{
		Client: pb.NewFileTransferServiceClient(conn),
		Dir:    dir,
	}, nil
}

func (c Client) TransferFiles() error {
	files, err := c.getAllFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		err := c.transferFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Client) transferFile(fileInfo os.FileInfo) error {
	file, err := c.getFile(fileInfo)
	if err != nil {
		return err
	}

	stream, err := c.Client.FileTransfer(context.Background())
	if err != nil {
		return errors.GetStreamError(err)
	}

	buff := make([]byte, MaxBuff)
	for {
		count, err := file.Read(buff)
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.FileReadError(err)
		}

		fileData := pb.FileData{
			Dir:      c.Dir,
			Filename: fileInfo.Name(),
			Mode:     int32(fileInfo.Mode()),
			Data:     buff[:count],
		}
		if err := stream.Send(&fileData); err != nil {
			return errors.StreamSendError(err)
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return errors.StreamCloseError(err)
	}

	return nil
}

func (c Client) getAllFiles() ([]os.FileInfo, error) {
	dir := c.Dir
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []os.FileInfo{}, errors.GetFilesError(dir, err)
	}
	return files, nil
}

func (c Client) getFilePath(filename string) string {
	return filepath.Join(c.Dir, filename)
}

func (c Client) getFile(fileInfo os.FileInfo) (*os.File, error) {
	fileName := fileInfo.Name()
	filePath := c.getFilePath(fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.OpenFileError(fileName, err)
	}
	return file, nil
}
