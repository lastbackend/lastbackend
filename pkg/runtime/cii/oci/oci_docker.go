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

package oci

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/types"
	encconfig "github.com/containers/ocicrypt/config"
	cryptUtils "github.com/containers/ocicrypt/utils"
	"github.com/containers/storage"
	"github.com/lastbackend/lastbackend/internal/pkg/models"
	"github.com/lastbackend/lastbackend/pkg/runtime/cii/oci/metrics"
	stg "github.com/lastbackend/lastbackend/pkg/runtime/cii/oci/storage"
	"github.com/pkg/errors"
)

const (
	DefaultLocalRegistry = "localhost"
)

var localRegistryPrefix = DefaultLocalRegistry + "/"

type Docker struct {
	// pullOperationsLock is used to synchronize pull operations.
	pullOperationsLock sync.Mutex

	config ConfigDocker
	// pullOperationsInProgress is used to avoid pulling the same image in parallel. Goroutines
	// will block on the pullResult.
	pullOperationsInProgress map[pullArguments]*pullOperation
	storageImageServer       stg.ImageServer

	SystemContext *types.SystemContext
	// DecryptionKeysPath is the path where keys for image decryption are stored.
	DecryptionKeysPath string `toml:"decryption_keys_path"`
}

func NewDocker(config ConfigDocker) (*Docker, error) {

	ctx := context.Background()

	options, err := storage.DefaultStoreOptions(false, 0)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create default image store options")
	}

	if config.RunRoot != "" {
		options.RunRoot = config.RunRoot
	}

	if config.Root != "" {
		options.GraphRoot = config.Root
	}

	if config.StorageDriver != "" {
		options.GraphDriverName = config.StorageDriver
	}

	store, err := storage.GetStore(options)
	if err != nil {
		return nil, err
	}

	systemContext := new(types.SystemContext)

	imageService, err := stg.GetImageService(ctx, systemContext, store, config.DefaultTransport, config.InsecureRegistries, config.Registries)
	if err != nil {
		return nil, err
	}

	d := new(Docker)
	d.config = config
	d.storageImageServer = imageService
	d.pullOperationsInProgress = make(map[pullArguments]*pullOperation)

	return d, nil
}

func (d *Docker) Auth(ctx context.Context, secret *models.SecretAuthData) (string, error) {
	config := models.AuthConfig{
		Username: secret.Username,
		Password: secret.Password,
	}
	js, err := json.Marshal(config)
	if err != nil {
		return models.EmptyString, err
	}
	return base64.URLEncoding.EncodeToString(js), nil
}

// DockerAuthConfig contains authorization information for connecting to a registry.
// the value of Username and Password can be empty for accessing the registry anonymously
type DockerAuthConfig struct {
	Username string
	Password string
	// IdentityToken can be used as an refresh_token in place of username and
	// password to obtain the bearer/access token in oauth2 flow. If identity
	// token is set, password should not be set.
	// Ref: https://docs.docker.com/registry/spec/auth/oauth/
	IdentityToken string
}

// pullOperation is used to synchronize parallel pull operations via the
// server's pullCache.  Goroutines can block the pullOperation's waitgroup and
// be released once the pull operation has finished.
type pullOperation struct {
	// wg allows for Goroutines trying to pull the same image to wait until the
	// currently running pull operation has finished.
	wg sync.WaitGroup
	// imageRef is the reference of the actually pulled image which will differ
	// from the input if it was a short name (e.g., alpine).
	imageRef string
	// err is the error indicating if the pull operation has succeeded or not.
	err error
}

// pullArguments are used to identify a pullOperation via an input image name and
// possibly specified credentials.
type pullArguments struct {
	image       string
	credentials types.DockerAuthConfig
}

