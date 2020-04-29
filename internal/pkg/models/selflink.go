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

package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	DefaultNamespace   = "default"
	ErrInvalidSelfLink = "invalid selflink"
)

type SelfLink interface {
	Parse(str string) error
	String() string
	Namespace() *NamespaceSelfLink
	Name() string
}

func NewSelfLink(namespace, kind, name string) SelfLink {
	switch kind {
	case KindNamespace:
		return NewNamespaceSelfLink(name)
	}

	return nil
}

type SelfLinkParent struct {
	Kind     string
	SelfLink SelfLink
}

type NamespaceSelfLink struct {
	string
}

func (sl *NamespaceSelfLink) Parse(namespace string) error {
	sl.string = namespace
	return nil
}

func (sl *NamespaceSelfLink) String() string {
	return sl.string
}

func (sl *NamespaceSelfLink) Parent() (string, SelfLink) {
	return EmptyString, nil
}

func (sl *NamespaceSelfLink) Namespace() *NamespaceSelfLink {
	return sl
}

func (sl *NamespaceSelfLink) Name() string {
	return sl.string
}

func (sl NamespaceSelfLink) MarshalJSON() ([]byte, error) {

	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")

	return buffer.Bytes(), nil
}

func (sl *NamespaceSelfLink) UnmarshalJSON(b []byte) error {

	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	return sl.Parse(link)
}

func NewNamespaceSelfLink(name string) *NamespaceSelfLink {

	var sl = new(NamespaceSelfLink)
	sl.string = name
	return sl
}

type ResourceSelfLink struct {
	string
	SelfLink
	namespace *NamespaceSelfLink
	name      string
	kind      string
}

func (sl *ResourceSelfLink) Parse(selflink string) error {

	parts := strings.Split(selflink, ":")
	sl.string = selflink
	if len(parts) < 2 {
		sl.namespace = NewNamespaceSelfLink(DefaultNamespace)
		sl.name = parts[0]
		return nil
	}

	sl.namespace = NewNamespaceSelfLink(parts[0])
	sl.name = parts[1]
	return nil
}

func (sl *ResourceSelfLink) String() string {
	return sl.string
}

func (sl *ResourceSelfLink) Namespace() *NamespaceSelfLink {
	return sl.namespace
}

func (sl *ResourceSelfLink) Name() string {
	return sl.name
}

func (sl ResourceSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *ResourceSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	return sl.Parse(link)
}

func NewResourceSelfLink(namespace, kind, resource string) *ResourceSelfLink {

	sl := new(ResourceSelfLink)

	link := fmt.Sprintf("%s:%s:%s", namespace, kind, resource)

	sl.string = link
	sl.namespace = NewNamespaceSelfLink(namespace)
	sl.name = resource
	sl.kind = kind

	return sl
}

type ServiceSelfLink struct {
	string
	SelfLink
	parent SelfLinkParent
	name   string
}

func (sl *ServiceSelfLink) Parse(selflink string) error {

	parts := strings.Split(selflink, ":")
	sl.string = selflink
	if len(parts) < 2 {
		sl.parent = SelfLinkParent{
			Kind:     KindNamespace,
			SelfLink: NewNamespaceSelfLink(DEFAULT_NAMESPACE),
		}
		sl.name = parts[0]
		return nil
	}

	sl.parent = SelfLinkParent{
		Kind:     KindNamespace,
		SelfLink: NewNamespaceSelfLink(parts[0]),
	}

	sl.name = parts[1]
	return nil
}

func (sl *ServiceSelfLink) String() string {
	return sl.string
}

func (sl *ServiceSelfLink) Parent() (string, SelfLink) {
	return sl.parent.Kind, sl.parent.SelfLink
}

func (sl *ServiceSelfLink) Namespace() *NamespaceSelfLink {
	return sl.parent.SelfLink.(*NamespaceSelfLink)
}

func (sl *ServiceSelfLink) Name() string {
	return sl.name
}

func (sl ServiceSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *ServiceSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	return sl.Parse(link)
}

