package links

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Bhupesh-V/areyouok/utils"
)

func GetFiles(filetype string, ignoredDirs []string) []string {
	var validFiles []string
	userPath := utils.GetUserDirectory()

	err := filepath.Walk(userPath,
		func(path string, info os.FileInfo, err error) error {
			utils.CheckErr(err)
			if info.IsDir() {
				if utils.In(info.Name(), ignoredDirs) {
					return filepath.SkipDir
				}
			}
			if strings.HasSuffix(filepath.Base(path), filetype) && !utils.In(info.Name(), ignoredDirs) {
				validFiles = append(validFiles, path)
			}
			return nil
		})
	utils.CheckErr(err)
	return validFiles
}
