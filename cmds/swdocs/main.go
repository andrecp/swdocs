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
	"strings"

	"github.com/andrecp/swdocs"

	log "github.com/sirupsen/logrus"
)

const (
	// Customizable via envvars.
	defaultPort          = "8087"
	defaultHTTPAddress   = "http://localhost"
	defaultTemplatesPath = "."
	defaultDbPath        = "swdev.sqlite"
	defaultLogLevel      = log.WarnLevel

	// Other constants
	subCommandHelp = `Missing or unsupported subcommand! You can use:
  * swdocs bootstrap mysoftware    # Creates a mysoftware.json to be modified and used with apply
  * swdocs apply mysoftware.json   # To create or update a swdoc for mysoftware
  * swdocs get mysoftware          # To get info about a swdoc called mysoftware
  * swdocs delete mysoftware       # To delete a swdoc called mysoftware
  * swdocs list                    # To list available swdocs, use --filter to filter.
  * swdocs serve                   # To run the swdoc server

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
		log.SetLevel(defaultLogLevel)
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
	bootstrapCmd := flag.NewFlagSet("bootstrap", flag.ExitOnError)

	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	fmtGetCmd := getCmd.String("format", "human", "The format of the output, options are 'json' and 'human'")

	applyCmd := flag.NewFlagSet("apply", flag.ExitOnError)
	userApplyCmd := applyCmd.String("user", "", "Override the user, useful for CI")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	filterListCmd := listCmd.String("filter", "%", "Filter by name, % is a wildcard.")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)

	httpAddr := os.Getenv("SWDOCS_HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = defaultHTTPAddress
	}

	port := os.Getenv("SWDOCS_PORT")
	if port == "" {
		port = defaultPort
	}

	baseURL := httpAddr + ":" + port

	// Verify we have at least one subcommand.
	if len(os.Args) < 2 {
		fmt.Println(subCommandHelp)
		os.Exit(1)
	}

	// Call the right subcommand.
	switch os.Args[1] {
	case "bootstrap":
		bootstrapCmd.Parse(os.Args[2:])
		err := bootstrapCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err.Error())
		}

		name := bootstrapCmd.Arg(0)
		if name == "" {
			fmt.Println("Must give one arg, SwDoc name, to bootstrap the file")
			os.Exit(1)
		}

		jsonBootstrap := `
	{
		"name": "$NAME",
		"description": "Describe $NAME here in a succint manner",
		  "sections": [
			{
				"header": "Each section of $NAME has a header to group similar links",
				"links": [
					{
						"url": "https://lmgtfy.app/?q=$NAME",
						"description": "Search google for $NAME"
					}
				]
			}
		  ]
	}`
		fileContents := strings.ReplaceAll(jsonBootstrap, "$NAME", name)
		fileName := name + ".json"
		ioutil.WriteFile(fileName, []byte(fileContents), 0644)
		fmt.Println(fileName + " created.")

	case "get":
		getCmd.Parse(os.Args[2:])
		err := getCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err.Error())
		}

		name := getCmd.Arg(0)
		if name == "" {
			fmt.Println("A name arg is required to get a SwDocs")
			os.Exit(1)
		}

		resp, err := http.Get(baseURL + "/api/v1/swdocs/" + name)
		if err != nil {
			log.Fatal(err.Error())
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err.Error())
		}

		if resp.StatusCode != 200 {
			fmt.Println(string(body))
			os.Exit(1)
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

		applyFilePath := applyCmd.Arg(0)
		if applyFilePath == "" {
			fmt.Println("A file path to a JSON is required as an argument to be apply.")
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

		jsonText, err := ioutil.ReadFile(applyFilePath)
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

		resp, err := http.Post(baseURL+"/api/v1/swdocs/apply", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Fatal(err.Error())
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println(string(body))

	case "list":
		listCmd.Parse(os.Args[2:])
		err := listCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err.Error())
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", baseURL+"/api/v1/swdocs/", nil)
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
			fmt.Println(swdoc.Name + " -> " + defaultHTTPAddress + "/" + swdoc.Name)
		}

	case "delete":
		err := deleteCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err.Error())
		}

		name := deleteCmd.Arg(0)
		if name == "" {
			fmt.Println("A name arg is required to delete a SwDocs")
			os.Exit(1)
		}

		client := &http.Client{}
		req, err := http.NewRequest("DELETE", baseURL+"/api/v1/swdocs/"+name, nil)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err.Error())
		}

		if resp.StatusCode != 200 {
			fmt.Println("Something went wrong, check server logs.")
			os.Exit(1)
		}
		fmt.Println("Ok.")

	case "serve":
		serveCmd.Parse(os.Args[2:])

		dbPath := os.Getenv("SWDOCS_DB_PATH")
		if dbPath == "" {
			dbPath = defaultDbPath
		}

		templatesPath := os.Getenv("SWDOCS_TEMPLATES_PATH")
		if templatesPath == "" {
			templatesPath = defaultTemplatesPath
		}

		// Create, initialize and run the app.
		c := swdocs.AppConfig{
			Port:          port,
			DbPath:        dbPath,
			TemplatesPath: templatesPath,
		}
		a := swdocs.App{Config: c}
		a.Initialize()
		a.Run()

	default:
		fmt.Println(subCommandHelp)
		os.Exit(1)
	}

}
