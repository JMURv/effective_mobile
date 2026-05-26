package ctrl

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/JMURv/effective-mobile/internal/dto"
	md "github.com/JMURv/effective-mobile/internal/models"
	"github.com/JMURv/effective-mobile/internal/repo"
	"github.com/JMURv/effective-mobile/tests/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestController_List(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := mocks.NewMockAppRepo(ctrlMock)
	mockCache := mocks.NewMockCacheService(ctrlMock)

	c := New(mockRepo, mockCache)
	ctx := context.Background()

	expected := []md.Subscription{
		{ID: uuid.New(), ServiceName: "netflix"},
	}

	tests := []struct {
		name    string
		setup   func()
		want    []md.Subscription
		wantErr bool
	}{
		{
			name: "CacheHit",
			setup: func() {
				mockCache.EXPECT().
					GetToStruct(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, key string, dst any) error {
						ptr := dst.(*[]md.Subscription)
						*ptr = expected
						return nil
					})
			},
			want: expected,
		},
		{
			name: "CacheMiss_RepoSuccess",
			setup: func() {
				mockCache.EXPECT().
					GetToStruct(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("miss"))

				mockRepo.EXPECT().
					List(gomock.Any()).
					Return(expected, nil)

				mockCache.EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), expected).
					Return().
					AnyTimes()
			},
			want: expected,
		},
		{
			name: "RepoError",
			setup: func() {
				mockCache.EXPECT().
					GetToStruct(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("miss"))

				mockRepo.EXPECT().
					List(gomock.Any()).
					Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			res, err := c.List(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestController_GetByID(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := mocks.NewMockAppRepo(ctrlMock)
	mockCache := mocks.NewMockCacheService(ctrlMock)

	c := New(mockRepo, mockCache)
	ctx := context.Background()

	id := uuid.New()

	sub := md.Subscription{
		ID:          id,
		ServiceName: "spotify",
	}

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "CacheHit",
			setup: func() {
				mockCache.EXPECT().
					GetToStruct(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, key string, dst any) error {
						ptr := dst.(*md.Subscription)
						*ptr = sub
						return nil
					})
			},
		},
		{
			name: "CacheMiss_RepoFound",
			setup: func() {
				mockCache.EXPECT().
					GetToStruct(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("miss"))

				mockRepo.EXPECT().
					GetByID(gomock.Any(), id).
					Return(sub, nil)

				mockCache.EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), sub).
					Return().
					AnyTimes()
			},
		},
		{
			name: "NotFound",
			setup: func() {
				mockCache.EXPECT().
					GetToStruct(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("miss"))

				mockRepo.EXPECT().
					GetByID(gomock.Any(), id).
					Return(md.Subscription{}, repo.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "RepoError",
			setup: func() {
				mockCache.EXPECT().
					GetToStruct(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("miss"))

				mockRepo.EXPECT().
					GetByID(gomock.Any(), id).
					Return(md.Subscription{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			_, err := c.GetByID(ctx, id)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestController_Create(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := mocks.NewMockAppRepo(ctrlMock)
	mockCache := mocks.NewMockCacheService(ctrlMock)

	c := New(mockRepo, mockCache)
	ctx := context.Background()

	req := dto.CreateSubscriptionRequest{
		ServiceName: "netflix",
	}

	sub := md.Subscription{
		ID: uuid.New(),
	}

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				mockRepo.EXPECT().
					Create(gomock.Any(), req).
					Return(sub, nil)

				mockCache.EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), sub).
					Return().AnyTimes()

				mockCache.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return()
			},
		},
		{
			name: "RepoError",
			setup: func() {
				mockRepo.EXPECT().
					Create(gomock.Any(), req).
					Return(md.Subscription{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := c.Create(ctx, req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestController_Update(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := mocks.NewMockAppRepo(ctrlMock)
	mockCache := mocks.NewMockCacheService(ctrlMock)

	c := New(mockRepo, mockCache)
	ctx := context.Background()

	id := uuid.New()

	req := dto.UpdateSubscriptionRequest{
		ServiceName: new("netflix"),
	}

	sub := md.Subscription{ID: id}

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				mockRepo.EXPECT().
					Update(gomock.Any(), id, req).
					Return(sub, nil)

				mockCache.EXPECT().
					Set(gomock.Any(), gomock.Any(), gomock.Any(), sub).
					Return().AnyTimes()

				mockCache.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return().AnyTimes()
			},
		},
		{
			name: "NotFound",
			setup: func() {
				mockRepo.EXPECT().
					Update(gomock.Any(), id, req).
					Return(md.Subscription{}, repo.ErrNotFound)
			},
			wantErr: true,
		},
		{
			name: "RepoError",
			setup: func() {
				mockRepo.EXPECT().
					Update(gomock.Any(), id, req).
					Return(md.Subscription{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := c.Update(ctx, id, req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestController_Delete(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := mocks.NewMockAppRepo(ctrlMock)
	mockCache := mocks.NewMockCacheService(ctrlMock)

	c := New(mockRepo, mockCache)
	ctx := context.Background()

	id := uuid.New()

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				mockRepo.EXPECT().
					Delete(gomock.Any(), id).
					Return(nil)

				mockCache.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return().AnyTimes()
			},
		},
		{
			name: "RepoError",
			setup: func() {
				mockRepo.EXPECT().
					Delete(gomock.Any(), id).
					Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := c.Delete(ctx, id)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestController_CalculateTotalCost(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := mocks.NewMockAppRepo(ctrlMock)
	mockCache := mocks.NewMockCacheService(ctrlMock)

	c := New(mockRepo, mockCache)
	ctx := context.Background()

	req := dto.CalculateTotalCostRequest{
		ServiceName: "netflix",
		UserID:      uuid.New(),
		From:        time.Now(),
		To:          time.Now().Add(1 * time.Hour),
	}

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				mockRepo.EXPECT().
					CalculateTotalCost(gomock.Any(), req).
					Return(int64(100), nil)
			},
		},
		{
			name: "RepoError",
			setup: func() {
				mockRepo.EXPECT().
					CalculateTotalCost(gomock.Any(), req).
					Return(int64(0), errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			res, err := c.CalculateTotalCost(ctx, req)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, int64(100), res.TotalCost)
		})
	}
}
