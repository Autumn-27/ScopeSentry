package dictionary

import (
	"context"
	"fmt"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	dictrepo "github.com/Autumn-27/ScopeSentry/internal/repositories/dictionary"
	nodessvc "github.com/Autumn-27/ScopeSentry/internal/services/node"
)

type ManageService interface {
	List(ctx context.Context) ([]models.DictionaryMeta, error)
	Create(ctx context.Context, name, category string, content []byte) error
	Download(ctx context.Context, id string) (filename string, data []byte, err error)
	Delete(ctx context.Context, ids []string) error
	Save(ctx context.Context, id string, content []byte) error
}

type manageService struct {
	repo        dictrepo.ManageRepository
	nodeService nodessvc.Service
}

func NewManageService() ManageService {
	return &manageService{
		repo:        dictrepo.NewManageRepository(),
		nodeService: nodessvc.NewService(),
	}
}

func (s *manageService) List(ctx context.Context) ([]models.DictionaryMeta, error) {
	return s.repo.List(ctx)
}

func (s *manageService) Create(ctx context.Context, name, category string, content []byte) error {
	// 重复性检查
	exists, err := s.repo.ExistsByNameCategory(ctx, name, category)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("duplication file name")
	}

	// 插入元数据，使用插入ID作为GridFS文件名
	id, err := s.repo.InsertMeta(ctx, name, category, fmt.Sprintf("%.2f", float64(len(content))/(1024*1024)))
	if err != nil {
		return err
	}
	if err := s.repo.UploadFile(ctx, id, content); err != nil {
		return err
	}

	// 刷新节点配置
	_ = s.nodeService.RefreshConfig(ctx, models.Message{ // 忽略错误不中断
		Name:    "all",
		Type:    "dictionary",
		Content: fmt.Sprintf("add:%s", id),
	})
	return nil
}

func (s *manageService) Download(ctx context.Context, id string) (string, []byte, error) {
	data, err := s.repo.DownloadFile(ctx, id)
	if err != nil {
		return "", nil, err
	}
	return id, data, nil
}

func (s *manageService) Delete(ctx context.Context, ids []string) error {
	if err := s.repo.DeleteMeta(ctx, ids); err != nil {
		return err
	}
	// 删除GridFS
	for _, id := range ids {
		if err := s.repo.DeleteFile(ctx, id); err != nil {
			return err
		}
		_ = s.nodeService.RefreshConfig(ctx, models.Message{
			Name:    "all",
			Type:    "dictionary",
			Content: fmt.Sprintf("delete:%s", id),
		})
	}
	return nil
}

func (s *manageService) Save(ctx context.Context, id string, content []byte) error {
	// 替换GridFS文件
	if err := s.repo.ReplaceFile(ctx, id, content); err != nil {
		return err
	}
	// 更新元数据大小
	if err := s.repo.UpdateMetaSize(ctx, id, fmt.Sprintf("%.2f", float64(len(content))/(1024*1024))); err != nil {
		return err
	}
	_ = s.nodeService.RefreshConfig(ctx, models.Message{
		Name:    "all",
		Type:    "dictionary",
		Content: fmt.Sprintf("add:%s", id),
	})
	return nil
}
