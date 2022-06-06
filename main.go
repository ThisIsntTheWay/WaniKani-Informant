package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/TwiN/go-color"
	"github.com/google/uuid"
)

type Configuration struct {
	ApiToken    string `json:"wkApiToken"`
	WebhookURL  string `json:"discordUrl"`
	LastReview  string `json:"lastReview"`
	PostOnError bool   `json:"postOnError"`
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
		reader := bufio.NewScanner(os.Stdin)
		fmt.Println(color.Colorize(color.Yellow, "[!] Configuration is missing."))

		text := ""
		fmt.Println("Please enter your WaniKani V2 API token.")
		for {
			fmt.Print(color.Colorize(color.Cyan, "> "))
			reader.Scan()
			text = reader.Text()

			// Validate input
			if !IsValidUUID(text) {
				fmt.Println(text)
				fmt.Println(color.Colorize(color.Red, "[!] Input is not a valid UUID, please retry."))
			} else {
				break
			}
		}

		cfg.ApiToken = text

		fmt.Println("Please enter your Discord Webhook URL.")
		fmt.Print(color.Colorize(color.Cyan, "> "))
		reader.Scan()
		text = reader.Text()

		cfg.WebhookURL = text

		cfg.LastReview = "2000-01-01T22:00:00.000000Z"
		cfg.PostOnError = true

		file, _ := json.MarshalIndent(cfg, "", " ")
		_ = ioutil.WriteFile("configuration.json", file, 0644)
	} else {
		file, _ := ioutil.ReadFile("configuration.json")
		_ = json.Unmarshal([]byte(file), &cfg)
	}

	// See what can graduate
	assignments := getAssignments(cfg.ApiToken)
	if assignments.DataUpdatedAt == "" {
		os.Exit(1)
	}

	reviews := getReviews(cfg.ApiToken)
	if reviews.DataUpdatedAt == "" {
		os.Exit(1)
	}

	fmt.Println(color.Colorize(color.Blue, "------------------"))

	var gradObject GraduationInfo
	gradObject.Counter = 0
	gradObject.RadGrads = 0
	gradObject.KanGrads = 0
	gradObject.VocGrads = 0

	for index, e := range reviews.SummaryData.Reviews {
		// Skip empty reviews
		/*fmt.Print(strconv.Itoa(index)+color.Colorize(color.Blue, " [i] Subject ID length: "), len(e.SubjectIds))
		fmt.Println()*/

		if len(e.SubjectIds) == 0 {
			continue
		} else {
			fmt.Print(color.Colorize(color.Yellow, strconv.Itoa(index)))
			fmt.Println("", e.AvailableAt)
		}

		haveFoundGraduatingReview := false
		graduatingReviewTotalItems := 0

		for i, referenceSubjectId := range e.SubjectIds {
			for _, assignmentCollectionElement := range assignments.Data {
				d := assignmentCollectionElement.Data
				if d.SubId == referenceSubjectId {
					fmt.Print(i)

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

	// Only POST if current time is equal to available at
	layout := "2006-01-02T15:04:05.000000Z"
	reviewTime, err := time.Parse(layout, gradObject.AvailableTime)
	reviewTimeLast, _ := time.Parse(layout, cfg.LastReview)
	nowTime := time.Now()

	fmt.Println()
	if err != nil {
		fmt.Println(color.Colorize(color.Red, "[!] Could not parse AvailableAt time:"))
		fmt.Println(err)

		postErrorToDiscord("Unable to determine AvailableAt time of review.", err.Error())
	} else {
		if nowTime.After(reviewTime) {
			// Check if review was posted already
			if reviewTime.Equal(reviewTimeLast) {
				fmt.Println(color.Colorize(color.Yellow, "Review already POSTed to Discord."))
			} else {
				fmt.Println(reviewTime)
				fmt.Println(color.Colorize(color.Yellow, "Attempting POST to Discord..."))
				r := postToDiscord(cfg.WebhookURL, gradObject)

				// Update time
				if r {
					cfg.LastReview = gradObject.AvailableTime
					file, _ := json.MarshalIndent(cfg, "", " ")
					_ = ioutil.WriteFile("configuration.json", file, 0644)
				}
			}
		} else {
			fmt.Println(color.Colorize(color.Yellow, "[i] NowTime not after ReviewTime:"), reviewTime)
			fmt.Println(color.Colorize(color.Yellow, "Will not do anything..."))
		}
	}

	fmt.Println()
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
