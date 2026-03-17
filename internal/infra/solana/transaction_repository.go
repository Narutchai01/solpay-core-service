// internal/infrastructure/solana/transaction_repository.go
package solana

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/models"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

type solanaTransactionRepository struct {
	client      *rpc.Client
	mintAddress string
}

func NewSolanaTransactionRepository(rpcURL string, mintAddress string) ports.SolanaClient {
	return &solanaTransactionRepository{
		client:      rpc.New(rpcURL),
		mintAddress: mintAddress,
	}
}

func parseBase58PublicKey(fieldName, value string) (solana.PublicKey, error) {
	publicKey, err := solana.PublicKeyFromBase58(value)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("invalid %s: %w", fieldName, err)
	}
	return publicKey, nil
}

func (r *solanaTransactionRepository) BuildUnsignedTransfer(ctx context.Context, req models.BuildTXUnsigned) (string, error) {
	receiverAddress := config.LoadConfig().RECEIVE_ADDRESS

	senderPubkey, err := parseBase58PublicKey("sender address", req.SenderAddress)
	if err != nil {
		return "", err
	}

	receiverPubkey, err := parseBase58PublicKey("receive address", receiverAddress)
	if err != nil {
		return "", err
	}

	mintPubkey, err := parseBase58PublicKey("mint token address", r.mintAddress)
	if err != nil {
		return "", err
	}

	senderATA, _, err := solana.FindAssociatedTokenAddress(senderPubkey, mintPubkey)
	if err != nil {
		return "", err
	}

	receiverATA, _, err := solana.FindAssociatedTokenAddress(receiverPubkey, mintPubkey)
	if err != nil {
		return "", err
	}

	// blockhash
	recent, err := r.client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return "", err
	}

	transferIx, err := token.NewTransferCheckedInstruction(
		req.Amount,
		req.Decimals,
		senderATA,
		mintPubkey,
		receiverATA,
		senderPubkey,
		[]solana.PublicKey{},
	).ValidateAndBuild()
	if err != nil {
		return "", err
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{transferIx},
		recent.Value.Blockhash,
		solana.TransactionPayer(senderPubkey),
	)
	if err != nil {
		return "", err
	}

	serialized, err := tx.MarshalBinary()
	if err != nil {
		return "", err
	}

	txHash := base64.StdEncoding.EncodeToString(serialized)
	return txHash, nil
}
