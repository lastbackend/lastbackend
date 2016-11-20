package mocks

import (
	"github.com/jarcoal/httpmock"
	"github.com/lastbackend/lastbackend/cmd/client/config"
)

func MockWhoami() string {
	token := "token"

	httpmock.Activate()

	httpmock.RegisterResponder("GET", config.Get().UserUrl,
		httpmock.NewStringResponder(200, `{"id":"some_id",
 											"username":"some_username",
											"email":"some_email",
											"gravatar":"some_gravatar",
											"balance":10,
											"organization":false,
											"profile":{
												"first_name":"some_first_name",
												"last_name":"some_last_name",
       											"company":"some_company"
       										},
											"created":"2014-01-16T07:38:28.45Z",
											"updated":"2014-01-16T07:38:28.45Z"}`))

	return token
}

func MockSignUp() (string, string, string) {
	username := "testname"
	email := "test@lb.com"
	password := "12345678"

	httpmock.Activate()

	httpmock.RegisterResponder("POST", config.Get().UserUrl,
		httpmock.NewStringResponder(200, `{"token": "token"}`))

	return username, email, password
}

func MockAuth() (string, string) {
	login := "testname"
	password := "12345678"

	httpmock.Activate()

	httpmock.RegisterResponder("POST", config.Get().AuthUserUrl,
		httpmock.NewStringResponder(200, `{"token": "token"}`))

	return login, password
}
