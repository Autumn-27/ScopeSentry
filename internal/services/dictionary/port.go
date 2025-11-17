package dictionary

import (
	"context"
	"fmt"

	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	dictrepo "github.com/Autumn-27/ScopeSentry-go/internal/repositories/dictionary"
	nodessvc "github.com/Autumn-27/ScopeSentry-go/internal/services/node"
)

type PortService interface {
	Get(ctx context.Context, search string, pageIndex, pageSize int) ([]models.PortDoc, int64, error)
	Add(ctx context.Context, name, value string) error
	Update(ctx context.Context, id, name, value string) error
	Delete(ctx context.Context, ids []string) error
}

type portService struct {
	repo        dictrepo.PortRepository
	nodeService nodessvc.Service
}

func NewPortService() PortService {
	return &portService{
		repo:        dictrepo.NewPortRepository(),
		nodeService: nodessvc.NewService(),
	}
}

func (s *portService) Get(ctx context.Context, search string, pageIndex, pageSize int) ([]models.PortDoc, int64, error) {
	total, err := s.repo.Count(ctx, search)
	if err != nil {
		return nil, 0, err
	}
	list, err := s.repo.Find(ctx, search, pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *portService) Add(ctx context.Context, name, value string) error {
	if value == "" {
		return fmt.Errorf("value is null")
	}
	if err := s.repo.Insert(ctx, name, value); err != nil {
		return err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "port"})
	return nil
}

func (s *portService) Update(ctx context.Context, id, name, value string) error {
	if err := s.repo.Update(ctx, id, name, value); err != nil {
		return err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "port"})
	return nil
}

func (s *portService) Delete(ctx context.Context, ids []string) error {
	if err := s.repo.Delete(ctx, ids); err != nil {
		return err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{Name: "all", Type: "port"})
	return nil
}