func NewServiceSelfLink(namespace, service string) *ServiceSelfLink {

	sl := new(ServiceSelfLink)

	link := fmt.Sprintf("%s:%s", namespace, service)

	sl.string = link
	sl.parent.Kind = KindNamespace
	sl.parent.SelfLink = NewNamespaceSelfLink(namespace)
	sl.name = service

	return sl
}

type JobSelfLink struct {
	string
	SelfLink
	parent SelfLinkParent
	name   string
}

func (sl *JobSelfLink) Parse(selflink string) error {

	parts := strings.Split(selflink, ":")

	sl.string = selflink
	if len(parts) < 2 {
		sl.parent = SelfLinkParent{
			Kind:     KindNamespace,
			SelfLink: NewNamespaceSelfLink(DEFAULT_NAMESPACE),
		}
		sl.name = parts[0]
		return nil
	}

	sl.parent = SelfLinkParent{
		Kind:     KindNamespace,
		SelfLink: NewNamespaceSelfLink(parts[0]),
	}

	sl.name = parts[1]
	return nil
}

func (sl *JobSelfLink) String() string {
	return sl.string
}

func (sl *JobSelfLink) Parent() (string, SelfLink) {
	return sl.parent.Kind, sl.parent.SelfLink
}

func (sl *JobSelfLink) Namespace() *NamespaceSelfLink {
	return sl.parent.SelfLink.(*NamespaceSelfLink)
}

func (sl *JobSelfLink) Name() string {
	return sl.name
}

func (sl JobSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *JobSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	return sl.Parse(link)
}

func NewJobSelfLink(namespace, service string) *JobSelfLink {

	sl := new(JobSelfLink)

	link := fmt.Sprintf("%s:%s", namespace, service)

	sl.string = link
	sl.parent.Kind = KindNamespace
	sl.parent.SelfLink = NewNamespaceSelfLink(namespace)
	sl.name = service

	return sl
}

type DeploymentSelfLink struct {
	string
	SelfLink
	namespace *NamespaceSelfLink
	parent    SelfLinkParent
	name      string
}

func (sl *DeploymentSelfLink) Parse(selflink string) error {

	parts := strings.Split(selflink, ":")

	sl.string = selflink
	sl.parent.Kind = KindService
	if len(parts) == 2 {
		sl.namespace = NewNamespaceSelfLink(DEFAULT_NAMESPACE)
		sl.parent.SelfLink = NewServiceSelfLink(DEFAULT_NAMESPACE, parts[0])
		sl.name = parts[1]
		return nil
	}

	if len(parts) == 1 {
		return errors.New(ErrInvalidSelfLink)
	}

	sl.namespace = NewNamespaceSelfLink(parts[0])
	sl.parent.SelfLink = NewServiceSelfLink(parts[0], parts[1])
	sl.name = parts[2]

	return nil
}

func (sl *DeploymentSelfLink) String() string {
	return sl.string
}

func (sl *DeploymentSelfLink) Namespace() *NamespaceSelfLink {
	return sl.namespace
}

func (sl *DeploymentSelfLink) Parent() (string, SelfLink) {
	return sl.parent.Kind, sl.parent.SelfLink
}

func (sl *DeploymentSelfLink) Name() string {
	return sl.name
}

func (sl DeploymentSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *DeploymentSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	return sl.Parse(link)
}

func NewDeploymentSelfLink(namespace, service, deployment string) *DeploymentSelfLink {

	sl := new(DeploymentSelfLink)

	link := fmt.Sprintf("%s:%s:%s", namespace, service, deployment)

	sl.string = link
	sl.namespace = NewNamespaceSelfLink(namespace)
	sl.parent.Kind = KindService
	sl.parent.SelfLink = NewServiceSelfLink(namespace, service)
	sl.name = deployment

	return sl
}

type EndpointSelfLink struct {
	string
	SelfLink
	namespace *NamespaceSelfLink
	parent    SelfLinkParent
	name      string
}

