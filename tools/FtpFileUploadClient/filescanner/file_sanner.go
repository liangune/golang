package filescanner

import (
	"fmt"
	"go/gopkg/logger/vglog"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileScanWalkFunc func(path string) error

const DefaultFileScanInterval = 100    // 默认扫描间隔, 单位微秒
const DefaultFileDirMaxAge = 7         // 默认保留7天文件夹
const DefaultFileDirTimerInterval = 12 // 12个小时检测一次

type FileScanner struct {
	srcPath        string           // 源文件夹
	destPath       string           // 目标文件夹
	isScanDestPath bool             // 是否扫描目标文件夹
	interval       time.Duration    // 扫描间隔
	scanWalkFunc   FileScanWalkFunc // 回调处理函数
}

func NewFileScanner(srcPath string, destPath string) *FileScanner {
	s := FileScanner{
		srcPath:        srcPath,
		destPath:       destPath,
		isScanDestPath: true,
		interval:       DefaultFileScanInterval,
		scanWalkFunc:   nil,
	}

	return &s
}

func (f *FileScanner) SetScanInterval(interval time.Duration) {
	if interval <= DefaultFileScanInterval {
		return
	}
	f.interval = interval
}

func (f *FileScanner) Start() {
	go f.scanner()
	go f.deleteEmptyDirTimer()
}

func (f *FileScanner) SetCallBackFunc(fn FileScanWalkFunc) {
	f.scanWalkFunc = fn
}

func (f *FileScanner) walkFunc(path string, info os.FileInfo, err error) error {
	if info == nil {
		//vglog.Error("scan file err: %v", err)
		return err
	}
	if info.IsDir() {
		return nil
	}

	var err1 error
	// .tump文件
	if strings.HasSuffix(path, ".tump") {
		err1 = os.Remove(path)
		if err1 != nil {
			vglog.Error("remove file %s error: %v", path, err1)
			return nil
		}
		vglog.Info("remove file %s", path)
		return nil
	}
	// .tmp文件
	if strings.HasSuffix(path, ".tmp") {
		err1 = os.Remove(path)
		if err1 != nil {
			vglog.Error("remove file %s error: %v", path, err1)
			return nil
		}
		vglog.Info("remove file %s", path)
		return nil
	}

	// .unimastest文件
	if strings.HasSuffix(path, ".unimastest") {
		err1 = os.Remove(path)
		if err1 != nil {
			vglog.Error("remove file %s error: %v", path, err1)
			return nil
		}
		vglog.Info("remove file %s", path)
		return nil
	}

	// 临时文件
	if strings.HasSuffix(path, ".unimastmp") {
		return nil
	}

	destFilePath, err1 := f.renameFile(path)
	if err1 != nil {
		vglog.Error("rename file: %v", err1)
		return nil
	}
	if f.scanWalkFunc != nil {
		f.scanWalkFunc(destFilePath)
	}
	return nil
}

func (f *FileScanner) walkFuncOnFirstScan(path string, info os.FileInfo, err error) error {
	if info == nil {
		//vglog.Error("scan file err: %v", err)
		return err
	}
	if info.IsDir() {
		return nil
	}

	if f.scanWalkFunc != nil {
		f.scanWalkFunc(path)
	}
	return nil
}

func (f *FileScanner) renameFile(path string) (string, error) {
	filePath, fileName := filepath.Split(path)
	// 先移动到临时路径
	destTempPath := f.destPath + string([]rune(filePath)[len([]rune(f.srcPath)):])
	destFilePath := destTempPath + fileName

	err := os.MkdirAll(destTempPath, os.ModePerm)
	vglog.Debug("Create dir : %s", destTempPath)
	if err != nil {
		vglog.Error("Create dir err: %s", err.Error())
		return "", err
	}
	err = os.Rename(path, destFilePath)
	vglog.Debug("src: %s, dest: %s", path, destTempPath)
	if err != nil {
		vglog.Error("rename file to %s from %s , error: %s", destFilePath, path, err)
		return "", err
	}

	return destFilePath, nil
}

func (f *FileScanner) scanner() {
	for {
		var err error
		if f.isScanDestPath {
			err = filepath.Walk(f.destPath, f.walkFuncOnFirstScan)
			f.isScanDestPath = false
		} else {
			err = filepath.Walk(f.srcPath, f.walkFunc)
		}

		if err != nil {
			vglog.Error("scan file path: %v", err)
		}
		time.Sleep(time.Microsecond * f.interval)
	}
}

func (f *FileScanner) deleteEmptyDirTimer() {
	for {
		f.deleteExpireDirectory(f.srcPath)
		f.deleteExpireDirectory(f.destPath)
		time.Sleep(time.Hour * DefaultFileDirTimerInterval)
	}
}

func (f *FileScanner) deleteExpireDirectory(path string) {
	dateFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	curTime := time.Now()
	duration := time.Duration(DefaultFileDirMaxAge*24) * time.Hour

	for _, file := range dateFiles {
		if file.IsDir() {
			date := file.Name()
			datetime := fmt.Sprintf("%s", date)
			objTime, err := time.Parse("20060102", datetime)
			if err != nil {
				vglog.Error("%v", err)
				continue
			}
			objTimeLoc, err := time.ParseInLocation("2006-01-02 15:04:05", objTime.Format("2006-01-02 15:04:05"), time.Local)
			if err != nil {
				vglog.Error("%v", err)
				continue
			}
			// 删除过期空文件夹
			if curTime.Unix() > (objTimeLoc.Unix() + int64(duration.Seconds())) {
				dateDirName := filepath.Join(path, file.Name())
				hourFiles, err := ioutil.ReadDir(dateDirName)
				if err != nil {
					vglog.Error("%v", err)
					continue
				}
				if len(hourFiles) == 0 {
					err = os.Remove(dateDirName)
					if err != nil {
						vglog.Error("%v", err)
					}
				} else {
					for _, hourFile := range hourFiles {
						hourDirName := filepath.Join(dateDirName, hourFile.Name())
						files, err := ioutil.ReadDir(hourDirName)
						if err != nil {
							vglog.Error("%v", err)
							continue
						}
						if len(files) == 0 {
							err = os.Remove(hourDirName)
							if err != nil {
								vglog.Error("%v", err)
							}
						}
					}
				}
			}
		}
	}

}
