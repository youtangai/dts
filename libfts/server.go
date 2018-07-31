package libfts

import (
	pb "github.com/youtangai/fts/proto"
	"golang.org/x/net/context"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	//DefFilePerm is
	DefFilePerm = 0755
)

//FileTransferService is
type FileTransferService struct {
	BaseDir    string
	folderName string
	fileName   string
	mode       uint32
}

//GetFolderInfo is
func (fts *FileTransferService) GetFolderInfo(ctx context.Context, info *pb.FolderInfo) (*pb.Res, error) {
	log.Println("start GetFolderInfo. foldername:", info.Name)
	if err := os.Chdir(fts.BaseDir); err != nil {
		return &pb.Res{Message: "failed to change base dir"}, err 
	}
	fts.folderName = info.Name
	log.Println("recieve folder name:", info.Name)
	if err := os.Mkdir(info.Name, DefFilePerm); err != nil {
		return &pb.Res{Message: "failed to create dir"}, err
	}
	path := filepath.Join(fts.BaseDir, info.Name)
	if err := os.Chdir(path); err != nil {
		return &pb.Res{Message: "failed to change dir"}, err
	}
	return &pb.Res{Message: "success"}, nil
}

//GetFileInfo is
func (fts *FileTransferService) GetFileInfo(ctx context.Context, info *pb.FileInfo) (*pb.Res, error) {
	log.Println("start GetFileInfo")
	fts.fileName = info.Name
	fts.mode = info.Mode
	log.Printf("recieve file info{ name:%s, mode:%d}\n", info.Name, info.Mode)
	return &pb.Res{Message: "success!!"}, nil
}

//TransferFile is
func (fts *FileTransferService) TransferFile(stream pb.FileTransferService_TransferFileServer) error {
	log.Println("start TransferFile")
	file, err := os.OpenFile(fts.fileName, os.O_WRONLY|os.O_CREATE, os.FileMode(fts.mode))
	if err != nil {
		log.Fatalf("cannot open file: %v\n", err)
	}
	defer file.Close()
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			if err = stream.SendAndClose(&pb.Res{Message: "file read done!"}); err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			return err
		}
		if _, err = file.Write(req.Data); err != nil {
			return err
		}
	}
}
