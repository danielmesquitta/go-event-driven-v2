package filestorage

import "context"

type Storage interface {
	StoreFile(ctx context.Context, fileID string, content string) error
}
