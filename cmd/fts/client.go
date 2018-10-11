package main

import (
	"github.com/urfave/cli"
	"github.com/youtangai/fts/client"
)

var clientCommand = cli.Command{
	Name:  "cli",
	Usage: "transfer files to fts server",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "host", Value: "127.0.0.1", Usage: "fts server ip"},
		cli.StringFlag{Name: "port", Value: "5050", Usage: "fts server port"},
	},
	ArgsUsage: "<dir> transfer dir path",
	Action:    execClient,
}

func execClient(ctx *cli.Context) error {
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
}
