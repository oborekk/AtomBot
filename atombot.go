package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/xml"
	"bytes"
	"time"
	"os"
)

type Link struct {
	Href	string	`xml:"href,attr"`
}

type Entry struct {
	Published	string `xml:"published"`
	Updated	string	`xml:"updated"`
	Links	Link	`xml:"link"`
	Title	string	`xml:"title"`
	Content	string	`xml:"content"`
}

type Name struct {
	Name	string	`xml:"name"`
}

type Feed struct {
	Title 	string 	`xml:"title"`
	Updated string 	`xml:"updated"`
	Author 	Name	`xml:"author"`
	Entries	[]Entry	`xml:"entry"`
}

func webexPost(entry Entry) {
	url := os.Getenv("CHATAPI")
	auth := os.Getenv("AUTH")
	room := os.Getenv("ROOM")

	message := fmt.Sprintf(`{"roomId":"%s","markdown":"### ⚠️ [%s](%s) %s"}`,
		room, entry.Title, entry.Links.Href, entry.Content)
	data := []byte(message)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Println("Error sendgin to chat ", err)
		return
	}

	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	bearer := fmt.Sprintf("Bearer %s", auth)
	req.Header.Add("Authorization", bearer)
	
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending to chat ", err)
		return
	}
	defer resp.Body.Close()
}

func testPost(entry Entry){
	fmt.Println(entry.Title)
}

func main() {
	var lastPubDate string
	rss := os.Getenv("RSS")

	for true {
		resp, err := http.Get(rss)
		if err != nil {
			log.Fatalln("Unable to get response: ", err)
		}
		defer resp.Body.Close()

		feed := Feed{}
		decoder := xml.NewDecoder(resp.Body)
		err = decoder.Decode(&feed)
		if err != nil {
			log.Printf("Error decoding xml response: ", err)
		}

		pubDate := feed.Updated
		var data Entry
		if lastPubDate != "" {
			time1, err1 := time.Parse("2006-01-02T15:04:05Z07:00", pubDate)
			if err1 != nil {
				log.Println("Could not parse time1:", err, pubDate)
			}
			time2, err2 := time.Parse("2006-01-02T15:04:05Z07:00", lastPubDate)
			if err2 != nil {
				log.Println("Could not parse time2:", err, lastPubDate)
			}
			
			if time2.Before(time1) {
				for _, entry := range feed.Entries {
					time3, err3 := time.Parse("2006-01-02T15:04:05Z07:00", entry.Updated)
					if err3 != nil {
						log.Println("Could not parse time3: ", err, entry.Updated)
					}

					if time3.After(time2){
						data = entry
						// testPost(data)
						webexPost(data)
					}
				}
			}
		// ** Else for tests: bulk posts current statuses
		// } else {
		// 	for _, entry := range feed.Entries {
		// 		data = entry
		// 		// testPost(data)
		// 		webexPost(data)
		// 	}
		}
		lastPubDate = pubDate
		time.Sleep(time.Second*10)
	}
}