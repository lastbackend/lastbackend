package storage

import (
	"github.com/lastbackend/lastbackend/libs/adapter"
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"
)

type IProfileService interface {
	Insert(db adapter.IDatabase, userID string) (*string, *e.Err)
	Update(db adapter.IDatabase, userID, firstName, lastName, company, city string) (bool, *e.Err)
	GetByUserID(db adapter.IDatabase, userID string) (*model.UserProfile, *e.Err)
}
