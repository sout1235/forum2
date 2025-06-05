package repository

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sout1235/forum2/backend/auth-service/internal/entity"
	"github.com/sout1235/forum2/backend/auth-service/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	user := &entity.User{ID: "1", Username: "testuser", Email: "test@example.com"}
	mockRepo.EXPECT().Create(gomock.Any(), user).Return(nil)
	err := mockRepo.Create(context.Background(), user)
	assert.NoError(t, err)
}

func TestUserRepository_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	user := &entity.User{ID: "1", Username: "testuser"}
	mockRepo.EXPECT().GetByID(gomock.Any(), "1").Return(user, nil)
	result, err := mockRepo.GetByID(context.Background(), "1")
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestUserRepository_GetByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	user := &entity.User{ID: "1", Username: "testuser"}
	mockRepo.EXPECT().GetByUsername(gomock.Any(), "testuser").Return(user, nil)
	result, err := mockRepo.GetByUsername(context.Background(), "testuser")
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestUserRepository_ExistsByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().ExistsByUsername(gomock.Any(), "testuser").Return(true, nil)
	exists, err := mockRepo.ExistsByUsername(context.Background(), "testuser")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().ExistsByEmail(gomock.Any(), "test@example.com").Return(true, nil)
	exists, err := mockRepo.ExistsByEmail(context.Background(), "test@example.com")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestUserRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	user := &entity.User{ID: "1", Username: "updateduser"}
	mockRepo.EXPECT().Update(gomock.Any(), user).Return(nil)
	err := mockRepo.Update(context.Background(), user)
	assert.NoError(t, err)
}

func TestUserRepository_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().DeleteUser(gomock.Any(), "1").Return(nil)
	err := mockRepo.DeleteUser(context.Background(), "1")
	assert.NoError(t, err)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	user := &entity.User{ID: "1", Email: "test@example.com"}
	mockRepo.EXPECT().GetByEmail(gomock.Any(), "test@example.com").Return(user, nil)
	result, err := mockRepo.GetByEmail(context.Background(), "test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}
