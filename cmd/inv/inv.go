package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Usage = "a kala cli for inventory"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Usage:  "use specified kala config",
			Value:  "~/.kala.toml",
			EnvVar: "KALA_CONFIG",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list the items within the given name or anchor",
			Action:  listCommand,
			Flags:   []cli.Flag{},
		},
	}

	app.Run(os.Args)
}

func Printlnf(f string, v ...interface{}) {
	fmt.Println(fmt.Sprintf(f, v...))
}
