package domain

import "context"

type EsInterface interface {
	GetLogs(ctx context.Context) error
}

type AlertRepo interface {
	Start(ctx context.Context) error
}
