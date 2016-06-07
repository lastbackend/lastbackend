package app

import (
	"fmt"
	"github.com/deployithq/deployit/daemon/env"
	"github.com/deployithq/deployit/utils"
	"os"
	"time"
)

type Layer struct {
	ID      string    `json:"id" yaml:"id"`
	Created time.Time `json:"created" yaml:"created"`
	Updated time.Time `json:"updated" yaml:"updated"`
}

func (l *Layer) CreateFromUrl(url string) error {
	return nil
}

func (l *Layer) CreateFromTarGz(path string, excludes []string) error {

	_, err := os.Stat(path)

	if os.IsNotExist(err) && len(excludes) == 0 {
		return nil
	}

	var update bool = true
	var layer string = utils.GenerateID()
	var tar_path string

	if !os.IsNotExist(err) {
		tar_path = fmt.Sprintf("%s/tmp/%s", env.Default_root_path, layer)

		// If there are no layers, then there is nothing to update
		if l.ID == "" {
			update = false
			l.Created = time.Now()
			tar_path = fmt.Sprintf("%s/apps/%s", env.Default_root_path, layer)
		}

		if err := utils.Ungzip(path, tar_path); err != nil {
			return err
		}
	}

	if update {
		src_path := fmt.Sprintf("%s/apps/%s", env.Default_root_path, l.ID)
		target_path := fmt.Sprintf("%s/apps/%s", env.Default_root_path, layer)
		if err := utils.Update(src_path, target_path, tar_path, excludes); err != nil {
			return err
		}

		if err := utils.RemoveDirs([]string{tar_path}); err != nil {
			return err
		}
	}

	l.Updated = time.Now()
	l.ID = layer

	return nil
}
