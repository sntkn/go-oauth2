package user

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestRepository_FindByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuery := NewMockUserRepository(ctrl)

	testCases := []struct {
		name     string
		id       string
		mockFunc func()
		want     *User
		wantErr  bool
	}{
		{
			name: "正常なケース",
			id:   "550e8400-e29b-41d4-a716-446655440000",
			mockFunc: func() {
				mockQuery.EXPECT().FindByID("550e8400-e29b-41d4-a716-446655440000").Return(&User{
					ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
				}, nil)
			},
			want: &User{
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

			repo := mockQuery
			got, err := repo.FindByID(tc.id)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}
