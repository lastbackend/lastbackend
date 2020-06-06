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
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/containerd/containerd/content"
	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const (
	manifestSizeLimit            = 8e6 // 8MB
	labelDockerSchema1EmptyLayer = "lstbknd.net/docker.schema1.empty-layer"
)

type blobState struct {
	diffID digest.Digest
	empty  bool
}

// Converter converts schema1 manifests to schema2 on fetch
type Converter struct {
	contentStore content.Store
	puller       puller

	pulledManifest *manifest

	mu         sync.Mutex
	blobMap    map[digest.Digest]blobState
	layerBlobs map[digest.Digest]ocispec.Descriptor
}

// NewConverter returns a new converter
func NewConverter(contentStore content.Store, puller puller) *Converter {
	return &Converter{
		contentStore: contentStore,
		puller:       puller,
		blobMap:      map[digest.Digest]blobState{},
		layerBlobs:   map[digest.Digest]ocispec.Descriptor{},
	}
}

// Handle fetching descriptors for a docker media type
func (c *Converter) Handle(ctx context.Context, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) {
	switch desc.MediaType {
	case MediaTypeDockerSchema1Manifest:
		if err := c.fetchManifest(ctx, desc); err != nil {
			return nil, err
		}

		m := c.pulledManifest
		if len(m.FSLayers) != len(m.History) {
			return nil, errors.New("invalid schema 1 manifest, history and layer mismatch")
		}
		descs := make([]ocispec.Descriptor, 0, len(c.pulledManifest.FSLayers))

		for i := range m.FSLayers {
			if _, ok := c.blobMap[c.pulledManifest.FSLayers[i].BlobSum]; !ok {
				empty, err := isEmptyLayer([]byte(m.History[i].V1Compatibility))
				if err != nil {
					return nil, err
				}

				// Do no attempt to download a known empty blob
				if !empty {
					descs = append([]ocispec.Descriptor{
						{
							MediaType: MediaTypeDockerSchema2LayerGzip,
							Digest:    c.pulledManifest.FSLayers[i].BlobSum,
							Size:      -1,
						},
					}, descs...)
				}
				c.blobMap[c.pulledManifest.FSLayers[i].BlobSum] = blobState{
					empty: empty,
				}
			}
		}
		return descs, nil
	case MediaTypeDockerSchema2LayerGzip:
		if c.pulledManifest == nil {
			return nil, errors.New("manifest required for schema 1 blob pull")
		}
		return nil, c.fetchBlob(ctx, desc)
	default:
		return nil, fmt.Errorf("%v not support for schema 1 manifests", desc.MediaType)
	}
}

// ConvertOptions provides options on converting a docker schema1 manifest.
type ConvertOptions struct {
	// ManifestMediaType specifies the media type of the manifest OCI descriptor.
	ManifestMediaType string

	// ConfigMediaType specifies the media type of the manifest config OCI
	// descriptor.
	ConfigMediaType string
}

// ConvertOpt allows configuring a convert operation.
type ConvertOpt func(context.Context, *ConvertOptions) error

// UseDockerSchema2 is used to indicate that a schema1 manifest should be
// converted into the media types for a docker schema2 manifest.
func UseDockerSchema2() ConvertOpt {
	return func(ctx context.Context, o *ConvertOptions) error {
		o.ManifestMediaType = MediaTypeDockerSchema2Manifest
		o.ConfigMediaType = MediaTypeDockerSchema2Config
		return nil
	}
}

