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

package metrics

import (
	"context"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
)

// UnaryInterceptor adds all necessary metrics to incoming gRPC requests
func UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// start values
		operationStart := time.Now()
		operation := filepath.Base(info.FullMethod)

		// the RPC
		resp, err := handler(ctx, req)

		// record the operation
		CRIOOperations.WithLabelValues(operation).Inc()
		CRIOOperationsLatency.WithLabelValues(operation).
			Observe(SinceInMicroseconds(operationStart))

		// record error metric if occurred
		if err != nil {
			CRIOOperationsErrors.WithLabelValues(operation).Inc()
		}

		return resp, err
	}
}
