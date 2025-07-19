package model

type Player struct {
	UserName string `json:"username"`
	UserID   string `json:"userId"`
	State    string `json:"state"`
}
