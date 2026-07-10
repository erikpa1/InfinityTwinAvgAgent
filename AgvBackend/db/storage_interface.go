package db

type IStorageController interface {
	GetFileBytes(container, fileName string) ([]byte, error)
	GetFileBytesNew(filePath string) ([]byte, error)
	GetFileString(container, fileName string) (string, error)
	GetFileStringNew(filePath string) (string, error)
	UploadFile(container, fileName string, _bytes []byte) error
	UploadFileNew(filePath string, _bytes []byte) error
	UploadFileString(container, fileName, content string) error
	UploadFileStringNew(filePath, content string) error
	UploadFileStringNewBase64(filePath, content string) error
	DeleteFolderNew(filePath string) error
	DeleteFolder(container, fileName string) error
	DeleteFileNew(filePath string) error
	MoveFileNew(oldFilePath, newFilePath string) error
	DownloadFile(url, container, fileName string, headers map[string]string) (bool, error)
	GetFileFolder(container, filePath string) string
	Exists(filePath string) bool
}
