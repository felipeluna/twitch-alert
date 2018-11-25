package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

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

// read the file that contains the credentials for talking with twitch.tv api
func readClientId() string {
	dat, err := ioutil.ReadFile("/Users/luna/.config/twitch-alert/client-id")
	if err != nil {
		fmt.Println(err)
	}
	return string(dat)

}

// make a get request to api.twitch.tv with the
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
	done := make(chan bool)
	go func() {
		for {
			result := notify()
			t := time.Now()
			fmt.Printf("%s - %t\n", t.String(), result)
			if result {
				done <- true
			}
			time.Sleep(1 * time.Minute)
		}
	}()
	<-done
}
func notify() bool {
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

	var result bool = false
	if len(d.Items) > 0 {
		result = true
		note.Title = streamer + " is Online"
		note.ContentImage = "./resources/PogChamp.png"
		note.Link = "http://www.twitch.com/" + streamer //or BundleID like: com.apple.Terminal
		//Then, push the notification
		err = note.Push()

		//If necessary, check error
		if err != nil {
			log.Println("Uh oh!")
		}
	}
	return result

}
