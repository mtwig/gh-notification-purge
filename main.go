package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/cli/go-gh"
	"net/http"
	"os"
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
	Url   string `json:"url"`
	Type  string `json:"type"`
	Title string `json:"title"`
}
type Notification struct {
	Id         string     `json:"id"`
	Subject    Subject    `json:"subject"`
	Reason     string     `json:"reason"`
	Repository Repository `json:"repository"`
	Url        string     `json:"url"`
}

// debug if enabled, debug information will be printed to the console
var debug = false

func main() {
	flag.BoolVar(&debug, "debug", false, "debug info while running")
	flag.Parse()

	httpClient, err := gh.HTTPClient(nil)
	if err != nil {
		appError(err, "unable to get http client")
	}
	var firstNotificationPageUrl = fmt.Sprintf("https://api.github.com/notifications?per_page=10?page=1?read=false")
	notificationResponse, err := httpClient.Get(firstNotificationPageUrl)
	if err != nil {
		appError(err, "error getting first page of unread notifications")
	}
	var notifications []Notification
	err = json.NewDecoder(notificationResponse.Body).Decode(&notifications)
	if err != nil {
		appError(err, "unable to decode notification response")
	}
	dbg(fmt.Sprintf("Github responded with %d notifications.\n", len(notifications)))

	// Track the unique set of URLs which need to be marked as read.
	// This will prevent sending multiple requests for the same thread
	// when the thread has multiple messages
	var uniqueThreads = make(map[string]Notification)
	for _, notification := range notifications {
		if notification.Subject.Type == "PullRequest" && notification.Reason == "review_requested" {
			dbg("Found a notification to mark as read\n")
			dbg(fmt.Sprintf("\tid: %s\n", notification.Id))
			dbg(fmt.Sprintf("\treason: %s\n", notification.Reason))
			dbg(fmt.Sprintf("\trepo: %s/%s\n", notification.Repository.Owner.Login, notification.Repository.Name))
			dbg(fmt.Sprintf("\turl: %s\n", notification.Subject.Url))
			dbg(fmt.Sprintf("Fetching details for pull request %s\n", notification.Subject.Url))
			pullResponse, err := httpClient.Get(notification.Subject.Url)
			if err != nil {
				appError(err, "unable pull request details")
			}
			dbg(fmt.Sprintf("\tResponse status: %s\n", pullResponse.Status))
			if pullResponse.StatusCode == 200 || pullResponse.StatusCode == 404 {
				var pull Pull
				err = json.NewDecoder(pullResponse.Body).Decode(&pull)
				if err != nil {
					appError(err, "unable to decode pull request")
				}
				dbg(fmt.Sprintf("\tPR state: %s\n", pull.State))
				if pull.State == "closed" {
					dbg(fmt.Sprintf("\tPR url: %s\n", notification.Subject.Url))
					uniqueThreads[notification.Url] = notification
				}
			} else {
				appError(errors.New(fmt.Sprintf("unhandled response code %d ", pullResponse.StatusCode)), "unhandled response code for pull request")
			}
		}
	}

	for url := range uniqueThreads {
		var currentThread = uniqueThreads[url]
		request, err := http.NewRequest("PATCH", url, nil)
		if err != nil {
			appError(err, "Unable to create request object.")
		}
		response, err := httpClient.Do(request)
		if err != nil {
			appError(err, fmt.Sprintf("Failed to mark notification %s as done.", currentThread.Subject.Title))
		}
		if response.StatusCode == 205 {
			info(fmt.Sprintf("Marked thread [%s] as read.\n", currentThread.Subject.Title))
		} else if response.StatusCode == 304 {
			warn("Thread was not marked as read\n")
		} else if response.StatusCode == 403 {
			appError(errors.New("unauthorized"), "authentication issue, check gh auth")
		} else {
			appError(errors.New(fmt.Sprintf("Response was %d", response.StatusCode)), "undocumented status code")
		}
		dbg(fmt.Sprintf("The response status was %d\n", response.StatusCode))
	}
}

// appError will print the msg string (if provided), and then print the error
// the code will then exit the CLI with exit code 1
func appError(err error, msg string) {
	if len(msg) > 0 {
		fmt.Printf("%s\n", msg)
	}
	fmt.Printf("%s\n", err)
	os.Exit(1)
}

// info print a message in an info format
func info(msg string) {
	fmt.Print(msg)
}

// info print a message in a warn format
func warn(msg string) {
	fmt.Printf(msg)
}

// if debug is enabled, print a message
func dbg(msg string) {
	if debug {
		fmt.Print(msg)
	}
}
