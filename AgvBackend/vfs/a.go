package vfs

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"turtle/core/lgr"
	"turtle/credentials"
)

type ListFileInfoStruct struct {
	FilePath string `json:"file_path"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
}

func GetExeFile() string {
	exePath, err := os.Executable()

	if err != nil {
		return "failedtogetworkdir"
	} else {
		return exePath
	}
}

func IsInDevelopment() bool {
	exePath, err := os.Executable()

	if IsLinux() {
		return false
	}

	if err != nil {
		lgr.Info("Error:", err)
		return false
	}
	exeDir := filepath.Dir(exePath)

	lgr.Info(exeDir)

	// Check if the binary is running from a temporary directory
	if strings.Contains(exeDir, "\\AppData\\") ||
		strings.Contains(exeDir, "src-go") ||
		strings.Contains(exeDir, "go-build") {
		return true
	} else {
		return false
	}
}

// FindAllFilesWithExtension - Finds all files with the given extension in a directory (optionally recursive)
func FindAllFilesWithExtension(folderPath string, extension string, recursive bool) ([]string, error) {
	var fileList []string
	if recursive {
		err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), "."+extension) {
				fileList = append(fileList, path)
			}
			return nil
		})
		return fileList, err
	}

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), "."+extension) {
			fileList = append(fileList, filepath.Join(folderPath, file.Name()))
		}
	}
	return fileList, nil
}

// FindAllFoldersWithExtension - Finds all folders containing files with the given extension
func FindAllFoldersWithExtension(folderPath string, extension string) ([]string, error) {
	var folderList []string
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), "."+extension) {
			folderList = append(folderList, filepath.Dir(path))
		}
		return nil
	})
	return folderList, err
}

// GetFileName - Returns file name from a path, optionally including the extension
func GetFileName(path string, includeExtension bool) string {
	if includeExtension {
		return filepath.Base(path)
	}
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

// Exists - Checks if the path exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetFileFolder - Returns the folder for a given file path
func GetFileFolder(filePath string) string {
	return filepath.Dir(filePath)
}

// GetFileFolder - Returns the folder for a given file path
func GetFilePath(filePath string) string {
	finalPath := filepath.Join(GetWorkingDirectory(), filePath)
	return finalPath
}

// GetFileFolder - Returns the folder for a given file path
func CreateFolder(folderPath string) {
	finalPath := filepath.Join(GetWorkingDirectory(), folderPath)
	err := os.MkdirAll(filepath.Dir(finalPath), 0755)
	if err != nil {
		lgr.Info("Create folder error:", err)
	}
}

// FindPreviewFromPostfix - Finds preview image by appending a postfix
func FindPreviewFromPostfix(filePath string, postFix string) string {
	folder := GetFileFolder(filePath)
	name := GetFileName(filePath, false)
	searched := filepath.Join(folder, name+postFix)
	if Exists(searched) {
		return searched
	}
	return ""
}

// GetPreviewPathInFolder - Returns the path of a preview.png in a folder
func GetPreviewPathInFolder(folder string) string {
	path := filepath.Join(folder, "preview.png")
	if Exists(path) {
		return path
	}
	return ""
}

// GetFirstFileWithExtension - Returns the first file with the given extension in a folder
func GetFirstFileWithExtension(folderPath string, extension string) string {
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), "."+extension) {
			return filepath.SkipDir
		}
		return nil
	})
	if err == nil {
		return ""
	}
	return err.Error()
}

// IsLinux - Check if the system is Linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

// GetWorkingDirectory - Returns the appropriate working directory based on OS
func GetWorkingDirectory() string {
	if IsLinux() {
		return credentials.LinuxWorkspace()
	} else if IsDarwin() {
		return credentials.GetDarwinWorkspace()
	}

	basePath := ""

	stPath := credentials.GetStoragePath()

	switch credentials.GetStoragePath() {
	case "LOCALAPPDATA":
		basePath = os.Getenv("LOCALAPPDATA")
	case "ProgramData":
		basePath = os.Getenv("ProgramData")
	default:
		basePath = stPath
	}

	return filepath.Join(basePath, credentials.GetAppName())
}

// GetFileFolderNew - Similar to GetFileFolder but normalizes path
func GetFileFolderNew(filePath string) string {
	return filepath.Dir(filepath.Clean(filePath))
}

// WriteFileToWD - Writes bytes to a file in the working directory
func WriteFileToWD(folder, filePath string, data []byte) error {
	wdPath := GetWorkingDirectory()
	finalPath := filepath.Join(wdPath, folder, filePath)

	err := os.MkdirAll(filepath.Dir(finalPath), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(finalPath, data, 0644)
}

func UploadFileChunk(folder, filePath string, data []byte, offset int) error {
	wdPath := GetWorkingDirectory()
	finalPath := filepath.Join(wdPath, folder, filePath)

	lgr.Error("Chunking: ", finalPath)

	err := os.MkdirAll(filepath.Dir(finalPath), 0755)
	if err != nil {
		return err
	}
	// Open the file in append mode, create if it doesn't exist
	file, err := os.OpenFile(finalPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file info to check its size
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Check if the offset is valid
	fileSize := fileInfo.Size()
	if offset > int(fileSize) {
		return fmt.Errorf("invalid offset: offset %d is greater than file size %d", offset, fileSize)
	}

	// If offset is specified, seek to that position
	if offset > 0 {
		if _, err := file.Seek(int64(offset), io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek to offset: %w", err)
		}
	} else {
		// If no offset provided, seek to end of file for appending
		if _, err := file.Seek(0, io.SeekEnd); err != nil {
			return fmt.Errorf("failed to seek to end of file: %w", err)
		}
	}

	// Write the data to the file
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	return nil
}

// WriteFileToWDNew - Similar to WriteFileToWD but without a folder prefix
func WriteFileToWDNew(filePath string, data []byte) error {
	finalPath := filepath.Join(GetWorkingDirectory(), filePath)
	err := os.MkdirAll(filepath.Dir(finalPath), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(finalPath, data, 0644)
}

// DeleteFolderNew - Deletes a folder in the working directory
func DeleteFolderNew(folderPath string) error {
	finalPath := filepath.Join(GetWorkingDirectory(), folderPath)
	return os.RemoveAll(finalPath)
}

// DeleteFolderNew - Deletes a folder in the working directory
func DeleteFolder(folderPath string, subfolder string) error {
	finalPath := filepath.Join(GetWorkingDirectory(), folderPath, subfolder)
	return os.RemoveAll(finalPath)
}

// DeleteFileNew - Deletes a file in the working directory
func DeleteFileNew(filePath string) error {
	finalPath := filepath.Join(GetWorkingDirectory(), filePath)
	return os.Remove(finalPath)
}

// MoveFileNew - Moves/renames a file in the working directory, creating the
// destination folder if needed
func MoveFileNew(oldFilePath, newFilePath string) error {
	wdPath := GetWorkingDirectory()
	oldFinalPath := filepath.Join(wdPath, oldFilePath)
	newFinalPath := filepath.Join(wdPath, newFilePath)

	if err := os.MkdirAll(filepath.Dir(newFinalPath), 0755); err != nil {
		return err
	}
	return os.Rename(oldFinalPath, newFinalPath)
}

// MakeDirs - Creates a directory (recursive)
func MakeDirs(folder string) error {
	wdPath := GetWorkingDirectory()
	finalPath := filepath.Join(wdPath, folder)
	return os.MkdirAll(finalPath, 0755)
}

// OpenWDFolder - Opens a folder in the working directory (on Windows)
func OpenWDFolder(folder string) error {
	finalPath := filepath.Join(GetWorkingDirectory(), folder)
	cmd := exec.Command(GetExploreCommand(), finalPath)

	result := cmd.Start()

	if result != nil {
		lgr.ErrorStack(result.Error())
	}

	return result

}

// OpenWDFolder - Opens a folder in the working directory (on Windows)
func OpenFolder(folder string) error {
	cmd := exec.Command(GetExploreCommand(), folder)

	result := cmd.Start()

	if result != nil {
		lgr.ErrorStack(result.Error())
	}

	return result

}

func GetExploreCommand() string {
	if IsLinux() {
		return "xdg-open"
	} else if IsDarwin() {
		return "open"
	} else {
		return "explorer"
	}
}

// WriteFileStringToWD - Writes a string to a file in the working directory
func WriteFileStringToWD(folder, filePath, content string) error {
	wdPath := GetWorkingDirectory()
	finalPath := filepath.Join(wdPath, folder, filePath)
	err := os.MkdirAll(filepath.Dir(finalPath), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(finalPath, []byte(content), 0644)
}

// WriteFileStringToWD - Writes a string to a file in the working directory
func WriteFileString(folder, filePath, content string) error {

	finalPath := filepath.Join(folder, filePath)
	err := os.MkdirAll(filepath.Dir(finalPath), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(finalPath, []byte(content), 0644)
}

func WriteFileStringToWDNew(filePath, content string) error {
	wdPath := GetWorkingDirectory()
	finalPath := filepath.Join(wdPath, filePath)
	err := os.MkdirAll(filepath.Dir(finalPath), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(finalPath, []byte(content), 0644)
}

func WriteFileStringToWDNewBase64(filePath, content string) error {
	wdPath := GetWorkingDirectory()
	finalPath := filepath.Join(wdPath, filePath)

	// Decode Base64 content
	decodedContent, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return fmt.Errorf("Failed to decode Base64 content: %w", err)
	}

	// Create necessary directories
	err = os.MkdirAll(filepath.Dir(finalPath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directories for path %s: %w", finalPath, err)
	}

	// Write the decoded content to the file
	err = os.WriteFile(finalPath, decodedContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", finalPath, err)
	}

	return nil
}

// ListWdFiles - Returns a list of all files in a folder
func ListFile(folderPath string) ([]string, error) {
	var fileList []string

	// Walk through the directory
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Return the error if accessing a file fails
		}

		// Check if the current item is a file (not a directory)
		if !info.IsDir() {
			replaced := strings.Replace(path, folderPath+"\\", "", -1)
			replaced = strings.Replace(replaced, "\\", "/", -1)
			fileList = append(fileList, replaced) // Add the file to the list
		}
		return nil
	})

	if err != nil {
		lgr.Error(err.Error())
		return nil, err // Return error if Walk fails
	}

	return fileList, nil
}

// ListWdFiles - Returns a list of all files in a folder
func ListWdFiles(folderPath string) ([]string, error) {
	var fileList []string

	prefix := filepath.Join(GetWorkingDirectory(), folderPath)

	// Walk through the directory
	err := filepath.Walk(prefix, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Return the error if accessing a file fails
		}

		// Check if the current item is a file (not a directory)
		if !info.IsDir() {

			replaced := strings.Replace(path, prefix+"\\", "", -1)
			replaced = strings.Replace(replaced, "\\", "/", -1)
			fileList = append(fileList, replaced) // Add the file to the list
		}
		return nil
	})

	if err != nil {

		lgr.Error(err.Error())
		return nil, err // Return error if Walk fails
	}

	return fileList, nil
}

// ListWdFiles - Returns a list of all files in a folder
func ListFilesWithInfo(folderPath string) ([]ListFileInfoStruct, error) {
	var fileList []ListFileInfoStruct

	prefix := filepath.Join(GetWorkingDirectory(), folderPath)

	// Walk through the directory
	err := filepath.Walk(prefix, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Return the error if accessing a file fails
		}

		// Check if the current item is a file (not a directory)
		if !info.IsDir() {

			replaced := strings.Replace(path, prefix+"\\", "", -1)
			replaced = strings.Replace(replaced, "\\", "/", -1)

			tmp := ListFileInfoStruct{}
			tmp.FilePath = replaced
			tmp.Size = info.Size()
			tmp.Modified = info.ModTime().String()

			fileList = append(fileList, tmp) // Add the file to the list
		}
		return nil
	})

	if err != nil {
		lgr.Error(err.Error())
		return nil, err // Return error if Walk fails
	}

	return fileList, nil
}

// GetFileBytesFromWDNew - Reads bytes from a file in the working directory
func GetFileBytesFromWDNew(filePath string) ([]byte, error) {
	finalPath := filepath.Join(GetWorkingDirectory(), filePath)
	return os.ReadFile(finalPath)
}
func GetFileBytesFromWD(folder string, filePath string) ([]byte, error) {
	finalPath := filepath.Join(GetWorkingDirectory(), folder, filePath)
	return os.ReadFile(finalPath)
}

// IsDesktop - Checks if the system is a Windows desktop
func IsDesktop() bool {
	return runtime.GOOS == "windows"
}

// GetFileStringFromWDNew - Reads a string from a file in the working directory
func GetFileStringFromWDNew(filePath string) (string, error) {
	finalPath := filepath.Join(GetWorkingDirectory(), filePath)
	data, err := os.ReadFile(finalPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func GetFileString(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func GetFileStringFromWD(folder string, filePath string) (string, error) {
	finalPath := filepath.Join(GetWorkingDirectory(), folder, filePath)
	data, err := os.ReadFile(finalPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func GetFolderSize(folder string) int64 {
	finalPath := filepath.Join(GetWorkingDirectory(), folder)

	var size int64

	err := filepath.Walk(finalPath, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() { // Only add file sizes
			size += info.Size()
		}
		return nil
	})
	if err != nil {
		return 0
	}
	return size
}

func DeleteFile(folderPath string, file string) error {
	finalPath := filepath.Join(GetWorkingDirectory(), folderPath, file)

	error := os.Remove(finalPath)

	lgr.Error("Going to delete: ", finalPath)

	if error != nil {
		lgr.Error(error.Error())
	}

	return error
}

func GetFilePathFromWD(container string, file string) string {
	finalPath := filepath.Join(GetWorkingDirectory(), file)
	return finalPath

}
