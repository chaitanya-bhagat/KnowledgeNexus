package documentupload

import (
	"context"
	"testing"
	"time"

	kbMock "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/base/mocks"
	docMock "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/document/mocks"
	uploadMock "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/documentupload/mocks"
	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/storage"
	objectMock "github.com/chaitanya-bhagat/knowledge-nexus/internals/storage/mocks"
	membershipMock "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/membership/mocks"
	tenantmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
	tenantMock "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/tenant/mocks"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestServiceInitiate(t *testing.T) {
	ctrl := gomock.NewController(t)

	uploadRepo := uploadMock.NewMockRepository(ctrl)
	tenantRepo := tenantMock.NewMockReader(ctrl)
	membershipRepo := membershipMock.NewMockReader(ctrl)
	documentRepo := docMock.NewMockRepository(ctrl)
	objectStore := objectMock.NewMockObjectStore(ctrl)
	kbRepo := kbMock.NewMockReader(ctrl)

	tenantID := uuid.New()
	userID := uuid.New()
	kbID := uuid.New()
	logger := zap.NewNop()

	tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(&tenantmodel.Tenant{
		ID:     tenantID,
		Status: tenantmodel.StatusActive,
	}, nil)

	membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(tenantmodel.Membership{
		TenantID: tenantID,
		UserID:   userID,
		Status:   tenantmodel.MembershipStatusActive,
	}, nil)

	kbRepo.EXPECT().GetKnowledgeBaseByID(gomock.Any(), kbID, tenantID).Return(kbmodel.KnowledgeBase{
		ID:        kbID,
		TenantID:  tenantID,
		CreatedBy: userID,
		Status:    kbmodel.KnowledgeBaseStatusActive,
	}, nil)

	// documentRepo.EXPECT().GetByID(gomock.Any(), tenantID, documentID).Return(kbmodel.Document{}, errors.New("document not found"))

	var createdDocument kbmodel.Document
	documentRepo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(kbmodel.Document{})).
		DoAndReturn(func(_ context.Context, doc kbmodel.Document) error {
			createdDocument = doc
			if doc.ID == uuid.Nil {
				t.Error("document ID must not be nil")
			}
			if doc.TenantID != tenantID {
				t.Errorf("document TenantID = %v, want %v", doc.TenantID, tenantID)
			}

			if doc.KbID != kbID {
				t.Errorf("document KnowledgeBaseID = %v, want %v", doc.KbID, kbID)
			}
			return nil
		})
	var createdUpload kbmodel.DocumentUpload

	uploadRepo.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(kbmodel.DocumentUpload{})).
		DoAndReturn(func(_ context.Context, u kbmodel.DocumentUpload) error {
			createdUpload = u
			if u.ID == uuid.Nil {
				t.Error("upload ID must not be nil")
			}

			if u.DocumentID != createdDocument.ID {
				t.Errorf("upload DocumentID = %v, want %v", u.DocumentID, createdDocument.ID)
			}

			if u.TenantID != tenantID {
				t.Errorf("upload TenantID = %v, want %v", u.TenantID, tenantID)
			}

			if u.CreatedBy != userID {
				t.Errorf("upload CreatedBy = %v, want %v", u.CreatedBy, userID)
			}

			if u.OriginalFileName != "judgment.pdf" {
				t.Errorf("filename = %q, want judgment.pdf", u.OriginalFileName)
			}

			if u.ContentType != "application/pdf" {
				t.Errorf("content type = %q, want application/pdf", u.ContentType)
			}

			if u.ExpectedSizeBytes != 1024 {
				t.Errorf("size = %d, want 1024", u.ExpectedSizeBytes)
			}

			if u.Status != kbmodel.DocumentUploadPending {
				t.Errorf("status = %v, want %v", u.Status, kbmodel.DocumentUploadPending)
			}

			if u.ObjectKey == "" {
				t.Error("object key must not be empty")
			}

			return nil
		})

	expiresAt := time.Now().UTC().Add(15 * time.Minute)

	objectStore.EXPECT().PresignPut(gomock.Any(), gomock.Any(), "application/pdf", gomock.Any()).
		DoAndReturn(func(_ context.Context, key string, _ string, _ time.Duration) (storage.PresignedUpload, error) {
			if key != createdUpload.ObjectKey {
				t.Errorf("presigned key = %q, persisted key = %q", key, createdUpload.ObjectKey)
			}

			return storage.PresignedUpload{
				URL:       "http://presigned.test/upload",
				ExpiresAt: expiresAt,
			}, nil
		})

	service := NewDocumentUploadService(uploadRepo, objectStore, tenantRepo, membershipRepo, documentRepo, kbRepo, logger)
	result, err := service.Initiate(context.Background(), kbmodel.InitiateInput{
		TenantID:   tenantID,
		DocumentID: nil,
		CreatedBy:  userID,

		Upload: kbmodel.UploadInput{
			OriginalFilename: "judgment.pdf",
			ContentType:      "application/pdf",
			SizeBytes:        1024,
		},

		Document: &kbmodel.DocumentInput{
			KnowledgeBaseID: kbID,
			Title:           "Kesavananda Bharati",
			DocumentType:    "judgment",
		},
	},
	)

	if err != nil {
		t.Fatalf("Initiate() error = %v", err)
	}

	if result.DocumentID != createdDocument.ID {
		t.Errorf("result DocumentID = %v, want %v", result.DocumentID, createdDocument.ID)
	}

	if createdUpload.DocumentID != result.DocumentID {
		t.Errorf("upload DocumentID = %v, result DocumentID = %v", createdUpload.DocumentID, result.DocumentID)
	}

	if result.UploadID != createdUpload.ID {
		t.Errorf("result UploadID = %v, want %v", result.UploadID, createdUpload.ID)
	}

	if result.URL != "http://presigned.test/upload" {
		t.Errorf("result URL = %q, want %q", result.URL, "http://presigned.test/upload")
	}

	if !result.ExpiresAt.Equal(expiresAt) {
		t.Errorf("result ExpiresAt = %v, want %v", result.ExpiresAt, expiresAt)
	}

}
