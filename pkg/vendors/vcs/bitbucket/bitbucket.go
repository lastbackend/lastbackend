package bitbucket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lastbackend/vendors/interfaces"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// const

const (
	API_V1_URL = "https://bitbucket.org"
	API_V2_URL = "https://api.bitbucket.org"
	TOKEN_URL  = "https://bitbucket.org/site/oauth2/access_token"
)

// Model

type BitBucket struct {
	proto interfaces.OAuth2
}

// IVendor

func GetClient(clientID, clientSecretID, redirectURI string) *BitBucket {
	return &BitBucket{proto: interfaces.OAuth2{ClientID: clientID, ClientSecret: clientSecretID, RedirectUri: redirectURI}}
}

func (BitBucket) GetVendorInfo() *interfaces.Vendor {
	return &interfaces.Vendor{Vendor: "bitbucket", Host: "bitbucket.org"}
}

// IOAuth2 func

func (g BitBucket) GetToken(code string) (*oauth2.Token, error) {

	token, err := g.getOAuth2Config().Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (g BitBucket) RefreshToken(token *oauth2.Token) (*oauth2.Token, bool, error) {

	var err error

	if !token.Expiry.Before(time.Now()) || token.RefreshToken == "" {
		return token, false, nil
	}

	restoredToken := &oauth2.Token{
		RefreshToken: token.RefreshToken,
	}

	c := g.getOAuth2Config().Client(oauth2.NoContext, restoredToken)

	newToken, err := c.Transport.(*oauth2.Transport).Source.Token()
	if err != nil {
		return nil, false, err
	}

	return newToken, true, nil
}

// IOAuth2 - Private functions

func (g BitBucket) getOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     g.proto.ClientID,
		ClientSecret: g.proto.ClientSecret,
		RedirectURL:  g.proto.RedirectUri,
		Endpoint: oauth2.Endpoint{
			TokenURL: TOKEN_URL,
		},
	}
}

// IVCS func

func (BitBucket) GetUser(token *oauth2.Token) (*interfaces.User, error) {

	var err error

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	payload := struct {
		Username string `json:"username"`
		ID       string `json:"uuid"`
	}{}

	var uri = fmt.Sprintf("%s/2.0/user", API_V2_URL)

	resUser, err := client.Get(uri)

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resUser.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	var user = new(interfaces.User)

	user.Username = payload.Username
	user.ServiceID = payload.ID

	emailsResponse := struct {
		Emails []struct {
			Email     string `json:"email"`
			Confirmed bool   `json:"is_confirmed"`
			Primary   bool   `json:"is_primary"`
		} `json:"values"`
	}{}

	uri = fmt.Sprintf("%s/2.0/user/emails", API_V2_URL)

	resEmails, err := client.Get(uri)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resEmails.Body).Decode(&emailsResponse)
	if err != nil {
		return nil, err
	}

	for _, email := range emailsResponse.Emails {
		if email.Confirmed == true && email.Primary == true {
			user.Email = email.Email
			break
		}
	}

	return user, nil
}

func (BitBucket) ListRepositories(token *oauth2.Token, username string, org bool) (*interfaces.VCSRepositories, error) {

	var res *http.Response
	var err error

	username = strings.ToLower(username)

	payload := struct {
		Repos []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Private     bool   `json:"is_private"`
		} `json:"values"`
	}{}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/2.0/repositories?role=owner", API_V2_URL)
	if org {
		uri = fmt.Sprintf("%s/2.0/repositories/%s?role=admin&pagelen=100", API_V2_URL, username)
	}

	res, err = client.Get(uri)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(res.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	var repositories = new(interfaces.VCSRepositories)

	for _, repo := range payload.Repos {
		repository := new(interfaces.VCSRepository)

		repository.Name = repo.Name
		repository.Description = repo.Description
		repository.Private = repo.Private
		*repositories = append(*repositories, *repository)
	}

	return repositories, nil
}

