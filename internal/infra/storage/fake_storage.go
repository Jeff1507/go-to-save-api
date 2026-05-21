package storage

import "io"

type FakeStorage struct {
	Uploaded    bool
	Deleted     bool
	UploadedKey string
	DeletedKey  string
}

func NewFakeStorage() *FakeStorage {
	return &FakeStorage{}
}

func (fs *FakeStorage) Upload(content io.Reader, contentType, path string) (string, error) {
	fs.Uploaded = true
	fs.UploadedKey = path + "/random-id" + contentType

	return fs.UploadedKey, nil
}

func (fs *FakeStorage) Delete(key string) error {
	if key == "" {
		return nil
	}

	fs.Deleted = true
	fs.DeletedKey = key

	return nil
}

func (fs *FakeStorage) URL(key string) string {
	return "http://localhost/uploads/" + key
}
