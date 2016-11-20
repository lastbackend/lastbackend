package cmd

type NewUserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenInfo struct {
	Token string `json:"token"`
}

type LoginInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type profileInfo struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Company   string `json:"company"`
}

type WhoamiInfo struct {
	Id           string      `json:"id"`
	Username     string      `json:"username"`
	Email        string      `json:"email"`
	Gravatar     string      `json:"gravatar"`
	Balance      float32     `json:"balance"`
	Organization bool        `json:"organization"`
	Profile      profileInfo `json:"profile"`
	Created      string      `json:"created"`
	Updated      string      `json:"updated"`
}
