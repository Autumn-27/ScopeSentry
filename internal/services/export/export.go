package export

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Autumn-27/ScopeSentry-go/internal/config"
	"github.com/Autumn-27/ScopeSentry-go/internal/logger"

	"github.com/Autumn-27/ScopeSentry-go/internal/models"
	"github.com/Autumn-27/ScopeSentry-go/internal/repositories/export"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Service 定义导出服务接口
type Service interface {
	CreateExportTask(ctx context.Context, task *models.ExportTask) error
	ProcessExportTask(ctx context.Context, run models.ExportRun)
	GetExportFields(ctx context.Context, index string) ([]string, error)
	GetExportRecords(ctx context.Context) ([]models.ExportTask, error)
	DeleteExportFiles(ctx context.Context, fileNames []string) error
	IsValidFileName(fileName string) bool
}

// ExportTask 导出任务结构体

type service struct {
	repository export.Repository
}

// NewService 创建导出服务实例
func NewService() Service {
	return &service{
		repository: export.NewRepository(),
	}
}

// CreateExportTask 创建导出任务
func (s *service) CreateExportTask(ctx context.Context, task *models.ExportTask) error {
	return s.repository.CreateExportTask(ctx, task)
}

// ProcessExportTask 处理导出任务
func (s *service) ProcessExportTask(ctx context.Context, run models.ExportRun) {

	// 在service层构建查询选项
	opts := options.Find().
		SetLimit(int64(run.Quantity))

	// 在service层处理投影
	projection := bson.M{}
	for _, field := range run.Field {
		projection[field] = 1
	}
	if run.Index == "asset" {
		projection["type"] = 1
	}
	opts.SetProjection(projection)

	// 获取游标
	cursor, err := s.repository.GetCollectionCursor(
		ctx,
		run.Index,
		run.Filter,
		opts,
	)
	if err != nil {
		err := s.repository.UpdateExportTaskStatus(ctx, run.FileName, 2, 0)
		if err != nil {
			logger.Error(fmt.Sprintf("Error updating export task status: %s", err.Error()))
			return
		}
		return
	}
	defer cursor.Close(ctx)

	// 导出数据
	filePath := filepath.Join(config.GlobalConfig.System.ExeDir, "files", "export", run.FileName)
	var fileSize float64
	if run.FileType == "xlsx" {
		fileSize, err = s.exportToXlsx(cursor, filePath, run.Field, run.Index)
	} else {
		fileSize, err = s.exportToJSON(cursor, filePath)
	}

	if err != nil {
		logger.Error(fmt.Sprintf("Error export: %s", err.Error()))
		err := s.repository.UpdateExportTaskStatus(ctx, run.FileName, 2, 0)
		if err != nil {
			logger.Error(fmt.Sprintf("Error updating export task status: %s", err.Error()))
			return
		}
		return
	}

	// 更新任务状态
	err = s.repository.UpdateExportTaskStatus(ctx, run.FileName, 1, fileSize)
	if err != nil {
		logger.Error(fmt.Sprintf("Error updating export task status: %s", err.Error()))
		return
	}
}

// GetExportFields 获取导出字段
func (s *service) GetExportFields(ctx context.Context, index string) ([]string, error) {
	// 从配置中获取字段列表
	fields, ok := models.Field[index]
	if !ok {
		return []string{}, fmt.Errorf("not found: %s", index)
	}
	return fields, nil
}

// GetExportRecords 获取导出记录
func (s *service) GetExportRecords(ctx context.Context) ([]models.ExportTask, error) {
	// 在 service 层构建查询选项
	opts := options.Find().
		SetSort(bson.D{{"create_time", -1}}).
		SetProjection(bson.D{
			{"_id", 0},
			{"id", bson.D{{"$toString", "$_id"}}},
			{"file_name", 1},
			{"end_time", 1},
			{"create_time", 1},
			{"data_type", 1},
			{"state", 1},
			{"file_size", 1},
			{"file_type", 1},
		})

	// 调用 repository 层执行查询
	return s.repository.GetExportRecords(ctx, bson.M{}, opts)
}

