# swdocs
Aggregate all of your software docs easily!


## Usage

```bash

# Create the .env file
> cat .dev.env
SWDOCS_DBPATH=/tmp/swdocs.db
SWDOCS_PORT=8087
SWDOCS_LOGLEVEL=debug

# Start the swdocs server and go to http://localhost:8087
> make run  # or source .dev.env && swdocs serve

# Apply either creates or updates an entry from a JSON file.
> swdocs apply --file tests/rabbitmq.json
```

### Working with sqlite

```sql
> sqlite3 /tmp/swdocs.db
>> .tables
>> .schema swdocs
```
