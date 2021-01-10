# Must Have

# Should have
* Better request logging
* Write tests
* handlers: Do not return the raw error as it can expose backend sensitive data

# Nice to have
* HTTPS
* Authentication
* Date shouldn't be in UTC for the clients (CLI/browser), for the browser with no javascript!
* Include metadata for docs (like in kubernetes) and allow people to build their own filters/searches based on custom metadata
* Search improvements -- Indexes to improve the queries, do not do a like % by default if no filter param is given
