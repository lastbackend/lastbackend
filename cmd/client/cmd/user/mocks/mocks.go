package mocks

import (
	"github.com/jarcoal/httpmock"
	"github.com/lastbackend/lastbackend/cmd/client/config"
)

func MockWhoamiOk() string {
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

func MockWhoamiBad() string {
	token := "token"

	httpmock.Activate()

	httpmock.RegisterResponder("GET", config.Get().UserUrl,
		httpmock.NewStringResponder(404, `{"code":404,
											"status":"USER_NOT_FOUND",
											"message":"user not found"}`))

	return token
}

func MockSignUpOk() (string, string, string) {
	username := "testname"
	email := "test@lb.com"
	password := "12345678"

	httpmock.Activate()

	httpmock.RegisterResponder("POST", config.Get().UserUrl,
		httpmock.NewStringResponder(200, `{"token": "token"}`))

	return username, email, password
}

func MockSignUpBadUsername() (string, string, string) {
	username := "tes"
	email := "test@lb.com"
	password := "12345678"

	httpmock.Activate()

	httpmock.RegisterResponder("POST", config.Get().UserUrl,
		httpmock.NewStringResponder(406, `{"code":406,
											"status":"BAD_PARAMETER_USERNAME",
											"message":"bad username parameter"}`))

	return username, email, password
}

func MockSignUpBadEmail() (string, string, string) {
	username := "testname"
	email := "test@lb"
	password := "12345678"

	httpmock.Activate()

	httpmock.RegisterResponder("POST", config.Get().UserUrl,
		httpmock.NewStringResponder(406, `{"code":406,
											"status":"BAD_PARAMETER_EMAIL",
											"message":"bad email parameter"}`))

	return username, email, password
}

func MockSignUpBadPassword() (string, string, string) {
	username := "testname"
	email := "test@lb.com"
	password := "12345"

	httpmock.Activate()

	httpmock.RegisterResponder("POST", config.Get().UserUrl,
		httpmock.NewStringResponder(406, `{"code":406,
											"status":"BAD_PARAMETER_PASSWORD",
											"message":"bad password parameter"}`))

	return username, email, password
}

func MockSignInOk() (string, string) {
	login := "testname"
	password := "12345678"

	httpmock.Activate()

	httpmock.RegisterResponder("POST", config.Get().AuthUserUrl,
		httpmock.NewStringResponder(200, `{"token": "token"}`))

	return login, password
}

func MockSignInBad() (string, string) {
	login := "testname"
	password := "12345678"

	httpmock.Activate()

	httpmock.RegisterResponder("POST", config.Get().AuthUserUrl,
		httpmock.NewStringResponder(401, `{"code":401,
											"status":"ACCESS_DENIED",
											"message":"access denied"}`))

	return login, password
}
