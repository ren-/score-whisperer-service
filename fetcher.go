package main

// Ideas:
// Calculate the standard deviation of top scores
// TOP pp, lowest PP
// Highest Accuracy, lowest accuracy
import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	API "github.com/ren-/osu/api"
)

// APIConnection stores the osu! API key and provides methods to consume API endpoints.
var APIConnection API.Config
var db *sqlx.DB

func main() {
	var err error
	err = APIConnection.SetAPIKey(os.Getenv("OSU_TOKEN"))

	if err != nil {
		fmt.Println(err)
	}

	db, err = sqlx.Connect("postgres", "host="+os.Getenv("DB_HOST")+" user="+os.Getenv("DB_USER")+" dbname="+os.Getenv("DB_DATABASE")+" sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	// Register messageCreate as a callback for the messageCreate events.

	fmt.Println("Service is now running.  Press CTRL-C to exit.")

	// start fetching top 300 players
	players := make(chan []string)
	sem := make(chan bool, 80)
	go func() {
		go getTopPlayersForCountry(300, "LT", players)
		var fetchedPlayers []string
		for {
			fetchedPlayers = <-players
			for _, element := range fetchedPlayers {

				go storeRecentPlays(element, sem)
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}