// DeleteExportFiles 删除导出文件
func (s *service) DeleteExportFiles(ctx context.Context, fileNames []string) error {
	// 删除文件
	for _, fileName := range fileNames {
		if s.IsValidFileName(fileName) {
			filePath := filepath.Join(config.GlobalConfig.System.ExeDir, "files", "export", fileName)
			if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
				logger.Error(fmt.Sprintf("删除文件失败: %s, 错误: %v", fileName, err))
				continue
			}
		}
	}

	// 调用 repository 层删除数据库记录
	return s.repository.DeleteExportTasks(ctx, fileNames)
}

// IsValidFileName 验证文件名是否有效
func (s *service) IsValidFileName(fileName string) bool {
	parts := strings.Split(fileName, ".")
	if len(parts) != 2 {
		return false
	}

	name := parts[0]
	fileType := parts[1]
	if fileType != "json" && fileType != "xlsx" {
		return false
	}
	matched, err := regexp.MatchString(`^[a-z0-9-]+$`, name)
	if err != nil {
		return false
	}
	return matched
}

// 生成Excel单元格坐标
func getCellName(col, row int) string {
	colName := ""
	for col > 0 {
		col--
		colName = string(rune('A'+col%26)) + colName
		col /= 26
	}
	if colName == "" {
		colName = "A"
	}
	return fmt.Sprintf("%s%d", colName, row)
}

// 使用游标导出到CSV
func (s *service) exportToXlsx(cursor *mongo.Cursor, filePath string, fields []string, index string) (float64, error) {
	// 创建新的Excel文件
	f := excelize.NewFile()
	defer f.Close()
	start := time.Now()
	if index == "asset" {
		// 分离HTTP和Other字段
		httpFields := []string{}
		otherFields := []string{}
		for _, field := range fields {
			if models.IsHTTPAssetField(field) {
				httpFields = append(httpFields, field)
			}
			if models.IsOtherAssetField(field) {
				otherFields = append(otherFields, field)
			}
		}

		// 创建HTTP工作表
		httpSheet := "HTTP Data"
		err := f.SetSheetName("Sheet1", httpSheet)
		if err != nil {
			logger.Error(fmt.Sprintf("Error setting sheet name: %s", err.Error()))
			return 0, err
		}

		// 创建Other工作表
		otherSheet := "Other Data"
		_, err = f.NewSheet(otherSheet)
		if err != nil {
			logger.Error(fmt.Sprintf("Error creating sheet: %s", err.Error()))
			return 0, err
		}

		// 写入表头
		for i, field := range httpFields {
			cell := getCellName(i+1, 1)
			f.SetCellValue(httpSheet, cell, field)
		}
		for i, field := range otherFields {
			cell := getCellName(i+1, 1)
			f.SetCellValue(otherSheet, cell, field)
		}

		// 处理数据
		httpRow := 2
		otherRow := 2
		for cursor.Next(context.Background()) {
			var doc bson.M
			if err := cursor.Decode(&doc); err != nil {
				logger.Error(fmt.Sprintf("Failed to decode document: %v, skipping...", err))
				continue
			}

			if doc["type"] == "http" {
				// 写入HTTP数据
				for i, field := range httpFields {
					cell := getCellName(i+1, httpRow)
					value := doc[field]
					if value == nil {
						if err := f.SetCellValue(httpSheet, cell, ""); err != nil {
							logger.Error(fmt.Sprintf("Failed to write empty value for field %s: %v, skipping...", field, err))
							continue
						}
					} else {
						if err := f.SetCellValue(httpSheet, cell, fmt.Sprintf("%v", value)); err != nil {
							logger.Error(fmt.Sprintf("Failed to write value for field %s: %v, skipping...", field, err))
							continue
						}
					}
				}
				httpRow++
			} else {
				// 写入Other数据
				for i, field := range otherFields {
					cell := getCellName(i+1, otherRow)
					value := doc[field]
					if value == nil {
						if err := f.SetCellValue(otherSheet, cell, ""); err != nil {
							logger.Error(fmt.Sprintf("Failed to write empty value for field %s: %v, skipping...", field, err))
							continue
						}
					} else {
						if err := f.SetCellValue(otherSheet, cell, fmt.Sprintf("%v", value)); err != nil {
							logger.Error(fmt.Sprintf("Failed to write value for field %s: %v, skipping...", field, err))
							continue
						}
					}
				}
				otherRow++
			}
		}

		// 检查游标错误
		if err := cursor.Err(); err != nil {
			return 0, fmt.Errorf("cursor error: %w", err)
		}

		// 保存文件
		if err := f.SaveAs(filePath); err != nil {
			return 0, fmt.Errorf("failed to save excel file: %w", err)
		}

		// 获取文件大小
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return 0, fmt.Errorf("failed to get file info: %w", err)
		}
		fileSize := float64(fileInfo.Size()) / (1024 * 1024)
		end := time.Now()
		logger.Info(fmt.Sprintf("Saved %v excel file size: %.2f MB time: %v number: %v", index, fileSize, end.Sub(start), httpRow+otherRow))
		return fileSize, nil
	} else {
		// 处理其他类型的数据
		sheet := "Sheet1"

		// 写入表头
		for i, field := range fields {
			cell := getCellName(i+1, 1)
			if err := f.SetCellValue(sheet, cell, field); err != nil {
				logger.Error(fmt.Sprintf("Failed to write header for field %s: %v, skipping...", field, err))
				continue
			}
		}

		// 处理数据
		row := 2
		for cursor.Next(context.Background()) {
			var doc bson.M
			if err := cursor.Decode(&doc); err != nil {
				logger.Error(fmt.Sprintf("Failed to decode document: %v, skipping...", err))
				continue
			}

			// 写入数据
			for i, field := range fields {
				cell := getCellName(i+1, row)
				value := doc[field]
				if value == nil {
					if err := f.SetCellValue(sheet, cell, ""); err != nil {
						logger.Error(fmt.Sprintf("Failed to write empty value for field %s: %v, skipping...", field, err))
						continue
					}
				} else {
					if err := f.SetCellValue(sheet, cell, fmt.Sprintf("%v", value)); err != nil {
						logger.Error(fmt.Sprintf("Failed to write value for field %s: %v, skipping...", field, err))
						continue
					}
				}
			}
			row++
		}

		// 检查游标错误
		if err := cursor.Err(); err != nil {
			return 0, fmt.Errorf("cursor error: %w", err)
		}

		// 保存文件
		if err := f.SaveAs(filePath); err != nil {
			return 0, fmt.Errorf("failed to save excel file: %w", err)
		}

		// 获取文件大小
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return 0, fmt.Errorf("failed to get file info: %w", err)
		}
		fileSize := float64(fileInfo.Size()) / (1024 * 1024)
		end := time.Now()
		logger.Info(fmt.Sprintf("Saved %v excel file size: %.2f MB time: %v number: %v", index, fileSize, end.Sub(start), row))
		return fileSize, nil
	}
}

