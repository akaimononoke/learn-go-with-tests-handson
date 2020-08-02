package main

import (
	"log"
	"net/http"

	"github.com/akaimononoke/learn-go-with-tests-handson/application/poker"
)

const dbFileName = "game.db.json"

func main() {
	store, closeDB, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer closeDB()

	server := poker.NewPlayerServer(store)

	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatalf("could not listen on port 8080 %v", err)
	}
}
