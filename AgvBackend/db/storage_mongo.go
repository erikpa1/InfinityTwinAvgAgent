package db

import (
	"context"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"turtle/core/lgr"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoFilesCollectionName is the single collection that holds every file's
// content across every container. There is no separate chunks collection:
// a file is one or more documents in this collection, each named after the
// file's uid (plus a "#<index>" suffix when it's split into multiple parts).
const MongoFilesCollectionName = "files"

// MongoChunkThresholdBytes is the size above which a file is split into
// multiple MongoChunkSizeBytes-sized parts instead of being stored as a
// single part. Kept safely under MongoDB's 16MB BSON document limit.
const MongoChunkThresholdBytes = 8 * 1024 * 1024

// MongoChunkSizeBytes is the size of each part once a file is chunked.
const MongoChunkSizeBytes = 4 * 1024 * 1024

// ErrMongoFileNotFound is returned when a requested file has no metadata
// document in its container collection.
var ErrMongoFileNotFound = errors.New("file not found")

// mongoFileMeta is the pointer document kept in the "container" collection
// for every file. It only records where to find the file's parts in the
// shared files collection - it never holds file content itself.
type mongoFileMeta struct {
	FileName  string    `bson:"fileName"`
	UID       string    `bson:"uid"`
	Size      int64     `bson:"size"`
	PartCount int       `bson:"partCount"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

// mongoFilePart is one document in the shared files collection. Every file,
// regardless of size, is stored as PartCount of these; small files simply
// have exactly one part. Name is the uid on its own for single-part files,
// or "<uid>#<index>" for files split into several parts, so every document
// in the collection is addressable and unique by name.
type mongoFilePart struct {
	Name  string `bson:"_id"`
	UID   string `bson:"uid"`
	Index int    `bson:"index"`
	Data  []byte `bson:"data"`
}

func partName(uid string, index int) string {
	if index == 0 {
		return uid
	}
	return uid + "#" + strconv.Itoa(index)
}

// MongoStorageClient implements IStorageController on top of MongoDB. Every
// container keeps its own pointer collection (fileName -> uid), and every
// file's actual bytes live in the single shared MongoFilesCollectionName
// collection as one or more parts named after that uid. Fetching a file is
// transparent to the caller either way - the client looks up the pointer and
// reassembles parts in order automatically.
type MongoStorageClient struct{}

func NewMongoStorageClient() *MongoStorageClient {
	return &MongoStorageClient{}
}

func (m *MongoStorageClient) _PathSplit(path string) (string, string) {
	split := strings.SplitN(path, "/", 2)
	if len(split) == 1 {
		return split[0], ""
	}
	return split[0], split[1]
}

func (m *MongoStorageClient) metaCollection(container string) *mongo.Collection {
	return DB.Col(container)
}

func (m *MongoStorageClient) filesCollection() *mongo.Collection {
	return DB.Col(MongoFilesCollectionName)
}

func (m *MongoStorageClient) deleteParts(uid string) error {
	_, err := m.filesCollection().DeleteMany(context.TODO(), bson.M{"uid": uid})
	return err
}

func (m *MongoStorageClient) upload(container, fileName string, _bytes []byte) error {
	var previous mongoFileMeta
	hasPrevious := false
	if err := m.metaCollection(container).FindOne(context.TODO(), bson.M{"fileName": fileName}).Decode(&previous); err == nil {
		hasPrevious = true
	} else if err != mongo.ErrNoDocuments {
		lgr.Error(err.Error())
		return err
	}

	uid := primitive.NewObjectID().Hex()

	chunkSize := len(_bytes)
	if len(_bytes) >= MongoChunkThresholdBytes {
		chunkSize = MongoChunkSizeBytes
	}
	if chunkSize == 0 {
		chunkSize = 1
	}

	parts := make([]interface{}, 0, (len(_bytes)/chunkSize)+1)
	for offset, index := 0, 0; offset < len(_bytes) || index == 0; index++ {
		end := offset + chunkSize
		if end > len(_bytes) {
			end = len(_bytes)
		}
		parts = append(parts, mongoFilePart{
			Name:  partName(uid, index),
			UID:   uid,
			Index: index,
			Data:  _bytes[offset:end],
		})
		offset = end
	}

	if _, err := m.filesCollection().InsertMany(context.TODO(), parts); err != nil {
		lgr.Error(err.Error())
		return err
	}

	meta := mongoFileMeta{
		FileName:  fileName,
		UID:       uid,
		Size:      int64(len(_bytes)),
		PartCount: len(parts),
		UpdatedAt: time.Now(),
	}

	_, err := m.metaCollection(container).UpdateOne(
		context.TODO(),
		bson.M{"fileName": fileName},
		bson.M{"$set": meta},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		lgr.Error(err.Error())
		_ = m.deleteParts(uid)
		return err
	}

	if hasPrevious && previous.UID != uid {
		if err := m.deleteParts(previous.UID); err != nil {
			lgr.Error(err.Error())
		}
	}

	return nil
}

func (m *MongoStorageClient) fetch(container, fileName string) ([]byte, error) {
	var meta mongoFileMeta
	err := m.metaCollection(container).FindOne(context.TODO(), bson.M{"fileName": fileName}).Decode(&meta)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrMongoFileNotFound
		}
		lgr.Error(err.Error())
		return nil, err
	}

	cursor, err := m.filesCollection().Find(
		context.TODO(),
		bson.M{"uid": meta.UID},
		options.Find().SetSort(bson.M{"index": 1}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	data := make([]byte, 0, meta.Size)
	for cursor.Next(context.TODO()) {
		var part mongoFilePart
		if err := cursor.Decode(&part); err != nil {
			return nil, err
		}
		data = append(data, part.Data...)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *MongoStorageClient) GetFileBytes(container, fileName string) ([]byte, error) {
	return m.fetch(container, fileName)
}

func (m *MongoStorageClient) GetFileBytesNew(filePath string) ([]byte, error) {
	container, fileName := m._PathSplit(filePath)
	return m.fetch(container, fileName)
}

func (m *MongoStorageClient) GetFileString(container, fileName string) (string, error) {
	data, err := m.fetch(container, fileName)
	return string(data), err
}

func (m *MongoStorageClient) GetFileStringNew(filePath string) (string, error) {
	container, fileName := m._PathSplit(filePath)
	data, err := m.fetch(container, fileName)
	return string(data), err
}

func (m *MongoStorageClient) UploadFile(container, fileName string, _bytes []byte) error {
	return m.upload(container, fileName, _bytes)
}

func (m *MongoStorageClient) UploadFileNew(filePath string, _bytes []byte) error {
	container, fileName := m._PathSplit(filePath)
	return m.upload(container, fileName, _bytes)
}

func (m *MongoStorageClient) UploadFileString(container, fileName, content string) error {
	return m.upload(container, fileName, []byte(content))
}

func (m *MongoStorageClient) UploadFileStringNew(filePath, content string) error {
	container, fileName := m._PathSplit(filePath)
	return m.upload(container, fileName, []byte(content))
}

func (m *MongoStorageClient) UploadFileStringNewBase64(filePath, content string) error {
	container, fileName := m._PathSplit(filePath)
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return err
	}
	return m.upload(container, fileName, decoded)
}

func (m *MongoStorageClient) deleteFile(container, fileName string) error {
	var meta mongoFileMeta
	err := m.metaCollection(container).FindOne(context.TODO(), bson.M{"fileName": fileName}).Decode(&meta)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrMongoFileNotFound
		}
		lgr.Error(err.Error())
		return err
	}

	if err := m.deleteParts(meta.UID); err != nil {
		return err
	}

	_, err = m.metaCollection(container).DeleteOne(context.TODO(), bson.M{"fileName": fileName})
	return err
}

func (m *MongoStorageClient) DeleteFileNew(filePath string) error {
	container, fileName := m._PathSplit(filePath)
	return m.deleteFile(container, fileName)
}

func (m *MongoStorageClient) MoveFileNew(oldFilePath, newFilePath string) error {
	oldContainer, oldFileName := m._PathSplit(oldFilePath)
	newContainer, newFileName := m._PathSplit(newFilePath)

	data, err := m.fetch(oldContainer, oldFileName)
	if err != nil {
		return err
	}

	if err := m.upload(newContainer, newFileName, data); err != nil {
		return err
	}

	return m.deleteFile(oldContainer, oldFileName)
}

func (m *MongoStorageClient) deleteFolder(container, folderName string) error {
	prefix := strings.TrimSuffix(folderName, "/")
	pattern := "^" + regexp.QuoteMeta(prefix) + "(/|$)"

	cursor, err := m.metaCollection(container).Find(context.TODO(), bson.M{
		"fileName": bson.M{"$regex": pattern},
	})
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	uids := make([]string, 0)
	for cursor.Next(context.TODO()) {
		var meta mongoFileMeta
		if err := cursor.Decode(&meta); err != nil {
			return err
		}
		uids = append(uids, meta.UID)
	}
	if err := cursor.Err(); err != nil {
		return err
	}

	if len(uids) > 0 {
		if _, err := m.filesCollection().DeleteMany(context.TODO(), bson.M{"uid": bson.M{"$in": uids}}); err != nil {
			return err
		}
	}

	_, err = m.metaCollection(container).DeleteMany(context.TODO(), bson.M{
		"fileName": bson.M{"$regex": pattern},
	})
	return err
}

func (m *MongoStorageClient) DeleteFolderNew(filePath string) error {
	container, folderName := m._PathSplit(filePath)
	return m.deleteFolder(container, folderName)
}

func (m *MongoStorageClient) DeleteFolder(container, fileName string) error {
	return m.deleteFolder(container, fileName)
}

func (m *MongoStorageClient) DownloadFile(url, container, fileName string, headers map[string]string) (bool, error) {
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

	err = m.upload(container, fileName, body)
	return err == nil, err
}

func (m *MongoStorageClient) GetFileFolder(container, filePath string) string {
	return filepath.Dir(filePath)
}

func (m *MongoStorageClient) Exists(filePath string) bool {
	container, fileName := m._PathSplit(filePath)
	count, err := m.metaCollection(container).CountDocuments(context.TODO(), bson.M{"fileName": fileName})
	if err != nil {
		lgr.Error(err.Error())
		return false
	}
	return count > 0
}
