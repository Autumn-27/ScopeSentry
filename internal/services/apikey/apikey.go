package apikey

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/repositories/apikey"
	"github.com/Autumn-27/ScopeSentry/internal/utils/random"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

const keyPrefix = "ssk_"

var (
	ErrInvalidKey = errors.New("invalid api key")
	ErrNotFound   = errors.New("api key not found")
)

type Service interface {
	Create(ctx context.Context, name, createdBy string) (*models.CreateApiKeyResponse, error)
	List(ctx context.Context) ([]models.ApiKey, error)
	Delete(ctx context.Context, id string) error
	Validate(ctx context.Context, rawKey string) (*models.ApiKey, error)
	EnsureIndexes(ctx context.Context) error
}

type service struct {
	repo apikey.Repository
}

func NewService() Service {
	return &service{repo: apikey.NewRepository()}
}

func (s *service) Create(ctx context.Context, name, createdBy string) (*models.CreateApiKeyResponse, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	rawKey := keyPrefix + random.GenerateRandomString(40)
	keyHash := hashKey(rawKey)
	displayPrefix := rawKey[:12] + "..."

	record := &models.ApiKey{
		Name:      name,
		KeyHash:   keyHash,
		KeyPrefix: displayPrefix,
		CreatedBy: createdBy,
	}
	if err := s.repo.Create(ctx, record); err != nil {
		return nil, err
	}

	return &models.CreateApiKeyResponse{
		ID:        record.ID.Hex(),
		Name:      record.Name,
		Key:       rawKey,
		KeyPrefix: displayPrefix,
		CreatedAt: record.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *service) List(ctx context.Context) ([]models.ApiKey, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) Validate(ctx context.Context, rawKey string) (*models.ApiKey, error) {
	rawKey = strings.TrimSpace(rawKey)
	if rawKey == "" || !strings.HasPrefix(rawKey, keyPrefix) {
		return nil, ErrInvalidKey
	}

	record, err := s.repo.FindByHash(ctx, hashKey(rawKey))
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, ErrInvalidKey
		}
		return nil, err
	}

	go func() {
		_ = s.repo.UpdateLastUsed(context.Background(), record.ID)
	}()

	return record, nil
}

func (s *service) EnsureIndexes(ctx context.Context) error {
	return s.repo.EnsureIndexes(ctx)
}

func hashKey(rawKey string) string {
	sum := sha256.Sum256([]byte(rawKey))
	return hex.EncodeToString(sum[:])
}

// IsAPIKeyFormat 判断 token 是否为 API Key 格式（非 JWT）
func IsAPIKeyFormat(token string) bool {
	return strings.HasPrefix(token, keyPrefix)
}
