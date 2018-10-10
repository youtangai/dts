package main

import (
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/youtangai/fts/libfts"
	pb "github.com/youtangai/fts/proto"
	"google.golang.org/grpc"
)

const (
	//CliArg is
	CliArg = 1
	//SrvArg is
	SrvArg = 1
)

func main() {
	app := cli.NewApp()
	app.Name = "fts"
	app.Commands = []cli.Command{
		{
			Name:  "cli",
			Usage: "transfer files to fts server",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "host", Value: "127.0.0.1", Usage: "fts server ip"},
				cli.StringFlag{Name: "port", Value: "5050", Usage: "fts server port"},
			},
			ArgsUsage: "<dir> transfer dir path",
			Action: func(ctx *cli.Context) error {
				//check arg num
				if ctx.NArg() < CliArg {
					log.Fatal("err: no much args num")
				}
				log.Println("call client command")
				libfts.TransferDir(ctx)
				return nil
			},
		},
		{
			Name:  "srv",
			Usage: "start fts server",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "host", Value: "127.0.0.1", Usage: "fts server ip"},
				cli.StringFlag{Name: "port", Value: "5050", Usage: "fts server port"},
			},
			ArgsUsage: "<dir> recieve path",
			Action: func(ctx *cli.Context) error {
				//check arg num
				if ctx.NArg() < SrvArg {
					log.Fatal("err: no much args num")
				}
				log.Println("call server command")
				host := ctx.String("host")
				port := ctx.String("port")
				dir := ctx.Args().Get(0)
				_, err := os.Stat(dir)
				if err != nil {
					if err := os.Mkdir(dir, 0755); err != nil {
						return err
					}
				}
				path, err := filepath.Abs(dir)
				if err != nil {
					return err
				}
				lis, err := net.Listen("tcp", host+":"+port)
				if err != nil {
					return err
				}
				log.Println("start listening on host:", host, "port:", port)
				log.Println("specify path:", path)
				grpcServer := grpc.NewServer()
				pb.RegisterFileTransferServiceServer(grpcServer, &libfts.FileTransferService{
					BaseDir: path,
				})
				log.Println("start grpcServer!!")
				if err := grpcServer.Serve(lis); err != nil {
					return err
				}
				return nil
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
