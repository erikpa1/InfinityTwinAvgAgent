package db

import (
	"encoding/base64"
	"turtle/core/lgr"
	"turtle/credentials"
)

// StorageController controls which client to use based on credentials
type StorageController struct {
	client IStorageController
}

func NewStorageController() *StorageController {
	sc := &StorageController{}

	switch credentials.GetStorageBackend() {
	case "mongo":
		sc.client = NewMongoStorageClient()
	default:
		sc.client = &DesktopClient{}
	}

	return sc
}

func (sc *StorageController) GetFileBase64(filePath string) (string, error) {
	_bytes, err := sc.client.GetFileBytesNew(filePath)
	if err != nil || len(_bytes) == 0 {
		lgr.Info("Unable to find file:", filePath)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(_bytes), nil
}

func (sc *StorageController) GetFileBytesNew(filePath string) ([]byte, error) {
	return sc.client.GetFileBytesNew(filePath)
}

func (sc *StorageController) GetFileBytes(container string, fileName string) ([]byte, error) {
	return sc.client.GetFileBytes(container, fileName)
}

func (sc *StorageController) GetFileString(container string, fileName string) (string, error) {
	return sc.client.GetFileString(container, fileName)
}

func (sc *StorageController) GetFileStringNew(filePath string) (string, error) {
	return sc.client.GetFileStringNew(filePath)
}

func (sc *StorageController) UploadFile(container string, fileName string, _bytes []byte) (bool, error) {
	err := sc.client.UploadFile(container, fileName, _bytes)
	return err == nil, err
}

func (sc *StorageController) UploadFileNew(filePath string, _bytes []byte) (bool, error) {
	err := sc.client.UploadFileNew(filePath, _bytes)
	return err == nil, err
}

func (sc *StorageController) UploadFileString(container, fileName, content string) (bool, error) {
	err := sc.client.UploadFileString(container, fileName, content)
	return err == nil, err
}

func (sc *StorageController) UploadFileStringNew(filePath, content string) (bool, error) {
	err := sc.client.UploadFileStringNew(filePath, content)
	return err == nil, err
}

func (sc *StorageController) UploadFileStringNewBase64(filePath, content string) (bool, error) {
	err := sc.client.UploadFileStringNewBase64(filePath, content)
	return err == nil, err
}

func (sc *StorageController) DeleteFolderNew(filePath string) error {
	return sc.client.DeleteFolderNew(filePath)
}

func (sc *StorageController) DeleteFolder(container, folderName string) error {
	return sc.client.DeleteFolder(container, folderName)
}

func (sc *StorageController) DeleteFileNew(filePath string) error {
	return sc.client.DeleteFileNew(filePath)
}

func (sc *StorageController) MoveFileNew(oldFilePath, newFilePath string) error {
	return sc.client.MoveFileNew(oldFilePath, newFilePath)
}

func (sc *StorageController) GetFileFolder(container, filePath string) string {
	return sc.client.GetFileFolder(container, filePath)
}

func (sc *StorageController) Exists(filePath string) bool {
	return sc.client.Exists(filePath)
}

var SC = NewStorageController()
