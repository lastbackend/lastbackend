package v1

import "github.com/lastbackend/lastbackend/pkg/apis/types"

func ToImageSpec(spec types.ImageSpec) ImageSpec {
	return ImageSpec{
		Name: spec.Name,
		Pull: spec.Pull,
		Auth: spec.Auth,
	}
}

func FromImageSpec(spec ImageSpec) types.ImageSpec {
	return types.ImageSpec{
		Name: spec.Name,
		Pull: spec.Pull,
		Auth: spec.Auth,
	}
}
