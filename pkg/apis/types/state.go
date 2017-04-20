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

package types

const StateWaiting   = "waiting"
const StateWarning   = "warning"
const StateReady     = "ready"


const StateCreated = "created"
const StateStarted = "started"
const StateStopped = "stopped"
const StateRestarted = "restarted"

const StateDestroy = "destroy"
const StateDestroyed = "destroyed"

const StateExited = "exited"
const StateRunning = "running"
const StateError = "error"

const EventStateStart = "start"
const EventStateStop = "stop"
const EventStateRestart = "restart"
const EventStateDestroy = "destroy"
const EventStateCreated = "created"
const EventStateKill = "kill"
