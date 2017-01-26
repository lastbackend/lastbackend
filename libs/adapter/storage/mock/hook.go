package mock

import (
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
	"time"
)

const hookTable string = "hooks"
const mockHookID string = "mocked"

// Service Build type for interface in interfaces folder
type HookMock struct {
	Mock *r.Mock
	storage.IHook
}

var hookMock = &model.Hook{
	ID:      mockHookID,
	User:    mockHookID,
	Service: mockHookID,
	Image:   "",
	Token:   "mocktoken",
	Created: time.Now(),
	Updated: time.Now(),
}

var hookImageMock = &model.Hook{
	ID:      mockHookID,
	User:    mockHookID,
	Service: "",
	Image:   mockHookID,
	Token:   "mocktoken",
	Created: time.Now(),
	Updated: time.Now(),
}

var hookServiceMock = &model.Hook{
	ID:      mockHookID,
	User:    mockHookID,
	Service: mockHookID,
	Image:   "",
	Token:   "mocktoken",
	Created: time.Now(),
	Updated: time.Now(),
}

// Get hooks by image
func (s *HookMock) GetByToken(token string) (*model.Hook, error) {

	var (
		err  error
		hook = new(model.Hook)
	)

	s.Mock.On(r.DB("test").Table(hookTable).Get(mockHookID)).Return(hookMock, nil)

	res, err := r.DB("test").Table(hookTable).Get(mockHookID).Run(s.Mock)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	err = res.One(hook)
	if err != nil {
		return nil, err
	}

	return hook, nil
}

// Get hooks by image
func (s *HookMock) ListByUser(id string) (*model.HookList, error) {
	return &model.HookList{}, nil
}

// Get hooks by image
func (s *HookMock) ListByImage(user, id string) (*model.HookList, error) {
	return &model.HookList{}, nil
}

// Get hooks by service
func (s *HookMock) ListByService(user, id string) (*model.HookList, error) {
	return &model.HookList{}, nil
}

// Insert new hook into storage
func (s *HookMock) Insert(hook *model.Hook) (*model.Hook, error) {

	var (
		err  error
		opts = r.InsertOpts{ReturnChanges: true}
	)

	s.Mock.On(r.DB("test").Table(hookTable).Insert(hookMock, opts)).Return(nil, nil)

	err = r.DB("test").Table(hookTable).Insert(hookMock, opts).Exec(s.Mock)
	if err != nil {
		return nil, err
	}

	return hook, nil
}

// Insert new hook into storage
func (s *HookMock) Delete(user, id string) error {

	var (
		err         error
		opts        = r.DeleteOpts{ReturnChanges: true}
		hook_filter = map[string]interface{}{
			"user": user,
			"id":   id,
		}
	)

	s.Mock.On(r.DB("test").Table(hookTable).Filter(hook_filter).Delete(opts)).Return(nil, nil)

	err = r.DB("test").Table(hookTable).Filter(hook_filter).Delete(opts).Exec(s.Mock)
	if err != nil {
		return err
	}

	return nil
}

func newHookMock(mock *r.Mock) *HookMock {
	s := new(HookMock)
	s.Mock = mock
	return s
}
