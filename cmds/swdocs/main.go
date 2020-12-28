package main

/*
* TODO handlers: Do not return the raw error as it can expose backend sensitive data
* TODO Write tests
* TODO Add comments / godoc
* TODO CRUD on docs
* TODO Equivalent of apply of kubernetes with a yml that people can get a doc, edit and apply or just apply to create
* TODO Versions of docs to allow for rolling back
* TODO Style the app with template inheritance (header/footer) and the /$SwDoc page
* TODO Implement a simple search functionality
 */

import (
	"flag"
	"fmt"
	"os"

	"github.com/andrecp/swdocs"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above if unset.
	loglevel, err := log.ParseLevel(os.Getenv("SWDOCS_LOGLEVEL"))
	if err != nil {
		log.SetLevel(log.WarnLevel)
	}
	log.SetLevel(loglevel)
}

func main() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	nameCreateCmd := createCmd.String("name", "", "The name of what you're documenting! Goes into the URL -> /name")

	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	linkCmd := flag.NewFlagSet("link", flag.ExitOnError)
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Missing subcommand, use one of: `create`, `edit`, `delete`, `link` or `serve`.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		createCmd.Parse(os.Args[2:])
		if *nameCreateCmd == "" {
			fmt.Println("--name is required when creating a new swdoc")
			os.Exit(1)
		}
	case "edit":
		editCmd.Parse(os.Args[2:])
	case "delete":
		deleteCmd.Parse(os.Args[2:])
	case "link":
		linkCmd.Parse(os.Args[2:])
	case "serve":
		serveCmd.Parse(os.Args[2:])
		a := swdocs.App{}
		a.Initialize(os.Getenv("SWDOCS_DBPATH"))

		a.Run(os.Getenv("SWDOCS_PORT"))

	default:
		fmt.Println("Must use one of the subcommands `create`, `edit`, `delete`, `link` or `serve`.")
		os.Exit(1)
	}

}
