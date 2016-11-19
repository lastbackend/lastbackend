package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/lastbackend/lastbackend/cmd/client/config"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	httpClient "github.com/lastbackend/lastbackend/libs/http/client"
	"os"
	"io/ioutil"
)

func Whoami(ctx *context.Context) {
	whoamiContent, err := WhoamiLogic(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(whoamiContent.Username, "", whoamiContent.Email, "", whoamiContent.Balance, "",
		whoamiContent.Organization, "", whoamiContent.Created, "", whoamiContent.Updated)
}

func WhoamiLogic(ctx *context.Context) (whoamiInfo, error) {
	var token string

	if ctx == context.Mock() {
		token = MockWhoami()
		defer httpmock.Deactivate()
	} else {
		tokenFile, err := os.Open(config.Get().StoragePath + "token")
		if err != nil{
			return whoamiInfo{}, err
		}
		defer tokenFile.Close()

		fileContent, err := ioutil.ReadAll(tokenFile)
		if err != nil {
			fmt.Println(err.Error())
		}
		token = string(fileContent)
	}

	data := tokenInfo{Token: token}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return whoamiInfo{}, err
	}

	resp := httpClient.Get(config.Get().UserUrl, jsonData, "Authorization", "Bearer " + token)

	var whoamiContent whoamiInfo
	err = json.Unmarshal(resp, &whoamiContent)
	if err != nil {
		return whoamiInfo{}, err
	}

	return whoamiContent, err
}

func MockWhoami() string {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbSI6ImxhdnJAbGI" +
		"uY29tIiwiZXhwIjoxNDg3NjExOTM5LCJqdGkiOjE0" +
		"ODc2MTE5MzksIm9pZCI6IiIsInVpZCI6ImU3Y2YyMTQxLTQzMT" +
		"ItNGJiNi05Yjc5LTUxNjE5Mzk2N2M2OCIsInVzZXIiOiJsYXZyIn0._gq" +
		"x4yNH29Qqphv3Rxu8RDKruaUh82mSd_5bnv-CaxA"

	httpmock.Activate()

	httpmock.RegisterResponder("POST", config.Get().AuthUserUrl,
		httpmock.NewStringResponder(200, `{"username": :"lavr", "email":"lavr@lb.com", "balance":10, "organization":false, "created":""2014-01-16T07:38:28.45Z",", "updated":""2014-01-16T07:38:28.45Z""}`))

	return token
}
