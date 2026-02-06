package spreadsheetapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients/spreadsheets"
)

type Client struct {
	// we are not mocking this client: it's pointless to use interface here
	clients *clients.Clients
}

func NewClient(clients *clients.Clients) *Client {
	if clients == nil {
		panic("NewClient: clients is nil")
	}

	return &Client{clients: clients}
}

func (c Client) AppendRow(ctx context.Context, spreadsheetName string, row []string) error {
	resp, err := c.clients.Spreadsheets.PostSheetsSheetRowsWithResponse(ctx, spreadsheetName, spreadsheets.PostSheetsSheetRowsJSONRequestBody{
		Columns: row,
	})
	if err != nil {
		return fmt.Errorf("failed to post row: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to post row: unexpected status code %d", resp.StatusCode())
	}

	return nil
}
