package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/youtangai/fts/client"
	"github.com/youtangai/fts/lib/errors"
	"github.com/youtangai/fts/lib/logging"
	"github.com/youtangai/fts/server"
)

const (
	ArgNum = 1
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
				checkArg(ctx)

				host, port, dir := getHostPortDir(ctx)

				cli, err := client.NewClient(dir, host, port)
				if err != nil {
					panic(err)
				}

				err = cli.TransferFiles()
				if err != nil {
					panic(err)
				}

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
				checkArg(ctx)

				host, port, dir := getHostPortDir(ctx)

				if !isDirExistDir(dir) {
					logging.NoDirInfo(dir)
					if err := os.Mkdir(dir, 0755); err != nil {
						return errors.MkdirError(dir, err)
					}
				}

				path, err := filepath.Abs(dir)
				if err != nil {
					return err
				}

				server := server.NewFileTransferServer(path, host, port)
				server.Run()

				return nil
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func isDirExistDir(dir string) bool {
	_, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return true
}

func checkArg(ctx *cli.Context) {
	if ctx.NArg() < ArgNum {
		errors.NoMuchArgNumPanic(ctx.NArg())
	}
}

func getHost(ctx *cli.Context) string {
	return ctx.String("host")
}

func getPort(ctx *cli.Context) string {
	return ctx.String("port")
}

func getDir(ctx *cli.Context) string {
	return ctx.Args().Get(0)
}

func getHostPortDir(ctx *cli.Context) (string, string, string) {
	return getHost(ctx), getPort(ctx), getDir(ctx)
}
