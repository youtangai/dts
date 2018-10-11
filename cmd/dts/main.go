package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/youtangai/dts/lib/errors"
)

var (
	version = "1.0.0"
)

const (
	ArgNum = 1
	usage  = `dts cli|srv -host=<IPADDR> -port=<PORTNUMBER> <DIR>`
)

func main() {
	app := cli.NewApp()
	app.Name = "dts"
	app.Version = version
	app.Usage = usage
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
