package constants

// 任务状态定义（与前端 task.status 字段对应）
//
// 历史实现只有 运行中/暂停/完成 三种状态，任务一入库就被写成 运行中，
// 但实际上扫描节点还没有领取任务，前端会把排队中的任务显示成“运行中”。
// 这里补上 待开始 状态：任务入库时为 待开始，出现扫描进度后才转为 运行中。
const (
	TaskStatusPending  = 0 // 待开始：任务已创建，扫描节点尚未产出进度
	TaskStatusRunning  = 1 // 运行中：已有扫描进度（或手动恢复/重测）
	TaskStatusPaused   = 2 // 暂停：被用户手动停止
	TaskStatusFinished = 3 // 运行完成：进度达到 100%
)
