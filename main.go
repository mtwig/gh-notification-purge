package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cli/go-gh"
)

type Pull struct {
	State string `json:"state"`
}
type Owner struct {
	Login string `json:"login"`
}
type Repository struct {
	Name  string `json:"name"`
	Owner Owner  `json:"owner"`
}
type Subject struct {
	Url  string `json:"url"`
	Type string `json:"type"`
}
type Notification struct {
	Id         string     `json:"id"`
	Subject    Subject    `json:"subject"`
	Reason     string     `json:"reason"`
	Repository Repository `json:"repository"`
	Url        string     `json:"url"`
}

func main() {

	httpClient, err := gh.HTTPClient(nil)
	if err != nil {
		log.Fatal(err)
	}

	restClient, err := gh.RESTClient(nil)
	if err != nil {
		log.Fatal(err)
	}

	var url = fmt.Sprintf("https://api.github.com/notifications?per_page=100?page=1")

	notificationResponse, err := httpClient.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Patch
	var notifications []Notification
	err = json.NewDecoder(notificationResponse.Body).Decode(&notifications)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Avoid setting

	fmt.Printf("The length is %d.\n", len(notifications))
	for _, notification := range notifications {
		if notification.Subject.Type == "PullRequest" && notification.Reason == "review_requested" {
			fmt.Printf("- ID %s.\n", notification.Id)
			fmt.Printf("  %s\n", notification.Reason)
			fmt.Printf("  %s/%s\n", notification.Repository.Owner.Login, notification.Repository.Name)
			fmt.Printf("  %s\n", notification.Subject.Url)

			pullResponse, err := httpClient.Get(notification.Subject.Url)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Status code for PR is %d\n", pullResponse.StatusCode)

			if pullResponse.StatusCode == 200 {
				var pull Pull
				err = json.NewDecoder(pullResponse.Body).Decode(&pull)
				fmt.Printf("  The state is %s\n", pull.State)
				if pull.State == "closed" {
					fmt.Printf("To close %s\n", notification.Url)

					//err = restClient.Patch(notification.Url, nil, nil)
					var response struct{}
					err = restClient.Patch(notification.Url, nil, &response)
					if err != nil {
						log.Fatal(err)
					}
					/// TODO: possible cases are 205, 304, 403
					// Check status???
					return
				}
			} else if pullResponse.StatusCode == 404 {
				// TODO: Just do a patch?
				fmt.Printf("Github responded with a 404! Fix this code!\n")
				return
			}

		}
	}
}
