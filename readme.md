# Chirpy

This is a project I am working on during the Learn HTTP Servers in Go course on Boot.dev

While this is not yet a production ready project I hope it can serve as a way to showcase some of the knowledge I have gained in the Go programming language and HTTP Servers.

## Set up

To run this project the following are needed:

- A recent version of Go
- The Goose migration tool (https://github.com/pressly/goose)
- Postgres
- The .env example file copied to .env and filled with appropriate data

Run the following script with your db string to get the app up and running:

```bash
cd /sql/schema/
// Use your own db string or this one when running postgres with default settings
goose postgres postgres://postgres:postgres@localhost:5432/chirpy up
cd -
go mod tidy
go run .
```
