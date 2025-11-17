package poc

import (
	"fmt"
	"strings"

	"github.com/Autumn-27/ScopeSentry-go/internal/logger"
	"github.com/Autumn-27/ScopeSentry-go/internal/services/node"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/repositories/poc"
	"github.com/Autumn-27/ScopeSentry-go/internal/utils/helper"
	"github.com/Autumn-27/ScopeSentry-go/internal/utils/random"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Service 定义POC服务接口
type Service interface {
	GetPocList(ctx *gin.Context, req *models.PocListRequest) (*models.PocListResponse, error)
	GetPocDetail(ctx *gin.Context, req *models.PocDetailRequest) (*models.Poc, error)
	UpdatePoc(ctx *gin.Context, req *models.PocUpdateRequest) error
	AddPoc(ctx *gin.Context, req *models.PocAddRequest) error
	DeletePoc(ctx *gin.Context, req *models.PocDeleteRequest) error
	GetAllPocData(ctx *gin.Context) ([]models.Poc, error)
	ImportPoc(ctx *gin.Context, filePath string) (*models.PocImportResponse, error)
}

type service struct {
	repo        poc.Repository
	nodeService node.Service
}

// NewService 创建新的POC服务实例
func NewService() Service {
	return &service{
		repo:        poc.NewRepository(),
		nodeService: node.NewService(),
	}
}

func (s *service) GetPocList(ctx *gin.Context, req *models.PocListRequest) (*models.PocListResponse, error) {
	query := bson.M{
		"name": bson.M{
			"$regex":   req.Search,
			"$options": "i",
		},
	}

	if len(req.Filter.Level) > 0 {
		query["level"] = bson.M{"$in": req.Filter.Level}
	}

	opts := options.Find().
		SetSort(bson.D{{"time", -1}}).
		SetSkip(int64((req.PageIndex - 1) * req.PageSize)).
		SetLimit(int64(req.PageSize))

	results, err := s.repo.FindWithPagination(ctx, query, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get poc list: %w", err)
	}

	total, err := s.repo.Count(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to count poc: %w", err)
	}

	return &models.PocListResponse{
		List:  results,
		Total: total,
	}, nil
}

func (s *service) GetPocDetail(ctx *gin.Context, req *models.PocDetailRequest) (*models.Poc, error) {
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	result, err := s.repo.FindOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, fmt.Errorf("failed to get poc detail: %w", err)
	}

	return result, nil
}

func (s *service) UpdatePoc(ctx *gin.Context, req *models.PocUpdateRequest) error {
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return fmt.Errorf("invalid id format: %w", err)
	}

	var pt models.PocTemplate
	err = yaml.Unmarshal([]byte(req.Content), &pt)
	if err != nil {
		return fmt.Errorf("failed to unmarshal poc template: %w", err)
	}
	ptTags := strings.Split(pt.Info.Tags, ",")
	tags := helper.RemoveArrayDuplicates(ptTags)
	update := bson.M{
		"$set": bson.M{
			"name":    pt.Info.Name,
			"content": req.Content,
			"level":   pt.Info.Severity,
			"tags":    tags,
		},
	}

	_, err = s.repo.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("failed to update poc: %w", err)
	}

	return nil
}

func (s *service) AddPoc(ctx *gin.Context, req *models.PocAddRequest) error {

	var pt models.PocTemplate
	err := yaml.Unmarshal([]byte(req.Content), &pt)
	if err != nil {
		return fmt.Errorf("failed to unmarshal poc template: %w", err)
	}
	id := pt.ID
	existing, err := s.repo.FindOne(ctx, bson.M{"id": id})
	if err == nil && existing != nil {
		return fmt.Errorf("poc already exists")
	}
	ptTags := strings.Split(pt.Info.Tags, ",")
	tags := helper.RemoveArrayDuplicates(ptTags)
	doc := models.Poc{
		Name:       pt.Info.Name,
		Content:    req.Content,
		Level:      pt.Info.Severity,
		TemplateId: id,
		Tags:       tags,
		Time:       helper.GetNowTimeString(),
	}

	res, err := s.repo.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to add poc: %w", err)
	}
	go func() {
		msg := models.Message{
			Name:    "all",
			Type:    "poc",
			Content: fmt.Sprintf(`add:%v`, res.InsertedID),
		}
		err = s.nodeService.RefreshConfig(ctx, msg)
		if err != nil {
			logger.Error("failed to refresh config", zap.Error(err))
		}
	}()
	return nil
}

