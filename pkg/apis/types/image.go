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

type ImageList []Image

type Image struct {
	imageMeta
	// Image Registry info
	Registry Registry `json:"registry"`
	// Image source info
	Source ImageSource `json:"source"`
}

type imageMeta struct{ ImageMeta }
type ImageMeta struct {
	meta

	// Add fields to expand the meta data
	// Example:
	// Note string `json:"note,omitempty"`
	// Uptime time.Time `json:"uptime"`

	BuildCount int `json:"build_count"`
}

type ImageSpec struct {
	Name     string
	Tag      string
	Registry Registry
}

type ImageSource struct {
	Hub   string `json:"hub"`
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	Tag   string `json:"tag"`
}

func (i *ImageSource) GenerateName() string {
	return ""
}
