package main

/* TODOS:
* handlers: Do not return the raw error as it can expose backend sensitive data
* Write tests
* Add comments / godoc
* CRUD on docs
* Equivalent of apply of kubernetes with a yml that people can get a doc, edit and apply or just apply to create
* Style the app with template inheritance (header/footer) and the /$SwDoc page
* Implement a simple search functionality
* The templates folder at runtime need to be configurable
* It is erroring silently when no .env file
* Test multiple writes, might need to tweak a bit sqlite (maximum of 1 conn, higher timeout) as per their docs, or, add a mux for writes.
 */

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/andrecp/swdocs"

	log "github.com/sirupsen/logrus"
)

const (
	httpAddress = "http://localhost:8087"
)

type (
	CreateRequest struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Sections    json.RawMessage `json:"sections"`
	}
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
	descriptionCreateCmd := createCmd.String("description", "", "A description of what it is.")
	sectionsCreateCmd := createCmd.String("sections", "", "JSON with the value, for example, '[{\"header\":\"Dashboards\",\"links\":[{\"url\":\"http://kibana.domain.com:5601\",\"description\":\"Kibana boards\"}]}]'")

	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Missing subcommand, use one of: `create`, `edit`, `delete`, or `serve`.")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "create":
		err := createCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err.Error())
		}

		if *nameCreateCmd == "" {
			fmt.Println("--name is required when creating a new swdoc")
			os.Exit(1)
		}
		requestBody, err := json.Marshal(CreateRequest{
			Name:        *nameCreateCmd,
			Description: *descriptionCreateCmd,
			Sections:    json.RawMessage(*sectionsCreateCmd),
		})
		if err != nil {
			log.Fatal(err.Error())
		}
		resp, err := http.Post(httpAddress+"/api/v1/swdocs", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Fatal(err.Error())
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Info(string(body))

	case "edit":
		editCmd.Parse(os.Args[2:])
	case "delete":
		deleteCmd.Parse(os.Args[2:])
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
