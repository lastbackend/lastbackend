package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/lastbackend/lastbackend/cmd/client/config"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	httpClient "github.com/lastbackend/lastbackend/libs/http/client"
	filesystem "github.com/lastbackend/lastbackend/libs/filesystem"
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
		homeDir, err := filesystem.GetHomeDir()
		if err != nil {
			return whoamiInfo{}, err
		}
		tokenFile, err := os.Open(homeDir + "/.lb/token")
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
	token := "eyJhbGciOiJIUzI1Ni" +
		"IsInR5cCI6IkpXVCJ9.eyJlbSI6ImxhdnJAb" +
		"GIuY29tIiwiZXhwIjoxNDg3NjA3MTExLCJqdGkiOjE0OD" +
		"c2MDcxMTEsIm9pZCI6IiIsInVpZCI6ImIzMjZjZjJlLTdmZTUtNDUzNS1h" +
		"NDg2LWEwY2I0Y2QzYTY5ZCIsInVzZXIiOiJsYXZyIn0.Xliv13Eko9xWqhcqx" +
		"tESLjfWuLuZYt5L4LARnawsfvw"

	httpmock.Activate()

	httpmock.RegisterResponder("POST", config.Get().AuthUserUrl,
		httpmock.NewStringResponder(200, `{"username": :"lavr", "email":"lavr@lb.com", "balance":10, "organization":false, "created":""2014-01-16T07:38:28.45Z",", "updated":""2014-01-16T07:38:28.45Z""}`))

	return token
}
