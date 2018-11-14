package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	gosxnotifier "github.com/deckarep/gosx-notifier"
)

type data struct {
	Items []struct {
		Id string `json:"id"`
	} `json:"data"`
	Pagination struct {
		cursor string `json:cursor`
	} `json:"pagination"`
}

func readClientId() string {
	dat, err := ioutil.ReadFile("/Users/luna/.config/twitch-alert/client-id")
	if err != nil {
		fmt.Println(err)
	}
	return string(dat)

}
func getRequestUser(user string, clientId string) *http.Request {
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/users?login="+user, nil)

	req.Header.Set("Client-ID", clientId)
	return req
}

func getRequestStream(userId string, clientId string) *http.Request {
	req, _ := http.NewRequest("GET", "https://api.twitch.tv/helix/streams?user_id="+userId, nil)

	req.Header.Set("Client-ID", clientId)
	return req
}
func main() {
	// ticker := time.NewTicker(5 * time.Second)
	// quit := make(chan struct{})
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			// do stuff
	// 			notify()
	// 		case <-quit:
	// 			ticker.Stop()
	// 			return
	// 		}
	// 	}
	// }()
	notify()
}
func notify() {
	// read id from file to get clientId
	clientId := readClientId()

	note := gosxnotifier.NewNotification("twitch noticator")
	note.Group = "com.github.felipeluna.twitch-notify"
	note.AppIcon = "~/.config/twich-alert/resources/twitch-logo.png"
	streamer := os.Args[1]

	// make requet with twitch username to get userId
	client := &http.Client{}
	req := getRequestUser(streamer, clientId)
	res, _ := client.Do(req)

	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	var d data
	err = json.Unmarshal(contents, &d)
	userId := d.Items[0].Id

	// user userId to make request and see if user is online
	req = getRequestStream(userId, clientId)
	res, _ = client.Do(req)

	contents, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	err = json.Unmarshal(contents, &d)

	if len(d.Items) == 0 {
		note.Title = streamer + " is Offline"
		note.ContentImage = "./resources/BibleThump.png"
	} else {
		note.Title = streamer + " is Online"
		note.ContentImage = "./resources/PogChamp.png"
		note.Link = "http://www.twitch.com/" + streamer //or BundleID like: com.apple.Terminal
	}

	//Then, push the notification
	err = note.Push()

	//If necessary, check error
	if err != nil {
		log.Println("Uh oh!")
	}

}
