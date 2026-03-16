// internal/infrastructure/solana/transaction_repository.go
package solana

import (
	"context"
	"encoding/base64"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/models"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

type solanaTransactionRepository struct {
	client      *rpc.Client
	mintAddress solana.PublicKey
}

func NewSolanaTransactionRepository(rpcURL string, mintAddress string) ports.SolanaClient {
	return &solanaTransactionRepository{
		client:      rpc.New(rpcURL),
		mintAddress: solana.MustPublicKeyFromBase58(mintAddress),
	}
}

func (r *solanaTransactionRepository) BuildUnsignedTransfer(ctx context.Context, req models.BuildTXUnsigned) (*string, error) {
	receiverAddress := config.LoadConfig().RECEIVE_ADDRESS
	senderPubkey := solana.MustPublicKeyFromBase58(req.SenderAddress)
	receiverPubkey := solana.MustPublicKeyFromBase58(receiverAddress)

	senderATA, _, err := solana.FindAssociatedTokenAddress(senderPubkey, r.mintAddress)
	if err != nil {
		return nil, err
	}

	receiverATA, _, err := solana.FindAssociatedTokenAddress(receiverPubkey, r.mintAddress)
	if err != nil {
		return nil, err
	}

	// blockhash
	recent, err := r.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}

	transferIx, err := token.NewTransferCheckedInstruction(
		req.Amount,
		req.Decimals,
		senderATA,
		r.mintAddress,
		receiverATA,
		senderPubkey,
		[]solana.PublicKey{},
	).ValidateAndBuild()
	if err != nil {
		return nil, err
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{transferIx},
		recent.Value.Blockhash,
		solana.TransactionPayer(senderPubkey),
	)
	if err != nil {
		return nil, err
	}

	serialized, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}

	txHash := base64.StdEncoding.EncodeToString(serialized)
	return &txHash, nil
}
