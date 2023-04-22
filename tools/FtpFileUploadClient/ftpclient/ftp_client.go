package ftpclient

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"go/gopkg/logger/vglog"
	"net"
	"os"
	"time"
)

const DefaultDialTimeout = time.Second * 10

// FileZilla Server max keepalive time is 2 minute
const DefaultDialKeepAlive = time.Minute * 2

type FtpClient struct {
	Addr          string
	User          string
	Password      string
	FtpConnection *ftp.ServerConn
}

func NewFtpClient(addr, user, password string) *FtpClient {
	client := FtpClient{
		Addr:          addr,
		User:          user,
		Password:      password,
		FtpConnection: nil,
	}
	return &client
}

func (client *FtpClient) Init() error {
	return nil
}

func (client *FtpClient) Connect() error {
	dialer := net.Dialer{
		KeepAlive: DefaultDialKeepAlive,
	}
	conn, err := ftp.Dial(client.Addr, ftp.DialWithTimeout(DefaultDialTimeout), ftp.DialWithDialer(dialer))
	if err != nil {
		vglog.Error("ftp client connect: %v", err)
		return err
	}

	client.FtpConnection = conn
	err = client.FtpConnection.Login(client.User, client.Password)
	if err != nil {
		vglog.Error("%v", err)
		return err
	}
	return nil
}

func (client *FtpClient) Close() error {
	if client.FtpConnection == nil {
		return fmt.Errorf("ftp connection is nil")
	}
	err := client.FtpConnection.Logout()
	if err != nil {
		vglog.Error("%v", err)
		return err
	}
	err = client.FtpConnection.Quit()
	if err != nil {
		vglog.Error("%v", err)
		return err
	}
	return nil
}

func (client *FtpClient) UploadFile(localFilePath, ftpPath, ftpFileName, tmpSuffix string) error {
	if client.FtpConnection == nil {
		err := client.Connect()
		if err != nil {
			return err
		}
	} else {
		err := client.FtpConnection.NoOp()
		if err != nil {
			err := client.Connect()
			if err != nil {
				return err
			}
		}
	}

	client.FtpConnection.MakeDir(ftpPath)
	client.FtpConnection.ChangeDir(ftpPath)

	file, err := os.Open(localFilePath)
	if err != nil {
		vglog.Error("Open file err: %v", err)
		return err
	}

	defer file.Close()
	ftpFileNameTmp := ftpFileName + tmpSuffix
	err = client.FtpConnection.Stor(ftpFileNameTmp, file)
	if err != nil {
		vglog.Error("ftp stor file err: %v", err)
		return err
	}
	err = client.FtpConnection.Rename(ftpFileNameTmp, ftpFileName)
	if err != nil {
		vglog.Error("ftp rename file err: %v", err)
		return err
	}

	return nil
}