func (sl *EndpointSelfLink) Parse(selflink string) error {

	parts := strings.Split(selflink, ":")

	sl.string = selflink
	sl.parent.Kind = KindService
	if len(parts) == 1 {
		sl.namespace = NewNamespaceSelfLink(DEFAULT_NAMESPACE)
		sl.parent.SelfLink = NewServiceSelfLink(DEFAULT_NAMESPACE, parts[0])
		sl.name = parts[0]
		return nil
	}

	sl.namespace = NewNamespaceSelfLink(parts[0])
	sl.parent.SelfLink = NewServiceSelfLink(parts[0], parts[1])
	sl.name = parts[1]

	return nil
}

func (sl *EndpointSelfLink) String() string {
	return sl.string
}

func (sl *EndpointSelfLink) Namespace() *NamespaceSelfLink {
	return sl.namespace
}

func (sl *EndpointSelfLink) Parent() (string, SelfLink) {
	return sl.parent.Kind, sl.parent.SelfLink
}

func (sl *EndpointSelfLink) Name() string {
	return sl.name
}

func (sl EndpointSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *EndpointSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	return sl.Parse(link)
}

func NewEndpointSelfLink(namespace, service string) *EndpointSelfLink {

	sl := new(EndpointSelfLink)

	link := fmt.Sprintf("%s:%s", namespace, service)

	sl.string = link
	sl.namespace = NewNamespaceSelfLink(namespace)
	sl.parent.Kind = KindService
	sl.parent.SelfLink = NewServiceSelfLink(namespace, service)
	sl.name = service

	return sl
}

type TaskSelfLink struct {
	string
	SelfLink
	namespace *NamespaceSelfLink
	parent    SelfLinkParent
	name      string
}

func (sl *TaskSelfLink) Parse(selflink string) error {

	parts := strings.Split(selflink, ":")

	sl.string = selflink
	sl.parent.Kind = KindService
	if len(parts) == 2 {
		sl.namespace = NewNamespaceSelfLink(DEFAULT_NAMESPACE)
		sl.parent.SelfLink = NewJobSelfLink(DEFAULT_NAMESPACE, parts[0])
		sl.name = parts[1]
		return nil
	}

	if len(parts) == 1 {
		return errors.New(ErrInvalidSelfLink)
	}

	sl.namespace = NewNamespaceSelfLink(parts[0])
	sl.parent.SelfLink = NewJobSelfLink(parts[0], parts[1])
	sl.name = parts[2]

	return nil
}

func (sl *TaskSelfLink) String() string {
	return sl.string
}

func (sl *TaskSelfLink) Namespace() *NamespaceSelfLink {
	return sl.namespace
}

func (sl *TaskSelfLink) Parent() (string, SelfLink) {
	return sl.parent.Kind, sl.parent.SelfLink
}

func (sl *TaskSelfLink) Name() string {
	return sl.name
}

func (sl TaskSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *TaskSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	return sl.Parse(link)
}

func NewTaskSelfLink(namespace, job, task string) *TaskSelfLink {

	sl := new(TaskSelfLink)

	link := fmt.Sprintf("%s:%s:%s", namespace, job, task)

	sl.string = link
	sl.namespace = NewNamespaceSelfLink(namespace)
	sl.parent.Kind = KindJob
	sl.parent.SelfLink = NewJobSelfLink(namespace, job)
	sl.name = task

	return sl
}

type PodSelfLink struct {
	string
	SelfLink
	namespace *NamespaceSelfLink
	parent    SelfLinkParent
	name      string
}

