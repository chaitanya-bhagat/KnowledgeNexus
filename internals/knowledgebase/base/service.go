package knowledgebase

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/identity"
	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	tenantmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/tenant"
)

type KnowledgeBaseService struct {
	repo       Repository
	tenantRepo tenant.Repository
	userRepo   identity.Repository
}

func NewKnowledgeBaseService(repo Repository, tenantRepo tenant.Repository, userRepo identity.Repository) *KnowledgeBaseService {
	return &KnowledgeBaseService{
		repo:       repo,
		tenantRepo: tenantRepo,
		userRepo:   userRepo,
	}
}

func (kbs *KnowledgeBaseService) CreateKnowledgeBase(ctx context.Context, kb kbmodel.KnowledgeBase) (kbmodel.KnowledgeBase, error) {

	if kb.TenantID == uuid.Nil {
		return kbmodel.KnowledgeBase{}, ErrInvalidTenantID
	}

	if kb.CreatedBy == uuid.Nil {
		return kbmodel.KnowledgeBase{}, ErrInvalidUserID
	}

	name := strings.TrimSpace(kb.Name)
	if name == "" {
		return kbmodel.KnowledgeBase{}, ErrInvalidKnowledgeBaseName
	}

	domainType := strings.ToLower(strings.TrimSpace(kb.DomainType))
	if domainType == "" {
		return kbmodel.KnowledgeBase{}, ErrInvalidDomainType
	}

	description := strings.TrimSpace(kb.Description)

	tenantDetails, err := kbs.tenantRepo.GetByID(ctx, kb.TenantID)
	if err != nil {
		return kbmodel.KnowledgeBase{}, err
	}
	if tenantDetails.Status == tenantmodel.StatusDisabled {
		return kbmodel.KnowledgeBase{}, ErrTenantDisabled
	}
	userDetails, err := kbs.userRepo.GetByID(ctx, kb.CreatedBy)
	if err != nil {
		return kbmodel.KnowledgeBase{}, err
	}
	if userDetails.Status == identity.StatusDisabled {
		return kbmodel.KnowledgeBase{}, ErrUserDisabled
	}
	now := time.Now().UTC()

	kb = kbmodel.KnowledgeBase{
		ID:          uuid.New(),
		Name:        name,
		TenantID:    kb.TenantID,
		CreatedBy:   kb.CreatedBy,
		Status:      kbmodel.KnowledgeBaseStatusActive,
		DomainType:  domainType,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := kbs.repo.CreateKnowledgeBase(ctx, kb); err != nil {
		return kbmodel.KnowledgeBase{}, err
	}
	return kb, nil
}

func (kbs *KnowledgeBaseService) GetKnowledgeBaseByID(ctx context.Context, kbID uuid.UUID, tenantID uuid.UUID) (kbmodel.KnowledgeBase, error) {
	if kbID == uuid.Nil {
		return kbmodel.KnowledgeBase{}, ErrInvalidKnowledgeBaseID
	}
	if tenantID == uuid.Nil {
		return kbmodel.KnowledgeBase{}, ErrInvalidTenantID
	}
	_, err := kbs.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return kbmodel.KnowledgeBase{}, err
	}

	return kbs.repo.GetKnowledgeBaseByID(ctx, kbID, tenantID)
}

func (kbs *KnowledgeBaseService) ListKnowledgeBasesByTenantID(ctx context.Context, tenantID uuid.UUID) ([]kbmodel.KnowledgeBase, error) {
	if tenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}
	_, err := kbs.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	return kbs.repo.ListKnowledgeBasesByTenantID(ctx, tenantID)
}

func (kbs *KnowledgeBaseService) UpdateKnowledgeBase(ctx context.Context, kbID uuid.UUID, tenantID uuid.UUID, kb *kbmodel.UpdateKnowledgeBase) (kbmodel.KnowledgeBase, error) {
	if kbID == uuid.Nil {
		return kbmodel.KnowledgeBase{}, ErrInvalidKnowledgeBaseID
	}
	if tenantID == uuid.Nil {
		return kbmodel.KnowledgeBase{}, ErrInvalidTenantID
	}

	name := strings.TrimSpace(kb.Name)
	if name == "" {
		return kbmodel.KnowledgeBase{}, ErrInvalidKnowledgeBaseName
	}
	kbDetails, err := kbs.repo.GetKnowledgeBaseByID(ctx, kbID, tenantID)
	if err != nil {
		return kbmodel.KnowledgeBase{}, err
	}

	kbDetails.Name = name
	kbDetails.Description = strings.TrimSpace(kb.Description)
	kbDetails.UpdatedAt = time.Now().UTC()

	err = kbs.repo.UpdateKnowledgeBase(ctx, &kbDetails)
	if err != nil {
		return kbmodel.KnowledgeBase{}, err
	}
	return kbDetails, nil
}

func (kbs *KnowledgeBaseService) Archive(ctx context.Context, tenantID uuid.UUID, kbID uuid.UUID) (kbmodel.KnowledgeBase, error) {
	if tenantID == uuid.Nil {
		return kbmodel.KnowledgeBase{}, ErrInvalidTenantID
	}
	if kbID == uuid.Nil {
		return kbmodel.KnowledgeBase{}, ErrInvalidKnowledgeBaseID
	}
	kbDetails, err := kbs.repo.GetKnowledgeBaseByID(ctx, kbID, tenantID)
	if err != nil {
		return kbmodel.KnowledgeBase{}, err
	}
	if kbDetails.Status == kbmodel.KnowledgeBaseStatusArchived {
		return kbDetails, nil
	}

	now := time.Now().UTC()
	err = kbs.repo.UpdateKnowledgeBaseStatus(ctx, tenantID, kbID, kbmodel.KnowledgeBaseStatusArchived, now)
	if err != nil {
		return kbmodel.KnowledgeBase{}, err
	}

	kbDetails.Status = kbmodel.KnowledgeBaseStatusArchived
	kbDetails.UpdatedAt = now

	return kbDetails, nil
}

func (kbs *KnowledgeBaseService) Activate(ctx context.Context, tenantID uuid.UUID, kbID uuid.UUID) (kbmodel.KnowledgeBase, error) {
	if tenantID == uuid.Nil {
		return kbmodel.KnowledgeBase{}, ErrInvalidTenantID
	}
	if kbID == uuid.Nil {
		return kbmodel.KnowledgeBase{}, ErrInvalidKnowledgeBaseID
	}

	tenantDetails, err := kbs.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {

		return kbmodel.KnowledgeBase{}, err
	}
	if tenantDetails.Status == tenantmodel.StatusDisabled {
		return kbmodel.KnowledgeBase{}, ErrTenantDisabled
	}

	kbDetails, err := kbs.repo.GetKnowledgeBaseByID(ctx, kbID, tenantID)
	if err != nil {
		return kbmodel.KnowledgeBase{}, err
	}

	if kbDetails.Status == kbmodel.KnowledgeBaseStatusActive {
		return kbDetails, nil
	}

	now := time.Now().UTC()

	err = kbs.repo.UpdateKnowledgeBaseStatus(ctx, tenantID, kbID, kbmodel.KnowledgeBaseStatusActive, now)
	if err != nil {
		return kbmodel.KnowledgeBase{}, err
	}
	kbDetails.Status = kbmodel.KnowledgeBaseStatusActive
	kbDetails.UpdatedAt = now

	return kbDetails, nil
}
