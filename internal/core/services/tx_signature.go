package services

import (
	"net/url"
	"path"
	"strings"

	"github.com/gagliardetto/solana-go"
)

func signatureFromTxHash(txHash string) string {
	normalized := strings.TrimSpace(txHash)
	if normalized == "" {
		return ""
	}

	if sig := signatureFromSignedTransactionBase64(normalized); sig != "" {
		return sig
	}

	if isSolanaSignature(normalized) {
		return normalized
	}

	if sig := signatureFromExplorerURL(normalized); sig != "" {
		return sig
	}

	return normalized
}

func signatureFromSignedTransactionBase64(base64Tx string) string {
	tx, err := solana.TransactionFromBase64(base64Tx)
	if err != nil || len(tx.Signatures) == 0 {
		return ""
	}

	sig := tx.Signatures[0]
	if sig.IsZero() {
		return ""
	}

	return sig.String()
}

func signatureFromExplorerURL(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}

	candidates := []string{
		strings.TrimSpace(parsed.Query().Get("signature")),
		strings.TrimSpace(parsed.Query().Get("sig")),
	}

	segment := strings.TrimSpace(path.Base(strings.TrimSuffix(parsed.Path, "/")))
	if segment != "" && segment != "." && segment != "/" {
		if decoded, err := url.PathUnescape(segment); err == nil {
			segment = decoded
		}
		candidates = append(candidates, segment)
	}

	for _, candidate := range candidates {
		if isSolanaSignature(candidate) {
			return candidate
		}
		if sig := signatureFromSignedTransactionBase64(candidate); sig != "" {
			return sig
		}
	}

	return ""
}

func isSolanaSignature(value string) bool {
	if strings.TrimSpace(value) == "" {
		return false
	}

	sig, err := solana.SignatureFromBase58(value)
	if err != nil {
		return false
	}

	return !sig.IsZero()
}
