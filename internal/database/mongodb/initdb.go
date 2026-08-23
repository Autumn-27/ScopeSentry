package mongodb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"github.com/Autumn-27/ScopeSentry/internal/utils/helper"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Autumn-27/ScopeSentry/internal/constants"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Autumn-27/ScopeSentry/internal/config"
	"github.com/Autumn-27/ScopeSentry/internal/utils/random"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func printProgressBar(step, total int, action string) {
	width := 50
	progress := float64(step) / float64(total)
	bar := int(progress * float64(width))

	fmt.Printf("\r[%s%s] %d/%d %s\n",
		strings.Repeat("=", bar),
		strings.Repeat(" ", width-bar),
		step, total, action)
}

func CreateDatabase() error {
	client := Client
	if client == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}

	// 获取数据库列表
	databases, err := client.ListDatabaseNames(context.Background(), bson.M{})
	if err != nil {
		return fmt.Errorf("failed to list databases: %v", err)
	}

	db := client.Database(config.GlobalConfig.MongoDB.Database)

	// 如果数据库不存在，创建数据库
	if !contains(databases, config.GlobalConfig.MongoDB.Database) {
		totalSteps := 13
		currentStep := 0

		// 创建用户集合
		collection := db.Collection("user")
		password, err := random.GeneratePassword(16)
		if err != nil {
			password = random.GenerateRandomString(16)
		}

		// 打印重要信息
		separator := strings.Repeat("=", 50)
		fmt.Printf("%s\n", separator)
		fmt.Println("✨✨✨ IMPORTANT NOTICE: Please review the User/Password below ✨✨✨")
		fmt.Println(separator)
		fmt.Printf("🔑 User/Password: ScopeSentry/%s\n", password)
		fmt.Println(separator)
		fmt.Println("✅ Ensure the User/Password is correctly copied!")
		fmt.Println("✅ The initialization password is stored in the file PASSWORD")

		// 保存密码到文件
		if err := os.WriteFile(filepath.Join(config.GlobalConfig.System.ExeDir, "PASSWORD"), []byte(password), 0644); err != nil {
			return fmt.Errorf("failed to write password file: %v", err)
		}

		// 加密密码
		hashedPassword := helper.Sha256Hex(password)
		// 创建用户
		_, err = collection.InsertOne(context.Background(), bson.M{
			"username": "ScopeSentry",
			"password": hashedPassword,
		})
		if err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		}
		currentStep++
		printProgressBar(currentStep, totalSteps, "Creating user")

		// 创建配置集合
		configCollection := db.Collection("config")

		// 插入系统配置
		_, err = configCollection.InsertOne(context.Background(), bson.M{
			"name":  "timezone",
			"value": "Asia/Shanghai",
			"type":  "system",
		})
		if err != nil {
			return fmt.Errorf("failed to insert timezone config: %v", err)
		}
		configCollection.InsertOne(context.Background(), bson.M{
			"name":  "ModulesConfig",
			"value": constants.ModulesConfig,
			"type":  "system",
		})

		configCollection.InsertOne(context.Background(), bson.M{
			"name":  "SubfinderApiConfig",
			"value": constants.SubfinderApiConfig,
			"type":  "subfinder",
		})

		configCollection.InsertOne(context.Background(), bson.M{
			"name":  "RadConfig",
			"value": constants.RadConfig,
			"type":  "rad",
		})

		// 插入指纹版本
		configCollection.InsertOne(context.Background(), bson.M{
			"name":  "FingerVersion",
			"value": "2026-01-20 23:50",
			"type":  "finger",
		})

		currentStep++
		printProgressBar(currentStep, totalSteps, "Setting timezone")

		// 创建通知配置
		_, err = configCollection.InsertOne(context.Background(), bson.M{
			"name":                          "notification",
			"dirScanNotification":           true,
			"portScanNotification":          true,
			"sensitiveNotification":         true,
			"subdomainTakeoverNotification": true,
			"pageMonNotification":           true,
			"subdomainNotification":         true,
			"vulNotification":               true,
			"type":                          "notification",
		})
		if err != nil {
			return fmt.Errorf("failed to insert notification config: %v", err)
		}
		currentStep++
		printProgressBar(currentStep, totalSteps, "Setting notification config")

		// 创建定时任务
		scheduledTasksCollection := db.Collection("ScheduledTasks")
		_, err = scheduledTasksCollection.InsertOne(context.Background(), bson.M{
			"id":      "page_monitoring",
			"name":    "Page Monitoring",
			"hour":    24,
			"node":    []string{},
			"allNode": true,
			"type":    "Page Monitoring",
			"state":   true,
		})
		if err != nil {
			return fmt.Errorf("failed to insert scheduled task: %v", err)
		}
		currentStep++
		printProgressBar(currentStep, totalSteps, "Creating scheduled tasks")

		// 创建通知集合
		err = db.CreateCollection(context.Background(), "notification")
		if err != nil {
			return fmt.Errorf("failed to create notification collection: %v", err)
		}
		currentStep++
		printProgressBar(currentStep, totalSteps, "Creating notification collection")

		// 创建字典集合
		dictionaryCollection := db.Collection("dictionary")
		// 插入目录扫描字典
		dirDict := constants.DirDict
		size := float64(len(dirDict)) / (1024 * 1024)
		result, err := dictionaryCollection.InsertOne(context.Background(), bson.M{
			"name":     "default",
			"category": "dir",
			"size":     fmt.Sprintf("%.2f", size),
		})
		if err != nil {
			return fmt.Errorf("failed to insert dir dictionary: %v", err)
		}
		// 使用GridFS存储字典内容
		if err := CreateGridFSFile(db, result.InsertedID.(primitive.ObjectID).Hex(), []byte(dirDict)); err != nil {
			return fmt.Errorf("failed to write dir dictionary: %v", err)
		}
		currentStep++
		printProgressBar(currentStep, totalSteps, "Creating dir dictionary")

		// 插入子域名字典
		domainDict := constants.DomainDict
		size = float64(len(domainDict)) / (1024 * 1024)
		result, err = dictionaryCollection.InsertOne(context.Background(), bson.M{
			"name":     "default",
			"category": "subdomain",
			"size":     fmt.Sprintf("%.2f", size),
		})
		if err != nil {
			return fmt.Errorf("failed to insert subdomain dictionary: %v", err)
		}
		if err := CreateGridFSFile(db, result.InsertedID.(primitive.ObjectID).Hex(), []byte(domainDict)); err != nil {
			return fmt.Errorf("failed to write subdomain dictionary: %v", err)
		}
		currentStep++
		printProgressBar(currentStep, totalSteps, "Creating subdomain dictionary")

		// 插入敏感信息规则
		sensitiveCollection := db.Collection("SensitiveRule")
		sensitiveData, _ := constants.GetSensitive()
		if len(sensitiveData) > 0 {
			_, err = sensitiveCollection.InsertMany(context.Background(), sensitiveData)
			if err != nil {
				return fmt.Errorf("failed to insert sensitive rules: %v", err)
			}
		}
		currentStep++
		printProgressBar(currentStep, totalSteps, "Creating sensitive rules")

		// 插入默认端口
		portCollection := db.Collection("PortDict")
		portData, _ := constants.GetPort()
		_, err = portCollection.InsertMany(context.Background(), portData)
		if err != nil {
			return fmt.Errorf("failed to insert port dictionary: %v", err)
		}
		currentStep++
		printProgressBar(currentStep, totalSteps, "Creating port dictionary")

		// 插入POC
		//pocCollection := db.Collection("PocList")
		//pocData := getPoc()
		//if len(pocData) > 0 {
		//	_, err = pocCollection.InsertMany(context.Background(), pocData)
		//	if err != nil {
		//		return fmt.Errorf("failed to insert POC: %v", err)
		//	}
		//}
		//currentStep++
		//printProgressBar(currentStep, totalSteps, "Creating POC")

		// 插入指纹规则
		fingerprintCollection := db.Collection("FingerprintRules")
		fingerprint, _ := constants.GetFingerprintData()
		if len(fingerprint) > 0 {
			_, err = fingerprintCollection.InsertMany(context.Background(), fingerprint)
			if err != nil {
				return fmt.Errorf("failed to insert fingerprint rules: %v", err)
			}
		}

		currentStep++
		printProgressBar(currentStep, totalSteps, "Creating fingerprint rules")

		// 创建默认插件
		pluginsCollection := db.Collection("plugins")
		var plgDocs []interface{}
		for _, p := range constants.Plugins {
			plgDocs = append(plgDocs, p)
		}
		_, err = pluginsCollection.InsertMany(context.Background(), plgDocs)
		if err != nil {
			return fmt.Errorf("failed to insert plugins: %v", err)
		}
		currentStep++
		printProgressBar(currentStep, totalSteps, "Creating plugins")

		// 创建默认扫描模板
		scanTemplateCollection := db.Collection("ScanTemplates")
		_, err = scanTemplateCollection.InsertOne(context.Background(), constants.ScanTemplateDefault)
		if err != nil {
			return fmt.Errorf("failed to insert scan template: %v", err)
		}
		currentStep++
		printProgressBar(currentStep, totalSteps, "Creating scan template")

		// 创建索引
		if err := createIndexes(db); err != nil {
			return fmt.Errorf("failed to create indexes: %v", err)
		}
		currentStep++
		printProgressBar(currentStep, totalSteps, "Creating indexes")

		fmt.Println() // 换行
		logger.Info("Project initialization successful")
	} else {
		// 数据库已存在，检查并更新必要的配置
		configCollection := db.Collection("config")
		var result bson.M
		err := configCollection.FindOne(context.Background(), bson.M{"name": "timezone"}).Decode(&result)
		if err != nil {
			return fmt.Errorf("failed to get timezone config: %v", err)
		}

		// 检查定时任务
		scheduledTasksCollection := db.Collection("ScheduledTasks")
		var taskResult bson.M
		err = scheduledTasksCollection.FindOne(context.Background(), bson.M{"id": "page_monitoring"}).Decode(&taskResult)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// 创建默认定时任务
				_, err = scheduledTasksCollection.InsertOne(context.Background(), bson.M{
					"id":    "page_monitoring",
					"name":  "Page Monitoring",
					"hour":  24,
					"type":  "Page Monitoring",
					"state": true,
				})
				if err != nil {
					return fmt.Errorf("failed to insert default scheduled task: %v", err)
				}
			} else {
				return fmt.Errorf("failed to check scheduled task: %v", err)
			}
		}
	}

	return nil
}