func (BitBucket) ListBranches(token *oauth2.Token, owner, repo string) (*interfaces.VCSBranches, error) {

	var branches = new(interfaces.VCSBranches)

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)

	payload := struct {
		Branches []struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"values"`
	}{}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/2.0/repositories/%s/%s/refs/branches?pagelen=100", API_V2_URL, owner, repo)

	res, err := client.Get(uri)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	for _, br := range payload.Branches {
		if br.Type == "branch" {
			branch := new(interfaces.VCSBranch)

			branch.Name = br.Name
			*branches = append(*branches, *branch)
		}
	}

	return branches, nil
}

func (BitBucket) GetLastCommitOfBranch(token *oauth2.Token, owner, repo, branch string) (*interfaces.Commit, error) {

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)
	branch = strings.ToLower(branch)

	payload := struct {
		Commits []struct {
			Hash   string `json:"hash"`
			Author struct {
				Raw  string `json:"raw"`
				User struct {
					Username string `json:"username"`
				} `json:"user"`
			} `json:"author"`
			Date    time.Time `json:"date"`
			Message string    `json:"message"`
		} `json:"values"`
	}{}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/2.0/repositories/%s/%s/commits/%s", API_V2_URL, owner, repo, branch)

	res, err := client.Get(uri)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	// Another regular expression: <(\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,6}\b)>$
	r, _ := regexp.Compile("<(.+)>$")

	if len(payload.Commits) == 0 {
		err = errors.New("Repo has no commits")
		return nil, err
	}

	if len(payload.Commits) == 0 {
		return nil, nil
	}

	var commit = new(interfaces.Commit)
	commit.Username = payload.Commits[0].Author.User.Username
	commit.Hash = payload.Commits[0].Hash
	commit.Message = payload.Commits[0].Message
	commit.Date = payload.Commits[0].Date
	commit.Email = r.FindStringSubmatch(payload.Commits[0].Author.Raw)[1]

	return commit, nil
}

func (BitBucket) GetReadme(token *oauth2.Token, owner string, repo string) (string, error) {

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/1.0/repositories/%s/%s/raw/master/README.md", API_V1_URL, owner, repo)

	res, err := client.Get(uri)
	if err != nil {
		return "", err
	}

	var content string

	if res.StatusCode == 200 {
		buf, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		content = string(buf)
	}

	return string(content), nil
}

func (BitBucket) PushPayload(data []byte) (*interfaces.VCSBranch, error) {

	var err error

	payload := struct {
		Push struct {
			Changes []struct {
				New struct {
					Name string `json:"name"`
				} `json:"new"`
				Commits []struct {
					Hash    string    `json:"hash"`
					Message string    `json:"message"`
					Date    time.Time `json:"date"`
					Author  struct {
						User struct {
							Username string `json:"username"`
						} `json:"user"`
						Raw string `json:"raw"`
					} `json:"author"`
				} `json:"commits"`
			} `json:"changes"`
		} `json:"push"`
	}{}

	if err = json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}

	r, _ := regexp.Compile("<(.+)>$")

	branch := new(interfaces.VCSBranch)
	branch.Name = payload.Push.Changes[0].New.Name
	branch.LastCommit = interfaces.Commit{
		Hash:     payload.Push.Changes[0].Commits[0].Hash,
		Date:     payload.Push.Changes[0].Commits[0].Date,
		Username: payload.Push.Changes[0].Commits[0].Author.User.Username,
		Message:  payload.Push.Changes[0].Commits[0].Message,
		Email:    r.FindStringSubmatch(payload.Push.Changes[0].Commits[0].Author.Raw)[1],
	}

	return branch, nil
}

func (BitBucket) CreateHook(token *oauth2.Token, hookID, owner, repo, host string) (*string, error) {

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)

	payload := struct {
		ID    string `json:"uuid"`
		Error struct {
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}{}

	body := struct {
		Description string   `json:"description"`
		URL         string   `json:"url"`
		Active      bool     `json:"active"`
		Events      []string `json:"events"`
	}{"web", fmt.Sprintf(`%s/hook/bitbucket/process/%s`, host, hookID), true, []string{"repo:push", "pullrequest:approved", "pullrequest:created"}}

	var buf io.ReadWriter
	buf = new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		return nil, err
	}

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/2.0/repositories/%s/%s/hooks", API_V2_URL, owner, repo)

	res, err := client.Post(uri, "application/json", buf)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	if payload.Error.Message != "" {
		return nil, errors.New(payload.Error.Message)
	}

	id := payload.ID

	return &id, nil
}

func (BitBucket) RemoveHook(token *oauth2.Token, hookID, owner, repo string) error {

	var err error

	repo = strings.ToLower(repo)
	owner = strings.ToLower(owner)

	client := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(token))

	var uri = fmt.Sprintf("%s/2.0/repositories/%s/%s/hooks/%s", API_V2_URL, owner, repo, hookID)

	req, err := http.NewRequest("DELETE", uri, nil)
	req.Header.Set("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