func (s *service) DeletePoc(ctx *gin.Context, req *models.PocDeleteRequest) error {
	var objIDs []primitive.ObjectID
	for _, id := range req.IDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return fmt.Errorf("invalid id format: %w", err)
		}
		objIDs = append(objIDs, objID)
	}

	_, err := s.repo.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": objIDs}})
	if err != nil {
		return fmt.Errorf("failed to delete poc: %w", err)
	}
	go func() {
		msg := models.Message{
			Name:    "all",
			Type:    "poc",
			Content: fmt.Sprintf(`delete:%v`, strings.Join(req.IDs, ",")),
		}
		err = s.nodeService.RefreshConfig(ctx, msg)
		if err != nil {
			logger.Error("failed to refresh config", zap.Error(err))
		}
	}()
	return nil
}

func (s *service) GetAllPocData(ctx *gin.Context) ([]models.Poc, error) {
	result, err := s.repo.GetAllPocData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all poc data: %w", err)
	}
	return result, nil
}

// ImportPoc 导入POC文件
func (s *service) ImportPoc(ctx *gin.Context, filePath string) (*models.PocImportResponse, error) {
	logger.Info("POC导入开始")

	// 生成随机文件名
	fileName := random.GenerateString(11)

	// 使用系统临时目录而不是相对路径
	tempDir := os.TempDir()

	// 解压路径也使用临时目录
	unzipPath := filepath.Join(tempDir, fileName)
	extractPath := unzipPath

	// 确保目录存在
	if err := os.MkdirAll(extractPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create extract directory: %w", err)
	}

	// 解压ZIP文件
	yamlFiles, err := s.extractZipFile(filePath, extractPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract zip file: %w", err)
	}

	// 获取现有POC的ID列表
	templateIds, err := s.repo.GetAllTemplateId(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing poc hashes: %w", err)
	}

	// 处理YAML文件
	successNum := 0
	errorNum := 0
	repeatNum := 0
	pocDataList := []models.Poc{}

	logger.Info(fmt.Sprintf("共%d个POC", len(yamlFiles)))

	for _, yamlFile := range yamlFiles {

		content, err := os.ReadFile(yamlFile)
		if err != nil {
			logger.Error(fmt.Sprintf("POC导入 读取文件失败: %s", yamlFile), zap.Error(err))
			errorNum++
			continue
		}
		// 解析YAML
		var pt models.PocTemplate
		if err := yaml.Unmarshal(content, &pt); err != nil {
			logger.Error(fmt.Sprintf("POC导入 解析YAML失败: %s", yamlFile), zap.Error(err))
			errorNum++
			continue
		}
		if pt.Info.Name == "" || pt.ID == "" {
			errorNum++
			continue
		}
		// 检查是否重复
		if helper.StringInSlice(pt.ID, templateIds) {
			repeatNum++
			continue
		}
		// 获取严重程度
		severity := "unknown"
		if pt.Info.Severity != "" {
			severity = strings.ToLower(pt.Info.Severity)
		}

		// 处理标签
		tags := []string{}
		if pt.Info.Tags != "" {
			tagList := strings.Split(pt.Info.Tags, ",")
			for _, tag := range tagList {
				tag = strings.TrimSpace(tag)
				if tag != "" {
					tags = append(tags, tag)
				}
			}
		}
		pocData := models.Poc{
			Name:       pt.Info.Name,
			Content:    string(content),
			TemplateId: pt.ID,
			Level:      severity,
			Time:       helper.GetNowTimeString(),
			Tags:       tags,
		}

		pocDataList = append(pocDataList, pocData)
	}

	// 批量插入POC数据
	if len(pocDataList) > 0 {
		result, err := s.repo.InsertMany(ctx, pocDataList)
		if err != nil {
			return nil, fmt.Errorf("failed to insert poc data: %w", err)
		}

		if result.InsertedIDs != nil {
			successNum = len(result.InsertedIDs)

			// 异步刷新配置
			go func() {
				insertedIDs := make([]string, len(result.InsertedIDs))
				for i, id := range result.InsertedIDs {
					insertedIDs[i] = id.(primitive.ObjectID).Hex()
				}

				msg := models.Message{
					Name:    "all",
					Type:    "poc",
					Content: fmt.Sprintf("add:%s", strings.Join(insertedIDs, ",")),
				}

				if err := s.nodeService.RefreshConfig(ctx, msg); err != nil {
					logger.Error("failed to refresh config", zap.Error(err))
				}
			}()
		}
	}

	// 清理临时文件
	s.cleanupTempFiles(filePath, extractPath)

	logger.Info(fmt.Sprintf("POC import completed: %d succeeded, %d duplicate, %d failed", successNum, repeatNum, errorNum))
	logger.Info("POC import finished")

	return &models.PocImportResponse{
		SuccessNum: successNum,
		ErrorNum:   errorNum,
		RepeatNum:  repeatNum,
		Message:    "import success",
	}, nil
}

