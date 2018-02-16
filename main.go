package main

import (
	"log"
	"os"
)

func main() {
	app, err := newApp(os.Args, os.Stdout, os.Stderr)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	res := app.run()
	log.Printf("[%s] OK?: %v (errors: %v)", app.site.Name, res.ok(), res.errors())
	os.Exit(res.statusCode)
}
