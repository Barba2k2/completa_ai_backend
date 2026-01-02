package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type FirebaseService struct {
	authClient    *auth.Client
	storageClient *storage.Client
	bucketName    string
}

func NewFirebaseService(projectID, bucketName string) (*FirebaseService, error) {
	ctx := context.Background()

	conf := &firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(ctx, conf, option.WithoutAuthentication())
	if err != nil {
		return nil, err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(bucketName, "gs://") {
		bucketName = strings.TrimPrefix(bucketName, "gs://")
	}
	bucketName = strings.TrimSuffix(bucketName, "/")

	if bucketName == "" && projectID != "" {
		bucketName = projectID + ".appspot.com"
	}

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &FirebaseService{
		authClient:    authClient,
		storageClient: storageClient,
		bucketName:    bucketName,
	}, nil
}

func (s *FirebaseService) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return s.authClient.VerifyIDToken(ctx, idToken)
}

func (s *FirebaseService) UploadProfileImage(
	ctx context.Context,
	userID string,
	content io.Reader,
	contentType string,
) (string, error) {
	if s.storageClient == nil {
		return "", fmt.Errorf("storage client not configured")
	}
	if s.bucketName == "" {
		return "", fmt.Errorf("storage bucket not configured")
	}

	if contentType == "" {
		contentType = "image/jpeg"
	}

	extension := ".jpg"
	switch contentType {
	case "image/png":
		extension = ".png"
	case "image/webp":
		extension = ".webp"
	case "image/jpeg", "image/jpg":
		extension = ".jpg"
	}

	objectPath := fmt.Sprintf("profile_images/%s%s", userID, extension)
	token := uuid.NewString()

	log.Printf("[FirebaseStorage] Uploading profile image user_id=%s bucket=%s object=%s", userID, s.bucketName, objectPath)

	writer := s.storageClient.Bucket(s.bucketName).Object(objectPath).NewWriter(ctx)
	writer.ContentType = contentType
	writer.Metadata = map[string]string{
		"firebaseStorageDownloadTokens": token,
	}

	if _, err := io.Copy(writer, content); err != nil {
		_ = writer.Close()
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	encodedPath := url.PathEscape(objectPath)
	downloadURL := fmt.Sprintf(
		"https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s",
		s.bucketName,
		encodedPath,
		token,
	)

	return downloadURL, nil
}