// extractZipFile 解压ZIP文件并返回YAML文件列表
func (s *service) extractZipFile(zipPath, extractPath string) ([]string, error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var yamlFiles []string

	// 确保解压路径是绝对路径
	absExtractPath, err := filepath.Abs(extractPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	for _, file := range reader.File {
		// 检查文件名是否包含路径穿越字符
		cleanName := filepath.Clean(file.Name)
		if strings.HasPrefix(cleanName, "..") || strings.Contains(cleanName, ".."+string(os.PathSeparator)) {
			logger.Error("Path traversal detected in ZIP file", zap.String("path", file.Name))
			continue
		}

		// 构建安全的文件路径
		filePath := filepath.Join(absExtractPath, file.Name)

		// 规范化路径
		filePath = filepath.Clean(filePath)

		// 安全检查：确保解压路径在目标目录内
		if !strings.HasPrefix(filePath, absExtractPath) {
			logger.Error("Path traversal detected in ZIP file", zap.String("path", file.Name), zap.String("resolved", filePath))
			continue
		}

		// 额外检查：确保路径不包含任何可疑的目录遍历
		if strings.Contains(filePath, "..") {
			logger.Error("Path traversal detected after resolution", zap.String("path", file.Name), zap.String("resolved", filePath))
			continue
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, 0755); err != nil {
				logger.Error("Failed to create directory", zap.String("path", filePath), zap.Error(err))
				continue
			}
			continue
		}

		// 创建父目录
		parentDir := filepath.Dir(filePath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			logger.Error("Failed to create parent directory", zap.String("path", parentDir), zap.Error(err))
			continue
		}

		// 解压文件
		rc, err := file.Open()
		if err != nil {
			logger.Error("Failed to open file in ZIP", zap.String("path", file.Name), zap.Error(err))
			continue
		}

		dst, err := os.Create(filePath)
		if err != nil {
			rc.Close()
			logger.Error("Failed to create file", zap.String("path", filePath), zap.Error(err))
			continue
		}

		_, err = io.Copy(dst, rc)
		rc.Close()
		dst.Close()

		if err != nil {
			logger.Error("Failed to copy file content", zap.String("path", filePath), zap.Error(err))
			continue
		}

		// 如果是YAML文件，添加到列表
		if strings.HasSuffix(strings.ToLower(file.Name), ".yaml") || strings.HasSuffix(strings.ToLower(file.Name), ".yml") {
			yamlFiles = append(yamlFiles, filePath)
		}
	}

	return yamlFiles, nil
}

// cleanupTempFiles 清理临时文件
func (s *service) cleanupTempFiles(zipPath, extractPath string) {
	// 删除ZIP文件
	if err := os.Remove(zipPath); err != nil {
		logger.Error("删除POC ZIP文件出错", zap.Error(err))
	}

	// 删除解压目录
	if err := os.RemoveAll(extractPath); err != nil {
		logger.Error("删除POC解压目录出错", zap.Error(err))
	}
}

// isSafeFilePath 检查文件路径是否安全
func (s *service) isSafeFilePath(filePath string) bool {
	// 检查路径是否包含危险字符
	if strings.Contains(filePath, "..") {
		return false
	}

	// 规范化路径
	cleanPath := filepath.Clean(filePath)

	// 检查规范化后的路径是否仍然包含危险字符
	if strings.Contains(cleanPath, "..") {
		return false
	}

	// if filepath.IsAbs(cleanPath) {
	//     return false
	// }

	return true
}
