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

	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.Alerter), store)
	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatalf("could not listen on port 8080 %v", err)
	}
}
