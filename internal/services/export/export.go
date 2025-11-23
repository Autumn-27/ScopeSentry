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

	"github.com/Autumn-27/ScopeSentry/internal/config"
	"github.com/Autumn-27/ScopeSentry/internal/logger"

	"github.com/Autumn-27/ScopeSentry/internal/models"
	"github.com/Autumn-27/ScopeSentry/internal/repositories/export"
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

// sanitizeExcelValue 清理Excel不允许的非法字符
// Excel XML不允许控制字符（0x00-0x1F，除了0x09(TAB)、0x0A(LF)、0x0D(CR)）
func sanitizeExcelValue(value interface{}) string {
	if value == nil {
		return ""
	}

	str := fmt.Sprintf("%v", value)
	// 过滤非法控制字符，保留TAB、LF、CR
	result := make([]rune, 0, len(str))
	for _, r := range str {
		// 允许：可打印字符(>=0x20)、TAB(0x09)、LF(0x0A)、CR(0x0D)
		if r >= 0x20 || r == 0x09 || r == 0x0A || r == 0x0D {
			result = append(result, r)
		} else {
			// 非法字符替换为空格或直接跳过
			// 这里选择跳过，如果需要保留位置可以用空格
		}
	}
	return string(result)
}

// 使用流式写入导出到Excel
func (s *service) exportToXlsx(cursor *mongo.Cursor, filePath string, fields []string, index string) (float64, error) {
	// 创建新的Excel文件
	f := excelize.NewFile()
	defer f.Close()
	start := time.Now()

	// 流式写入刷新间隔（每5000行刷新一次，减少内存占用）
	const flushInterval = 1000

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

		// 创建流式写入器
		httpStreamWriter, err := f.NewStreamWriter(httpSheet)
		if err != nil {
			return 0, fmt.Errorf("failed to create stream writer for HTTP sheet: %w", err)
		}

		otherStreamWriter, err := f.NewStreamWriter(otherSheet)
		if err != nil {
			return 0, fmt.Errorf("failed to create stream writer for Other sheet: %w", err)
		}

		// 写入表头 - 转换为interface{}切片
		httpHeader := make([]interface{}, len(httpFields))
		for i, field := range httpFields {
			httpHeader[i] = field
		}
		if err := httpStreamWriter.SetRow("A1", httpHeader); err != nil {
			return 0, fmt.Errorf("failed to write HTTP header: %w", err)
		}

		otherHeader := make([]interface{}, len(otherFields))
		for i, field := range otherFields {
			otherHeader[i] = field
		}
		if err := otherStreamWriter.SetRow("A1", otherHeader); err != nil {
			return 0, fmt.Errorf("failed to write Other header: %w", err)
		}

		// 处理数据
		// 注意：excelize流式写入要求行号必须严格连续递增，不能跳过任何行号
		// 表头已写入第1行，数据从第2行开始
		httpRowIdx := 2  // 从2开始（1是表头）
		otherRowIdx := 2 // 从2开始（1是表头）
		httpCount := 0   // HTTP实际成功写入的行数
		otherCount := 0  // Other实际成功写入的行数

		for cursor.Next(context.Background()) {
			var doc bson.M
			if err := cursor.Decode(&doc); err != nil {
				logger.Error(fmt.Sprintf("Failed to decode document: %v, skipping...", err))
				continue
			}

			if doc["type"] == "http" {
				// 准备HTTP数据行，清理非法字符
				rowData := make([]interface{}, len(httpFields))
				for i, field := range httpFields {
					value := doc[field]
					rowData[i] = sanitizeExcelValue(value)
				}

				// 使用流式写入 - 行号必须严格连续递增
				cell, err := excelize.CoordinatesToCellName(1, httpRowIdx)
				if err != nil {
					logger.Error(fmt.Sprintf("Failed to convert coordinates for HTTP row %d: %v", httpRowIdx, err))
					httpRowIdx++ // 即使失败也要递增行号，保持连续性
					continue
				}

				if err := httpStreamWriter.SetRow(cell, rowData); err != nil {
					logger.Error(fmt.Sprintf("Failed to write HTTP row %d: %v", httpRowIdx, err))
					// 尝试写入空行以保持行号连续性
					emptyRow := make([]interface{}, len(httpFields))
					for i := range emptyRow {
						emptyRow[i] = ""
					}
					if retryErr := httpStreamWriter.SetRow(cell, emptyRow); retryErr != nil {
						return 0, fmt.Errorf("excelize stream writer corrupted at HTTP row %d: %w", httpRowIdx, retryErr)
					}
				} else {
					httpCount++
				}
				httpRowIdx++ // 无论成功失败，行号都必须递增
			} else {
				// 准备Other数据行，清理非法字符
				rowData := make([]interface{}, len(otherFields))
				for i, field := range otherFields {
					value := doc[field]
					rowData[i] = sanitizeExcelValue(value)
				}

				// 使用流式写入 - 行号必须严格连续递增
				cell, err := excelize.CoordinatesToCellName(1, otherRowIdx)
				if err != nil {
					logger.Error(fmt.Sprintf("Failed to convert coordinates for Other row %d: %v", otherRowIdx, err))
					otherRowIdx++ // 即使失败也要递增行号，保持连续性
					continue
				}

				if err := otherStreamWriter.SetRow(cell, rowData); err != nil {
					logger.Error(fmt.Sprintf("Failed to write Other row %d: %v", otherRowIdx, err))
					// 尝试写入空行以保持行号连续性
					emptyRow := make([]interface{}, len(otherFields))
					for i := range emptyRow {
						emptyRow[i] = ""
					}
					if retryErr := otherStreamWriter.SetRow(cell, emptyRow); retryErr != nil {
						return 0, fmt.Errorf("excelize stream writer corrupted at Other row %d: %w", otherRowIdx, retryErr)
					}
				} else {
					otherCount++
				}
				otherRowIdx++ // 无论成功失败，行号都必须递增
			}

		}

		// 检查游标错误
		if err := cursor.Err(); err != nil {
			return 0, fmt.Errorf("cursor error: %w", err)
		}

		// 最终刷新
		if err := httpStreamWriter.Flush(); err != nil {
			return 0, fmt.Errorf("failed to flush HTTP stream: %w", err)
		}
		if err := otherStreamWriter.Flush(); err != nil {
			return 0, fmt.Errorf("failed to flush Other stream: %w", err)
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
		totalRows := httpCount + otherCount
		logger.Info(fmt.Sprintf("Saved %v excel file size: %.2f MB time: %v number: %v (HTTP: %v, Other: %v)", index, fileSize, end.Sub(start), totalRows, httpCount, otherCount))
		return fileSize, nil
	} else {
		// 处理其他类型的数据
		sheet := "Sheet1"

		// 创建流式写入器
		streamWriter, err := f.NewStreamWriter(sheet)
		if err != nil {
			return 0, fmt.Errorf("failed to create stream writer: %w", err)
		}

		// 写入表头
		header := make([]interface{}, len(fields))
		for i, field := range fields {
			header[i] = field
		}
		if err := streamWriter.SetRow("A1", header); err != nil {
			return 0, fmt.Errorf("failed to write header: %w", err)
		}

		// 处理数据
		row := 2
		rowCount := 0 // 实际写入的行数

		for cursor.Next(context.Background()) {
			var doc bson.M
			if err := cursor.Decode(&doc); err != nil {
				logger.Error(fmt.Sprintf("Failed to decode document: %v, skipping...", err))
				continue
			}

			// 准备数据行，清理非法字符
			rowData := make([]interface{}, len(fields))
			for i, field := range fields {
				value := doc[field]
				rowData[i] = sanitizeExcelValue(value)
			}

			// 使用流式写入
			cell, err := excelize.CoordinatesToCellName(1, row)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to convert coordinates for row %d: %v", row, err))
				row++ // 即使失败也要递增行号，保持连续性
				continue
			}

			if err := streamWriter.SetRow(cell, rowData); err != nil {
				logger.Error(fmt.Sprintf("Failed to write row %d: %v", row, err))
				// 尝试写入空行以保持行号连续性
				emptyRow := make([]interface{}, len(fields))
				for i := range emptyRow {
					emptyRow[i] = ""
				}
				if retryErr := streamWriter.SetRow(cell, emptyRow); retryErr != nil {
					return 0, fmt.Errorf("excelize stream writer corrupted at row %d: %w", row, retryErr)
				}
			} else {
				rowCount++
			}
			row++ // 无论成功失败，行号都必须递增
		}

		// 检查游标错误
		if err := cursor.Err(); err != nil {
			return 0, fmt.Errorf("cursor error: %w", err)
		}

		// 最终刷新
		if err := streamWriter.Flush(); err != nil {
			return 0, fmt.Errorf("failed to flush stream: %w", err)
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
		logger.Info(fmt.Sprintf("Saved %v excel file size: %.2f MB time: %v number: %v", index, fileSize, end.Sub(start), rowCount))
		return fileSize, nil
	}
}

// 使用流式编码导出到JSON
func (s *service) exportToJSON(cursor *mongo.Cursor, filePath string) (float64, error) {
	start := time.Now()
	file, err := os.Create(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 使用JSON编码器进行流式编码，直接写入文件，避免创建临时字节数组
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false) // 不转义HTML字符，提高性能

	// 使用游标逐条读取数据
	row := 0
	for cursor.Next(context.Background()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			logger.Error(fmt.Sprintf("Failed to decode document: %v, skipping...", err))
			continue
		}

		// 使用流式编码直接写入文件，避免创建临时字节数组
		if err := encoder.Encode(doc); err != nil {
			logger.Error(fmt.Sprintf("Failed to encode document: %v, skipping...", err))
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
