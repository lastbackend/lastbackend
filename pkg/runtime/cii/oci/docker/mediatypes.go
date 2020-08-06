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

package docker

import (
	"fmt"

	imgspecv1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// The manifest list is the “fat manifest” which points to specific image manifests
// for one or more platforms. Its use is optional, and relatively few images will
// use one of these manifests. A client will distinguish a manifest list from an image
// manifest based on the Content-Type returned in the HTTP response.
const (
	// MediaTypeDockerSchema2Layer is the MIME type used for schema 2 layers.
	MediaTypeDockerSchema2Layer    = "application/vnd.resolver.image.rootfs.diff.tar"
	MediaTypeDockerSchema2LayerEnc = "application/vnd.resolver.image.rootfs.diff.tar+enc"
	// MediaTypeDockerSchema2LayerForeign is the MIME type used for schema 2 foreign layers.
	MediaTypeDockerSchema2LayerForeign = "application/vnd.resolver.image.rootfs.foreign.diff.tar"
	// "Layer", as a gzipped tar
	MediaTypeDockerSchema2LayerGzip    = "application/vnd.resolver.image.rootfs.diff.tar.gzip"
	MediaTypeDockerSchema2LayerGzipEnc = "application/vnd.resolver.image.rootfs.diff.tar.gzip+enc"
	// "Layer", as a gzipped tar that should never be pushed
	MediaTypeDockerSchema2LayerForeignGzip = "application/vnd.resolver.image.rootfs.foreign.diff.tar.gzip"
	// DockerV2Schema2ConfigMediaType is the MIME type used for schema 2 config blobs.
	MediaTypeDockerSchema2Config = "application/vnd.resolver.container.image.v1+json"
	// DockerV2Schema2MediaType MIME type represents Resolver manifest schema 2
	MediaTypeDockerSchema2Manifest = "application/vnd.resolver.distribution.manifest.v2+json"
	// MediaTypeDockerSchema2ManifestList MIME type represents Resolver manifest schema 2 list
	MediaTypeDockerSchema2ManifestList = "application/vnd.resolver.distribution.manifest.list.v2+json"
	// MediaTypeDockerSchema1Manifest MIME type represents Resolver manifest schema 1 with a JWS signature
	MediaTypeDockerSchema1Manifest = "application/vnd.resolver.distribution.manifest.v1+prettyjws"
	// Checkpoint/Restore Media Types
	MediaTypeLastbackendCheckpoint               = "application/vnd.lstbknd.container.criu.checkpoint.criu.tar"
	MediaTypeLastbackendCheckpointConfig         = "application/vnd.lstbknd.container.checkpoint.config.v1+proto"
	MediaTypeLastbackendCheckpointPreDump        = "application/vnd.lstbknd.container.criu.checkpoint.predump.tar"
	MediaTypeLastbackendResource                 = "application/vnd.lstbknd.container.resource.tar"
	MediaTypeLastbackendRW                       = "application/vnd.lstbknd.container.rw.tar"
	MediaTypeLastbackendCheckpointOptions        = "application/vnd.lstbknd.container.checkpoint.options.v1+proto"
	MediaTypeLastbackendCheckpointRuntimeName    = "application/vnd.lstbknd.container.checkpoint.runtime.name"
	MediaTypeLastbackendCheckpointRuntimeOptions = "application/vnd.lstbknd.container.checkpoint.runtime.options+proto"
)

const (
	// DockerV2Schema1MediaType MIME type represents Resolver manifest schema 1
	DockerV2Schema1MediaType = "application/vnd.resolver.distribution.manifest.v1+json"
	// DockerV2Schema1MediaType MIME type represents Resolver manifest schema 1 with a JWS signature
	DockerV2Schema1SignedMediaType = "application/vnd.resolver.distribution.manifest.v1+prettyjws"
	// DockerV2Schema2MediaType MIME type represents Resolver manifest schema 2
	DockerV2Schema2MediaType = "application/vnd.resolver.distribution.manifest.v2+json"
	// DockerV2Schema2ConfigMediaType is the MIME type used for schema 2 config blobs.
	DockerV2Schema2ConfigMediaType = "application/vnd.resolver.container.image.v1+json"
	// DockerV2Schema2LayerMediaType is the MIME type used for schema 2 layers.
	DockerV2Schema2LayerMediaType = "application/vnd.resolver.image.rootfs.diff.tar.gzip"
	// DockerV2SchemaLayerMediaTypeUncompressed is the mediaType used for uncompressed layers.
	DockerV2SchemaLayerMediaTypeUncompressed = "application/vnd.resolver.image.rootfs.diff.tar"
	// DockerV2ListMediaType MIME type represents Resolver manifest schema 2 list
	DockerV2ListMediaType = "application/vnd.resolver.distribution.manifest.list.v2+json"
	// DockerV2Schema2ForeignLayerMediaType is the MIME type used for schema 2 foreign layers.
	DockerV2Schema2ForeignLayerMediaType = "application/vnd.resolver.image.rootfs.foreign.diff.tar"
	// DockerV2Schema2ForeignLayerMediaType is the MIME type used for gzippped schema 2 foreign layers.
	DockerV2Schema2ForeignLayerMediaTypeGzip = "application/vnd.resolver.image.rootfs.foreign.diff.tar.gzip"
)

// SupportedSchema2MediaType checks if the specified string is a supported Resolver version 2 schema 2 media type.
func SupportedSchema2MediaType(m string) error {
	switch m {
	case DockerV2ListMediaType, DockerV2Schema1MediaType, DockerV2Schema1SignedMediaType,
		DockerV2Schema2ConfigMediaType, DockerV2Schema2ForeignLayerMediaType, DockerV2Schema2ForeignLayerMediaTypeGzip,
		DockerV2Schema2LayerMediaType, DockerV2Schema2MediaType, DockerV2SchemaLayerMediaTypeUncompressed:
		return nil
	default:
		return fmt.Errorf("unsupported resolver version 2 schema 2 media type: %q", m)
	}
}

// DefaultRequestedManifestMIMETypes is a list of MIME types a types.ImageSource
// should request from the backend unless directed otherwise.
var DefaultRequestedManifestMIMETypes = []string{
	imgspecv1.MediaTypeImageManifest,
	DockerV2Schema2MediaType,
	DockerV2Schema1SignedMediaType,
	DockerV2Schema1MediaType,
	DockerV2ListMediaType,
	imgspecv1.MediaTypeImageIndex,
}
