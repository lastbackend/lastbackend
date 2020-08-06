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

package cntx

import (
	"context"
	"errors"
	"time"
)

var (
	ErrValueNotSet    = errors.New("value not set")
	ErrValueNotBool   = errors.New("non-bool value")
	ErrValueNotString = errors.New("non-string value")
	ErrValueNotInt    = errors.New("non-int value")
	ErrNotCastType    = errors.New("typecast error")
)

func NewContext() context.Context {
	return context.Background()
}

func NewContextTODO() context.Context {
	return context.TODO()
}

func NewContextWithCancel() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func NewContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

func NewContextWithDeadline(d time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), d)
}

func IsSetInContext(ctx context.Context, key string) bool {
	v := ctx.Value(key)
	if v == nil {
		return false
	}
	return true
}

func SetContext(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetBoolFromContext(ctx context.Context, key string) bool {
	v := ctx.Value(key)
	if v == nil {
		return false
	}
	if val, ok := v.(bool); ok {
		return val
	}
	return false
}

func GetStringFromContext(ctx context.Context, key string) string {
	v := ctx.Value(key)
	if v == nil {
		return ""
	}
	if val, ok := v.(string); ok {
		return val
	}
	return ""
}

func GetIntFromContext(ctx context.Context, key string) int {
	v := ctx.Value(key)
	if v == nil {
		return 0
	}
	if val, ok := v.(int); ok {
		return val
	}
	return 0
}

func GetDurationFromContext(ctx context.Context, key string) time.Duration {
	v := ctx.Value(key)
	if v == nil {
		return 0
	}
	if val, ok := v.(time.Duration); ok {
		return val
	}
	return 0
}

func GetStringSliceFromContext(ctx context.Context, key string) []string {
	v := ctx.Value(key)
	if v == nil {
		return make([]string, 0)
	}
	if val, ok := v.([]string); ok {
		return val
	}
	return make([]string, 0)
}

func GetInt64SliceSliceFromContext(ctx context.Context, key string) []int64 {
	v := ctx.Value(key)
	if v == nil {
		return make([]int64, 0)
	}
	if val, ok := v.([]int64); ok {
		return val
	}
	return make([]int64, 0)
}
