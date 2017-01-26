package storage

import (
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
)

const HookTable string = "hooks"

// Service Build type for interface in interfaces folder
type HookStorage struct {
	Session *r.Session
	storage.IHook
}

// Get hooks by image
func (s *HookStorage) GetByToken(token string) (*model.Hook, error) {

	var (
		err          error
		hook         = new(model.Hook)
		token_filter = r.Row.Field("token").Eq(token)
	)

	res, err := r.Table(HookTable).Filter(token_filter).Run(s.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	err = res.All(hook)
	if err != nil {
		return nil, err
	}

	return hook, nil
}

// Get hooks by image
func (s *HookStorage) ListByUser(id string) (*model.HookList, error) {

	var (
		err         error
		hooks       = new(model.HookList)
		user_filter = r.Row.Field("user").Eq(id)
	)

	res, err := r.Table(HookTable).Filter(user_filter).Run(s.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	err = res.All(hooks)
	if err != nil {
		return nil, err
	}

	return hooks, nil
}

// Get hooks by image
func (s *HookStorage) ListByImage(user, id string) (*model.HookList, error) {

	var (
		err         error
		hookList    = new(model.HookList)
		hook_filter = map[string]interface{}{
			"user":  user,
			"image": id,
		}
	)

	res, err := r.Table(HookTable).Filter(hook_filter).Run(s.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	err = res.All(hookList)
	if err != nil {
		return nil, err
	}

	return hookList, nil
}

// Get hooks by service
func (s *HookStorage) ListByService(user, id string) (*model.HookList, error) {

	var (
		err         error
		hookList    = new(model.HookList)
		hook_filter = map[string]interface{}{
			"user":    user,
			"service": id,
		}
	)

	res, err := r.Table(HookTable).Filter(hook_filter).Run(s.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	err = res.All(hookList)
	if err != nil {
		return nil, err
	}

	return hookList, nil
}

// Insert new hook into storage
func (s *HookStorage) Insert(hook *model.Hook) (*model.Hook, error) {

	res, err := r.Table(HookTable).Insert(hook, r.InsertOpts{ReturnChanges: true}).RunWrite(s.Session)
	if err != nil {
		return nil, err
	}

	hook.ID = res.GeneratedKeys[0]

	return hook, nil
}

// Insert new hook into storage
func (s *HookStorage) Remove(id string) error {

	var opts = r.DeleteOpts{}

	err := r.Table(HookTable).Get(id).Delete(opts).Exec(s.Session)
	if err != nil {
		return err
	}

	return nil
}

func newHookStorage(session *r.Session) *HookStorage {
	r.TableCreate(HookTable, r.TableCreateOpts{}).Run(session)
	s := new(HookStorage)
	s.Session = session
	return s
}
