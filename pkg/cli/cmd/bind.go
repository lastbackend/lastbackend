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

package cmd

import (
	ws "github.com/lastbackend/lastbackend/pkg/cli/cmd/namespace"
	sr "github.com/lastbackend/lastbackend/pkg/cli/cmd/service"
)

func init() {

	// ------------------------- COMMAND -------------------------

	// ----- root -----

	RootCmd.AddCommand(versionCmd, workspace, service, repo)

	// ----- workspace -----

	workspace.AddCommand(ws.CreateWorkspace, ws.ListWorkspace, ws.InfoWorkspace, ws.RemoveWorkspace, ws.SelectWorkspace)

	// ----- service -----

	service.AddCommand(sr.ListService, sr.InfoService, sr.ScaleService, sr.UpdateService, sr.RemoveService, sr.CreateService, sr.LogsService)

	// ------------------------- FLAGS -------------------------

	// ----- workspace -----

	ws.CreateWorkspace.Flags().StringVarP(&ws.Desc, "desc", "d", "", "Set description")

	// ----- service -----

	sr.CreateService.Flags().Int64VarP(&sr.Memory, "memory", "m", 0, "Set memory")
	sr.CreateService.Flags().StringVarP(&sr.Sources, "sources", "s", "", "Set sources")
	sr.UpdateService.Flags().Int64VarP(&sr.Memory, "memory", "m", 0, "Set memory")
}
