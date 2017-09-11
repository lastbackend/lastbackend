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

package routes

import (
	c "github.com/lastbackend/lastbackend/pkg/api/context"
	"github.com/lastbackend/lastbackend/pkg/api/vendors/views/v1"
	"github.com/lastbackend/lastbackend/pkg/common/errors"
	"github.com/lastbackend/lastbackend/pkg/common/types"
	"github.com/lastbackend/lastbackend/pkg/storage/store"
	"github.com/lastbackend/lastbackend/pkg/util/http/utils"
	"github.com/lastbackend/lastbackend/pkg/vendors"
	"github.com/lastbackend/lastbackend/pkg/vendors/docker"
	"github.com/lastbackend/lastbackend/pkg/vendors/interfaces"
	"net/http"
	"github.com/lastbackend/lastbackend/pkg/log"
)

func VendorsH(w http.ResponseWriter, r *http.Request) {
	var (
		storage = c.Get().GetStorage()
	)

	log.Debug("Get vendor services handler")

	vendors, err := storage.Vendor().List(r.Context())
	if err != nil && err.Error() != store.ErrKeyNotFound {
		log.Error(err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewVendorList(vendors).ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func VendorConnectH(w http.ResponseWriter, r *http.Request) {

	var (
		storage = c.Get().GetStorage()
		params  = utils.Vars(r)
		vendor  = params[`vendor`]
		token   = params[`token`]
		client  interfaces.IVCS
	)

	log.Debug("Connect service handler")

	switch vendor {
	case "github":
		client = vendors.GetGitHub(token)
	case "bitbucket":
		client = vendors.GetBitBucket(token)
	case "gitlab":
		client = vendors.GetGitLab(token)
	default:
		log.Error("vendor is not supported yet")
		errors.BadParameter("vendor").Http(w)
		return
	}

	serviceUser, err := client.GetUser()
	if err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
		return
	}

	vendorInfo := client.VendorInfo()

	if err := storage.Vendor().Insert(r.Context(), serviceUser.Username, vendorInfo.Name, vendorInfo.Host, serviceUser.ServiceID, vendorInfo.Token); err != nil {
		log.Error(err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	vendors, err := storage.Vendor().List(r.Context())
	if err != nil {
		log.Error(err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewVendorList(vendors).ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func VendorDisconnectH(w http.ResponseWriter, r *http.Request) {

	var (
		storage = c.Get().GetStorage()
		params  = utils.Vars(r)
		vendor  = params[`vendor`]
	)

	log.Debug("Disconnect service handler")

	if err := storage.Vendor().Remove(r.Context(), vendor); err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
		return
	}

	vendors, err := storage.Vendor().List(r.Context())
	if err != nil {
		log.Error(err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := v1.NewVendorList(vendors).ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func VCSRepositoryListH(w http.ResponseWriter, r *http.Request) {

	var (
		storage = c.Get().GetStorage()
		client  interfaces.IVCS
		params  = utils.Vars(r)
		vendor  = params[`vendor`]
	)

	vndr, err := storage.Vendor().Get(r.Context(), vendor)
	if err != nil {
		log.Error(err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	switch vendor {
	case "github":
		client = vendors.GetGitHub(vndr.Token.AccessToken)
	case "bitbucket":
		client = vendors.GetBitBucket(vndr.Token.AccessToken)
	case "gitlab":
		client = vendors.GetGitLab(vndr.Token.AccessToken)
	default:
		log.Error("vendor is not supported yet")
		errors.BadParameter("vendor").Http(w)
		return
	}

	repos, err := client.ListRepositories(vndr.Username, false)
	if err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
		return
	}

	rp := types.VCSRepositoryList{}
	rp.Convert(repos)
	response, err := rp.ToJson()

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func VCSBranchListH(w http.ResponseWriter, r *http.Request) {

	var (
		storage = c.Get().GetStorage()
		client  interfaces.IVCS
		params  = utils.Vars(r)
		vendor  = params[`vendor`]
		repo    = r.URL.Query().Get(`repo`)
	)

	vndr, err := storage.Vendor().Get(r.Context(), vendor)
	if err != nil {
		log.Error(err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	switch vendor {
	case "github":
		client = vendors.GetGitHub(vndr.Token.AccessToken)
	case "bitbucket":
		client = vendors.GetBitBucket(vndr.Token.AccessToken)
	case "gitlab":
		client = vendors.GetGitLab(vndr.Token.AccessToken)
	default:
		log.Error("vendor is not supported yet")
		errors.BadParameter("vendor").Http(w)
		return
	}

	branches, err := client.ListBranches(vndr.Username, repo)
	if err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
		return
	}

	br := types.VCSBranchList{}
	br.Convert(branches)
	response, err := br.ToJson()

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func DockerRepositorySearchH(w http.ResponseWriter, r *http.Request) {

	var (
		params = r.URL.Query()
		name   = params.Get("name")
	)

	log.Debug("Search docker repository handler")

	repoListModel, err := docker.GetRepository(name)
	if err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := repoListModel.ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}

func DockerRepositoryTagListH(w http.ResponseWriter, r *http.Request) {

	var (
		params = r.URL.Query()
		owner  = params.Get("owner")
		name   = params.Get("name")
	)

	log.Debug("List docker repository tags handler")

	tagListModel, err := docker.ListTag(owner, name)
	if err != nil {
		log.Error(err)
		errors.HTTP.InternalServerError(w)
		return
	}

	response, err := tagListModel.ToJson()
	if err != nil {
		log.Error("Error: convert struct to json", err.Error())
		errors.HTTP.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(response); err != nil {
		log.Error("Error: write response", err.Error())
		return
	}
}