func (d *Docker) Pull(ctx context.Context, spec *models.ImageManifest, out io.Writer) (*models.Image, error) {

	pullArgs := pullArguments{image: spec.Name}

	// We use the server's pullOperationsInProgress to record which images are
	// currently being pulled. This allows for avoiding pulling the same image
	// in parallel. Hence, if a given image is currently being pulled, we queue
	// into the pullOperation's waitgroup and wait for the pulling goroutine to
	// unblock us and re-use its results.
	pullOp, pullInProcess := func() (pullOp *pullOperation, inProgress bool) {
		d.pullOperationsLock.Lock()
		defer d.pullOperationsLock.Unlock()
		pullOp, inProgress = d.pullOperationsInProgress[pullArgs]
		if !inProgress {
			pullOp = &pullOperation{}
			d.pullOperationsInProgress[pullArgs] = pullOp
			pullOp.wg.Add(1)
		}
		return pullOp, inProgress
	}()

	if !pullInProcess {
		pullOp.err = errors.New("pull image was aborted")
		defer func() {
			d.pullOperationsLock.Lock()
			delete(d.pullOperationsInProgress, pullArgs)
			pullOp.wg.Done()
			d.pullOperationsLock.Unlock()
		}()
		pullOp.imageRef, pullOp.err = d.pullImage(ctx, &pullArgs)
	} else {
		// Wait for the pull operation to finish.
		pullOp.wg.Wait()
	}

	if pullOp.err != nil {
		return nil, pullOp.err
	}

	return nil, nil
}

func (d *Docker) Remove(ctx context.Context, image string) error {
	return nil
}

func (d *Docker) Push(ctx context.Context, spec *models.ImageManifest, out io.Writer) (*models.Image, error) {
	return nil, nil
}

func (d *Docker) Build(ctx context.Context, stream io.Reader, spec *models.SpecBuildImage, out io.Writer) (*models.Image, error) {
	return nil, nil
}

func (d *Docker) List(ctx context.Context, filters ...string) ([]*models.Image, error) {
	return nil, nil
}

func (d *Docker) Inspect(ctx context.Context, id string) (*models.Image, error) {
	return nil, nil
}

func (d *Docker) Subscribe(ctx context.Context) (chan *models.Image, error) {
	return nil, nil
}

func (d *Docker) Close() error {
	return nil
}

