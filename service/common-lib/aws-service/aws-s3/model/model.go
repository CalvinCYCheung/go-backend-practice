package model

type S3Object struct {
	BucketName  string
	KeyName     string
	FileContent []byte
}
