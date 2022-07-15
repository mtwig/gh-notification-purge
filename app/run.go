package app

import (
	"encoding/json"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/cli/go-gh"
	nlog "github.com/mtwig/gh-notification-purge/log"
	"strings"
)

func (app application) Run() (err error) {
	notifications, err := fetchNotifications()
	if err != nil {
		nlog.Debugln("Error fetch notification " + fmt.Sprintf("%s", err))
		return
	}
	var numDone = 0
	for _, notification := range notifications {
		if notification.NotificationReason != "review_requested" {
			continue
		}
		if notification.NotificationType != "PullRequest" {
			continue
		}

		pull, pullErr := fetchPull(notification.Pull)
		if pullErr != nil {
			err = pullErr
			return
		}

		if !(pull.Merged == true || pull.State == "closed") {
			continue
		}

		if app.printSubjects {
			nlog.ConsoleC(color.Bold, fmt.Sprintf(" 📮 %s\n", notification.NotificationTitle))
		}
		if app.dryRun {
			nlog.Console("Would run `")
			nlog.ConsoleC(color.Cyan, "gh api -X PATCH "+notification.NotificationUrl)
			nlog.ConsoleC(color.Blue, "`\n")
		} else {
			_, stdErr, patchErr := gh.Exec("api",
				"-X", "PATCH",
				notification.NotificationUrl)
			if patchErr != nil {
				nlog.ConsoleC(color.Red, "\t😲 "+stdErr.String()+"\n")
			}
			numDone++
		}

	}

	nlog.ConsoleC(color.Green, fmt.Sprintf("🎉 Marked %d notifications as read!\n", numDone))

	return
}

func fetchNotifications() (details []NotificationDetail, err error) {
	stdOut, _, err := gh.Exec("api",
		"notifications?per_page=100",
		"--paginate",
		"--jq",
		".[] | {url:.url,"+
			"reason:.reason,"+
			"type:.subject.type,"+
			"pull:.subject.url,"+
			"repository:.repository.full_name,"+
			"title:.subject.title}")
	if err != nil {
		return
	}
	var responseLines = strings.Split(stdOut.String(), "\n")

	for _, response := range responseLines {
		if len(response) == 0 {
			continue
		}
		var notification NotificationDetail
		err = json.Unmarshal([]byte(response), &notification)
		if err != nil {
			nlog.Debugln("Unable to unmarshal json")
		}
		details = append(details, notification)
	}
	return
}

func fetchPull(url string) (pull Pull, err error) {
	stdOut, _, err := gh.Exec("api",
		url,
		"--jq",
		"{state,merged,user:.user.login,commits,additions,deletions,changed_files,comments}")
	if err != nil {
		return
	}
	err = json.Unmarshal(stdOut.Bytes(), &pull)
	return
}

// user.login, commits, additions, deletions, changed
