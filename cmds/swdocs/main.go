package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"

	"github.com/andrecp/swdocs"

	log "github.com/sirupsen/logrus"
)

const (
	httpAddress = "http://localhost:8087"
)

type (
	applyRequest struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		User        string          `json:"user"`
		Sections    json.RawMessage `json:"sections"`
	}
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Display Debug Level severity log or above if unset.
	loglevel, err := log.ParseLevel(os.Getenv("SWDOCS_LOGLEVEL"))
	if err != nil {
		log.SetLevel(log.DebugLevel)
	}
	log.SetLevel(loglevel)
}

func main() {

	applyCmd := flag.NewFlagSet("apply", flag.ExitOnError)
	userApplyCmd := applyCmd.String("user", "", "Override the user, useful for CI")
	filePathApplyCmd := applyCmd.String("file", "", "The JSON file you want to apply the changes from")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	nameDeleteCmd := deleteCmd.String("name", "", "The name of the SwDoc you want to delete")

	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Missing subcommand, use one of: `apply`, `delete` or `serve`.")
		os.Exit(1)
	}
	switch os.Args[1] {
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

		var username string
		if *userApplyCmd == "" {
			user, err := user.Current()
			if err != nil {
				log.Fatal(err.Error())
			}
			username = user.Username
		} else {
			username = *userApplyCmd
		}

		jsonText, err := ioutil.ReadFile(*filePathApplyCmd)
		if err != nil {
			log.Fatal(err.Error())
		}

		r := applyRequest{}
		err = json.Unmarshal(jsonText, &r)
		if err != nil {
			log.Fatal(err.Error())
		}

		r.User = username

		requestBody, err := json.Marshal(r)
		if err != nil {
			log.Fatal(err.Error())
		}

		resp, err := http.Post(httpAddress+"/api/v1/swdocs/apply", "application/json", bytes.NewBuffer(requestBody))
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

	case "serve":
		serveCmd.Parse(os.Args[2:])
		a := swdocs.App{}
		a.Initialize(os.Getenv("SWDOCS_DBPATH"))

		a.Run(os.Getenv("SWDOCS_PORT"))

	default:
		fmt.Println("Must use one of the subcommands `apply`, `delete` or `serve`.")
		os.Exit(1)
	}

}
