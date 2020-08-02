package main

import (
	"log"
	"net/http"
	"os"

	"github.com/akaimononoke/learn-go-with-tests-handson/application"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("opening %s %v", dbFileName, err)
	}

	store, err := application.NewFileSystemPlayerStore(db)
	if err != nil {
		log.Fatalf("problem creating file system player store: %v", err)
	}
	server := application.NewPlayerServer(store)

	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatalf("could not listen on port 8080 %v", err)
	}
}
