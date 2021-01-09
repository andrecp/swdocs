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
	httpAddress    = "http://localhost:8087"
	subCommandHelp = `Missing or unsupported subcommand! You can use:
  * swdocs apply -f file.json to create or update an entry
  * swdocs get $swdocname to get info about a swdoc
  * swdocs list to get available swdocs
  * swdocs delete $swdocname to delete a swdoc
  * swdocs serve to run the swdoc server

Every subcommand supports --help.
	`
)

func init() {
	// Log as JSON instead of the default ASCII formatter
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Display Warn level log if unset
	loglevel := os.Getenv("SWDOCS_LOGLEVEL")
	if loglevel == "" {
		log.SetLevel(log.WarnLevel)
	} else {
		loglevel, err := log.ParseLevel(loglevel)
		if err != nil {
			panic(err)
		}
		log.SetLevel(loglevel)
	}

}

func main() {

	// Declare command line subcommands and options.
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	nameGetCmd := getCmd.String("name", "", "The name of the SwDoc you want to get")
	fmtGetCmd := getCmd.String("format", "human", "The format of the output, options are 'json' and 'human'")
	urlGetCmd := getCmd.String("url", httpAddress, "The URL to make the request to, in format of http://NAME:PORT")

	applyCmd := flag.NewFlagSet("apply", flag.ExitOnError)
	userApplyCmd := applyCmd.String("user", "", "Override the user, useful for CI")
	filePathApplyCmd := applyCmd.String("file", "", "The JSON file you want to apply the changes from")
	urlApplyCmd := applyCmd.String("url", httpAddress, "The URL to make the request to, in format of http://NAME:PORT")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	filterListCmd := listCmd.String("filter", "%", "Filter by name, % is a wildcard.")
	urlListrCmd := listCmd.String("url", httpAddress, "The URL to make the request to, in format of http://NAME:PORT")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	nameDeleteCmd := deleteCmd.String("name", "", "The name of the SwDoc you want to delete")
	urlDeleteCmd := deleteCmd.String("url", httpAddress, "The URL to make the request to, in format of http://NAME:PORT")

	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)

	// Verify we have at least one subcommand.
	if len(os.Args) < 2 {
		fmt.Println(subCommandHelp)
		os.Exit(1)
	}

	// Call the right subcommand.
	switch os.Args[1] {
	case "get":
		getCmd.Parse(os.Args[2:])
		err := getCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err.Error())
		}

		if *nameGetCmd == "" {
			fmt.Println("--name is required to get a SwDocs")
			os.Exit(1)
		}

		resp, err := http.Get(*urlGetCmd + "/api/v1/swdocs/" + *nameGetCmd)
		if err != nil {
			log.Fatal(err.Error())
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err.Error())
		}

		r := swdocs.SwDoc{}
		err = json.Unmarshal(body, &r)
		if err != nil {
			log.Fatal(err.Error())
		}

		if *fmtGetCmd == "human" {
			fmt.Println("Name: " + string(r.Name))
			fmt.Println("Description: " + string(r.Description))
			fmt.Println("Last updated by: " + string(r.User))
			fmt.Println("Last updated on: " + r.Updated.ToString())
			fmt.Println("")
			for _, section := range r.Sections {
				fmt.Println(section.Header)
				for _, link := range section.Links {
					fmt.Println(" * " + link.Description + " (" + link.URL + ")")
				}

			}
		} else if *fmtGetCmd == "json" {
			json, err := json.MarshalIndent(r, "", "  ")
			if err != nil {
				log.Fatal(err.Error())
			}
			fmt.Println(string(json))
		} else {
			fmt.Println("Unsupported format, options are 'json' and 'human'")
			os.Exit(1)
		}

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

		r := swdocs.SwDoc{}
		err = json.Unmarshal(jsonText, &r)
		if err != nil {
			log.Fatal(err.Error())
		}

		r.User = username

		requestBody, err := json.Marshal(r)
		if err != nil {
			log.Fatal(err.Error())
		}

		resp, err := http.Post(*urlApplyCmd+"/api/v1/swdocs/apply", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Fatal(err.Error())
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Info(string(body))

	case "list":
		listCmd.Parse(os.Args[2:])
		err := listCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err.Error())
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", *urlListrCmd+"/api/v1/swdocs/", nil)
		if err != nil {
			log.Fatal(err.Error())
		}

		q := req.URL.Query()
		q.Add("filter", *filterListCmd)
		req.URL.RawQuery = q.Encode()

		resp, err := client.Do(req)
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err.Error())
		}

		r := []swdocs.SwDoc{}
		err = json.Unmarshal(body, &r)
		if err != nil {
			log.Fatal(err.Error())
		}
		for _, swdoc := range r {
			fmt.Println(swdoc.Name)
		}

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
		req, err := http.NewRequest("DELETE", *urlDeleteCmd+"/api/v1/swdocs/"+*nameDeleteCmd, nil)
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
		fmt.Println(subCommandHelp)
		os.Exit(1)
	}

}
