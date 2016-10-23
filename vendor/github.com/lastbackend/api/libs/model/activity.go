package model

type Activity struct {
	ID					string `json:"id"`
	UserID      string `json:"user_id"`
	EntityID    string `json:"entity_id"`
	Name        string `json:"name"`
	Action      string `json:"action"`
	Message     string `json:"message"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Created     string `json:"created"`
}

type Activities []Activity