// Convert a docker manifest to an OCI descriptor
func (c *Converter) Convert(ctx context.Context, opts ...ConvertOpt) (ocispec.Descriptor, error) {
	co := ConvertOptions{
		ManifestMediaType: ocispec.MediaTypeImageManifest,
		ConfigMediaType:   ocispec.MediaTypeImageConfig,
	}
	for _, opt := range opts {
		if err := opt(ctx, &co); err != nil {
			return ocispec.Descriptor{}, err
		}
	}

	history, diffIDs, err := c.schema1ManifestHistory()
	if err != nil {
		return ocispec.Descriptor{}, errors.Wrap(err, "schema 1 conversion failed")
	}

	var img ocispec.Image
	if err := json.Unmarshal([]byte(c.pulledManifest.History[0].V1Compatibility), &img); err != nil {
		return ocispec.Descriptor{}, errors.Wrap(err, "failed to unmarshal image from schema 1 history")
	}

	img.History = history
	img.RootFS = ocispec.RootFS{
		Type:    "layers",
		DiffIDs: diffIDs,
	}

	b, err := json.MarshalIndent(img, "", "   ")
	if err != nil {
		return ocispec.Descriptor{}, errors.Wrap(err, "failed to marshal image")
	}

	config := ocispec.Descriptor{
		MediaType: co.ConfigMediaType,
		Digest:    digest.Canonical.FromBytes(b),
		Size:      int64(len(b)),
	}

	layers := make([]ocispec.Descriptor, len(diffIDs))
	for i, diffID := range diffIDs {
		layers[i] = c.layerBlobs[diffID]
	}

	manifest := ocispec.Manifest{
		Versioned: specs.Versioned{
			SchemaVersion: 2,
		},
		Config: config,
		Layers: layers,
	}

	mb, err := json.MarshalIndent(manifest, "", "   ")
	if err != nil {
		return ocispec.Descriptor{}, errors.Wrap(err, "failed to marshal image")
	}

	desc := ocispec.Descriptor{
		MediaType: co.ManifestMediaType,
		Digest:    digest.Canonical.FromBytes(mb),
		Size:      int64(len(mb)),
	}

	labels := map[string]string{}
	labels["containerd.io/gc.ref.content.0"] = manifest.Config.Digest.String()
	for i, ch := range manifest.Layers {
		labels[fmt.Sprintf("containerd.io/gc.ref.content.%d", i+1)] = ch.Digest.String()
	}

	ref := MakeRefKey(ctx, desc)
	if err := content.WriteBlob(ctx, c.contentStore, ref, bytes.NewReader(mb), desc, content.WithLabels(labels)); err != nil {
		return ocispec.Descriptor{}, errors.Wrap(err, "failed to write image manifest")
	}

	ref = MakeRefKey(ctx, config)
	if err := content.WriteBlob(ctx, c.contentStore, ref, bytes.NewReader(b), config); err != nil {
		return ocispec.Descriptor{}, errors.Wrap(err, "failed to write image config")
	}

	return desc, nil
}

// ReadStripSignature reads in a schema1 manifest and returns a byte array
// with the "signatures" field stripped
func ReadStripSignature(schema1Blob io.Reader) ([]byte, error) {
	b, err := ioutil.ReadAll(io.LimitReader(schema1Blob, manifestSizeLimit)) // limit to 8MB
	if err != nil {
		return nil, err
	}

	return stripSignature(b)
}

func (c *Converter) fetchManifest(ctx context.Context, desc ocispec.Descriptor) error {
	rc, err := c.puller.pull(ctx, desc)
	if err != nil {
		return err
	}

	b, err := ReadStripSignature(rc)
	rc.Close()
	if err != nil {
		return err
	}

	var m manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	c.pulledManifest = &m

	return nil
}

func (c *Converter) fetchBlob(ctx context.Context, desc ocispec.Descriptor) error {
	var (
		ref            = MakeRefKey(ctx, desc)
		calc           = newBlobStateCalculator()
		compressMethod = Gzip
	)

	// size may be unknown, set to zero for content ingest
	ingestDesc := desc
	if ingestDesc.Size == -1 {
		ingestDesc.Size = 0
	}

	cw, err := content.OpenWriter(ctx, c.contentStore, content.WithRef(ref), content.WithDescriptor(ingestDesc))
	if err != nil {
		if !IsAlreadyExists(err) {
			return err
		}

		reuse, err := c.reuseLabelBlobState(ctx, desc)
		if err != nil {
			return err
		}

		if reuse {
			return nil
		}

		ra, err := c.contentStore.ReaderAt(ctx, desc)
		if err != nil {
			return err
		}
		defer ra.Close()

		r, err := DecompressStream(content.NewReader(ra))
		if err != nil {
			return err
		}

		compressMethod = r.GetCompression()
		_, err = io.Copy(calc, r)
		r.Close()
		if err != nil {
			return err
		}
	} else {
		defer cw.Close()

		rc, err := c.puller.pull(ctx, desc)
		if err != nil {
			return err
		}
		defer rc.Close()

		eg, _ := errgroup.WithContext(ctx)
		pr, pw := io.Pipe()

		eg.Go(func() error {
			r, err := DecompressStream(pr)
			if err != nil {
				return err
			}

			compressMethod = r.GetCompression()
			_, err = io.Copy(calc, r)
			r.Close()
			pr.CloseWithError(err)
			return err
		})

		eg.Go(func() error {
			defer pw.Close()

			return content.Copy(ctx, cw, io.TeeReader(rc, pw), ingestDesc.Size, ingestDesc.Digest)
		})

		if err := eg.Wait(); err != nil {
			return err
		}
	}

	if desc.Size == -1 {
		info, err := c.contentStore.Info(ctx, desc.Digest)
		if err != nil {
			return errors.Wrap(err, "failed to get blob info")
		}
		desc.Size = info.Size
	}

	if compressMethod == Uncompressed {
		desc.MediaType = MediaTypeDockerSchema2Layer
	}

	state := calc.State()

	cinfo := content.Info{
		Digest: desc.Digest,
		Labels: map[string]string{
			"containerd.io/uncompressed": state.diffID.String(),
			labelDockerSchema1EmptyLayer: strconv.FormatBool(state.empty),
		},
	}

	if _, err := c.contentStore.Update(ctx, cinfo, "labels.containerd.io/uncompressed", fmt.Sprintf("labels.%s", labelDockerSchema1EmptyLayer)); err != nil {
		return errors.Wrap(err, "failed to update uncompressed label")
	}

	c.mu.Lock()
	c.blobMap[desc.Digest] = state
	c.layerBlobs[state.diffID] = desc
	c.mu.Unlock()

	return nil
}

