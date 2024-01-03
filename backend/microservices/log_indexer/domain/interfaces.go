package domain

import "context"

type EsInterface interface {
	IndexLog(ctx context.Context, log Log) error
}