func (sl *PodSelfLink) Parse(selflink string) error {

	sl.string = selflink
	parts := strings.Split(selflink, ":")
	if len(parts) == 4 {

		sl.namespace = NewNamespaceSelfLink(parts[0])

		// get parent prefix from selflink

		if strings.HasPrefix(parts[2], "d_") {
			parts[2] = strings.TrimPrefix(parts[2], "d_")
			sl.parent.Kind = KindDeployment
			sl.parent.SelfLink = NewDeploymentSelfLink(parts[0], parts[1], parts[2])
		}

		if strings.HasPrefix(parts[2], "t_") {
			parts[2] = strings.TrimPrefix(parts[2], "t_")
			sl.parent.Kind = KindTask
			sl.parent.SelfLink = NewTaskSelfLink(parts[0], parts[1], parts[2])
		}

		sl.name = parts[3]
		return nil
	}

	if len(parts) == 3 {
		sl.namespace = NewNamespaceSelfLink(DEFAULT_NAMESPACE)

		// get parent prefix from selflink

		sl.parent.Kind = KindDeployment
		sl.parent.SelfLink = NewDeploymentSelfLink(parts[0], parts[1], parts[2])

		sl.parent.Kind = KindTask
		sl.parent.SelfLink = NewTaskSelfLink(parts[0], parts[1], parts[2])

		sl.name = parts[2]
	}

	if len(parts) == 2 {
		sl.namespace = NewNamespaceSelfLink(parts[0])
		sl.name = parts[1]
	}

	if len(parts) < 2 {
		return errors.New(ErrInvalidSelfLink)
	}

	return nil
}

func (sl *PodSelfLink) String() string {
	return sl.string
}

func (sl *PodSelfLink) Namespace() *NamespaceSelfLink {
	return sl.namespace
}

func (sl *PodSelfLink) Parent() (string, SelfLink) {
	return sl.parent.Kind, sl.parent.SelfLink
}

func (sl *PodSelfLink) Name() string {
	return sl.name
}

func (sl PodSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *PodSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	return sl.Parse(link)
}

func NewPodSelfLink(kind, parent, name string) (*PodSelfLink, error) {

	sl := new(PodSelfLink)

	switch kind {
	case KindDeployment:

		prefix := "d"

		dsl := DeploymentSelfLink{}
		if err := dsl.Parse(parent); err != nil {
			return nil, err
		}
		_, ssl := dsl.Parent()
		service := ssl.Name()

		sl.namespace = dsl.Namespace()
		sl.parent.Kind = KindDeployment
		sl.parent.SelfLink = &dsl
		sl.string = fmt.Sprintf("%s:%s:%s_%s:%s", dsl.Namespace().String(), service, prefix, dsl.Name(), name)

	case KindTask:
		prefix := "t"
		tsl := TaskSelfLink{}
		if err := tsl.Parse(parent); err != nil {
			return nil, err
		}
		_, ssl := tsl.Parent()
		job := ssl.Name()
		sl.namespace = tsl.Namespace()
		sl.parent.Kind = KindDeployment
		sl.parent.SelfLink = &tsl
		sl.string = fmt.Sprintf("%s:%s:%s_%s:%s", tsl.Namespace().String(), job, prefix, tsl.Name(), name)
	}

	return sl, nil
}

type ConfigSelfLink struct {
	string
	SelfLink
	parent SelfLinkParent
	name   string
}

func (sl *ConfigSelfLink) Parse(selflink string) error {

	parts := strings.Split(selflink, ":")

	sl.string = selflink
	if len(parts) < 2 {
		sl.parent = SelfLinkParent{
			Kind:     KindNamespace,
			SelfLink: NewNamespaceSelfLink(DEFAULT_NAMESPACE),
		}
		sl.name = parts[0]
		return nil
	}

	sl.parent = SelfLinkParent{
		Kind:     KindNamespace,
		SelfLink: NewNamespaceSelfLink(parts[0]),
	}

	sl.name = parts[1]
	return nil
}

func (sl *ConfigSelfLink) String() string {
	return sl.string
}

func (sl *ConfigSelfLink) Parent() (string, SelfLink) {
	return sl.parent.Kind, sl.parent.SelfLink
}

func (sl *ConfigSelfLink) Namespace() *NamespaceSelfLink {
	return sl.parent.SelfLink.(*NamespaceSelfLink)
}

func (sl *ConfigSelfLink) Name() string {
	return sl.name
}

