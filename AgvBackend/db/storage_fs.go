package db

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"turtle/core/lgr"
	"turtle/vfs"
)

// _DesktopClient defines methods for working with local storage
type DesktopClient struct{}

func (DesktopClient) GetFileBytes(container, fileName string) ([]byte, error) {
	return vfs.GetFileBytesFromWD(container, fileName)
}

func (DesktopClient) GetFileBytesNew(filePath string) ([]byte, error) {
	return vfs.GetFileBytesFromWDNew(filePath)
}

func (DesktopClient) GetFileString(container, fileName string) (string, error) {
	return vfs.GetFileStringFromWD(container, fileName)
}

func (DesktopClient) GetFileStringNew(filePath string) (string, error) {
	return vfs.GetFileStringFromWDNew(filePath)
}

func (DesktopClient) UploadFile(container, fileName string, _bytes []byte) error {
	return vfs.WriteFileToWD(container, fileName, _bytes)
}

func (DesktopClient) UploadFileNew(filePath string, _bytes []byte) error {
	return vfs.WriteFileToWDNew(filePath, _bytes)
}

func (DesktopClient) DeleteFileNew(filePath string) error {
	return vfs.DeleteFileNew(filePath)
}

func (DesktopClient) MoveFileNew(oldFilePath, newFilePath string) error {
	return vfs.MoveFileNew(oldFilePath, newFilePath)
}

func (DesktopClient) UploadFileString(container, fileName, content string) error {
	return vfs.WriteFileStringToWD(container, fileName, content)
}

func (DesktopClient) UploadFileStringNew(filePath, content string) error {
	return vfs.WriteFileStringToWDNew(filePath, content)
}

func (DesktopClient) UploadFileStringNewBase64(filePath, content string) error {
	return vfs.WriteFileStringToWDNewBase64(filePath, content)
}

func (DesktopClient) DeleteFolderNew(filePath string) error {
	return vfs.DeleteFolderNew(filePath)
}

func (DesktopClient) DeleteFolder(container, fileName string) error {
	return vfs.DeleteFolder(container, fileName)
}

func (DesktopClient) DownloadFile(url, container, fileName string, headers map[string]string) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		lgr.Info("Error:", err)
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	err = vfs.WriteFileToWD(container, fileName, body)
	return err == nil, err
}

func (DesktopClient) GetFileFolder(container, filePath string) string {
	return vfs.GetFileFolder(filepath.Join(container, filePath))
}

func (DesktopClient) Exists(filePath string) bool {
	return vfs.Exists(filePath)
}
