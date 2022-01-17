package models

type Result struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}
