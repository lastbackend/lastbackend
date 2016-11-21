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
	"strconv"
)

func Whoami() {
	if cfg.Token == "" {
		return
	}

	whoamiContent, err, _ := WhoamiLogic()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var header []string = []string{"Username", "Email", "Balance", "Organization", "Created", "Updated"}
	var data [][]string

	organization := strconv.FormatBool(whoamiContent.Organization)
	balance := strconv.FormatFloat(float64(whoamiContent.Balance), 'f', 2, 32)
	d := []string{
		whoamiContent.Username, whoamiContent.Email, balance,
		organization, whoamiContent.Created[:10], whoamiContent.Updated[:10]}
	data = append(data, d)
	d = d[:0]

	table.PrintTable(header, data, []string{})
}

func WhoamiLogic() (structs.WhoamiInfo, error, string) {
	var token string

	if ctx == context.Mock() {
		if ctx.Info.Version == "OK" {
			token = mock.MockWhoamiOk()
		} else if ctx.Info.Version == "BAD" {
			token = mock.MockWhoamiBad()
		}
		defer httpmock.Deactivate()
	} else {
		token = config.Get().Token
	}

	data := structs.TokenInfo{Token: token}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return structs.WhoamiInfo{}, err, ""
	}

	resp, status := httpClient.Get(cfg.UserUrl, jsonData, "Authorization", "Bearer "+token)
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
	fmt.Printf("Whoami failed: %s", httpError.Message)

	return structs.WhoamiInfo{}, nil, httpError.Message
}
