package utils

import (
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"net/http"
)

func StartCorn(endpoint string) {
	c := cron.New()
	_, err := c.AddFunc("*/1 * * * *", func() {
		log.Printf("Cron job : Calling publish news ")
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			log.Println("Error creating new request : ", err)
		}
		req.Header.Set("X-Cron-Job", "true")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("error calling publish news ")
			return

		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println("error closing body")
			}
		}(resp.Body)
		log.Println("response status : ", resp.Status)
	})
	if err != nil {
		log.Fatalf("error calling publish news %v", err)
	}
	c.Start()
	log.Println("Cron job started")
}
