package main

import (
	"fmt"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
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
	fmt.Printf("%+v\n\n", parsedThing)
	fmt.Println(" ------------------------ ")
	fmt.Println("|  Checking Comments... | ")
	fmt.Println(" ------------------------ ")
	for _, cGroup := range parsedThing.Comments {
		for _, cGroupComment := range cGroup.List {
			if strings.ToLower(cGroupComment.Text) == "# hello sentinel" {
				color.Green("+ Howdy little terraform script :)")
			}
		}
	}
	fmt.Print("\u2714 Check complete\n\n")

	switch n := parsedThing.Node.(type) {
	case *ast.ObjectList:
		for _, item := range n.Items {
			fmt.Println(reflect.TypeOf(item.Val))
			for _, k := range item.Keys {
				fmt.Printf("%+v\n", k)
				if k.Token.Type.String() == token.IDENT.String() && k.Token.Text == "locals" {
					fmt.Println("FOUND LOCAL VARIABLES DECL")
					checkLocalVariables(item)
				}
			}
			switch x := item.Val.(type) {
			case *ast.ObjectType:
				for _, oItem := range x.List.Items {
					for _, oItemKey := range oItem.Keys {
						fmt.Printf("%+v\n", oItemKey)
					}
				}
			}
		}
	default:
		fmt.Printf("unknown type: %T\n", n)
	}
}

func checkLocalVariables(objectItem *ast.ObjectItem) {
	fmt.Println(" ------------------------------- ")
	fmt.Println("|  Checking Local Variables... | ")
	fmt.Println(" ------------------------------- ")
	for _, key := range objectItem.Keys {
		fmt.Println(key.Token.String())
	}
}
