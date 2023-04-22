package utils

import (
	"go/gopkg/logger/vglog"
	"go/tools/FtpFileUploadClient/filehandle"
)

func ScanDirectoryWalkFunc(path string) error {
	fileHandler := filehandle.GetFileHandler()
	fileHandler.Dispatch(path)
	vglog.Access("path: %s", path)
	return nil
}
