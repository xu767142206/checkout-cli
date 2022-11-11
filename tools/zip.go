package tools

import (
	"github.com/mholt/archiver/v3"
)

func Unpack(fileName string,  savePath string) error {
	return archiver.Unarchive(fileName, savePath)
}
