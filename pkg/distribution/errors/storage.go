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

package errors

import "github.com/lastbackend/lastbackend/pkg/storage/types"

type storage struct {}

func (storage) IsErrEntityExists (err error) bool {
	return err.Error() == types.ErrEntityExists
}

func (storage) IsErrOperationFailure (err error) bool {
	return err.Error() == types.ErrOperationFailure
}

func (storage) IsErrEntityNotFound (err error) bool {
	return err.Error() == types.ErrEntityNotFound
}

func (storage) IsErrStructArgIsNil (err error) bool {
	return err.Error() == types.ErrStructArgIsNil
}

func (storage) IsErrStructOutIsNil (err error) bool {
	return err.Error() == types.ErrStructOutIsNil
}

func (storage) IsErrStructArgIsInvalid (err error) bool {
	return err.Error() == types.ErrStructArgIsInvalid
}
func (storage) IsErrStructOutIsInvalid (err error) bool {
	return err.Error() == types.ErrStructOutIsInvalid
}

func Storage() storage {
	return storage{}
}
