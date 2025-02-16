package configs

import (
	"os"
	"regexp"
)

func GetRootPath() string {
	projectDirName := "AIYapper"
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	return string(projectName.Find([]byte(currentWorkDirectory)))
}