func (c *Converter) reuseLabelBlobState(ctx context.Context, desc ocispec.Descriptor) (bool, error) {
	cinfo, err := c.contentStore.Info(ctx, desc.Digest)
	if err != nil {
		return false, errors.Wrap(err, "failed to get blob info")
	}
	desc.Size = cinfo.Size

	diffID, ok := cinfo.Labels["containerd.io/uncompressed"]
	if !ok {
		return false, nil
	}

	emptyVal, ok := cinfo.Labels[labelDockerSchema1EmptyLayer]
	if !ok {
		return false, nil
	}

	isEmpty, err := strconv.ParseBool(emptyVal)
	if err != nil {
		return false, nil
	}

	bState := blobState{empty: isEmpty}

	if bState.diffID, err = digest.Parse(diffID); err != nil {
		return false, nil
	}

	// NOTE: there is no need to read header to get compression method
	// because there are only two kinds of methods.
	if bState.diffID == desc.Digest {
		desc.MediaType = MediaTypeDockerSchema2Layer
	} else {
		desc.MediaType = MediaTypeDockerSchema2LayerGzip
	}

	c.mu.Lock()
	c.blobMap[desc.Digest] = bState
	c.layerBlobs[bState.diffID] = desc
	c.mu.Unlock()
	return true, nil
}

func (c *Converter) schema1ManifestHistory() ([]ocispec.History, []digest.Digest, error) {
	if c.pulledManifest == nil {
		return nil, nil, errors.New("missing schema 1 manifest for conversion")
	}
	m := *c.pulledManifest

	if len(m.History) == 0 {
		return nil, nil, errors.New("no history")
	}

	history := make([]ocispec.History, len(m.History))
	diffIDs := []digest.Digest{}
	for i := range m.History {
		var h v1History
		if err := json.Unmarshal([]byte(m.History[i].V1Compatibility), &h); err != nil {
			return nil, nil, errors.Wrap(err, "failed to unmarshal history")
		}

		blobSum := m.FSLayers[i].BlobSum

		state := c.blobMap[blobSum]

		history[len(history)-i-1] = ocispec.History{
			Author:     h.Author,
			Comment:    h.Comment,
			Created:    &h.Created,
			CreatedBy:  strings.Join(h.ContainerConfig.Cmd, " "),
			EmptyLayer: state.empty,
		}

		if !state.empty {
			diffIDs = append([]digest.Digest{state.diffID}, diffIDs...)

		}
	}

	return history, diffIDs, nil
}

type fsLayer struct {
	BlobSum digest.Digest `json:"blobSum"`
}

type history struct {
	V1Compatibility string `json:"v1Compatibility"`
}

type manifest struct {
	FSLayers []fsLayer `json:"fsLayers"`
	History  []history `json:"history"`
}

type v1History struct {
	Author          string    `json:"author,omitempty"`
	Created         time.Time `json:"created"`
	Comment         string    `json:"comment,omitempty"`
	ThrowAway       *bool     `json:"throwaway,omitempty"`
	Size            *int      `json:"Size,omitempty"` // used before ThrowAway field
	ContainerConfig struct {
		Cmd []string `json:"Cmd,omitempty"`
	} `json:"container_config,omitempty"`
}

// isEmptyLayer returns whether the v1 compatibility history describes an
// empty layer. A return value of true indicates the layer is empty,
// however false does not indicate non-empty.
func isEmptyLayer(compatHistory []byte) (bool, error) {
	var h v1History
	if err := json.Unmarshal(compatHistory, &h); err != nil {
		return false, err
	}

	if h.ThrowAway != nil {
		return *h.ThrowAway, nil
	}
	if h.Size != nil {
		return *h.Size == 0, nil
	}

	// If no `Size` or `throwaway` field is given, then
	// it cannot be determined whether the layer is empty
	// from the history, return false
	return false, nil
}

