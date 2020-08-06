package v1

import (
	"context"
	"github.com/lastbackend/lastbackend/internal/server/api/v1"
)

func (V1) NamespaceCreate(ctx context.Context, req *v1.NamespaceRequest) (*v1.NamespaceResponse, error) {
	return new(v1.NamespaceResponse), nil
}

func (v V1) NamespaceGet(ctx context.Context, query *v1.NamespaceQuery) (*v1.NamespaceResponse, error) {
	panic("implement me")
}

func (v V1) NamespaceList(ctx context.Context, filter *v1.NamespaceFilter) (*v1.NamespaceResponse, error) {
	panic("implement me")
}

func (v V1) NamespaceStatus(ctx context.Context, query *v1.NamespaceQuery) (*v1.NamespaceResponse, error) {
	panic("implement me")
}

func (v V1) NamespacePatch(ctx context.Context, request *v1.NamespaceRequest) (*v1.NamespaceResponse, error) {
	panic("implement me")
}

func (v V1) NamespaceReplace(ctx context.Context, request *v1.NamespaceRequest) (*v1.NamespaceResponse, error) {
	panic("implement me")
}

func (v V1) NamespaceDelete(ctx context.Context, query *v1.NamespaceQuery) (*v1.NamespaceResponse, error) {
	panic("implement me")
}
