//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
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

package errors

import (
	"errors"
)

const (
	ErrEntityExists          = "entity exists"
	ErrOperationFailure      = "operation failure"
	ErrEntityNotFound        = "entity not found"
	ErrStructArgIsNil        = "input structure is nil"
	ErrStructOutIsNil        = "output structure is nil"
	ErrStructArgIsInvalid    = "input structure is invalid"
	ErrStructOutIsInvalid    = "output structure is invalid"
	ErrStructOutIsNotPointer = "output structure is not pointer"
)

type storage struct{}

func (storage) IsErrEntityExists(err error) bool {
	return err.Error() == ErrEntityExists
}

func (storage) NewErrEntityExists() error {
	return errors.New(ErrEntityExists)
}

func (storage) IsErrOperationFailure(err error) bool {
	return err.Error() == ErrOperationFailure
}

func (storage) NewErrOperationFailure() error {
	return errors.New(ErrOperationFailure)
}

func (storage) IsErrEntityNotFound(err error) bool {
	return err.Error() == ErrEntityNotFound
}

func (storage) NewErrEntityNotFound() error {
	return errors.New(ErrEntityNotFound)
}

func (storage) IsErrStructArgIsNil(err error) bool {
	return err.Error() == ErrStructArgIsNil
}

func (storage) NewErrStructArgIsNil() error {
	return errors.New(ErrStructArgIsNil)
}

func (storage) IsErrStructOutIsNil(err error) bool {
	return err.Error() == ErrStructOutIsNil
}

func (storage) NewErrStructOutIsNil() error {
	return errors.New(ErrStructOutIsNil)
}

func (storage) IsErrStructArgIsInvalid(err error) bool {
	return err.Error() == ErrStructArgIsInvalid
}

func (storage) NewErrStructArgIsInvalid() error {
	return errors.New(ErrStructArgIsInvalid)
}

func (storage) IsErrStructOutIsInvalid(err error) bool {
	return err.Error() == ErrStructOutIsInvalid
}

func (storage) NewErrStructOutIsInvalid() error {
	return errors.New(ErrStructOutIsInvalid)
}
func (storage) NewErrStructOutIsNotPointer() error {
	return errors.New(ErrStructOutIsNotPointer)
}

func Storage() storage {
	return storage{}
}
