package documentversion

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/identity"
	knowledgebase "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/base"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/document"
	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant"
	tenantmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
	"go.uber.org/zap"
)

type DocVersionService struct {
	docVersionRepo Repository
	docRepo        document.Repository
	kbRepo         knowledgebase.Repository
	tenantRepo     tenant.Repository
	userRepo       identity.Repository
	logger         *zap.Logger
}

func NewDocumentVersionService(docVersionRepo Repository, docRepo document.Repository, kbRepo knowledgebase.Repository, tenantRepo tenant.Repository, userRepo identity.Repository, logger *zap.Logger) *DocVersionService {
	return &DocVersionService{
		docVersionRepo: docVersionRepo,
		docRepo:        docRepo,
		kbRepo:         kbRepo,
		tenantRepo:     tenantRepo,
		userRepo:       userRepo,
		logger:         logger,
	}
}

func (dvs *DocVersionService) Create(ctx context.Context, input kbmodel.DocumentVersionInput) (kbmodel.DocumentVersion, error) {
	if input.TenantID == uuid.Nil {
		return kbmodel.DocumentVersion{}, ErrInvalidTenantID
	}
	if input.DocumentID == uuid.Nil {
		return kbmodel.DocumentVersion{}, ErrInvalidDocumentID
	}
	if input.CreatedBy == uuid.Nil {
		return kbmodel.DocumentVersion{}, ErrInvalidUserID
	}

	objectKey := strings.TrimSpace(input.ObjectKey)
	if objectKey == "" {
		return kbmodel.DocumentVersion{}, ErrInvalidObjectKey
	}

	fileName := strings.TrimSpace(input.OriginalFileName)
	if fileName == "" {
		return kbmodel.DocumentVersion{}, ErrInvalidOriginalFileName
	}

	if input.SizeBytes <= 0 {
		return kbmodel.DocumentVersion{}, ErrInvalidSize
	}
	// Check Tenant is active

	tenantDetail, err := dvs.tenantRepo.GetByID(ctx, input.TenantID)
	if err != nil {
		if errors.Is(err, tenant.ErrNotFound) {
			return kbmodel.DocumentVersion{}, ErrTenantNotFound
		}
		return kbmodel.DocumentVersion{}, err
	}
	if tenantDetail.Status == tenantmodel.StatusDisabled {
		return kbmodel.DocumentVersion{}, ErrTenantDisabled
	}

	// Check document
	docDetail, err := dvs.docRepo.GetByID(ctx, input.TenantID, input.DocumentID)
	if err != nil {
		if errors.Is(err, document.ErrDocumentNotFound) {
			return kbmodel.DocumentVersion{}, ErrDocumentNotFound
		}
		return kbmodel.DocumentVersion{}, err
	}
	if docDetail.Status == kbmodel.DocumentStatusArchived {
		return kbmodel.DocumentVersion{}, ErrDocumentArchived
	}

	kbDetail, err := dvs.kbRepo.GetKnowledgeBaseByID(ctx, docDetail.KbID, input.TenantID)
	if err != nil {
		if errors.Is(err, document.ErrKnowledgeBaseNotFound) {
			return kbmodel.DocumentVersion{}, ErrKnowledgeBaseNotFound
		}
		return kbmodel.DocumentVersion{}, err
	}
	if kbDetail.Status == kbmodel.KnowledgeBaseStatusArchived {
		return kbmodel.DocumentVersion{}, ErrKnowledgeBaseArchived
	}

	now := time.Now().UTC()

	version := kbmodel.DocumentVersion{
		ID:               uuid.New(),
		TenantID:         input.TenantID,
		DocumentID:       input.DocumentID,
		VersionNumber:    0,
		ObjectKey:        objectKey,
		OriginalFileName: fileName,
		ContentType:      strings.TrimSpace(input.ContentType),
		SizeBytes:        input.SizeBytes,
		Checksum:         strings.TrimSpace(input.Checksum),
		Status:           kbmodel.DocumentVersionStatusUploaded,
		FailureReason:    "",
		CreatedBy:        input.CreatedBy,
		CreatedAt:        now,
	}
	err = dvs.docVersionRepo.Create(ctx, version)
	if err != nil {
		return kbmodel.DocumentVersion{}, err
	}
	return version, nil
}

