package mock

import (
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
)

const HookTable string = "hooks"

// Service Build type for interface in interfaces folder
type HookMock struct {
	Mock *r.Mock
	storage.IHook
}

// Get hooks by image
func (s *HookMock) GetByToken(token string) (*model.Hook, error) {
	return nil, nil
}

// Get hooks by image
func (s *HookMock) ListByUser(id string) (*model.HookList, error) {
	return nil, nil
}

// Get hooks by image
func (s *HookMock) ListByImage(user, id string) (*model.HookList, error) {
	return nil, nil
}

// Get hooks by service
func (s *HookMock) ListByService(user, id string) (*model.HookList, error) {
	return nil, nil
}

// Insert new hook into storage
func (s *HookMock) Insert(hook *model.Hook) (*model.Hook, error) {
	return nil, nil
}

// Insert new hook into storage
func (s *HookMock) Delete(user, id string) error {
	return nil
}

func newHookMock(mock *r.Mock) *HookMock {
	s := new(HookMock)
	s.Mock = mock
	return s
}
