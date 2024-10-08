package user

import (
	"testing"

	"github.com/go-errors/errors"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/api/internal/infrastructure/db/model"
	"github.com/sntkn/go-oauth2/api/internal/modules/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestService_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuery := domain.NewMockRepository(ctrl)

	testCases := []struct {
		name     string
		id       string
		mockFunc func()
		want     *domain.User
		wantErr  bool
	}{
		{
			name: "正常なケース",
			id:   "550e8400-e29b-41d4-a716-446655440000",
			mockFunc: func() {
				mockQuery.EXPECT().FindByID("550e8400-e29b-41d4-a716-446655440000").Return(&model.User{
					ID: "550e8400-e29b-41d4-a716-446655440000",
				}, nil)
			},
			want: &domain.User{
				ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			},
			wantErr: false,
		},
		{
			name: "ユーザーが見つからない場合",
			id:   "non-existent-id",
			mockFunc: func() {
				mockQuery.EXPECT().FindByID("non-existent-id").Return(nil, gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "無効なUUID",
			id:   "invalid-uuid",
			mockFunc: func() {
				mockQuery.EXPECT().FindByID("invalid-uuid").Return(nil, errors.New("invalid-uuid"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			service := NewUsecase(mockQuery)

			got, err := service.FindUser(tc.id)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}
