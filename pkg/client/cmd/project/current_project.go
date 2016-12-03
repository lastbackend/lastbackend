package project

import (
	"errors"
	tab "github.com/crackcomm/go-clitable"
	"github.com/lastbackend/lastbackend/libs/model"
	"github.com/lastbackend/lastbackend/pkg/client/context"
)

func Current() error {
	var ctx = context.Get()

	var project = new(model.Project)
	//var name string
	//err := ctx.Storage.Get(".current", &name)
	err := ctx.Storage.Get("project", &project)
	if err != nil {
		return errors.New(err.Error())
	}

	table := tab.New([]string{"ID", "Name", "Created", "Updated"})
	table.AddRow(map[string]interface{}{
		"ID":      project.ID,
		"Name":    project.Name,
		"Created": project.Created.String()[:10],
		"Updated": project.Updated.String()[:10],
	})
	table.Markdown = true
	table.Print()

	return nil
}
