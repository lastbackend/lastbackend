//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2017] Last.Backend LLC
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

package handler
//
//import (
//	"net/http"
//	c "github.com/lastbackend/lastbackend/pkg/daemon/context"
//	"github.com/lastbackend/lastbackend/pkg/template"
//)
//
//func TemplateListH(w http.ResponseWriter, _ *http.Request) {
//
//	var (
//		er             error
//		ctx            = c.Get()
//		response_empty = func() {
//			w.WriteHeader(http.StatusOK)
//			_, er = w.Write([]byte("[]"))
//			if er != nil {
//				ctx.Log.Error("Error: write response", er.Error())
//				return
//			}
//			return
//		}
//	)
//
//	templates, err := template.List()
//	if err != nil {
//		ctx.Log.Error(err.Error())
//		response_empty()
//		return
//	}
//
//	if templates == nil {
//		response_empty()
//		return
//	}
//
//	response, err := templates.ToJson()
//	if er != nil {
//		ctx.Log.Error(err.Error())
//		response_empty()
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	_, er = w.Write(response)
//	if er != nil {
//		ctx.Log.Error("Error: write response", er.Error())
//		return
//	}
//}