func (d *Docker) pullImage(ctx context.Context, pullArgs *pullArguments) (string, error) {
	var err error

	// A shallow copy we can modify
	sourceCtx := *d.SystemContext

	if pullArgs.credentials.Username != "" {
		sourceCtx.DockerAuthConfig = &pullArgs.credentials
	}

	decryptConfig, err := getDecryptionKeys(d.DecryptionKeysPath)
	if err != nil {
		return "", err
	}

	var (
		images []string
		pulled string
	)

	images, err = d.storageImageServer.ResolveNames(d.SystemContext, pullArgs.image)
	if err != nil {
		return "", err
	}

	for _, img := range images {
		var tmpImg types.ImageCloser
		tmpImg, err = d.storageImageServer.PrepareImage(&sourceCtx, img)
		if err != nil {
			// We're not able to find the image remotely, check if it's
			// available locally, but only for localhost/ prefixed ones.
			// This allows pulling localhost/ prefixed images even if the
			// `imagePullPolicy` is set to `Always`.
			if strings.HasPrefix(img, localRegistryPrefix) {
				if _, err := d.storageImageServer.ImageStatus(
					d.SystemContext, img,
				); err == nil {
					pulled = img
					break
				}
			}
			fmt.Println(ctx, "error preparing image %s: %v", img, err)
			continue
		}
		defer tmpImg.Close()

		var storedImage *stg.ImageResult
		storedImage, err = d.storageImageServer.ImageStatus(d.SystemContext, img)
		if err == nil {
			tmpImgConfigDigest := tmpImg.ConfigInfo().Digest
			if tmpImgConfigDigest.String() == "" {
				// this means we are playing with a schema1 image, in which
				// case, we're going to repull the image in any case
				fmt.Println(ctx, "image config digest is empty, re-pulling image")
			} else if tmpImgConfigDigest.String() == storedImage.ConfigDigest.String() {
				fmt.Println(ctx, "image %s already in store, skipping pull", img)
				pulled = img

				// Skipped bytes metrics
				if storedImage.Size != nil {
					counter, err := metrics.CRIOImagePullsByNameSkipped.GetMetricWithLabelValues(img)
					if err != nil {
						fmt.Println(ctx, "Unable to write image pull name (skipped) metrics: %v", err)
					} else {
						counter.Add(float64(*storedImage.Size))
					}
				}

				break
			}
			fmt.Println(ctx, "image in store has different ID, re-pulling %s", img)
		}

		// Pull by collecting progress metrics
		progress := make(chan types.ProgressProperties)
		go func() {
			for p := range progress {
				if p.Artifact.Size > 0 {
					fmt.Println(ctx, "ImagePull (%v): %s (%s): %v bytes (%.2f%%)",
						p.Event, img, p.Artifact.Digest, p.Offset,
						float64(p.Offset)/float64(p.Artifact.Size)*100,
					)
				} else {
					fmt.Println(ctx, "ImagePull (%v): %s (%s): %v bytes",
						p.Event, img, p.Artifact.Digest, p.Offset,
					)
				}

				// Metrics for every digest
				digestCounter, err := metrics.CRIOImagePullsByDigest.GetMetricWithLabelValues(
					img, p.Artifact.Digest.String(), p.Artifact.MediaType,
					fmt.Sprintf("%d", p.Artifact.Size),
				)
				if err != nil {
					fmt.Println(ctx, "Unable to write image pull digest metrics: %v", err)
				} else {
					digestCounter.Add(float64(p.OffsetUpdate))
				}

				// Metrics for the overall image
				nameCounter, err := metrics.CRIOImagePullsByName.GetMetricWithLabelValues(
					img, fmt.Sprintf("%d", imageSize(tmpImg)),
				)
				if err != nil {
					fmt.Println(ctx, "Unable to write image pull name metrics: %v", err)
				} else {
					nameCounter.Add(float64(p.OffsetUpdate))
				}
			}
		}()

		_, err = d.storageImageServer.PullImage(d.SystemContext, img, &copy.Options{
			SourceCtx:        &sourceCtx,
			DestinationCtx:   d.SystemContext,
			OciDecryptConfig: decryptConfig,
			ProgressInterval: time.Second,
			Progress:         progress,
		})
		if err != nil {
			fmt.Println(ctx, "error pulling image %s: %v", img, err)
			continue
		}
		pulled = img
		break
	}

	if pulled == "" && err != nil {
		return "", err
	}

	status, err := d.storageImageServer.ImageStatus(d.SystemContext, pulled)
	if err != nil {
		return "", err
	}
	imageRef := status.ID
	if len(status.RepoDigests) > 0 {
		imageRef = status.RepoDigests[0]
	}

	return imageRef, nil
}

func decodeDockerAuth(s string) (user, password string, err error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", "", err
	}
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		// if it's invalid just skip, as docker does
		return "", "", nil
	}
	user = parts[0]
	password = strings.Trim(parts[1], "\x00")
	return user, password, nil
}

func imageSize(img types.ImageCloser) (size int64) {
	for _, layer := range img.LayerInfos() {
		if layer.Size > 0 {
			size += layer.Size
		} else {
			return -1
		}
	}

	configSize := img.ConfigInfo().Size
	if configSize >= 0 {
		size += configSize
	} else {
		return -1
	}

	return size
}

// getDecryptionKeys reads the keys from the given directory
func getDecryptionKeys(keysPath string) (*encconfig.DecryptConfig, error) {
	if _, err := os.Stat(keysPath); os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("skipping non-existing decryption_keys_path: %s", keysPath))
		return nil, nil
	}

	base64Keys := make([]string, 0)
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Handle symlinks
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			return errors.New("Symbolic links not supported in decryption keys paths")
		}

		privateKey, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		sEnc := base64.StdEncoding.EncodeToString(privateKey)
		base64Keys = append(base64Keys, sEnc)

		return nil
	}

	if err := filepath.Walk(keysPath, walkFn); err != nil {
		return nil, err
	}

	sortedDc, err := cryptUtils.SortDecryptionKeys(strings.Join(base64Keys, ","))
	if err != nil {
		return nil, err
	}

	return encconfig.InitDecryption(sortedDc).DecryptConfig, nil
}