// 使用游标导出到JSON
func (s *service) exportToJSON(cursor *mongo.Cursor, filePath string) (float64, error) {
	start := time.Now()
	file, err := os.Create(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 使用游标逐条读取数据
	row := 0
	for cursor.Next(context.Background()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			logger.Error(fmt.Sprintf("Failed to decode document: %v, skipping...", err))
			continue
		}

		// 将文档转换为JSON字符串
		jsonData, err := json.Marshal(doc)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to marshal document: %v, skipping...", err))
			continue
		}

		// 写入一行数据
		if _, err := file.Write(jsonData); err != nil {
			logger.Error(fmt.Sprintf("Failed to write document: %v, skipping...", err))
			continue
		}
		if _, err := file.Write([]byte("\n")); err != nil {
			logger.Error(fmt.Sprintf("Failed to write newline: %v, skipping...", err))
			continue
		}
		row++
	}

	// 检查游标错误
	if err := cursor.Err(); err != nil {
		return 0, fmt.Errorf("cursor error: %w", err)
	}

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}

	fileSize := float64(fileInfo.Size()) / (1024 * 1024)
	end := time.Now()
	logger.Info(fmt.Sprintf("Saved JSON file size: %.2f MB time: %v number: %v", fileSize, end.Sub(start), row))
	return fileSize, nil
}
