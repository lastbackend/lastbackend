package v1

import "github.com/lastbackend/lastbackend/pkg/apis/types"

func NewSpec(node *types.Node) *Spec {
	return ToNodeSpec(node.Spec)
}
