package discordInterface

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"ThisIsntTheWay/wk-informant/app/structs"

	"github.com/TwiN/go-color"
	"github.com/go-resty/resty/v2"
)

/* -----------
	FUNCTIONS
   ----------- */
func PostToDiscord(url string, gradObj structs.GraduationInfo) bool {
	// Sanity check
	if gradObj.Counter == 0 {
		fmt.Println(color.Colorize(color.Red, "No items to graduate, will not POST."))
		return false
	}

	gradTemplate, _ := ioutil.ReadFile("json/msgGraduationTemplate.json")

	var obj structs.WebhookMessage
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
	// List applicable to webhooks embed titles
	for index, filterElement := range [...]string{"Radicals", "Kanji", "Vocab"} {
		x := 0
		s := ""

		switch index {
		case 0:
			x = gradObj.RadGrads
			s = "!radGrad!"
			break
		case 1:
			x = gradObj.KanGrads
			s = "!kanGrad!"
			break
		case 2:
			x = gradObj.VocGrads
			s = "!vocGrad!"
			break
		}

		// Remove slice if it doesn't contain any elements to graduate
		// Otherwise, adjust template
		if x == 0 {
			for i, e := range obj.Embeds {
				if e.Title == filterElement {
					obj.Embeds = append(obj.Embeds[:i], obj.Embeds[i+1:]...)
				}
			}
		} else {
			for i, e := range obj.Embeds {
				if e.Title == filterElement {
					obj.Embeds[i].Description = strings.Replace(
						obj.Embeds[i].Description, s, strconv.Itoa(x),
						-1)
				}
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

func PostErrorToDiscord(errHeader string, errMsg string) bool {
	// Get config
	var cfg structs.Configuration

	file, _ := ioutil.ReadFile("configuration.json")
	_ = json.Unmarshal([]byte(file), &cfg)

	if !cfg.PostOnError {
		return false
	}

	// Read template
	errMsgTemplate, _ := ioutil.ReadFile("json/msgErrorTemplate.json")

	var obj structs.WebhookMessage
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
