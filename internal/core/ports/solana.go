package ports

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/models"
)

type SolanaClient interface {
	BuildUnsignedTransfer(ctx context.Context, req models.BuildTXUnsigned) (*string, error)
}
