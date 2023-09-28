package models

type Conference struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Division string `json:"division"`
	Gender   string `json:"gender"`
}
