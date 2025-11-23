package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Autumn-27/ScopeSentry/internal/utils/helper"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Autumn-27/ScopeSentry/internal/config"
	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/repositories/user"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUsernameExists  = errors.New("username already exists")
	ErrInvalidUserData = errors.New("invalid user data")
)

type Service interface {
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, user *models.User) error
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, user *models.User) error
	DeleteUser(ctx context.Context, id string) error
	ChangePassword(ctx context.Context, id string, newPassword string) error
}

type service struct {
	userRepo user.Repository
}

func NewService() Service {
	return &service{
		userRepo: user.NewRepository(),
	}
}

func (s *service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrUserNotFound
	}

	hashPwd := helper.Sha256Hex(password)
	if hashPwd != user.Password {
		return "", ErrInvalidPassword
	}

	// 生成JWT token
	now := time.Now()
	expire := now.Add(config.GlobalConfig.JWT.Expire)

	claims := jwt.MapClaims{
		"userID":   user.ID.Hex(),
		"username": user.Username,
		"role":     user.Role,
		"exp":      expire.Unix(),
		"iat":      now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return tokenString, nil
}

func (s *service) Register(ctx context.Context, user *models.User) error {
	// 验证用户数据
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return ErrInvalidUserData
	}

	// 检查用户名是否已存在
	existingUser, err := s.userRepo.FindByUsername(ctx, user.Username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return ErrUsernameExists
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.userRepo.Create(ctx, user)
}

func (s *service) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.userRepo.FindAll(ctx)
}

func (s *service) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *service) UpdateUser(ctx context.Context, id string, user *models.User) error {
	// 检查用户是否存在
	existingUser, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return ErrUserNotFound
	}

	// 如果更新了密码，需要重新加密
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	return s.userRepo.Update(ctx, id, user)
}

func (s *service) DeleteUser(ctx context.Context, id string) error {
	// 检查用户是否存在
	existingUser, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return ErrUserNotFound
	}

	return s.userRepo.Delete(ctx, id)
}

// ChangePassword 仅更新密码字段，使用与登录一致的 SHA256 存储
func (s *service) ChangePassword(ctx context.Context, id string, newPassword string) error {
	if newPassword == "" {
		return ErrInvalidUserData
	}
	hashed := helper.Sha256Hex(newPassword)
	fields := bson.M{
		"password":   hashed,
		"updated_at": time.Now(),
	}
	return s.userRepo.UpdateFields(ctx, id, fields)
}
