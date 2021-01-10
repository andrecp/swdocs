# swdocs
Aggregate all of your software docs easily!

* Very simple to get up and running, it is a single binary;
* Run for yourself in your machine or in a server for your company;
* Offers a web ui (without javascript) and a CLI to navigate your docs;
* Have all the relevant links to your code base *with* your code base;
* Edit the links from where you edit your codem, let your CI system update your docs when you push new code.

## Install

```bash
> go get -u github.com/andrecp/swdocs/cmd/swdocs
```

You can also download this repo and run `make build`.

## Screenshots

This is what the UI looks like with a single swdoc on it from the [tests](tests/rabbitmq.json)

<img src="images/swdocs-home.png" width="800" /> ![The home page](images/swdocs-home.png)
<img src="images/swdoc-page.png" width="800" /> ![Visualizing a SwDoc for rabbitmq](images/swdoc-page.png)

## Usage

You can use SwDocs from a browser or from the CLI. You start by running the server and creating some docs through its `apply` CLI and then you can visualize and search from your favorite browser or through CLI operations like `get` and `list`.

### Running the server

```bash

# Create the .env file
> cat .dev.env
SWDOCS_DBPATH=/tmp/swdocs.db
SWDOCS_PORT=8087
SWDOCS_LOGLEVEL=debug

# Start the swdocs server and go to http://localhost:8087
> make run  # or source .dev.env && swdocs serve

```

### Creating and updating a SwDoc

```bash
# Apply either creates or updates an entry from a JSON file.
> swdocs apply --file tests/rabbitmq.json

# By default apply will use the process owner username, but, you can override it too for CIs (say from jenkins to parse the commit owner username)
> swdocs apply --file tests/rabbitmq.json --user ken
```

### Getting and listing SwDocs

```bash

# List every SwDoc in the database.
> swdocs list

# List every SwDoc in the database containing the word rabbit.
> swdocs list --filter rabbit%

# Get the URLs for a swdoc from the terminal
> swdocs get --name rabbitmq

# Or the JSON
> swdocs get --name rabbitmq --format json
```

### Deleting a SwDoc

```
> swdocs delete --name rabbitmq
```

## SwDoc definition

The structs are defined in [model.go](model.go), an example of a JSON to be inserted is

```json
{
    "name": "rabbitmq",
    "description": "A broker for your messages! AMQP!",
      "sections": [
        {
            "header": "Guides",
            "links": [
                {
                    "url": "https://www.rabbitmq.com/getstarted.html",
                    "description": "Official get started, great guide!"
                }
            ]
        }
      ]
  }
```

## Working with sqlite

The database gets created the first time the program runs.

```sql
> sqlite3 /tmp/swdocs.sqlite
>> .tables
>> .schema swdocs
```

You likely want to backup your .sqlite file periodically!
