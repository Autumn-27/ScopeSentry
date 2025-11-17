package constants

// 角色定义
const (
	RoleSuperAdmin = "super_admin" // 超级管理员
	RoleAdmin      = "admin"       // 管理员
	RoleUser       = "user"        // 普通用户
	RoleGuest      = "guest"       // 访客
)

// 权限前缀定义
const (
	PrefixUser    = "user"    // 用户管理
	PrefixRole    = "role"    // 角色管理
	PrefixAsset   = "asset"   // 资产管理
	PrefixDomain  = "domain"  // 子域名管理
	PrefixApp     = "app"     // 应用管理
	PrefixTask    = "task"    // 任务管理
	PrefixProject = "project" // 项目管理
	PrefixSystem  = "system"  // 系统管理
)

// 操作类型定义
const (
	ActionCreate = "create" // 创建
	ActionRead   = "read"   // 读取
	ActionUpdate = "update" // 更新
	ActionDelete = "delete" // 删除
	ActionList   = "list"   // 列表
	ActionExport = "export" // 导出
	ActionImport = "import" // 导入
)

// 权限定义
const (
	// 用户管理权限
	PermissionUserCreate = PrefixUser + ":" + ActionCreate
	PermissionUserRead   = PrefixUser + ":" + ActionRead
	PermissionUserUpdate = PrefixUser + ":" + ActionUpdate
	PermissionUserDelete = PrefixUser + ":" + ActionDelete
	PermissionUserList   = PrefixUser + ":" + ActionList

	// 角色管理权限
	PermissionRoleCreate = PrefixRole + ":" + ActionCreate
	PermissionRoleRead   = PrefixRole + ":" + ActionRead
	PermissionRoleUpdate = PrefixRole + ":" + ActionUpdate
	PermissionRoleDelete = PrefixRole + ":" + ActionDelete
	PermissionRoleList   = PrefixRole + ":" + ActionList

	// 资产管理权限
	PermissionAssetCreate = PrefixAsset + ":" + ActionCreate
	PermissionAssetRead   = PrefixAsset + ":" + ActionRead
	PermissionAssetUpdate = PrefixAsset + ":" + ActionUpdate
	PermissionAssetDelete = PrefixAsset + ":" + ActionDelete
	PermissionAssetList   = PrefixAsset + ":" + ActionList
	PermissionAssetExport = PrefixAsset + ":" + ActionExport
	PermissionAssetImport = PrefixAsset + ":" + ActionImport

	// 子域名管理权限
	PermissionDomainCreate = PrefixDomain + ":" + ActionCreate
	PermissionDomainRead   = PrefixDomain + ":" + ActionRead
	PermissionDomainUpdate = PrefixDomain + ":" + ActionUpdate
	PermissionDomainDelete = PrefixDomain + ":" + ActionDelete
	PermissionDomainList   = PrefixDomain + ":" + ActionList
	PermissionDomainExport = PrefixDomain + ":" + ActionExport
	PermissionDomainImport = PrefixDomain + ":" + ActionImport

	// 应用管理权限
	PermissionAppCreate = PrefixApp + ":" + ActionCreate
	PermissionAppRead   = PrefixApp + ":" + ActionRead
	PermissionAppUpdate = PrefixApp + ":" + ActionUpdate
	PermissionAppDelete = PrefixApp + ":" + ActionDelete
	PermissionAppList   = PrefixApp + ":" + ActionList

	// 任务管理权限
	PermissionTaskCreate = PrefixTask + ":" + ActionCreate
	PermissionTaskRead   = PrefixTask + ":" + ActionRead
	PermissionTaskUpdate = PrefixTask + ":" + ActionUpdate
	PermissionTaskDelete = PrefixTask + ":" + ActionDelete
	PermissionTaskList   = PrefixTask + ":" + ActionList
	PermissionTaskExport = PrefixTask + ":" + ActionExport

	// 项目管理权限
	PermissionProjectCreate = PrefixProject + ":" + ActionCreate
	PermissionProjectRead   = PrefixProject + ":" + ActionRead
	PermissionProjectUpdate = PrefixProject + ":" + ActionUpdate
	PermissionProjectDelete = PrefixProject + ":" + ActionDelete
	PermissionProjectList   = PrefixProject + ":" + ActionList

	// 系统管理权限
	PermissionSystemConfig = PrefixSystem + ":config"
	PermissionSystemLog    = PrefixSystem + ":log"
	PermissionSystemBackup = PrefixSystem + ":backup"
)

// 角色权限映射
var RolePermissions = map[string][]string{
	RoleSuperAdmin: {
		// 用户管理
		PermissionUserCreate,
		PermissionUserRead,
		PermissionUserUpdate,
		PermissionUserDelete,
		PermissionUserList,

		// 角色管理
		PermissionRoleCreate,
		PermissionRoleRead,
		PermissionRoleUpdate,
		PermissionRoleDelete,
		PermissionRoleList,

		// 资产管理
		PermissionAssetCreate,
		PermissionAssetRead,
		PermissionAssetUpdate,
		PermissionAssetDelete,
		PermissionAssetList,
		PermissionAssetExport,
		PermissionAssetImport,

		// 子域名管理
		PermissionDomainCreate,
		PermissionDomainRead,
		PermissionDomainUpdate,
		PermissionDomainDelete,
		PermissionDomainList,
		PermissionDomainExport,
		PermissionDomainImport,

		// 应用管理
		PermissionAppCreate,
		PermissionAppRead,
		PermissionAppUpdate,
		PermissionAppDelete,
		PermissionAppList,

		// 任务管理
		PermissionTaskCreate,
		PermissionTaskRead,
		PermissionTaskUpdate,
		PermissionTaskDelete,
		PermissionTaskList,
		PermissionTaskExport,

		// 项目管理
		PermissionProjectCreate,
		PermissionProjectRead,
		PermissionProjectUpdate,
		PermissionProjectDelete,
		PermissionProjectList,

		// 系统管理
		PermissionSystemConfig,
		PermissionSystemLog,
		PermissionSystemBackup,
	},
	RoleAdmin: {
		// 用户管理
		PermissionUserRead,
		PermissionUserList,

		// 资产管理
		PermissionAssetRead,
		PermissionAssetList,
		PermissionAssetExport,

		// 子域名管理
		PermissionDomainRead,
		PermissionDomainList,
		PermissionDomainExport,

		// 应用管理
		PermissionAppRead,
		PermissionAppList,

		// 任务管理
		PermissionTaskCreate,
		PermissionTaskRead,
		PermissionTaskList,
		PermissionTaskExport,

		// 项目管理
		PermissionProjectCreate,
		PermissionProjectRead,
		PermissionProjectUpdate,
		PermissionProjectList,
	},
	RoleUser: {
		// 资产管理
		PermissionAssetRead,
		PermissionAssetList,

		// 子域名管理
		PermissionDomainRead,
		PermissionDomainList,

		// 应用管理
		PermissionAppRead,
		PermissionAppList,

		// 任务管理
		PermissionTaskRead,
		PermissionTaskList,

		// 项目管理
		PermissionProjectRead,
		PermissionProjectList,
	},
	RoleGuest: {
		// 只读权限
		PermissionAssetRead,
		PermissionDomainRead,
		PermissionAppRead,
		PermissionTaskRead,
		PermissionProjectRead,
	},
} 