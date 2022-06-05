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

	for index, e := range reviews.SummaryData.Reviews {
		// Skip empty reviews
		if len(e.SubjectIds) == 0 {
			continue
		} else {
			fmt.Print(color.Colorize(color.Yellow, strconv.Itoa(index)))
			fmt.Println("", e.AvailableAt)
		}

		// ToDo: Detect SRS stages for subject IDs
		for _, referenceSubjectId := range e.SubjectIds {
			for _, assignmentCollectionElement := range assignments.Data {
				//fmt.Println(color.Colorize(color.Red, "------------"))

				e := assignmentCollectionElement.Data
				if e.SubId == referenceSubjectId {
					// Items with SRS stage 4, meaning Apprentice 4, can succeed to Guru 1 -> Meaning "passed"
					if e.SrsStage == 4 {
						fmt.Print(color.Colorize(color.Green, "--> "))
						fmt.Println(e.SubId, "can graduate.")
					}
				}
				//fmt.Println()
			}
		}
	}
}
