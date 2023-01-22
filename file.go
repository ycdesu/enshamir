package enshamir

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func WriteIfNotExisted(path string, data []byte, perm os.FileMode) error {
	existed, err := IsFilePathExisted(path)
	if err != nil {
		return fmt.Errorf("unable to accesss file %s: %w", path, err)
	}
	if existed {
		return fmt.Errorf("file %s is existed", path)
	}

	return os.WriteFile(path, data, perm)
}

func IsFilePathExisted(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}
