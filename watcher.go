package main

import (
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/ren-/osu/api"

	"github.com/PuerkitoBio/goquery"
)

func getTopPlayersForCountry(top int, country string, players chan []string) []string {
	for {
		var playersTemp []string
		pagesToScrape := int(top / 50)

		for page := 1; page <= pagesToScrape; page++ {

			doc, err := goquery.NewDocument("https://osu.ppy.sh/p/pp/?c=" + country + "&m=0&s=3&o=1&f=&page=" + strconv.Itoa(page))
			if err != nil {
				log.Fatal(err)
			}

			doc.Find(".beatmapListing a").Each(func(i int, s *goquery.Selection) {
				playersTemp = append(playersTemp, s.Text())
			})
			time.Sleep(100 * time.Millisecond)

		}
		players <- playersTemp
		time.Sleep(60 * time.Second)
	}
}

func storeRecentPlays(username string, sem chan bool) {
	songs, err := APIConnection.GetRecentPlays(url.QueryEscape(username), api.OSU, 100)

	if err != nil {
		//fmt.Println(err)
	}

	for i := range songs {
		song := &songs[i]
		song.Username = username
	}
	for _, element := range songs {
		sem <- true
		stmt, err := db.PrepareNamed("INSERT INTO plays(beatmap_id, score, max_combo, count50, count100, count300, count_miss, count_katu, count_geki, perfect, enabled_mods, user_id, date, rank, username) VALUES(:beatmap_id, :score, :max_combo, :count50, :count100, :count300, :count_miss, :count_katu, :count_geki, :perfect, :enabled_mods, :user_id, :date, :rank, :username) ON CONFLICT DO NOTHING")
		if err != nil {
			log.Fatal(err)
		}
		res, err := stmt.Exec(&element)
		if err != nil || res == nil {
			log.Fatal(err)
		}

		stmt.Close()
		<-sem
		//_, err = db.NamedQuery(`INSERT INTO plays(beatmap_id, score, max_combo, count50, count100, count300, count_miss, count_katu, count_geki, perfect, enabled_mods, user_id, date, rank, username) VALUES(:beatmap_id, :score, :max_combo, :count50, :count100, :count300, :count_miss, :count_katu, :count_geki, :perfect, :enabled_mods, :user_id, :date, :rank, :username)`, &element)

	}
}
