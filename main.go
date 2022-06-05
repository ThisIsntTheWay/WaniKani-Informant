package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/TwiN/go-color"
)

type Configuration struct {
	ApiToken   string `json:"wkApiToken"`
	WebhookURL string `json:"discordUrl"`
}

type GraduationInfo struct {
	Counter       int
	RadGrads      int
	KanGrads      int
	VocGrads      int
	TotItems      int
	AvailableTime string
}

func main() {
	cfg := Configuration{
		ApiToken:   "",
		WebhookURL: "",
	}

	// Create config file
	if _, err := os.Stat("configuration.json"); errors.Is(err, os.ErrNotExist) {
		// Read from console
		reader := bufio.NewReader(os.Stdin)
		fmt.Println(color.Colorize(color.Yellow, "[!] Configuration is missing."))

		fmt.Println("Please enter your WaniKani V2 API token.")
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		text = strings.TrimSuffix(text, "\r\n")
		cfg.ApiToken = text

		fmt.Println("Please enter your Discord Webhook URL.")
		fmt.Print("> ")
		text, _ = reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		text = strings.TrimSuffix(text, "\r\n")
		cfg.WebhookURL = text

		file, _ := json.MarshalIndent(cfg, "", " ")
		_ = ioutil.WriteFile("configuration.json", file, 0644)
	} else {
		file, _ := ioutil.ReadFile("configuration.json")
		_ = json.Unmarshal([]byte(file), &cfg)
	}

	// See what can graduate
	assignments := getAssignments(cfg.ApiToken)
	reviews := getReviews(cfg.ApiToken)

	fmt.Println(color.Colorize(color.Blue, "------------------"))

	var gradObject GraduationInfo
	gradObject.Counter = 0
	gradObject.RadGrads = 0
	gradObject.KanGrads = 0
	gradObject.VocGrads = 0

	for index, e := range reviews.SummaryData.Reviews {
		// Skip empty reviews
		if len(e.SubjectIds) == 0 {
			continue
		} else {
			fmt.Print(color.Colorize(color.Yellow, strconv.Itoa(index)))
			fmt.Println("", e.AvailableAt)
		}

		haveFoundGraduatingReview := false
		graduatingReviewTotalItems := 0

		for i, referenceSubjectId := range e.SubjectIds {
			fmt.Print(i)

			for _, assignmentCollectionElement := range assignments.Data {
				d := assignmentCollectionElement.Data
				if d.SubId == referenceSubjectId {
					// Items with SRS stage 4, meaning Apprentice 4, can succeed to Guru 1 -> Meaning "passed"
					if d.SrsStage == 4 {
						fmt.Print(color.Colorize(color.Green, " --> "))
						fmt.Println(d.SubId, "can graduate, is SubType", d.SubType)

						haveFoundGraduatingReview = true

						gradObject.Counter++
						switch d.SubType {
						case "radical":
							gradObject.RadGrads++
							break
						case "kanji":
							gradObject.KanGrads++
							break
						case "vocabulary":
							gradObject.VocGrads++
							break
						default:
							fmt.Println(color.Colorize(color.Red, " -> Subtype unknown"))
						}
					} else {
						fmt.Print(color.Colorize(color.Red, " --> "))
						fmt.Println(d.SubId, "cannot graduate, is SubType", d.SubType)
					}
				}
			}

			// Reset total items counter if this is not a graduating review
			if !haveFoundGraduatingReview && (i == (len(e.SubjectIds) - 1)) {
				graduatingReviewTotalItems = 0
			} else {
				graduatingReviewTotalItems = i + 1
				gradObject.AvailableTime = e.AvailableAt
			}
		}

		// Abort if the first review has graduating items
		if haveFoundGraduatingReview {
			gradObject.TotItems = graduatingReviewTotalItems
			break
		}
	}

	fmt.Println(gradObject)
	postToDiscord(cfg.WebhookURL, gradObject)
}
