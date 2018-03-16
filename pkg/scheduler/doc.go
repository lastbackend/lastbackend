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

package scheduler

// watch service pods
// generate new spec after new pod creation
// allocate node for new spec

// watch nodes online states, if node goes offline
// more than 30 seconds, move specs to another node

// node/<hostname>/alive
// node/<hostname>/spec/pod

// watch builders online states, if builder goes offline
// more than 30 seconds, move build to another builder

// builder/elected:<hostname>
// builder/builds/<id>:<state>
