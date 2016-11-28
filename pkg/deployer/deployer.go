package deployer

import (
	"github.com/lastbackend/lastbackend/libs/model"
)

type Deployer struct{}

func (Deployer) DeployFromTemplate(service string, template *model.Template) error {
	return nil
}
