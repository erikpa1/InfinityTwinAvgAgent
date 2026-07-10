package db

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"turtle/core/lgr"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// Azure Storage Quickstart Sample - Demonstrate how to upload, list, download, and delete blobs.
//
// Documentation References:
// - What is a Storage Account - https://docs.microsoft.com/azure/storage/common/storage-create-storage-account
// - Blob Service Concepts - https://docs.microsoft.com/rest/api/storageservices/Blob-Service-Concepts
// - Blob Service Go SDK API - https://godoc.org/github.com/Azure/azure-storage-blob-go
// - Blob Service REST API - https://docs.microsoft.com/rest/api/storageservices/Blob-Service-REST-API
// - Scalability and performance targets - https://docs.microsoft.com/azure/storage/common/storage-scalability-targets
// - Azure Storage Performance and Scalability checklist https://docs.microsoft.com/azure/storage/common/storage-performance-checklist
// - Storage Emulator - https://docs.microsoft.com/azure/storage/common/storage-use-emulator

func handleError(err error) {
	if err != nil {
		lgr.Error(err.Error())
	}
}

var CONN_STR = ""

func TestStorage() {
	fmt.Printf("Azure Blob storage quick start sample\n")

	// TODO: replace <storage-account-name> with your actual storage account name

	ctx := context.TODO()

	client, err := azblob.NewClientFromConnectionString(CONN_STR, nil)
	handleError(err)

	// Create the container
	containerName := "quickstart-sample-container"
	fmt.Printf("Creating a container named %s\n", containerName)
	_, err = client.CreateContainer(ctx, containerName, nil)
	handleError(err)

	data := []byte("\nHello, world! This is a blob.\n")
	blobName := "sample-blob"

	// Upload to data to blob storage
	fmt.Printf("Uploading a blob named %s\n", blobName)
	_, err = client.UploadBuffer(ctx, containerName, blobName, data, &azblob.UploadBufferOptions{})
	handleError(err)

	// List the blobs in the container

	lgr.Info("Listing the blobs in the container:")

	pager := client.NewListBlobsFlatPager(containerName, &azblob.ListBlobsFlatOptions{
		Include: azblob.ListBlobsInclude{Snapshots: true, Versions: true},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		handleError(err)

		for _, blob := range resp.Segment.BlobItems {
			lgr.Info(*blob.Name)
		}
	}

	// Download the blob
	get, err := client.DownloadStream(ctx, containerName, blobName, nil)
	handleError(err)

	downloadedData := bytes.Buffer{}
	retryReader := get.NewRetryReader(ctx, &azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(retryReader)
	handleError(err)

	err = retryReader.Close()
	handleError(err)

	// Print the content of the blob we created
	lgr.Info("Blob contents:")
	lgr.Info(downloadedData.String())

	fmt.Printf("Press enter key to delete resources and exit the application.\n")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	fmt.Printf("Cleaning up.\n")

	// Delete the blob
	fmt.Printf("Deleting the blob " + blobName + "\n")

	_, err = client.DeleteBlob(ctx, containerName, blobName, nil)
	handleError(err)

	// Delete the container
	fmt.Printf("Deleting the container " + containerName + "\n")
	_, err = client.DeleteContainer(ctx, containerName, nil)
	handleError(err)
}

// Dummy AzureStorageExecutioner to replace later with actual implementation
type AzureStorageExecutioner struct{}

func (a *AzureStorageExecutioner) GetContainerHandler(container string) *AzureStorageExecutioner {
	return a
}

func (a *AzureStorageExecutioner) GetFile(fileName string) ([]byte, error) {
	return []byte{}, nil
}

func (a *AzureStorageExecutioner) GetFileString(fileName string) (string, error) {
	return "", nil
}

func (a *AzureStorageExecutioner) SaveFile(fileName string, data []byte) (bool, error) {
	return true, nil
}

func (a *AzureStorageExecutioner) SaveFileString(fileName, content string) (bool, error) {
	return true, nil
}

func (a *AzureStorageExecutioner) DeleteFolder(folderName string) error {
	return nil
}

func (a *AzureStorageExecutioner) DeleteFile(fileName string) error {
	return nil
}

func (a *AzureStorageExecutioner) Exists(fileName string) bool {
	return true
}

// StorageClient defines methods for interacting with cloud storage
type StorageClient struct {
	conn *AzureStorageExecutioner
}

func (sc *StorageClient) _PathSplit(path string) (string, string) {
	split := strings.SplitN(path, "/", 2)
	return split[0], split[1]
}

func (sc *StorageClient) GetFileBytesNew(filePath string) ([]byte, error) {
	container, file := sc._PathSplit(filePath)
	return sc.conn.GetContainerHandler(container).GetFile(file)
}

func (sc *StorageClient) GetFileBytes(container, fileName string) ([]byte, error) {
	return sc.conn.GetContainerHandler(container).GetFile(fileName)
}

func (sc *StorageClient) GetFileString(container, fileName string) (string, error) {
	return sc.conn.GetContainerHandler(container).GetFileString(fileName)
}

func (sc *StorageClient) GetFileStringNew(filePath string) (string, error) {
	container, file := sc._PathSplit(filePath)
	return sc.conn.GetContainerHandler(container).GetFileString(file)
}

func (sc *StorageClient) UploadFile(container, fileName string, _bytes []byte) (bool, error) {
	return sc.conn.GetContainerHandler(container).SaveFile(fileName, _bytes)
}

func (sc *StorageClient) UploadFileNew(filePath string, _bytes []byte) (bool, error) {
	container, file := sc._PathSplit(filePath)
	return sc.conn.GetContainerHandler(container).SaveFile(file, _bytes)
}

func (sc *StorageClient) UploadFileString(container, fileName, content string) (bool, error) {
	return sc.conn.GetContainerHandler(container).SaveFileString(fileName, content)
}

func (sc *StorageClient) UploadFileStringNew(filePath, content string) (bool, error) {
	container, file := sc._PathSplit(filePath)
	return sc.conn.GetContainerHandler(container).SaveFileString(file, content)
}

func (sc *StorageClient) DeleteFolderNew(filePath string) error {
	container, folder := sc._PathSplit(filePath)
	return sc.conn.GetContainerHandler(container).DeleteFolder(folder)
}

func (sc *StorageClient) DeleteFolder(container, folderName string) error {
	return sc.conn.GetContainerHandler(container).DeleteFolder(folderName)
}

func (sc *StorageClient) DeleteFileNew(filePath string) error {
	container, file := sc._PathSplit(filePath)
	return sc.conn.GetContainerHandler(container).DeleteFile(file)
}

func (sc *StorageClient) MoveFileNew(oldFilePath, newFilePath string) error {
	oldContainer, oldFile := sc._PathSplit(oldFilePath)
	newContainer, newFile := sc._PathSplit(newFilePath)

	handler := sc.conn.GetContainerHandler(oldContainer)
	data, err := handler.GetFile(oldFile)
	if err != nil {
		return err
	}

	if _, err := sc.conn.GetContainerHandler(newContainer).SaveFile(newFile, data); err != nil {
		return err
	}

	return handler.DeleteFile(oldFile)
}

func (sc *StorageClient) GetFileFolder(container, filePath string) string {
	return filepath.Dir(filePath)
}

func (sc *StorageClient) Exists(filePath string) bool {
	container, file := sc._PathSplit(filePath)
	return sc.conn.GetContainerHandler(container).Exists(file)
}
