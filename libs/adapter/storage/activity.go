package storage

import (
	"github.com/lastbackend/lastbackend/libs/interface/storage"
	"github.com/lastbackend/lastbackend/libs/model"
	r "gopkg.in/dancannon/gorethink.v2"
	"time"
)

const ActivityTable string = "activities"

// Activity Service type for interface in interfaces folder
type ActivityStorage struct {
	Session *r.Session
	storage.IActivity
}

func (s *ActivityStorage) Insert(activity *model.Activity) (*model.Activity, error) {

	var (
		err  error
		opts = r.InsertOpts{ReturnChanges: true}
	)

	activity.Created = time.Now()
	activity.Updated = time.Now()

	res, err := r.Table(ActivityTable).Insert(activity, opts).RunWrite(s.Session)
	if err != nil {
		return nil, err
	}

	activity.ID = res.GeneratedKeys[0]

	return activity, nil
}

func (s *ActivityStorage) ListProjectActivity(user, project string) (*model.ActivityList, error) {

	var (
		err                 error
		serviceActivityList = new(model.ActivityList)
		activity_filter     = map[string]interface{}{
			"user":    user,
			"project": project,
		}
	)

	res, err := r.Table(ActivityTable).Filter(activity_filter).Run(s.Session)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(serviceActivityList)

	return serviceActivityList, nil
}

func (s *ActivityStorage) ListServiceActivity(user, service string) (*model.ActivityList, error) {

	var (
		err                 error
		serviceActivityList = new(model.ActivityList)
		activity_filter     = map[string]interface{}{
			"user":    user,
			"service": service,
		}
	)

	res, err := r.Table(ActivityTable).Filter(activity_filter).Run(s.Session)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.IsNil() {
		return nil, nil
	}

	res.One(serviceActivityList)

	return serviceActivityList, nil
}

func (s *ActivityStorage) RemoveByProject(user, project string) error {

	var (
		err            error
		opts           = r.DeleteOpts{ReturnChanges: true}
		activity_filter = map[string]interface{}{
			"user":    user,
			"project": project,
		}
	)

	_, err = r.Table(ActivityTable).Filter(activity_filter).Delete(opts).RunWrite(s.Session)
	if err != nil {
		return err
	}

	return nil
}

func (s *ActivityStorage) RemoveByService(user, service string) error {

	var (
		err            error
		opts           = r.DeleteOpts{ReturnChanges: true}
		activity_filter = map[string]interface{}{
			"user":    user,
			"service": service,
		}
	)

	_, err = r.Table(ActivityTable).Filter(activity_filter).Delete(opts).RunWrite(s.Session)
	if err != nil {
		return err
	}

	return nil
}

func newActivityStorage(session *r.Session) *ActivityStorage {
	r.TableCreate(ActivityTable, r.TableCreateOpts{}).Run(session)
	s := new(ActivityStorage)
	s.Session = session
	return s
}
