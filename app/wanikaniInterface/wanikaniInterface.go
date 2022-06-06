package wanikaniInterface

import (
	"encoding/json"
	"fmt"

	"ThisIsntTheWay/wk-informant/app/discordInterface"
	"ThisIsntTheWay/wk-informant/app/structs"

	"github.com/TwiN/go-color"
	"github.com/go-resty/resty/v2"
)

/* -----------
	FUNCTIONS
   ----------- */
func GetReviews(apiToken string) structs.Summary {
	client := resty.New()

	resp, err := client.R().
		EnableTrace().
		SetAuthToken(apiToken).
		Get("https://api.wanikani.com/v2/summary")

	// Explore response object
	if resp.StatusCode() != 200 {
		fmt.Println(color.Colorize(color.Red, "[!] Error accessing review summary."))
		fmt.Println("Status Code:", resp.StatusCode())
		fmt.Println("Error      :", err)

		e := discordInterface.PostErrorToDiscord("Unable to access review summary", ("Error: " + string(resp.Body())))
		if e {
			fmt.Println(color.Colorize(color.Yellow, "[i] Error posted to Discord."))
		}

		return structs.Summary{}
	} else {
		var obj structs.Summary
		json.Unmarshal(resp.Body(), &obj)

		return obj
	}
}

func GetAssignments(apiToken string) structs.Assignment {
	client := resty.New()
	resp, err := client.R().
		SetAuthToken(apiToken).
		Get("https://api.wanikani.com/v2/assignments?srs_stages=1,2,3,4")

	if resp.StatusCode() != 200 {
		fmt.Println(color.Colorize(color.Red, "[!] Error accessing assignments."))
		fmt.Println("Status Code:", resp.StatusCode())
		fmt.Println("Error      :", err)

		e := discordInterface.PostErrorToDiscord("Unable to access assignments", ("Error: " + string(resp.Body())))
		if e {
			fmt.Println(color.Colorize(color.Yellow, "[i] Error posted to Discord."))
		}

		return structs.Assignment{}
	} else {
		var obj structs.Assignment
		json.Unmarshal(resp.Body(), &obj)

		return obj
	}
}