func (sl ConfigSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *ConfigSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	return sl.Parse(link)
}

func NewConfigSelfLink(namespace, config string) *ConfigSelfLink {

	sl := new(ConfigSelfLink)

	link := fmt.Sprintf("%s:%s", namespace, config)

	sl.string = link
	sl.parent.Kind = KindNamespace
	sl.parent.SelfLink = NewNamespaceSelfLink(namespace)
	sl.name = config

	return sl
}

type SecretSelfLink struct {
	string
	SelfLink
	parent SelfLinkParent
	name   string
}

func (sl *SecretSelfLink) Parse(selflink string) error {

	parts := strings.Split(selflink, ":")

	sl.string = selflink
	if len(parts) < 2 {
		sl.parent = SelfLinkParent{
			Kind:     KindNamespace,
			SelfLink: NewNamespaceSelfLink(DEFAULT_NAMESPACE),
		}
		sl.name = parts[0]
		return nil
	}

	sl.parent = SelfLinkParent{
		Kind:     KindNamespace,
		SelfLink: NewNamespaceSelfLink(parts[0]),
	}

	sl.name = parts[1]
	return nil
}

func (sl *SecretSelfLink) String() string {
	return sl.string
}

func (sl *SecretSelfLink) Parent() (string, SelfLink) {
	return sl.parent.Kind, sl.parent.SelfLink
}

func (sl *SecretSelfLink) Namespace() *NamespaceSelfLink {
	return sl.parent.SelfLink.(*NamespaceSelfLink)
}

func (sl *SecretSelfLink) Name() string {
	return sl.name
}

func (sl SecretSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *SecretSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	return sl.Parse(link)
}

func NewSecretSelfLink(namespace, secret string) *SecretSelfLink {

	sl := new(SecretSelfLink)

	link := fmt.Sprintf("%s:%s", namespace, secret)

	sl.string = link
	sl.parent.Kind = KindNamespace
	sl.parent.SelfLink = NewNamespaceSelfLink(namespace)
	sl.name = secret

	return sl
}

type VolumeSelfLink struct {
	string
	SelfLink
	parent SelfLinkParent
	name   string
}

func (sl *VolumeSelfLink) Parse(selflink string) error {

	parts := strings.Split(selflink, ":")

	sl.string = selflink
	if len(parts) < 2 {
		sl.parent = SelfLinkParent{
			Kind:     KindNamespace,
			SelfLink: NewNamespaceSelfLink(DEFAULT_NAMESPACE),
		}
		sl.name = parts[0]
		return nil
	}

	sl.parent = SelfLinkParent{
		Kind:     KindNamespace,
		SelfLink: NewNamespaceSelfLink(parts[0]),
	}

	sl.name = parts[1]
	return nil
}

func (sl *VolumeSelfLink) String() string {
	return sl.string
}

func (sl *VolumeSelfLink) Parent() (string, SelfLink) {
	return sl.parent.Kind, sl.parent.SelfLink
}

func (sl *VolumeSelfLink) Namespace() *NamespaceSelfLink {
	return sl.parent.SelfLink.(*NamespaceSelfLink)
}

func (sl *VolumeSelfLink) Name() string {
	return sl.name
}

func (sl VolumeSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *VolumeSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	return sl.Parse(link)
}

func NewVolumeSelfLink(namespace, volume string) *VolumeSelfLink {

	sl := new(VolumeSelfLink)

	link := fmt.Sprintf("%s:%s", namespace, volume)

	sl.string = link
	sl.parent.Kind = KindNamespace
	sl.parent.SelfLink = NewNamespaceSelfLink(namespace)
	sl.name = volume

	return sl
}

type RouteSelfLink struct {
	string
	SelfLink
	parent SelfLinkParent
	name   string
}

