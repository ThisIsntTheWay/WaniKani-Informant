package main

import (
	"encoding/json"
	"fmt"

	"github.com/TwiN/go-color"
	"github.com/go-resty/resty/v2"
)

// REVIEW STRUCTURE
type Reviews struct {
	SubjectIds  []int  `json:"subject_ids"`
	AvailableAt string `json:"available_at"`
}

type SummaryData struct {
	NextReviewsAt string     `json:"next_reviews_at"`
	Reviews       []*Reviews `json:"reviews"`
}

type Summary struct {
	DataUpdatedAt string      `json:"data_updated_at"`
	SummaryData   SummaryData `json:"data"`
}

// ASSIGNMENT STRUCTURE
type AssignmentsSubData struct {
	SubId       int    `json:"subject_id"`
	SubType     string `json:"subject_type"`
	SrsStage    int    `json:"srs_stage"`
	AvailableAt string `json:"available_at"`
}

type AssignmentsData struct {
	Id        int                `json:"id"`
	UpdatedAt string             `json:"data_updated_at"`
	Data      AssignmentsSubData `json:"data"`
}

type Assignment struct {
	TotalCount    int                `json:"total_count"`
	DataUpdatedAt string             `json:"data_updated_at"`
	Data          []*AssignmentsData `json:"data"`
}

/* -----------
	FUNCTIONS
   ----------- */
func getReviews(apiToken string) Summary {
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
		return Summary{}
	} else {
		var obj Summary

		json.Unmarshal(resp.Body(), &obj)
		/*
			fmt.Println("DataUpdatedAt     :", obj.DataUpdatedAt)
			fmt.Println("NextReviewsAt     :", obj.SummaryData.Reviews)
		*/

		return obj
	}
}

func getAssignments(apiToken string) Assignment {
	client := resty.New()
	resp, err := client.R().
		SetAuthToken(apiToken).
		Get("https://api.wanikani.com/v2/assignments?srs_stages=1,2,3,4")

	if resp.StatusCode() != 200 {
		fmt.Println(color.Colorize(color.Red, "[!] Error accessing assignments."))
		fmt.Println("Status Code:", resp.StatusCode())
		fmt.Println("Error      :", err)
		return Assignment{}
	} else {
		var obj Assignment
		json.Unmarshal(resp.Body(), &obj)
		/*
			fmt.Printf("obj is: %v\n", obj)
			fmt.Println("  Status Code:", resp.StatusCode())
			fmt.Println("  Header:", resp.Header())

			fmt.Println("TCount            :", obj.TotalCount)
			fmt.Println("DataUpdatedAt     :", obj.DataUpdatedAt)
			fmt.Println("Data1             :", obj.Data[0].Id)
			fmt.Println("Data2             :", obj.Data[1].Id)
		*/

		return obj
	}
}