func createIndexes(db *mongo.Database) error {
	// 创建页面监控索引
	pageMonitoring := db.Collection("PageMonitoring")
	_, err := pageMonitoring.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"url", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create PageMonitoring index: %v", err)
	}

	// 创建页面监控内容索引
	pageMonitoringBody := db.Collection("PageMonitoringBody")
	_, err = pageMonitoringBody.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"md5", 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create PageMonitoringBody index: %v", err)
	}

	// 创建资产索引
	asset := db.Collection("asset")
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{"time", -1}}},
		{Keys: bson.D{{"url", 1}}},
		{Keys: bson.D{{"ip", 1}}},
		{Keys: bson.D{{"port", 1}}},
		{Keys: bson.D{{"host", 1}, {"port", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{"project", 1}, {"time", -1}}},
		{Keys: bson.D{{"project", 1}}},
		{Keys: bson.D{{"taskName", 1}}},
		{Keys: bson.D{{"rootDomain", 1}}},
		// 1.8 新增索引
		{Keys: bson.D{{"tags", 1}}},
		{Keys: bson.D{{"technologies", 1}}},
		{Keys: bson.D{{"faviconmmh3", 1}}},
		{Keys: bson.D{{"bodyhash", 1}}},
		// 1.8.1新增索引
		{Keys: bson.D{{"host", 1}}},

		{Keys: bson.D{{"rootDomain", 1}, {"time", -1}}}, // 优化按 rootDomain 查询并按 time 排序的性能
	}
	_, err = asset.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create asset indexes: %v", err)
	}

	// 创建ip资产临时集合索引
	ipAssetTmp := db.Collection("IPAssetTmp")
	indexes = []mongo.IndexModel{
		{Keys: bson.D{{"ip", 1}, {"port", 1}, {"domain", 1}}, Options: options.Index().SetUnique(true)},
	}
	_, err = ipAssetTmp.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create ipAssetTmp indexes: %v", err)
	}

	// 创建ip资产集合索引
	ipAsset := db.Collection("IPAsset")
	indexes = []mongo.IndexModel{
		{Keys: bson.D{{"time", -1}}},
		{Keys: bson.D{{"ports.server.domain", 1}}},
		{Keys: bson.D{{"ports.port", 1}}},
		{Keys: bson.D{{"ip", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{"ports.server.service", 1}}},
		{Keys: bson.D{{"ports.server.webServer", 1}}},
		{Keys: bson.D{{"rootDomain", 1}}},
		{Keys: bson.D{{"ports.server.technologies", 1}}},
		{Keys: bson.D{{"project", 1}}},
		{Keys: bson.D{{"taskName", 1}}},
	}
	_, err = ipAsset.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create ipAsset indexes: %v", err)
	}

	// 创建icon集合索引
	icon := db.Collection("icon")
	indexes = []mongo.IndexModel{
		{Keys: bson.D{{"fav3", 1}}, Options: options.Index().SetUnique(true)},
	}
	_, err = icon.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create icon indexes: %v", err)
	}

	// 创建http res body集合索引
	httpBody := db.Collection("HttpBody")
	indexes = []mongo.IndexModel{
		{Keys: bson.D{{"hash", 1}}, Options: options.Index().SetUnique(true)},
	}
	_, err = httpBody.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create httpBody indexes: %v", err)
	}

	// 创建http截图索引
	httpScreenshot := db.Collection("screenshot")
	indexes = []mongo.IndexModel{
		{Keys: bson.D{{"hash", 1}}, Options: options.Index().SetUnique(true)},
	}
	_, err = httpScreenshot.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create httpBody indexes: %v", err)
	}

	// 创建子域名索引
	subdomain := db.Collection("subdomain")
	indexes = []mongo.IndexModel{
		{Keys: bson.D{{"project", 1}}},
		{Keys: bson.D{{"taskName", 1}}},
		{Keys: bson.D{{"rootDomain", 1}}},
		{Keys: bson.D{{"time", 1}}},
	}
	_, err = subdomain.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create subdomain indexes: %v", err)
	}

	// 创建URL扫描索引
	urlScan := db.Collection("UrlScan")
	indexes = []mongo.IndexModel{
		{Keys: bson.D{{"project", 1}}},
		{Keys: bson.D{{"taskName", 1}}},
		{Keys: bson.D{{"rootDomain", 1}}},
	}
	_, err = urlScan.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create UrlScan indexes: %v", err)
	}

	// 创建敏感信息body索引
	SensitiveBody := db.Collection("SensitiveBody")
	_, err = SensitiveBody.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{"md5", 1}},
	})
	if err != nil {
		fmt.Errorf("failed to create SensitiveBody index: %v", err)
	}
	// 创建爬虫索引
	crawler := db.Collection("crawler")
	indexes = []mongo.IndexModel{
		{Keys: bson.D{{"project", 1}}},
		{Keys: bson.D{{"taskName", 1}}},
		{Keys: bson.D{{"rootDomain", 1}}},
	}
	_, err = crawler.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create crawler indexes: %v", err)
	}

	// 创建敏感信息结果索引
	sensitiveResult := db.Collection("SensitiveResult")
	indexes = []mongo.IndexModel{
		{Keys: bson.D{{"project", 1}}},
		{Keys: bson.D{{"taskName", 1}}},
		{Keys: bson.D{{"rootDomain", 1}}},
	}
	_, err = sensitiveResult.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create SensitiveResult indexes: %v", err)
	}

	// 创建目录扫描结果索引
	dirScanResult := db.Collection("DirScanResult")
	indexes = []mongo.IndexModel{
		{Keys: bson.D{{"project", 1}}},
		{Keys: bson.D{{"taskName", 1}}},
		{Keys: bson.D{{"rootDomain", 1}}},
	}
	_, err = dirScanResult.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create DirScanResult indexes: %v", err)
	}

	// 创建漏洞索引
	vulnerability := db.Collection("vulnerability")
	indexes = []mongo.IndexModel{
		{Keys: bson.D{{"project", 1}}},
		{Keys: bson.D{{"taskName", 1}}},
		{Keys: bson.D{{"rootDomain", 1}}},
		{Keys: bson.D{{"hash", 1}}},
	}
	_, err = vulnerability.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create vulnerability indexes: %v", err)
	}

	// 创建vulDetail索引
	vulnerabilityDetail := db.Collection("vulnerabilityDetail")
	indexes = []mongo.IndexModel{
		{Keys: bson.D{{"hash", 1}}},
	}
	_, err = vulnerabilityDetail.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create vulnerabilityDetail indexes: %v", err)
	}

	// 创建根域名索引
	rootDomain := db.Collection("RootDomain")
	indexes = []mongo.IndexModel{
		{Keys: bson.D{{"project", 1}}},
		{Keys: bson.D{{"taskName", 1}}},
		{Keys: bson.D{{"domain", 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{"time", 1}}},
	}
	_, err = rootDomain.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create RootDomain indexes: %v", err)
	}

	// 创建应用索引
	app := db.Collection("app")
	indexes = []mongo.IndexModel{
		{Keys: bson.D{{"project", 1}}},
		{Keys: bson.D{{"taskName", 1}}},
		{Keys: bson.D{{"time", 1}}},
		{Keys: bson.D{{"name", 1}}},
	}
	_, err = app.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create app indexes: %v", err)
	}

	// 创建MP索引
	mp := db.Collection("mp")
	indexes = []mongo.IndexModel{
		{Keys: bson.D{{"project", 1}}},
		{Keys: bson.D{{"taskName", 1}}},
		{Keys: bson.D{{"time", 1}}},
		{Keys: bson.D{{"name", 1}}},
	}
	_, err = mp.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create mp indexes: %v", err)
	}

	// 创建task相关索引
	task := db.Collection("task")
	indexes = []mongo.IndexModel{
		{Keys: bson.D{{"creatTime", -1}}},
		{Keys: bson.D{{"name", 1}}},
		{Keys: bson.D{{"progress", 1}}},
	}
	_, err = task.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		return fmt.Errorf("failed to create task indexes: %v", err)
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
