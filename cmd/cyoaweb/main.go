package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	cyoa "github.com/Ed-cred/choose_adventure"
)

func main() {
	filename := flag.String("file", "gopher.json", "The JSON file to read for the story.")
	flag.Parse()
	log.Println("Using the story from", *filename)
	f, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	story := make(map[string]cyoa.Chapter)
	err = json.NewDecoder(f).Decode(&story)
	if err != nil {
		log.Fatal("failed to decode json", err)
	}
	log.Printf("%+v", story)

	// f, err := os.ReadFile(*file)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
