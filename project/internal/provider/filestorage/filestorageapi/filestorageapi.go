package filestorageapi

import (
	"context"
	"fmt"
	"net/http"

	"tickets/internal/pkg/log"
	"tickets/internal/provider/filestorage"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
)

type Client struct {
	clients *clients.Clients
}

func NewClient(clients *clients.Clients) *Client {
	if clients == nil {
		panic("NewClient: clients is nil")
	}

	return &Client{clients: clients}
}

func (c Client) StoreFile(ctx context.Context, fileID string, content string) error {
	resp, err := c.clients.Files.PutFilesFileIdContentWithTextBodyWithResponse(ctx, fileID, content)
	if err != nil {
		return fmt.Errorf("failed to store file: %w", err)
	}

	switch resp.StatusCode() {
	case http.StatusCreated:
		return nil
	case http.StatusConflict:
		log.New(ctx).With("file", fileID).Info("file already exists")
		return nil
	default:
		return fmt.Errorf("unexpected status code for PUT files/%s/content: %d", fileID, resp.StatusCode())
	}
}

var _ filestorage.Storage = &Client{}
