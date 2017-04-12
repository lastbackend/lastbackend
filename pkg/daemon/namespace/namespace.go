package namespace

import (
	"context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
	ctx "github.com/lastbackend/lastbackend/pkg/daemon/context"
	"github.com/lastbackend/lastbackend/pkg/daemon/namespace/routes/request"
	"github.com/lastbackend/lastbackend/pkg/util/validator"
)

type namespace struct {
	Context context.Context
}

func New(ctx context.Context) *namespace {
	return &namespace{ctx}
}

func (ns *namespace) List() (*types.NamespaceList, error) {
	var storage = ctx.Get().GetStorage()
	return storage.Namespace().List(ns.Context)
}

func (ns *namespace) Get(id string) (*types.Namespace, error) {
	var (
		err     error
		storage = ctx.Get().GetStorage()
		n       *types.Namespace
	)

	if validator.IsUUID(id) {
		n, err = storage.Namespace().GetByID(ns.Context, id)
	} else {
		n, err = storage.Namespace().GetByName(ns.Context, id)
	}

	return n, err
}

func (ns *namespace) Create(rq *request.RequestNamespaceCreateS) (*types.Namespace, error) {

	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
		n       *types.Namespace
	)

	n, err = storage.Namespace().Insert(ns.Context, rq.Name, rq.Description)
	if err != nil {
		log.Errorf("Error: insert project to db: %s", err.Error())
		return n, err
	}

	return n, nil
}

func (ns *namespace) Update(n *types.Namespace) (*types.Namespace, error) {
	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)

	n, err = storage.Namespace().Update(ns.Context, n)
	if err != nil {
		log.Errorf("Error: update project to db: %s", err.Error())
		return n, err
	}

	return n, nil
}

func (ns *namespace) Remove(id string) error {
	var (
		err     error
		log     = ctx.Get().GetLogger()
		storage = ctx.Get().GetStorage()
	)
	err = storage.Namespace().Remove(ns.Context, id)
	if err != nil {
		log.Error("Error: remove project from db", err)
		return err
	}

	return nil
}