func (dvs *DocVersionService) GetByID(ctx context.Context, tenantID uuid.UUID, versionID uuid.UUID) (kbmodel.DocumentVersion, error) {
	if tenantID == uuid.Nil {
		return kbmodel.DocumentVersion{}, ErrInvalidTenantID
	}
	if versionID == uuid.Nil {
		return kbmodel.DocumentVersion{}, ErrInvalidVersionID
	}
	return dvs.docVersionRepo.GetByID(ctx, tenantID, versionID)
}

func (dvs *DocVersionService) GetListByDocumentID(ctx context.Context, tenantID uuid.UUID, docID uuid.UUID) ([]kbmodel.DocumentVersion, error) {
	if tenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if docID == uuid.Nil {
		return nil, ErrInvalidDocumentID
	}
	docDetail, err := dvs.docRepo.GetByID(ctx, tenantID, docID)
	if err != nil {
		if errors.Is(err, document.ErrDocumentNotFound) {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}
	if docDetail.Status == kbmodel.DocumentStatusArchived {
		return nil, ErrDocumentArchived
	}
	return dvs.docVersionRepo.GetList(ctx, tenantID, docID)
}

func (dvs *DocVersionService) MarkProcessing(ctx context.Context, tenantID uuid.UUID, versionID uuid.UUID) error {
	if tenantID == uuid.Nil {
		return ErrInvalidTenantID
	}
	if versionID == uuid.Nil {
		return ErrInvalidVersionID
	}
	versionDetail, err := dvs.docVersionRepo.GetByID(ctx, tenantID, versionID)
	if err != nil {
		return err
	}
	if versionDetail.Status == kbmodel.DocumentVersionStatusProcessing {
		return nil
	}
	if versionDetail.Status != kbmodel.DocumentVersionStatusUploaded {
		return ErrInvalidStatus
	}
	return dvs.docVersionRepo.Update(ctx, tenantID, versionID, kbmodel.DocumentVersionStatusProcessing, "")
}

func (dvs *DocVersionService) MarkReady(ctx context.Context, tenantID uuid.UUID, versionID uuid.UUID) error {
	if tenantID == uuid.Nil {
		return ErrInvalidTenantID
	}
	if versionID == uuid.Nil {
		return ErrInvalidVersionID
	}
	versionDetail, err := dvs.docVersionRepo.GetByID(ctx, tenantID, versionID)
	if err != nil {
		return err
	}
	if versionDetail.Status == kbmodel.DocumentVersionStatusReady {
		return nil
	}
	if versionDetail.Status != kbmodel.DocumentVersionStatusProcessing {
		return ErrInvalidStatus
	}
	return dvs.docVersionRepo.Update(ctx, tenantID, versionID, kbmodel.DocumentVersionStatusReady, "")
}

func (dvs *DocVersionService) MarkFailed(ctx context.Context, tenantID uuid.UUID, versionID uuid.UUID, failureReason string) error {
	if tenantID == uuid.Nil {
		return ErrInvalidTenantID
	}
	if versionID == uuid.Nil {
		return ErrInvalidVersionID
	}
	versionDetail, err := dvs.docVersionRepo.GetByID(ctx, tenantID, versionID)
	if err != nil {
		return err
	}
	if versionDetail.Status == kbmodel.DocumentVersionStatusReady {
		return ErrInvalidStatus
	}
	reason := strings.TrimSpace(failureReason)
	return dvs.docVersionRepo.Update(ctx, tenantID, versionID, kbmodel.DocumentVersionStatusReady, reason)
}
