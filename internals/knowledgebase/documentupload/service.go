package documentupload

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/document"
	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/storage"
	"go.uber.org/zap"
)

type DocumentUploadService struct {
	uploadRepo     Repository
	objectStore    storage.ObjectStore
	documentReader document.Reader
	logger         *zap.Logger
}

var uploadExpiry = 15 * time.Minute

func NewDocumentUploadService(uploadRepo Repository, objectStore storage.ObjectStore, documentReader document.Reader, logger *zap.Logger) *DocumentUploadService {
	return &DocumentUploadService{
		uploadRepo:     uploadRepo,
		objectStore:    objectStore,
		documentReader: documentReader,
		logger:         logger,
	}
}

func (dus *DocumentUploadService) Initiate(ctx context.Context, input kbmodel.InitiateInput) (kbmodel.InitiateResult, error) {
	if input.TenantID == uuid.Nil {
		return kbmodel.InitiateResult{}, ErrInvalidTenantID
	}
	if input.DocumentID == uuid.Nil {
		return kbmodel.InitiateResult{}, ErrInvalidDocumentID
	}
	if input.CreatedBy == uuid.Nil {
		return kbmodel.InitiateResult{}, ErrInvalidUserID
	}

	filename := safeFilename(input.OriginalFilename)
	if filename == "" {
		return kbmodel.InitiateResult{}, ErrInvalidOriginalFileName
	}
	if input.SizeBytes <= 0 {
		return kbmodel.InitiateResult{}, ErrInvalidSize
	}

	document, err := dus.documentReader.GetByID(ctx, input.TenantID, input.DocumentID)
	if err != nil {
		return kbmodel.InitiateResult{}, fmt.Errorf("failed to get document: %w", err)
	}
	if document.Status == kbmodel.DocumentStatusArchived {
		return kbmodel.InitiateResult{}, ErrDocumentArchived
	}
	uploadID := uuid.New()

	objectKey := fmt.Sprintf("tenants/%s/documents/%s/upload/%s/%s", input.TenantID.String(), input.DocumentID.String(), uploadID.String(), filename)

	now := time.Now().UTC()

	u := kbmodel.DocumentUpload{
		ID:                uploadID,
		TenantID:          input.TenantID,
		DocumentID:        input.DocumentID,
		ObjectKey:         objectKey,
		OriginalFileName:  filename,
		ContentType:       input.ContentType,
		ExpectedSizeBytes: input.SizeBytes,
		Status:            kbmodel.DocumentUploadPending,
		CreatedBy:         input.CreatedBy,
		ExpiredAt:         now.Add(uploadExpiry),
		CreatedAt:         now,
	}
	if err := dus.uploadRepo.Create(ctx, u); err != nil {
		return kbmodel.InitiateResult{}, fmt.Errorf("upload persistance failed: %w", err)
	}

	presigned, err := dus.objectStore.PresignPut(ctx, objectKey, input.ContentType, uploadExpiry)
	if err != nil {
		return kbmodel.InitiateResult{}, fmt.Errorf("failed to get presigned url: %w", err)
	}
	return kbmodel.InitiateResult{UploadID: uploadID, URL: presigned.URL, ExpiresAt: presigned.ExpiresAt}, nil
}

func safeFilename(name string) string {
	name = strings.ReplaceAll(name, "\\", "/")

	parts := strings.Split(name, "/")
	name = parts[len(parts)-1]

	return strings.TrimSpace(name)
}
