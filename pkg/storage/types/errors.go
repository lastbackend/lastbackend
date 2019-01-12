//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2019] Last.Backend LLC
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

package types

import "github.com/lastbackend/lastbackend/pkg/distribution/errors"

const (
	ErrEntityExists          = errors.ErrEntityExists
	ErrOperationFailure      = errors.ErrOperationFailure
	ErrEntityNotFound        = errors.ErrEntityNotFound
	ErrStructArgIsNil        = errors.ErrStructArgIsNil
	ErrStructOutIsNil        = errors.ErrStructOutIsNil
	ErrStructArgIsInvalid    = errors.ErrStructArgIsInvalid
	ErrStructOutIsInvalid    = errors.ErrStructOutIsInvalid
	ErrStructOutIsNotPointer = errors.ErrStructOutIsNotPointer
)
