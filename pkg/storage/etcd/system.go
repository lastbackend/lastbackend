//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"github.com/lastbackend/lastbackend/pkg/log"
	"github.com/lastbackend/lastbackend/pkg/storage/storage"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
)

const systemStorage = "system"
const systemLeadTTL = 11

// App Service type for interface in interfaces folder
type SystemStorage struct {
	storage.System
}

func (s *SystemStorage) ProcessSet(ctx context.Context, process *types.Process) error {

	log.V(logLevel).Debugf("Storage: System: set process: %#v", process)

	if process == nil {
		err := errors.New("process can not be empty")
		log.V(logLevel).Errorf("Storage: System: set process err: %s", err.Error())
		return err
	}

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: System: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	keyMember := keyCreate(systemStorage, process.Meta.Kind, "process", process.Meta.Hostname)
	if err := client.Upsert(ctx, keyMember, process, nil, systemLeadTTL); err != nil {
		log.V(logLevel).Errorf("Storage: System: upsert process err: %s", err.Error())
		return err
	}

	return nil
}

func (s *SystemStorage) Elect(ctx context.Context, process *types.Process) (bool, error) {

	log.V(logLevel).Debugf("Storage: System: elect process: %#v", process)

	if process == nil {
		err := errors.New("process can not be empty")
		log.V(logLevel).Errorf("Storage: System: elect process err: %s", err.Error())
		return false, err
	}

	var (
		id   string
		err  error
		lead bool
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: System: create client err: %s", err.Error())
		return false, err
	}
	defer destroy()

	key := keyCreate(systemStorage, process.Meta.Kind, "lead")
	err = client.Get(ctx, key, &id)
	if err != nil && (err.Error() != store.ErrEntityNotFound) {
		log.V(logLevel).Errorf("Storage: System: get process lead info err: %s", err.Error())
		return false, err
	}

	if id != "" {
		return false, nil
	}

	if err.Error() == store.ErrEntityNotFound {
		err = client.Create(ctx, key, &process.Meta.ID, nil, systemLeadTTL)
		if err != nil && err.Error() != store.ErrEntityExists {
			log.V(logLevel).Errorf("Storage: System: create process ttl err: %s", err.Error())
			return false, err
		}
		lead = true
	}

	return lead, nil
}

func (s *SystemStorage) ElectUpdate(ctx context.Context, process *types.Process) error {

	log.V(logLevel).Debugf("Storage: System: elect update process: %#v", process)

	if process == nil {
		err := errors.New("process can not be empty")
		log.V(logLevel).Errorf("Storage: System: elect update process err: %s", err.Error())
		return err
	}

	var (
		id  string
		err error
	)

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: System: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(systemStorage, process.Meta.Kind, "lead")
	err = client.Get(ctx, key, &id)
	if err != nil && err.Error() != store.ErrEntityNotFound {
		log.V(logLevel).Errorf("Storage: System: get process lead err: %s", err.Error())
		return err
	}

	if id != process.Meta.ID {
		return nil
	}

	if err := client.Update(ctx, key, &process.Meta.ID, nil, systemLeadTTL); err != nil {
		log.V(logLevel).Errorf("Storage: System: update process ttl err: %s", err.Error())
		return err
	}

	return nil
}

func (s *SystemStorage) ElectWait(ctx context.Context, process *types.Process, lead chan bool) error {

	log.V(logLevel).Debug("Storage: System: elect wait process")

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("Storage: System: create client err: %s", err.Error())
		return err
	}
	defer destroy()

	key := keyCreate(systemStorage, process.Meta.Kind, "lead")
	cb := func(action, key string, val []byte) {
		if action == "create" || action == "update" {
			var id string
			if err := json.Unmarshal(val, &id); err != nil {
				//TODO: return error and start loop over
			}

			if id == process.Meta.ID {
				lead <- true
			} else {
				lead <- false
			}
		}

		if action == "delete" {
			_, err := s.Elect(ctx, process)
			if err != nil {
				log.V(logLevel).Errorf("Storage: System: elect process err: %s", err.Error())
			}
		}
	}

	if err := client.Watch(ctx, key, "", cb); err != nil {
		log.V(logLevel).Errorf("Storage: System: watch process err: %s", err.Error())
		return err
	}

	return nil
}

// Clear system storage
func (s *SystemStorage) Clear(ctx context.Context) error {

	log.V(logLevel).Debugf("storage:etcd:system:> clear")

	client, destroy, err := getClient(ctx)
	if err != nil {
		log.V(logLevel).Errorf("storage:etcd:system:> clear err: %s", err.Error())
		return err
	}
	defer destroy()

	if err := client.DeleteDir(ctx, systemStorage); err != nil {
		log.V(logLevel).Errorf("storage:etcd:system:> clear err: %s", err.Error())
		return err
	}

	return nil
}

func newSystemStorage() *SystemStorage {
	s := new(SystemStorage)
	return s
}
