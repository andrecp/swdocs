package main

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
	nameCreateCmd := createCmd.String("name", "", "The name of the SwDoc you're creating! Goes into the URL -> /:name")
	descriptionCreateCmd := createCmd.String("description", "", "A description of the SwDoc is.")
	sectionsCreateCmd := createCmd.String("sections", "", "JSON with the value, for example, '[{\"header\":\"Dashboards\",\"links\":[{\"url\":\"http://kibana.domain.com:5601\",\"description\":\"Kibana boards\"}]}]'")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	nameDeleteCmd := deleteCmd.String("name", "", "The name of the SwDoc you want to delete")

	applyCmd := flag.NewFlagSet("apply", flag.ExitOnError)
	filePathApplyCmd := applyCmd.String("file", "", "The JSON file you want to apply the changes from")

	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Missing subcommand, use one of: `create`, `apply`, `delete`, or `serve`.")
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
	case "delete":
		err := deleteCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err.Error())
		}

		if *nameDeleteCmd == "" {
			fmt.Println("--name is required when deleting an existing swdoc")
			os.Exit(1)
		}

		client := &http.Client{}
		req, err := http.NewRequest("DELETE", httpAddress+"/api/v1/swdocs/"+*nameDeleteCmd, nil)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err.Error())
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Info(string(body))

	case "apply":
		applyCmd.Parse(os.Args[2:])
		err := applyCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err.Error())
		}

		if *filePathApplyCmd == "" {
			fmt.Println("--file is required to apply a file to SwDocs")
			os.Exit(1)
		}

		jsonData, err := ioutil.ReadFile(*filePathApplyCmd)
		if err != nil {
			log.Fatal(err.Error())
		}

		resp, err := http.Post(httpAddress+"/api/v1/swdocs/apply", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal(err.Error())
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Info(string(body))

	case "serve":
		serveCmd.Parse(os.Args[2:])
		a := swdocs.App{}
		a.Initialize(os.Getenv("SWDOCS_DBPATH"))

		a.Run(os.Getenv("SWDOCS_PORT"))

	default:
		fmt.Println("Must use one of the subcommands `create`, `apply`, `delete`, `link` or `serve`.")
		os.Exit(1)
	}

}
