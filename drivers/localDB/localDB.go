package localDB

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"io/ioutil"
)

type LDB struct {
	path string
	mode os.FileMode
}

func Init(rootPath string) (*LDB, error) {

	conn := new(LDB)
	conn.path = path.Dir(fmt.Sprintf("%s/storage/", rootPath))
	conn.mode = os.FileMode(666)

	if err := os.MkdirAll(conn.path, conn.mode); err != nil {
		return conn, err
	}

	return conn, nil
}

func (ldb *LDB) Read(key string, i interface{}) error {

	if key == "" {
		return nil
	}

	source, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", ldb.path, key))
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(source, i); err != nil {
		return err
	}

	return nil
}

func (ldb *LDB) Write(key string, i interface{}) error {

	payload, _ := yaml.Marshal(i)

	filepath := fmt.Sprintf("%s/%s", ldb.path, key)

	var _, err = os.Stat(filepath)

	if !os.IsNotExist(err) {
		err = os.Remove(filepath)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filepath, os.O_CREATE |os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	_, err = file.Write(payload)
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}

	file.Close()

	return nil
}

func (ldb *LDB) Remove(key string) error {

	filepath := fmt.Sprintf("%s/%s", ldb.path, key)

	var err = os.Remove(filepath)
	if err != nil {
		return err
	}

	return nil
}
