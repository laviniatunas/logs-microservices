package domain

import "context"

type EsInterface interface {
	IndexLog(ctx context.Context, log Log) error
}

type AlertsInterface interface {
	TriggerAlert(ctx context.Context, logBytes []byte) error
}
