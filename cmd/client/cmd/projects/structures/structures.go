package structures

type Project struct {
	Id   string           `json:"id"`
	Name string 		`json:"name"`
	User string 		`json:"user"`
	Description string 	`json:"description"`
	Created string 		`json:"created"`
	Updated string 		`json:"updated"`


}

type ProjList []Project
