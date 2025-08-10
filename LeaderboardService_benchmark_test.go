package leaderboard

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ... existing test functions ...

// 基准测试：插入N个用户的性能测试
func BenchmarkLeaderboard_InsertUsers(b *testing.B) {
	// 初始化随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	lb := NewLeaderboard()
	baseTimestamp := int32(time.Now().Unix())
	userCount := 1_000_000 // 100万用户

	// 预热：创建测试数据（避免数据生成影响性能测量）
	playerIDs := make([]string, userCount)
	scores := make([]int32, userCount)
	for i := 0; i < userCount; i++ {
		playerIDs[i] = fmt.Sprintf("player-%07d", i)
		scores[i] = int32(r.Intn(1_000_000)) // 0-999999随机分数
	}

	b.ResetTimer()   // 重置计时器，排除数据准备时间
	b.ReportAllocs() // 启用内存分配统计

	// 批量插入测试数据
	for i := 0; i < userCount; i++ {
		lb.updateScore(playerIDs[i], scores[i], baseTimestamp-int32(i%86400)) // 时间戳略微随机化
	}

	b.StopTimer()

	// 验证数据完整性
	assert.Equal(b, int64(userCount), lb.zset.Length(), "插入用户数量不匹配")

	// 测试查询性能 - Top 100
	b.StartTimer()
	top100 := lb.getTopN(100)
	b.StopTimer()
	assert.Len(b, top100, 100, "Top 100查询结果数量错误")

	// 测试随机查询性能
	b.StartTimer()
	for i := 0; i < 1000; i++ {
		// 随机选择1000个用户查询排名
		idx := r.Intn(userCount)
		lb.getPlayerRank(playerIDs[idx])
	}
	b.StopTimer()

	// 测试查询性能 - 排名范围
	b.StartTimer()
	for i := 0; i < 1000; i++ {
		// 随机选择1000个用户查询排名范围
		idx := r.Intn(userCount)
		lb.getPlayerRankRange(playerIDs[idx], 100)
	}
	b.StopTimer()

}
