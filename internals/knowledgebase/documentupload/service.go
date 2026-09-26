package documentupload

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"

	knowledgebase "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/base"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/document"
	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/storage"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/membership"
	tenantmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/tenant"
	"go.uber.org/zap"
)

type DocumentUploadService struct {
	uploadRepo       Repository
	objectStore      storage.ObjectStore
	tenantReader     tenant.Reader
	membershipReader membership.Reader
	documentRepo     document.Repository
	kbRepo           knowledgebase.Reader
	logger           *zap.Logger
}

var uploadExpiry = 15 * time.Minute

func NewDocumentUploadService(uploadRepo Repository, objectStore storage.ObjectStore, tenantReader tenant.Reader, membershipReader membership.Reader, documentRepo document.Repository, kbRepo knowledgebase.Reader, logger *zap.Logger) *DocumentUploadService {
	return &DocumentUploadService{
		uploadRepo:       uploadRepo,
		objectStore:      objectStore,
		tenantReader:     tenantReader,
		membershipReader: membershipReader,
		documentRepo:     documentRepo,
		kbRepo:           kbRepo,
		logger:           logger,
	}
}

func (dus *DocumentUploadService) Initiate(ctx context.Context, input kbmodel.InitiateInput) (kbmodel.InitiateResult, error) {
	if input.TenantID == uuid.Nil {
		return kbmodel.InitiateResult{}, ErrInvalidTenantID
	}

	if input.CreatedBy == uuid.Nil {
		return kbmodel.InitiateResult{}, ErrInvalidUserID
	}

	filename := safeFilename(input.Upload.OriginalFilename)
	if filename == "" {
		return kbmodel.InitiateResult{}, ErrInvalidOriginalFileName
	}
	if input.Upload.SizeBytes <= 0 {
		return kbmodel.InitiateResult{}, ErrInvalidSize
	}

	tenantDetail, err := dus.tenantReader.GetByID(ctx, input.TenantID)
	if err != nil {
		return kbmodel.InitiateResult{}, fmt.Errorf("failed to get tenant: %w", err)
	}

	if tenantDetail.Status == tenantmodel.StatusDisabled {
		return kbmodel.InitiateResult{}, ErrInActiveTenant
	}

	membershipDetail, err := dus.membershipReader.GetMembership(ctx, input.TenantID, input.CreatedBy)
	if err != nil {
		return kbmodel.InitiateResult{}, fmt.Errorf("failed to get membership: %w", err)
	}
	if membershipDetail.Status == tenantmodel.MembershipStatusDisabled {
		return kbmodel.InitiateResult{}, ErrUsersMembershipDisabled
	}
	var documentID uuid.UUID
	if input.DocumentID != nil {
		documentID = *input.DocumentID
		if documentID == uuid.Nil {
			return kbmodel.InitiateResult{}, ErrInvalidDocumentID
		}
		document, err := dus.documentRepo.GetByID(ctx, input.TenantID, documentID)
		if err != nil {
			return kbmodel.InitiateResult{}, fmt.Errorf("failed to get document: %w", err)
		}
		if document.Status == kbmodel.DocumentStatusArchived {
			return kbmodel.InitiateResult{}, ErrDocumentArchived
		}
	} else {
		if input.Document == nil {
			return kbmodel.InitiateResult{}, ErrInvalidDocumentMetadata
		}
		if input.Document.KnowledgeBaseID == uuid.Nil {
			return kbmodel.InitiateResult{}, ErrInvalidKBID
		}
		title := strings.TrimSpace(input.Document.Title)
		if title == "" {
			return kbmodel.InitiateResult{}, ErrInvalidTitle
		}
		documentType := strings.TrimSpace(input.Document.DocumentType)
		if documentType == "" {
			return kbmodel.InitiateResult{}, ErrInvalidDocumentType
		}

		kb, err := dus.kbRepo.GetKnowledgeBaseByID(ctx, input.Document.KnowledgeBaseID, input.TenantID)
		if err != nil {
			return kbmodel.InitiateResult{}, err
		}
		if kb.Status == kbmodel.KnowledgeBaseStatusArchived {
			return kbmodel.InitiateResult{}, ErrKBArchived
		}

		documentID = uuid.New()

		doc := kbmodel.Document{
			ID:           documentID,
			TenantID:     input.TenantID,
			KbID:         input.Document.KnowledgeBaseID,
			Title:        input.Document.Title,
			DocumentType: input.Document.DocumentType,
			CreatedBy:    input.CreatedBy,
		}
		err = dus.documentRepo.Create(ctx, doc)
		if err != nil {
			return kbmodel.InitiateResult{}, fmt.Errorf("failed to create document: %w", err)
		}
	}

	uploadID := uuid.New()

	objectKey := fmt.Sprintf("tenants/%s/documents/%s/upload/%s/%s", input.TenantID.String(), documentID.String(), uploadID.String(), filename)

	now := time.Now().UTC()

	u := kbmodel.DocumentUpload{
		ID:                uploadID,
		TenantID:          input.TenantID,
		DocumentID:        documentID,
		ObjectKey:         objectKey,
		OriginalFileName:  filename,
		ContentType:       input.Upload.ContentType,
		ExpectedSizeBytes: input.Upload.SizeBytes,
		Status:            kbmodel.DocumentUploadPending,
		CreatedBy:         input.CreatedBy,
		ExpiredAt:         now.Add(uploadExpiry),
		CreatedAt:         now,
	}
	if err := dus.uploadRepo.Create(ctx, u); err != nil {
		return kbmodel.InitiateResult{}, fmt.Errorf("upload persistance failed: %w", err)
	}

	presigned, err := dus.objectStore.PresignPut(ctx, objectKey, input.Upload.ContentType, uploadExpiry)
	if err != nil {
		return kbmodel.InitiateResult{}, fmt.Errorf("failed to get presigned url: %w", err)
	}
	return kbmodel.InitiateResult{DocumentID: documentID, UploadID: uploadID, URL: presigned.URL, ExpiresAt: presigned.ExpiresAt}, nil
}

func safeFilename(filename string) string {
	filename = strings.TrimSpace(filename)

	if filename == "" {
		return ""
	}

	filename = strings.ReplaceAll(filename, "\\", "/")
	filename = path.Base(filename)

	if filename == "." || filename == "/" {
		return ""
	}

	return filename
}
