package spreadsheet

import "context"

type API interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}
