package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/youtangai/fts/lib/errors"
)

const (
	ArgNum = 1
)

func main() {
	app := cli.NewApp()
	app.Name = "fts"
	app.Commands = []cli.Command{
		clientCommand,
		serverCommand,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
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
