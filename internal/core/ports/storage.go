package ports

type Storage interface {
	UploadFile(bucketName string, filePath string, fileContent []byte) (string, error)
}
