package models

type Order struct {
	Zone  string `json:"zone"`
	Row   int64  `json:"row"`
	Seat  int64  `json:"seat"`
	Email string `json:"email"`
}
