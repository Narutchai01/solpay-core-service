package supabase

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"

	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseStorage struct {
	SUPABASE_PRIVATE_KEY string
	SUPABASE_URL         string
}

func NewSupabaseStorage(privateKey, url string) *SupabaseStorage {
	return &SupabaseStorage{
		SUPABASE_PRIVATE_KEY: strings.TrimSpace(privateKey),
		SUPABASE_URL:         strings.TrimSpace(url),
	}
}

func (s *SupabaseStorage) UploadFile(bucketName string, filePath string, fileContent []byte) (string, error) {
	bucketName = strings.TrimSpace(bucketName)
	objectPath := strings.TrimLeft(strings.TrimSpace(filePath), "/")
	apiKey := strings.TrimSpace(s.SUPABASE_PRIVATE_KEY)
	projectBaseURL, storageAPIURL := normalizeSupabaseBaseURL(strings.TrimSpace(s.SUPABASE_URL))

	if bucketName == "" {
		return "", errors.New("bucket name is required")
	}
	if objectPath == "" {
		return "", errors.New("file path is required")
	}
	if len(fileContent) == 0 {
		return "", errors.New("file content is empty")
	}
	if apiKey == "" {
		return "", errors.New("supabase private key is required")
	}
	if projectBaseURL == "" || storageAPIURL == "" {
		return "", errors.New("supabase url is required")
	}

	objectPath = strings.TrimPrefix(objectPath, bucketName+"/")

	storageClient := storage_go.NewClient(storageAPIURL, apiKey, nil)
	contentType := http.DetectContentType(fileContent)

	options := storage_go.FileOptions{
		ContentType: &contentType,
	}

	reader := bytes.NewReader(fileContent)
	_, err := storageClient.UploadFile(bucketName, objectPath, reader, options)
	if err != nil {
		if strings.Contains(err.Error(), "Invalid Compact JWS") {
			return "", fmt.Errorf("invalid supabase key format: use a JWT service_role/anon key for SUPABASE_PRIVATE_KEY: %w", err)
		}
		return "", err
	}

	if bucketName == "slip" {
		publicUrlResp := storageClient.GetPublicUrl(bucketName, objectPath)
		return publicUrlResp.SignedURL, nil
	}

	signedUrlResp, err := storageClient.CreateSignedUrl(bucketName, objectPath, 604800) // 7 days (604800 seconds)
	if err != nil {
		return "", fmt.Errorf("failed to create signed url: %w", err)
	}

	url := storageAPIURL + signedUrlResp.SignedURL

	return url, nil
}

func normalizeSupabaseBaseURL(raw string) (projectBaseURL, storageAPIURL string) {
	base := strings.TrimRight(strings.TrimSpace(raw), "/")
	if base == "" {
		return "", ""
	}

	if strings.HasSuffix(base, "/storage/v1") {
		return strings.TrimSuffix(base, "/storage/v1"), base
	}

	return base, base + "/storage/v1"
}
