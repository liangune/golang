package filehandle

import (
	"fmt"
	"github.com/spf13/viper"
	"go/gopkg/logger/vglog"
	"go/gopkg/workerpool"
	"os"
	"path/filepath"
)

type FileTask struct {
	path string
}

func (t *FileTask) Execute(w workerpool.WorkerInterface) error {
	tmpPath := viper.GetString("ftp.tmpPath")
	ftpServerPath := viper.GetString("ftp.ftpServerPath")
	tmpSuffix := viper.GetString("ftp.tmpSuffix")

	srcDir, srcFileName := filepath.Split(t.path)
	remoteFtpPath := string([]rune(srcDir)[len([]rune(tmpPath)):])
	remoteFtpPath = ftpServerPath + remoteFtpPath

	handle := w.GetHandle()
	if handle != nil {
		client, ok := handle.(*FileUploader)
		if !ok {
			err := fmt.Errorf("handle interface transforms into *FileUploader errror")
			vglog.Error("%v", err)
			return err
		}
		err := client.UploadFile(t.path, remoteFtpPath, srcFileName, tmpSuffix)
		if err != nil {
			vglog.Error("%v", err)
			return err
		} else {
			destName := remoteFtpPath + srcFileName
			vglog.Notice("Successfully transferred ftp file: %s, local file: %s,", destName, t.path)
			os.Remove(t.path)
		}
	}

	return nil
}
