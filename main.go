package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	"github.com/hashicorp/hcl/hcl/ast"

	"github.com/hashicorp/hcl"
	"github.com/micro/cli"
)

var app = cli.NewApp()

var files = []string{}

func command() {
	app.Commands = []cli.Command{
		{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "Specify the HCL file you would like to analyze.",
			Action: func(c *cli.Context) {
				fmt.Println("specified file: ", c.Args().First())
				return
			},
		},
	}
}

func info() {
	app.Name = "HCL Sentinel CLI"
	app.Usage = "HCL file analysis for identifying common issues."
	app.Author = "Matt Williams"
	app.Version = "0.0.1"
}

func main() {

	command()
	info()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	dat, err := ioutil.ReadFile("./test.tf")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(dat))

	parsedThing, err := hcl.ParseBytes(dat)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(reflect.TypeOf(parsedThing.Node))
	switch n := parsedThing.Node.(type) {
	case *ast.ObjectList:
		for i, item := range n.Items {
			fmt.Print(reflect.TypeOf(n.Items[i]))
			fmt.Print(reflect.TypeOf(item))
		}
	}
}