type signature struct {
	Signatures []jsParsedSignature `json:"signatures"`
}

type jsParsedSignature struct {
	Protected string `json:"protected"`
}

type protectedBlock struct {
	Length int    `json:"formatLength"`
	Tail   string `json:"formatTail"`
}

// joseBase64UrlDecode decodes the given string using the standard base64 url
// decoder but first adds the appropriate number of trailing '=' characters in
// accordance with the jose specification.
// http://tools.ietf.org/html/draft-ietf-jose-json-web-signature-31#section-2
func joseBase64UrlDecode(s string) ([]byte, error) {
	switch len(s) % 4 {
	case 0:
	case 2:
		s += "=="
	case 3:
		s += "="
	default:
		return nil, errors.New("illegal base64url string")
	}
	return base64.URLEncoding.DecodeString(s)
}

func stripSignature(b []byte) ([]byte, error) {
	var sig signature
	if err := json.Unmarshal(b, &sig); err != nil {
		return nil, err
	}
	if len(sig.Signatures) == 0 {
		return nil, errors.New("no signatures")
	}
	pb, err := joseBase64UrlDecode(sig.Signatures[0].Protected)
	if err != nil {
		return nil, errors.Wrapf(err, "could not decode %s", sig.Signatures[0].Protected)
	}

	var protected protectedBlock
	if err := json.Unmarshal(pb, &protected); err != nil {
		return nil, err
	}

	if protected.Length > len(b) {
		return nil, errors.New("invalid protected length block")
	}

	tail, err := joseBase64UrlDecode(protected.Tail)
	if err != nil {
		return nil, errors.Wrap(err, "invalid tail base 64 value")
	}

	return append(b[:protected.Length], tail...), nil
}

type blobStateCalculator struct {
	empty    bool
	digester digest.Digester
}

func newBlobStateCalculator() *blobStateCalculator {
	return &blobStateCalculator{
		empty:    true,
		digester: digest.Canonical.Digester(),
	}
}

func (c *blobStateCalculator) Write(p []byte) (int, error) {
	if c.empty {
		for _, b := range p {
			if b != 0x00 {
				c.empty = false
				break
			}
		}
	}
	return c.digester.Hash().Write(p)
}

func (c *blobStateCalculator) State() blobState {
	return blobState{
		empty:  c.empty,
		diffID: c.digester.Digest(),
	}
}

type refKeyPrefix struct{}

// MakeRefKey returns a unique reference for the descriptor. This reference can be
// used to lookup ongoing processes related to the descriptor. This function
// may look to the context to namespace the reference appropriately.
func MakeRefKey(ctx context.Context, desc ocispec.Descriptor) string {
	if v := ctx.Value(refKeyPrefix{}); v != nil {
		values := v.(map[string]string)
		if prefix := values[desc.MediaType]; prefix != "" {
			return prefix + "-" + desc.Digest.String()
		}
	}

	switch mt := desc.MediaType; {
	case mt == MediaTypeDockerSchema2Manifest || mt == ocispec.MediaTypeImageManifest:
		return "manifest-" + desc.Digest.String()
	case mt == MediaTypeDockerSchema2ManifestList || mt == ocispec.MediaTypeImageIndex:
		return "index-" + desc.Digest.String()
	case IsLayerType(mt):
		return "layer-" + desc.Digest.String()
	case IsKnownConfig(mt):
		return "config-" + desc.Digest.String()
	default:
		return "unknown-" + desc.Digest.String()
	}
}

// parseMediaTypes splits the media type into the base type and
// an array of sorted extensions
func parseMediaTypes(mt string) (string, []string) {
	if mt == "" {
		return "", []string{}
	}

	s := strings.Split(mt, "+")
	ext := s[1:]
	sort.Strings(ext)

	return s[0], ext
}

// IsLayerTypes returns true if the media type is a layer
func IsLayerType(mt string) bool {
	if strings.HasPrefix(mt, "application/vnd.oci.image.layer.") {
		return true
	}

	// Parse Docker media types, strip off any + suffixes first
	base, _ := parseMediaTypes(mt)
	switch base {
	case MediaTypeDockerSchema2Layer, MediaTypeDockerSchema2LayerGzip,
		MediaTypeDockerSchema2LayerForeign, MediaTypeDockerSchema2LayerForeignGzip:
		return true
	}
	return false
}

// IsKnownConfig returns true if the media type is a known config type
func IsKnownConfig(mt string) bool {
	switch mt {
	case MediaTypeDockerSchema2Config, ocispec.MediaTypeImageConfig,
		MediaTypeLastbackendCheckpoint, MediaTypeLastbackendCheckpointConfig:
		return true
	}
	return false
}
