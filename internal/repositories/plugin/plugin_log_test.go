package plugin

import (
	"context"
	"strings"
	"testing"

	redisdriver "github.com/redis/go-redis/v9"
)

// fakeLogRedis 是 pluginLogRedis 的内存实现，用于在不依赖真实 Redis 的情况下验证日志读取逻辑。
type fakeLogRedis struct {
	types map[string]string
	lists map[string][]string
	sets  map[string][]string
}

func (f *fakeLogRedis) Type(_ context.Context, key string) *redisdriver.StatusCmd {
	cmd := redisdriver.NewStatusCmd(context.Background())
	if tp, ok := f.types[key]; ok {
		cmd.SetVal(tp)
	} else {
		cmd.SetVal("none")
	}
	return cmd
}

func (f *fakeLogRedis) LRange(_ context.Context, key string, _, _ int64) *redisdriver.StringSliceCmd {
	cmd := redisdriver.NewStringSliceCmd(context.Background())
	cmd.SetVal(f.lists[key])
	return cmd
}

func (f *fakeLogRedis) SMembers(_ context.Context, key string) *redisdriver.StringSliceCmd {
	cmd := redisdriver.NewStringSliceCmd(context.Background())
	cmd.SetVal(f.sets[key])
	return cmd
}

func (f *fakeLogRedis) Del(_ context.Context, _ ...string) *redisdriver.IntCmd {
	cmd := redisdriver.NewIntCmd(context.Background())
	cmd.SetVal(1)
	return cmd
}

func (f *fakeLogRedis) Keys(_ context.Context, _ string) *redisdriver.StringSliceCmd {
	cmd := redisdriver.NewStringSliceCmd(context.Background())
	cmd.SetVal(nil)
	return cmd
}

func newFakeLogRepository(f *fakeLogRedis) *repository {
	return &repository{redisClient: f}
}

// 新版日志以 list 存储，必须按写入顺序返回。
func TestGetLogsReadsListKey(t *testing.T) {
	repo := newFakeLogRepository(&fakeLogRedis{
		types: map[string]string{"logs:server_plugins:abc": "list"},
		lists: map[string][]string{"logs:server_plugins:abc": {"[1] first", "[2] second"}},
	})

	got, err := repo.GetLogs(context.Background(), "logs:server_plugins:abc")
	if err != nil {
		t.Fatalf("GetLogs returned error: %v", err)
	}
	if want := "[1] first\n[2] second"; got != want {
		t.Fatalf("GetLogs = %q, want %q", got, want)
	}
}

// 历史版本的日志 key 是 set，LRange 会报 WRONGTYPE（issue #290），这里必须兼容读取。
func TestGetLogsReadsLegacySetKey(t *testing.T) {
	repo := newFakeLogRepository(&fakeLogRedis{
		types: map[string]string{"logs:plugins:httpx:abc": "set"},
		sets:  map[string][]string{"logs:plugins:httpx:abc": {"[2] second", "[1] first"}},
	})

	got, err := repo.GetLogs(context.Background(), "logs:plugins:httpx:abc")
	if err != nil {
		t.Fatalf("GetLogs returned error for legacy set key: %v", err)
	}
	if want := "[1] first\n[2] second"; got != want {
		t.Fatalf("GetLogs = %q, want %q", got, want)
	}
}

// key 不存在时应返回空日志而不是错误。
func TestGetLogsMissingKeyReturnsEmpty(t *testing.T) {
	repo := newFakeLogRepository(&fakeLogRedis{})

	got, err := repo.GetLogs(context.Background(), "logs:plugins:httpx:missing")
	if err != nil {
		t.Fatalf("GetLogs returned error for missing key: %v", err)
	}
	if got != "" {
		t.Fatalf("GetLogs = %q, want empty string", got)
	}
}

// 其它类型（例如 hash）应返回带类型的明确错误，而不是 Redis 的 WRONGTYPE 原始报错。
func TestGetLogsRejectsUnsupportedType(t *testing.T) {
	repo := newFakeLogRepository(&fakeLogRedis{
		types: map[string]string{"logs:plugins:httpx:abc": "hash"},
	})

	_, err := repo.GetLogs(context.Background(), "logs:plugins:httpx:abc")
	if err == nil {
		t.Fatal("GetLogs accepted a hash key, want error")
	}
	if !strings.Contains(err.Error(), "unsupported redis type") || !strings.Contains(err.Error(), "hash") {
		t.Fatalf("GetLogs error = %v, want it to name the unsupported type", err)
	}
}

func TestGetLogsEmptyKeyReturnsEmpty(t *testing.T) {
	repo := newFakeLogRepository(&fakeLogRedis{})

	got, err := repo.GetLogs(context.Background(), "")
	if err != nil || got != "" {
		t.Fatalf("GetLogs(\"\") = %q, %v; want empty string, nil", got, err)
	}
}
