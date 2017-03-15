package mock

import (
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
	"time"
)

const mockActivityID = "mocked"
const activityTable string = "activities"

// Project Service type for interface in interfaces folder
type ActivityMock struct {
	Mock *r.Mock
	storage.IActivity
}

var activityMock = &model.Activity{
	ID:      mockActivityID,
	User:    mockUserID,
	Project: mockProjectID,
	Service: mockServiceID,
	Name:    "mockname",
	Event:   "mockevent",
	Created: time.Now(),
	Updated: time.Now(),
}

func (s *ActivityMock) Insert(_ *model.Activity) (*model.Activity, error) {
	var err error
	var opts = r.InsertOpts{ReturnChanges: true}

	s.Mock.On(r.DB("test").Table(activityTable).Insert(activityMock, opts)).Return(nil, nil)

	err = r.DB("test").Table(activityTable).Insert(activityMock, opts).Exec(s.Mock)
	if err != nil {
		return nil, err
	}

	return activityMock, nil
}

func (s *ActivityMock) ListProjectActivity(user, project string) (*model.ActivityList, error) {

	var err error
	var activity_filter = map[string]interface{}{
		"user":    user,
		"project": project,
	}
	var activityList = new(model.ActivityList)

	s.Mock.On(r.DB("test").Table(activityTable).Filter(activity_filter)).Return(activityMock, nil)

	res, err := r.DB("test").Table(activityTable).Filter(activity_filter).Run(s.Mock)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(activityList)

	return activityList, nil
}

func (s *ActivityMock) ListServiceActivity(user, service string) (*model.ActivityList, error) {

	var err error
	var activity_filter = map[string]interface{}{
		"user":    user,
		"service": service,
	}
	var activityList = new(model.ActivityList)

	s.Mock.On(r.DB("test").Table(activityTable).Filter(activity_filter)).Return(activityMock, nil)

	res, err := r.DB("test").Table(activityTable).Filter(activity_filter).Run(s.Mock)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(activityList)

	return activityList, nil
}

func (s *ActivityMock) RemoveByProject(user, project string) error {

	var err error
	var activity_filter = map[string]interface{}{
		"user":    user,
		"project": project,
	}
	var opts = r.DeleteOpts{ReturnChanges: true}

	s.Mock.On(r.DB("test").Table(activityTable).Filter(activity_filter).Delete(opts)).Return(nil, nil)

	err = r.DB("test").Table(activityTable).Filter(activity_filter).Delete(opts).Exec(s.Mock)
	if err != nil {
		return err
	}

	return nil
}

func (s *ActivityMock) RemoveByService(user, service string) error {

	var err error
	var activity_filter = map[string]interface{}{
		"user":    user,
		"service": service,
	}
	var opts = r.DeleteOpts{ReturnChanges: true}

	s.Mock.On(r.DB("test").Table(activityTable).Filter(activity_filter).Delete(opts)).Return(nil, nil)

	err = r.DB("test").Table(activityTable).Filter(activity_filter).Delete(opts).Exec(s.Mock)
	if err != nil {
		return err
	}

	return nil
}

func newActivityMock(mock *r.Mock) *ActivityMock {
	s := new(ActivityMock)
	s.Mock = mock
	return s
}
