package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/go-resty/resty/v2"
)

/* -----------
	STRUCTS
   ----------- */
type EmbedItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       int    `json:"color"`
}
type WebhookMessage struct {
	Content     string      `json:"content"`
	Username    string      `json:"username"`
	AvatarUrl   string      `json:"avatar_url"`
	Embeds      []EmbedItem `json:"embeds"`
	Attachments []string
}

/* -----------
	FUNCTIONS
   ----------- */
func postToDiscord(url string, gradObj GraduationInfo) bool {
	// Sanity check
	if gradObj.Counter == 0 {
		fmt.Println(color.Colorize(color.Red, "No items to graduate, will not POST."))
		return false
	}

	gradTemplate, _ := ioutil.ReadFile("json/msgGraduationTemplate.json")

	var obj WebhookMessage
	json.Unmarshal(gradTemplate, &obj)

	/*
		fmt.Println(color.Colorize(color.Red, "1 ----------------"))
		fmt.Println(obj)*/

	// Adjust template
	obj.Content = strings.Replace(
		obj.Content, "!revTotal!", strconv.Itoa(gradObj.TotItems), -1)

	obj.Content = strings.Replace(
		obj.Content, "!reviewTime!", gradObj.AvailableTime, -1)

	// The wierd for loops are required as indexes could move when being modified
	f := "Radicals"
	if gradObj.RadGrads == 0 {
		for i, e := range obj.Embeds {
			if e.Title == f {
				obj.Embeds = append(obj.Embeds[:i], obj.Embeds[i+1:]...)
			}
		}
	} else {
		for i, e := range obj.Embeds {
			if e.Title == f {
				obj.Embeds[i].Description = strings.Replace(
					obj.Embeds[i].Description, "!radGrad!", strconv.Itoa(gradObj.RadGrads),
					-1)
			}
		}
	}

	f = "Kanji"
	if gradObj.KanGrads == 0 {
		for i, e := range obj.Embeds {
			if e.Title == f {
				obj.Embeds = append(obj.Embeds[:i], obj.Embeds[i+1:]...)
			}
		}
	} else {
		for i, e := range obj.Embeds {
			if e.Title == f {
				obj.Embeds[i].Description = strings.Replace(
					obj.Embeds[i].Description, "!kanGrad!", strconv.Itoa(gradObj.KanGrads),
					-1)
			}
		}
	}

	f = "Vocab"
	if gradObj.VocGrads == 0 {
		for i, e := range obj.Embeds {
			if e.Title == f {
				obj.Embeds = append(obj.Embeds[:i], obj.Embeds[i+1:]...)
			}
		}
	} else {
		for i, e := range obj.Embeds {
			if e.Title == f {
				obj.Embeds[i].Description = strings.Replace(
					obj.Embeds[i].Description, "!vocGrad!", strconv.Itoa(gradObj.VocGrads),
					-1)
			}
		}
	}

	client := resty.New()
	resp, _ := client.R().
		SetBody(obj).
		Post(url)

	if resp.IsError() {
		fmt.Println(color.Colorize(color.Red, "[!] Error POSTing to Discord."))
		fmt.Println(color.Colorize(color.Gray, "URL: "+url))
		fmt.Println("> Status Code:", resp.StatusCode())
		fmt.Println("> Response   :", resp)
		return false
	} else {
		fmt.Println(color.Colorize(color.Green, "Successfully POSTed to Discord."))
		return true
	}
}

func postErrorToDiscord(errHeader string, errMsg string) bool {
	// Get config
	var cfg Configuration

	file, _ := ioutil.ReadFile("configuration.json")
	_ = json.Unmarshal([]byte(file), &cfg)

	if !cfg.PostOnError {
		return false
	}

	// Read template
	errMsgTemplate, _ := ioutil.ReadFile("json/msgErrorTemplate.json")

	var obj WebhookMessage
	json.Unmarshal(errMsgTemplate, &obj)

	// Adjust template
	obj.Embeds[0].Description = strings.Replace(
		obj.Embeds[0].Description, "!errMsgContent!", errMsg, -1)
	obj.Embeds[0].Title = strings.Replace(
		obj.Embeds[0].Title, "!errMsgheader!", errHeader, -1)

	// POST
	client := resty.New()
	resp, _ := client.R().
		SetBody(obj).
		Post(cfg.WebhookURL)

	if resp.IsError() {
		fmt.Println(color.Colorize(color.Red, "[!] Could not POST error to discord."))
		fmt.Println(color.Colorize(color.Gray, "URL: "+cfg.WebhookURL))
		fmt.Println("> Status Code:", resp.StatusCode())
		fmt.Println("> Response   :", resp)
		return false
	} else {
		return true
	}
}