func (sl *RouteSelfLink) Parse(selflink string) error {

	parts := strings.Split(selflink, ":")

	sl.string = selflink
	if len(parts) < 2 {
		sl.parent = SelfLinkParent{
			Kind:     KindNamespace,
			SelfLink: NewNamespaceSelfLink(DEFAULT_NAMESPACE),
		}
		sl.name = parts[0]
		return nil
	}

	sl.parent = SelfLinkParent{
		Kind:     KindNamespace,
		SelfLink: NewNamespaceSelfLink(parts[0]),
	}

	sl.name = parts[1]
	return nil
}

func (sl *RouteSelfLink) String() string {
	return sl.string
}

func (sl *RouteSelfLink) Parent() (string, SelfLink) {
	return sl.parent.Kind, sl.parent.SelfLink
}

func (sl *RouteSelfLink) Namespace() *NamespaceSelfLink {
	return sl.parent.SelfLink.(*NamespaceSelfLink)
}

func (sl *RouteSelfLink) Name() string {
	return sl.name
}

func (sl RouteSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *RouteSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	return sl.Parse(link)
}

func NewRouteSelfLink(namespace, route string) *RouteSelfLink {

	sl := new(RouteSelfLink)

	link := fmt.Sprintf("%s:%s", namespace, route)

	sl.string = link
	sl.parent.Kind = KindNamespace
	sl.parent.SelfLink = NewNamespaceSelfLink(namespace)
	sl.name = route

	return sl
}

type SubnetSelfLink struct {
	string
	cidr string
}

func (sl *SubnetSelfLink) Parse(selflink string) {
	sl.cidr = selflink
	sl.string = selflink
}

func (sl *SubnetSelfLink) String() string {
	return strings.Replace(sl.string, "/", ":", -1)
}

func (sl *SubnetSelfLink) Hostname() string {
	return sl.cidr
}

func (sl SubnetSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *SubnetSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	sl.Parse(link)
	return nil
}

func NewSubnetSelfLink(cidr string) *SubnetSelfLink {

	sl := new(SubnetSelfLink)

	link := fmt.Sprintf("%s", cidr)

	sl.string = strings.Replace(link, "/", ":", -1)
	sl.cidr = cidr

	return sl
}

type NodeSelfLink struct {
	string
	hostname string
}

func (sl *NodeSelfLink) Parse(selflink string) {
	sl.hostname = selflink
	sl.string = selflink
}

func (sl *NodeSelfLink) String() string {
	return sl.string
}

func (sl *NodeSelfLink) Hostname() string {
	return sl.hostname
}

func (sl NodeSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *NodeSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	sl.Parse(link)
	return nil
}

func NewNodeSelfLink(hostname string) *NodeSelfLink {

	sl := new(NodeSelfLink)

	link := fmt.Sprintf("%s", hostname)

	sl.string = link
	sl.hostname = hostname

	return sl
}

type IngressSelfLink struct {
	string
	hostname string
}

func (sl *IngressSelfLink) Parse(selflink string) {
	sl.hostname = selflink
	sl.string = selflink
}

func (sl *IngressSelfLink) String() string {
	return sl.string
}

func (sl *IngressSelfLink) Hostname() string {
	return sl.hostname
}

func (sl IngressSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *IngressSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	sl.Parse(link)
	return nil
}

func NewIngressSelfLink(hostname string) *IngressSelfLink {

	sl := new(IngressSelfLink)

	link := fmt.Sprintf("%s", hostname)

	sl.string = link
	sl.hostname = hostname

	return sl
}

type DiscoverySelfLink struct {
	string
	hostname string
}

func (sl *DiscoverySelfLink) Parse(selflink string) {
	sl.hostname = selflink
	sl.string = selflink
}

func (sl *DiscoverySelfLink) String() string {
	return sl.string
}

func (sl *DiscoverySelfLink) Hostname() string {
	return sl.hostname
}

func (sl DiscoverySelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *DiscoverySelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	sl.Parse(link)
	return nil
}

func NewDiscoverySelfLink(hostname string) *DiscoverySelfLink {

	sl := new(DiscoverySelfLink)

	link := fmt.Sprintf("%s", hostname)

	sl.string = link
	sl.hostname = hostname

	return sl
}

