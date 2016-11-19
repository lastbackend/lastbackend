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
