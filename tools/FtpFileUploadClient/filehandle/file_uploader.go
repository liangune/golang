package filehandle

import (
	"github.com/spf13/viper"
	"go/gopkg/workerpool"
	"go/tools/FtpFileUploadClient/ftpclient"
)

type FileUploader struct {
	ftpclient.FtpClient
}

func NewFileUploader() workerpool.HandleInterface {
	addr := viper.GetString("ftp.addr")
	user := viper.GetString("ftp.user")
	password := viper.GetString("ftp.password")

	c := FileUploader{
		ftpclient.FtpClient{
			Addr:          addr,
			User:          user,
			Password:      password,
			FtpConnection: nil,
		},
	}

	return &c
}

func (u *FileUploader) Init() error {
	return u.FtpClient.Init()
}

func (u *FileUploader) UploadFile(localFilePath, ftpPath, ftpFileName, tmpSuffix string) error {
	return u.FtpClient.UploadFile(localFilePath, ftpPath, ftpFileName, tmpSuffix)
}
