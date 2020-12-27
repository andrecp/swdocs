# swdocs
Aggregate all of your software docs easily!


## Usage

```bash

# Start the mdoc server
swdocs serve --port 8080

# Create a new mdoc page about the zeromq software package.
swdocs create zeromq --title "A nice abstraction for sockets" --description "ZeroMQ is a high-performance asynchronous messaging library, aimed at use in distributed or concurrent applications. It provides a message queue, but unlike message-oriented middleware, a ZeroMQ system can run without a dedicated message broker. "

# Edit some metadata
swdocs edit zeromq 

# Link to some documents
swdocs link zeromq "http://confluence.com/ZeroMQ+Internal+Guide" "Internal docs" "Our internal page"
swdocs link zeromq "https://zeromq.org/" "Upstream site" "Homepage"
swdocs link zeromq "https://zguide.zeromq.org/" "Upstream site" "A good guide for developing zeromq"
```
