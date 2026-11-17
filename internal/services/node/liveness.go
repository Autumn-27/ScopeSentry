package node

import (
	"os"
	"strconv"
	"time"
)

// nodeTimeLayout 节点写入 Redis 的时间格式（本地时间字符串，无时区信息）
const nodeTimeLayout = "2006-01-02 15:04:05"

// DefaultNodeTimeout 默认超过多久没有心跳判定离线
const DefaultNodeTimeout = 300 * time.Second

// NodeTimeout 允许用 NODE_TIMEOUT_SECONDS 覆盖（节点多时可以把窗口放宽）
func NodeTimeout() time.Duration {
	if v := os.Getenv("NODE_TIMEOUT_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return time.Duration(n) * time.Second
		}
	}
	return DefaultNodeTimeout
}

// Liveness 节点存活判定结果
type Liveness struct {
	Alive       bool
	TimeDiff    time.Duration // 服务器时间 - 节点上报时间（按最接近的解释）
	ParseFailed bool          // 时间字符串无法解析（此时不判离线，避免误杀）
	ClockSkew   bool          // 节点时钟比服务器快（TimeDiff 为负）
	UsedUTC     bool          // 判定时用了 UTC 解释（节点与服务器时区不一致）
}

// ClassifyLiveness 根据节点状态与心跳时间判断是否在线。
//
// 修 #254：节点上报的是「本地时间字符串」，服务器按自身时区解析；节点与服务器
// 时区/时钟不一致时，差值会大得离谱（甚至为负），导致节点被误判离线并写死成
// state=3。这里对 loc 与 UTC 两种解释各算一次，取与服务器时间最接近的那个：
//   - 只差一点点（< timeout）→ 在线
//   - 节点时间在未来（|diff| 很小）→ 在线，标记 ClockSkew
//   - 两种解释都超过 timeout → 离线
//
// 解析失败时保留节点（Alive=true + ParseFailed），不让一个格式问题把节点从列表里抹掉。
func ClassifyLiveness(state, updateTime string, now time.Time, loc *time.Location, timeout time.Duration) Liveness {
	if state != "1" {
		return Liveness{Alive: false}
	}
	if loc == nil {
		loc = time.UTC
	}

	best, bestAbs, ok, usedUTC := time.Duration(0), time.Duration(1<<62), false, false
	for _, candidate := range []*time.Location{loc, time.UTC} {
		t, err := time.ParseInLocation(nodeTimeLayout, updateTime, candidate)
		if err != nil {
			continue
		}
		diff := now.Sub(t)
		abs := diff
		if abs < 0 {
			abs = -abs
		}
		if abs < bestAbs {
			best, bestAbs, ok = diff, abs, true
			usedUTC = candidate == time.UTC && loc != time.UTC
		}
	}
	if !ok {
		return Liveness{Alive: true, ParseFailed: true}
	}

	res := Liveness{TimeDiff: best, UsedUTC: usedUTC, ClockSkew: best < 0}
	res.Alive = bestAbs <= timeout
	return res
}
