package domain

import "context"

type EsInterface interface {
	GetLogs(ctx context.Context) error
}
