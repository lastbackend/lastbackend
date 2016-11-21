package structures

type Project struct {
	//Id   string           `json:"id"`
	//User string		`json:"user"`
	Name string 		`json:"project_name"`
	Description string 	`json:"description"`
	Created string 		`json:"time"`
	Updated string 		`json:"time"`


}
type ProjList struct {
	Proj []Project `json:"projects"`
}

/*
	Name string `json:"project_name" gorethink:"project_name"`
	Description string `json:"description" gorethink:"description"`
	Created string `json:"time" gorethink:"time"`
	Updated string `json:"time" gorethink:"time"`
 */
