* handlers: Do not return the raw error as it can expose backend sensitive data
* Write tests
* Add comments / godoc
* CRUD on docs
* Equivalent of apply of kubernetes with a yml that people can get a doc, edit and apply or just apply to create
* Style the app with template inheritance (header/footer) and the /$SwDoc page
* The templates folder at runtime need to be configurable
* It is erroring silently when no .env file
* Test multiple writes, might need to tweak a bit sqlite (maximum of 1 conn, higher timeout) as per their docs, or, add a mux for writes.
