package namespace

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	ctx "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/daemon/namespace/routes/request"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
)

type Namespace struct {
}

func New() *Namespace {
	return new(Namespace)
}

func (ns *Namespace) List(c context.Context) (*types.NamespaceList, error) {
	var storage = ctx.Get().GetStorage()
	return storage.Namespace().List(c)
}

func (ns *Namespace) Get(c context.Context, id string) (*types.Namespace, error) {
	var (
		err     error
		storage = ctx.Get().GetStorage()
		n       *types.Namespace
	)

	if validator.IsUUID(id) {
		n, err = storage.Namespace().GetByID(c, id)
	} else {
		n, err = storage.Namespace().GetByName(c, id)
	}

	return n, err
}

func (ns *Namespace) Create(c context.Context, rq *request.RequestNamespaceCreateS) (*types.Namespace, error) {

	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
		n       *types.Namespace
	)

	n, err = storage.Namespace().Insert(c, rq.Name, rq.Description)
	if err != nil {
		log.Errorf("Error: insert project to db: %s", err.Error())
		return n, err
	}

	return n, nil
}

func (ns *Namespace) Update(c context.Context, n *types.Namespace) (*types.Namespace, error) {
	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	n, err = storage.Namespace().Update(c, n)
	if err != nil {
		log.Errorf("Error: update project to db: %s", err.Error())
		return n, err
	}

	return n, nil
}

func (ns *Namespace) Remove(c context.Context, id string) error {
	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)
	err = storage.Namespace().Remove(c, id)
	if err != nil {
		log.Error("Error: remove project from db", err)
		return err
	}

	return nil
}
