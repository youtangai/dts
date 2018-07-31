package libfts

import (
	pb "github.com/youtangai/fts/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
	"github.com/urfave/cli"
	"path/filepath"
)

//Client is file transfer service's client
type Client struct {
	client pb.FileTransferServiceClient
}

const (
	//MaxBuff is max buffer file transfer 
	MaxBuff = 16384
)

func (c *Client)runGetFolderInfo(folderName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	folderInfo := &pb.FolderInfo{
		Name: folderName,
	}
	res, err := c.client.GetFolderInfo(ctx, folderInfo)
	if err != nil {
		return err
	}
	log.Printf("return message: %s\n", res.Message)
	return nil
}

func (c *Client)runGetFileInfo(name string, size int64, mode uint32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fileInfo := &pb.FileInfo{
		Name: name,
		Mode: mode,
	}
	res, err := c.client.GetFileInfo(ctx, fileInfo)
	if err != nil {
		return err
	}
	log.Printf("return message: %s\n", res.Message)
	return nil
}

func (c *Client)runTransferFile(file *os.File) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := c.client.TransferFile(ctx)
	buff := make([]byte, MaxBuff)
	for {
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


//TransferDir
func TransferDir(context *cli.Context) {
	host := context.String("host")
	port := context.String("port")
	dir := context.Args().Get(0)
    path, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal("failed to get abs path", path)
	}

	log.Println("request to", host, ":",  port)
	log.Println("dir:", dir, "path:", path)

	conn, err := grpc.Dial(host+":"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	ftsCli := &Client{
		client: pb.NewFileTransferServiceClient(conn),
	}
	
	ftsCli.runGetFolderInfo(dir)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("failed to read dir: %+v\n", err)
	}

	err = os.Chdir(dir)
	if err != nil {
		log.Fatalf("failed to change dir: %v\n", err)
	}

	for _, file := range files {
		ftsCli.runGetFileInfo(file.Name(), file.Size(), uint32(file.Mode()))
		f, err := os.Open(file.Name())
		if err != nil {
			log.Fatalf("cannot open file: %s err: %v\n", file.Name(), err)
		}
		ftsCli.runTransferFile(f)
	}
}
