package io

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// WriteToFile will create or edit an existing a os.File in `outputDir` with `data` given. If `append` is true,
// `data` given will be appended to the file
func WriteToFile(outputDir, name, data string, append bool) (*os.File, error) {
	outfileName := strings.Join([]string{outputDir, name}, string(filepath.Separator))

	err := ensureDirectoryForFile(outfileName)
	if err != nil {
		return nil, err
	}

	f, err := createOrOpenFile(outfileName, append)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("file_path", outfileName).Msg("file saved")
	return f, nil
}

func createOrOpenFile(filename string, append bool) (*os.File, error) {
	if append {
		return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	}
	return os.Create(filename)
}

func ensureDirectoryForFile(file string) error {
	baseDir := path.Dir(file)
	_, err := os.Stat(baseDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(baseDir, 0755)
}
