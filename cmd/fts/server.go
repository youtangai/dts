package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/youtangai/fts/lib/errors"
	"github.com/youtangai/fts/lib/logging"
	"github.com/youtangai/fts/lib/util"
	"github.com/youtangai/fts/server"
)

var serverCommand = cli.Command{
	Name:  "srv",
	Usage: "start fts server",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "host", Value: "127.0.0.1", Usage: "fts server ip"},
		cli.StringFlag{Name: "port", Value: "5050", Usage: "fts server port"},
	},
	ArgsUsage: "<dir> recieve path",
	Action:    execServer,
}

func execServer(ctx *cli.Context) error {
	checkArg(ctx)

	host, port, dir := getHostPortDir(ctx)

	if !util.IsExistDir(dir) {
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
}
