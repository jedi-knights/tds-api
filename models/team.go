package models

type Team struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Url            string `json:"url"`
	Gender         string `json:"gender"`
	Division       string `json:"division"`
	ConferenceId   int    `json:"conference_id"`
	ConferenceName string `json:"conference_name"`
	ConferenceUrl  string `json:"conference_url"`
}
