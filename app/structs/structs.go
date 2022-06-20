package structs

/* -----------
	General
   ----------- */
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

/* -----------
	Cache
   ----------- */
type Cache struct {
	GradObjects  []int  `json:"GradObjects"`
	LastReviewId string `json:"LastReviewId"`
}

/* -----------
	WaniKani
   ----------- */
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
	Webhooks
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
