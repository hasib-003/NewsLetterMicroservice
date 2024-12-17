package utils

import (
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"net/http"
)

func StartCorn(endpoint string) {
	c := cron.New()
	_, err := c.AddFunc("*/10 * * * *", func() {
		log.Printf("Cron job : Calling publish news ")
		resp, err := http.Get(endpoint)
		if err != nil {
			log.Println("error calling publish news ")

		} else {
			log.Println("Cron job response ", resp.Status)
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Println(err)
				}
			}(resp.Body)
		}
	})
	if err != nil {
		log.Fatalf("error calling publish news %v", err)
	}
	c.Start()
	log.Println("Cron job started")
}
