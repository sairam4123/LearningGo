package main

type Status string

const (
	Pending   Status = "pending"
	Completed Status = "completed"
	Archived  Status = "archived"
)

func (s Status) String() string {
	return string(s)
}


type Todo struct {
	Message string `json:"message"`
	Status  Status `json:"status"`
	ID      int    `json:"id"`
}
