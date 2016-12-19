package storage

import (
	e "github.com/lastbackend/lastbackend/libs/errors"
	"github.com/lastbackend/lastbackend/libs/model"

	"github.com/lastbackend/lastbackend/libs/interface/storage"
	r "gopkg.in/dancannon/gorethink.v2"
)

const HookTable string = "hooks"

// Service Build type for interface in interfaces folder
type HookStorage struct {
	Session *r.Session
	storage.IHook
}

// Get hooks by image
func (s *HookStorage) GetByToken(token string) (*model.Hook, *e.Err) {

	var err error
	var hook = new(model.Hook)
	var token_filter = r.Row.Field("token").Eq(token)
	res, err := r.Table(HookTable).Filter(token_filter).Run(s.Session)
	if err != nil {
		return nil, e.New("hook").Unknown(err)
	}

	res.One(hook)

	defer res.Close()
	return hook, nil
}

// Get hooks by image
func (s *HookStorage) GetByUser(id string) (*model.HookList, *e.Err) {

	var err error
	var hooks = new(model.HookList)
	var user_filter = r.Row.Field("user").Eq(id)
	res, err := r.Table(HookTable).Filter(user_filter).Run(s.Session)
	if err != nil {
		return nil, e.New("hook").Unknown(err)
	}

	res.All(hooks)

	defer res.Close()
	return hooks, nil
}

// Get hooks by image
func (s *HookStorage) ListByImage(user, id string) (*model.HookList, *e.Err) {

	var err error
	var hooks = new(model.HookList)
	var image_filter = r.Row.Field("image").Eq(id)
	var user_filter = r.Row.Field("user").Eq(user)
	res, err := r.Table(HookTable).Filter(image_filter).Filter(user_filter).Run(s.Session)
	if err != nil {
		return nil, e.New("hook").Unknown(err)
	}

	res.All(hooks)

	defer res.Close()
	return hooks, nil
}

// Get hooks by service
func (s *HookStorage) ListByService(user, id string) (*model.HookList, *e.Err) {

	var err error
	var hooks = new(model.HookList)
	var service_filter = r.Row.Field("service").Eq(id)
	var user_filter = r.Row.Field("user").Eq(user)
	res, err := r.Table(HookTable).Filter(service_filter).Filter(user_filter).Run(s.Session)
	if err != nil {
		return nil, e.New("hook").Unknown(err)
	}

	res.All(hooks)

	defer res.Close()
	return hooks, nil
}

// Insert new hook into storage
func (s *HookStorage) Insert(hook *model.Hook) (*model.Hook, *e.Err) {

	res, err := r.Table(HookTable).Insert(hook, r.InsertOpts{ReturnChanges: true}).RunWrite(s.Session)
	if err != nil {
		return nil, e.New("hook").Unknown(err)
	}

	hook.ID = res.GeneratedKeys[0]

	return hook, nil
}

// Insert new hook into storage
func (s *HookStorage) Delete(user, id string) *e.Err {
	var user_filter = r.Row.Field("user").Eq(user)
	_, err := r.Table(HookTable).Get(id).Filter(user_filter).Delete().Run(s.Session)
	if err != nil {
		return e.New("hook").NotFound(err)
	}

	return nil
}

func newHookStorage(session *r.Session) *HookStorage {
	r.TableCreate(HookTable, r.TableCreateOpts{}).Run(session)
	s := new(HookStorage)
	s.Session = session
	return s
}
