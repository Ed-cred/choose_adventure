package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	cyoa "github.com/Ed-cred/choose_adventure"
)

func main() {
	port := flag.Int("port", 3000, "The port to start the CYOA web app on")
	filename := flag.String("file", "gopher.json", "The JSON file to read for the story.")
	flag.Parse()
	log.Println("Using the story from", *filename)
	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	story, err := cyoa.JsonStoryFromFile(f)
	if err != nil {
		log.Fatal(err)
	}

	h := cyoa.NewHandler(story, nil)
	log.Printf("Strting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