type ExporterSelfLink struct {
	string
	hostname string
}

func (sl *ExporterSelfLink) Parse(selflink string) {
	sl.hostname = selflink
	sl.string = selflink
}

func (sl *ExporterSelfLink) String() string {
	return sl.string
}

func (sl *ExporterSelfLink) Hostname() string {
	return sl.hostname
}

func (sl ExporterSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *ExporterSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	sl.Parse(link)
	return nil
}

func NewExporterSelfLink(hostname string) *ExporterSelfLink {

	sl := new(ExporterSelfLink)

	link := fmt.Sprintf("%s", hostname)

	sl.string = link
	sl.hostname = hostname

	return sl
}

type APISelfLink struct {
	string
	hostname string
}

func (sl *APISelfLink) Parse(selflink string) {
	sl.hostname = selflink
	sl.string = selflink
}

func (sl *APISelfLink) String() string {
	return sl.string
}

func (sl *APISelfLink) Hostname() string {
	return sl.hostname
}

func (sl APISelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *APISelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	sl.Parse(link)
	return nil
}

func NewAPISelfLink(hostname string) *APISelfLink {

	sl := new(APISelfLink)

	link := fmt.Sprintf("%s", hostname)

	sl.string = link
	sl.hostname = hostname

	return sl
}

type ControllerSelfLink struct {
	string
	hostname string
}

func (sl *ControllerSelfLink) Parse(selflink string) {
	sl.hostname = selflink
	sl.string = selflink
}

func (sl *ControllerSelfLink) String() string {
	return sl.string
}

func (sl *ControllerSelfLink) Hostname() string {
	return sl.hostname
}

func (sl ControllerSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *ControllerSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	sl.Parse(link)
	return nil
}

func NewControllerSelfLink(hostname string) *ControllerSelfLink {

	sl := new(ControllerSelfLink)

	link := fmt.Sprintf("%s", hostname)

	sl.string = link
	sl.hostname = hostname

	return sl
}

type ProcessSelfLink struct {
	string
	name string
	pid  int
	kind string
}

func (sl *ProcessSelfLink) Parse(selflink string) error {

	var err error

	parts := strings.Split(selflink, ":")
	if len(parts) != 3 {
		return errors.New(ErrInvalidSelfLink)
	}

	sl.kind = parts[0]
	sl.name = parts[1]

	sl.pid, err = strconv.Atoi(parts[2])
	if err != nil {
		return err
	}

	sl.string = selflink
	return nil
}

func (sl *ProcessSelfLink) String() string {
	return sl.string
}

func (sl *ProcessSelfLink) Hostname() string {
	return sl.name
}

func (sl ProcessSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	buffer.WriteString(sl.string)
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *ProcessSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	return sl.Parse(link)
}

func NewProcessSelfLink(kind, name string, pid int) *ProcessSelfLink {

	sl := new(ProcessSelfLink)

	link := fmt.Sprintf("%s:%s:%d", kind, name, pid)

	sl.string = link
	sl.kind = kind
	sl.name = name
	sl.pid = pid

	return sl
}

type ClusterSelfLink struct {
	string
	name string
}

func (sl *ClusterSelfLink) Parse(selflink string) {
	sl.name = selflink
	sl.string = selflink
}

func (sl *ClusterSelfLink) String() string {
	return sl.string
}

func (sl ClusterSelfLink) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("\"")
	if sl.string != EmptyString {
		buffer.WriteString(sl.string)
	}
	buffer.WriteString("\"")
	return buffer.Bytes(), nil
}

func (sl *ClusterSelfLink) UnmarshalJSON(b []byte) error {
	var link string
	if err := json.Unmarshal(b, &link); err != nil {
		return err
	}

	sl.Parse(link)
	return nil
}

func NewClusterSelfLink(name string) *ClusterSelfLink {

	sl := new(ClusterSelfLink)

	sl.string = name
	sl.name = name

	return sl
}
