package request

import (
	"encoding/json"
	"github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/errors"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
	"io"
	"io/ioutil"
	"strings"
)

type RequestNamespaceCreateS struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *RequestNamespaceCreateS) DecodeAndValidate(reader io.Reader) *errors.Err {

	var (
		log = context.Get().GetLogger()
	)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Error(err)
		return errors.New("user").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("project").IncorrectJSON(err)
	}

	if s.Name == "" {
		return errors.New("project").BadParameter("name")
	}

	s.Name = strings.ToLower(s.Name)

	if len(s.Name) < 4 && len(s.Name) > 64 && !validator.IsProjectName(s.Name) {
		return errors.New("project").BadParameter("name")
	}

	return nil
}

type RequestNamespaceUpdateS struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *RequestNamespaceUpdateS) DecodeAndValidate(reader io.Reader) *errors.Err {

	var (
		log = context.Get().GetLogger()
	)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Error(err)
		return errors.New("user").Unknown(err)
	}

	err = json.Unmarshal(body, s)
	if err != nil {
		return errors.New("project").IncorrectJSON(err)
	}

	if s.Name == "" {
		return errors.New("project").BadParameter("name")
	}

	s.Name = strings.ToLower(s.Name)

	if len(s.Name) < 4 && len(s.Name) > 64 && !validator.IsProjectName(s.Name) {
		return errors.New("project").BadParameter("name")
	}

	return nil
}
