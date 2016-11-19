package cmd

type newUserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenInfo struct {
	Token string `json:"token"`
}

type loginInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type whoamiInfo struct {
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	Balance      float64 `json:"balance"`
	Organization bool    `json:"organization"`
	Created      string  `json:"created"`
	Updated      string  `json:"updated"`
}
