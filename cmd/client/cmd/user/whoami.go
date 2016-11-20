package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jarcoal/httpmock"
	mock "github.com/lastbackend/lastbackend/cmd/client/cmd/user/mocks"
	structs "github.com/lastbackend/lastbackend/cmd/client/cmd/user/structs"
	"github.com/lastbackend/lastbackend/cmd/client/config"
	"github.com/lastbackend/lastbackend/cmd/client/context"
	httpClient "github.com/lastbackend/lastbackend/libs/http/client"
	"github.com/lastbackend/lastbackend/libs/table"
	"io/ioutil"
	"os"
	"strconv"
)

func Whoami(ctx *context.Context) {
	whoamiContent, err, _ := WhoamiLogic(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if whoamiContent.Id == "" {
		return
	}

	var header []string = []string{"Username", "Email", "Balance", "Organization", "Created", "Updated"}
	var data [][]string

	organization := strconv.FormatBool(whoamiContent.Organization)
	balance := strconv.FormatFloat(float64(whoamiContent.Balance), 'f', 2, 32)
	d := []string{
		whoamiContent.Username, whoamiContent.Email, balance,
		organization, whoamiContent.Created, whoamiContent.Updated}
	data = append(data, d)
	d = d[:0]

	table.PrintTable(header, data, []string{})
}

func WhoamiLogic(ctx *context.Context) (structs.WhoamiInfo, error, string) {
	var token string

	if ctx == context.Mock() {
		if ctx.Info.Version == "OK" {
			token = mock.MockWhoamiOk()
		} else if ctx.Info.Version == "BAD" {
			token = mock.MockWhoamiBad()
		}
		defer httpmock.Deactivate()
	} else {
		tokenFile, err := os.Open(config.Get().StoragePath + "token")
		if err != nil {
			return structs.WhoamiInfo{}, err, ""
		}
		defer tokenFile.Close()

		fileContent, err := ioutil.ReadAll(tokenFile)
		if err != nil {
			return structs.WhoamiInfo{}, err, ""
		}
		token = string(fileContent)
	}

	data := structs.TokenInfo{Token: token}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return structs.WhoamiInfo{}, err, ""
	}

	resp, status := httpClient.Get(config.Get().UserUrl, jsonData, "Authorization", "Bearer "+token)
	if status == 200 {
		var whoamiContent structs.WhoamiInfo
		err = json.Unmarshal(resp, &whoamiContent)
		if err != nil {
			return structs.WhoamiInfo{}, err, ""
		}

		return whoamiContent, err, ""
	}

	var httpError structs.ErrorJson
	err = json.Unmarshal(resp, &httpError)
	if err != nil {
		return structs.WhoamiInfo{}, err, ""
	}
	fmt.Printf("Whoami failed: %s\n", httpError.Message)

	return structs.WhoamiInfo{}, nil, httpError.Message
}
