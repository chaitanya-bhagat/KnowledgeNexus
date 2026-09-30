package document

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/identity"
	knowledgebase "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/base"
	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	tenantmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/tenant"
	"go.uber.org/zap"
)

type DocumentService struct {
	documentRepo Repository
	tenanrRepo   tenant.Repository
	kbRepo       knowledgebase.Repository
	userRepo     identity.Repository
	logger       *zap.Logger
}

func NewDocumentService(documentRepo Repository, tenanrRepo tenant.Repository, kbRepo knowledgebase.Repository, userRepo identity.Repository, logger *zap.Logger) *DocumentService {
	return &DocumentService{
		documentRepo: documentRepo,
		tenanrRepo:   tenanrRepo,
		kbRepo:       kbRepo,
		userRepo:     userRepo,
		logger:       logger,
	}
}

func (ds *DocumentService) Create(ctx context.Context, input kbmodel.CreateDocument) (kbmodel.Document, error) {
	if input.TenantID == uuid.Nil {
		return kbmodel.Document{}, ErrInvalidTenantID
	}
	if input.KnowledgeBaseID == uuid.Nil {
		return kbmodel.Document{}, ErrInvalidKnowledgeBaseID
	}
	if input.CreatedBy == uuid.Nil {
		return kbmodel.Document{}, ErrInvalidUserID
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		return kbmodel.Document{}, ErrInvalidTitle
	}

	documentType := strings.ToLower(strings.TrimSpace(input.DocumentType))

	tenantDetails, err := ds.tenanrRepo.GetByID(ctx, input.TenantID)
	if err != nil {
		if errors.Is(err, tenant.ErrNotFound) {
			return kbmodel.Document{}, ErrTenantNotFound
		}
		return kbmodel.Document{}, err
	}

	if tenantDetails.Status == tenantmodel.StatusDisabled {
		return kbmodel.Document{}, ErrTenantDisabled
	}

	kbDetails, err := ds.kbRepo.GetKnowledgeBaseByID(ctx, input.KnowledgeBaseID, input.TenantID)
	if err != nil {
		if errors.Is(err, knowledgebase.ErrKnowledgeBaseNotFound) {
			return kbmodel.Document{}, ErrKnowledgeBaseNotFound
		}
		return kbmodel.Document{}, err
	}

	if kbDetails.Status == kbmodel.KnowledgeBaseStatusArchived {
		return kbmodel.Document{}, ErrKnowledgeBaseArchived
	}

	userDetails, err := ds.userRepo.GetByID(ctx, input.CreatedBy)
	if err != nil {
		if errors.Is(err, identity.ErrNotFound) {
			return kbmodel.Document{}, ErrUserNotFound
		}
		return kbmodel.Document{}, err
	}
	if userDetails.Status == identity.StatusDisabled {
		return kbmodel.Document{}, ErrUserNotFound
	}

	now := time.Now().UTC()

	doc := kbmodel.Document{
		ID:           uuid.New(),
		Title:        title,
		TenantID:     input.TenantID,
		CreatedBy:    input.CreatedBy,
		KbID:         input.KnowledgeBaseID,
		DocumentType: documentType,
		Status:       kbmodel.DocumentStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := ds.documentRepo.Create(ctx, doc); err != nil {
		return kbmodel.Document{}, err
	}

	return doc, nil
}

func (ds *DocumentService) GetByID(ctx context.Context, tenantID uuid.UUID, docID uuid.UUID) (kbmodel.Document, error) {
	if tenantID == uuid.Nil {
		return kbmodel.Document{}, ErrInvalidTenantID
	}
	if docID == uuid.Nil {
		return kbmodel.Document{}, ErrInvalidDocumentID
	}
	return ds.documentRepo.GetByID(ctx, tenantID, docID)
}

func (ds *DocumentService) ListByKnowledgeBase(ctx context.Context, kbID uuid.UUID, tenantID uuid.UUID) ([]kbmodel.Document, error) {
	if tenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	if kbID == uuid.Nil {
		return nil, ErrInvalidDocumentID
	}
	_, err := ds.kbRepo.GetKnowledgeBaseByID(ctx, kbID, tenantID)
	if err != nil {
		if errors.Is(err, knowledgebase.ErrKnowledgeBaseNotFound) {
			return nil, ErrKnowledgeBaseNotFound
		}
		return nil, err
	}
	return ds.documentRepo.GetList(ctx, tenantID, kbID)
}

func (ds *DocumentService) Update(ctx context.Context, tenantID uuid.UUID, docID uuid.UUID, input kbmodel.UpdateDocumentInput) (kbmodel.Document, error) {
	if tenantID == uuid.Nil {
		return kbmodel.Document{}, ErrInvalidTenantID
	}
	if docID == uuid.Nil {
		return kbmodel.Document{}, ErrInvalidDocumentID
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return kbmodel.Document{}, ErrInvalidTitle
	}

	tenantDetails, err := ds.tenanrRepo.GetByID(ctx, tenantID)
	if err != nil {
		if errors.Is(err, tenant.ErrNotFound) {
			return kbmodel.Document{}, ErrTenantNotFound
		}
		return kbmodel.Document{}, err
	}

	if tenantDetails.Status == tenantmodel.StatusDisabled {
		return kbmodel.Document{}, ErrTenantDisabled
	}

	doc, err := ds.documentRepo.GetByID(ctx, tenantID, docID)
	if err != nil {
		if errors.Is(err, ErrDocumentNotFound) {
			return kbmodel.Document{}, ErrDocumentNotFound
		}
		return kbmodel.Document{}, err
	}

	kbDetails, err := ds.kbRepo.GetKnowledgeBaseByID(ctx, doc.KbID, tenantID)
	if err != nil {
		if errors.Is(err, knowledgebase.ErrKnowledgeBaseNotFound) {
			return kbmodel.Document{}, ErrKnowledgeBaseNotFound
		}
		return kbmodel.Document{}, err
	}

	if kbDetails.Status == kbmodel.KnowledgeBaseStatusArchived {
		return kbmodel.Document{}, ErrKnowledgeBaseArchived
	}

	doc.Title = title
	doc.UpdatedAt = time.Now()
	doc.DocumentType = strings.ToLower(strings.TrimSpace(input.DocumentType))

	if err := ds.documentRepo.Update(ctx, doc); err != nil {
		return kbmodel.Document{}, err
	}
	return doc, nil
}

func (ds *DocumentService) Archive(ctx context.Context, tenantID uuid.UUID, docID uuid.UUID) (kbmodel.Document, error) {
	if tenantID == uuid.Nil {
		return kbmodel.Document{}, ErrInvalidTenantID
	}
	if docID == uuid.Nil {
		return kbmodel.Document{}, ErrInvalidDocumentID
	}

	doc, err := ds.documentRepo.GetByID(ctx, tenantID, docID)
	if err != nil {
		return kbmodel.Document{}, err
	}

	if doc.Status == kbmodel.DocumentStatusArchived {
		return doc, nil
	}

	now := time.Now().UTC()

	if err := ds.documentRepo.UpdateStatus(ctx, tenantID, docID, kbmodel.DocumentStatusArchived, now); err != nil {
		return kbmodel.Document{}, err
	}

	doc.Status = kbmodel.DocumentStatusArchived
	doc.UpdatedAt = now
	return doc, nil
}

func (ds *DocumentService) Activate(ctx context.Context, tenantID uuid.UUID, docID uuid.UUID) (kbmodel.Document, error) {
	if tenantID == uuid.Nil {
		return kbmodel.Document{}, ErrInvalidTenantID
	}
	if docID == uuid.Nil {
		return kbmodel.Document{}, ErrInvalidDocumentID
	}

	tenantDetails, err := ds.tenanrRepo.GetByID(ctx, tenantID)
	if err != nil {
		if errors.Is(err, tenant.ErrNotFound) {
			return kbmodel.Document{}, ErrTenantNotFound
		}
		return kbmodel.Document{}, err
	}

	if tenantDetails.Status == tenantmodel.StatusDisabled {
		return kbmodel.Document{}, ErrTenantDisabled
	}

	doc, err := ds.documentRepo.GetByID(ctx, tenantID, docID)
	if err != nil {
		return kbmodel.Document{}, err
	}

	kbDetails, err := ds.kbRepo.GetKnowledgeBaseByID(ctx, doc.KbID, tenantID)
	if err != nil {
		if errors.Is(err, knowledgebase.ErrKnowledgeBaseNotFound) {
			return kbmodel.Document{}, ErrKnowledgeBaseNotFound
		}
		return kbmodel.Document{}, err
	}

	if kbDetails.Status == kbmodel.KnowledgeBaseStatusArchived {
		return kbmodel.Document{}, ErrKnowledgeBaseArchived
	}

	if doc.Status == kbmodel.DocumentStatusActive {
		return doc, nil
	}

	now := time.Now().UTC()

	if err := ds.documentRepo.UpdateStatus(ctx, tenantID, docID, kbmodel.DocumentStatusActive, now); err != nil {
		return kbmodel.Document{}, err
	}

	doc.Status = kbmodel.DocumentStatusActive
	doc.UpdatedAt = now
	return doc, nil
}
