# swdocs
Aggregate all of your software docs easily!


## Usage

```bash

Create a .env file

# > cat .dev.env
export SWDOCS_DBPATH=/tmp/swdocs.db
export SWDOCS_PORT=8087

# Start the swdocs server
make run

# Or
source .dev.env && swdocs serve

# Create a new mdoc page about the zeromq software package.
swdocs create --name zeromq --description "ZeroMQ is a high-performance asynchronous messaging library, aimed at use in distributed or concurrent applications. It provides a message queue, but unlike message-oriented middleware, a ZeroMQ system can run without a dedicated message broker. "

# Edit some metadata
swdocs edit zeromq 
```

### Working with sqlite

```sql
> sqlite3 /tmp/swdocs.db
>> .tables
>> .schema swdocs
```
