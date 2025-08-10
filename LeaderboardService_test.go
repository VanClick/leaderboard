package leaderboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeaderboard_UpdateScore(t *testing.T) {
	lb := NewLeaderboard()
	playerId := "player1"
	initialScore := int32(100)
	timestamp := int32(1620000000)

	// 测试更新分数
	lb.UpdateScore(playerId, initialScore, timestamp)

	// 验证分数是否正确更新
	info := lb.GetPlayerRank(playerId)
	assert.Equal(t, playerId, info.PlayerId)
	assert.Equal(t, int32(1), info.Rank)
	assert.Equal(t, initialScore, info.Score)

	// 测试更新更高分数
	newScore := int32(200)
	newTimestamp := int32(1630000000)
	lb.UpdateScore(playerId, newScore, newTimestamp)
	info = lb.GetPlayerRank(playerId)
	assert.Equal(t, newScore, info.Score)
}

func TestLeaderboard_GetPlayerRank(t *testing.T) {
	lb := NewLeaderboard()
	playerId := "player1"

	// 测试不存在的玩家
	info := lb.GetPlayerRank(playerId)
	assert.Equal(t, int32(0), info.Rank)
	assert.Equal(t, int32(0), info.Score)

	// 添加玩家后测试
	lb.UpdateScore(playerId, 100, 1620000000)
	info = lb.GetPlayerRank(playerId)
	assert.Equal(t, int32(1), info.Rank)
	assert.Equal(t, int32(100), info.Score)
}

func TestLeaderboard_GetTopN(t *testing.T) {
	lb := NewLeaderboard()

	// 添加测试数据
	lb.UpdateScore("player1", 300, 1620000000)
	lb.UpdateScore("player2", 200, 1620000000)
	lb.UpdateScore("player3", 100, 1620000000)
	lb.UpdateScore("player4", 200, 1630000000)
	lb.UpdateScore("player5", 200, 1610000000)

	// 测试获取前2名
	result := lb.GetTopN(5)
	assert.Len(t, result, 5)
	assert.Equal(t, "player1", result[0].PlayerId)
	assert.Equal(t, int32(300), result[0].Score)
	assert.Equal(t, "player5", result[1].PlayerId)
	assert.Equal(t, int32(200), result[1].Score)
	assert.Equal(t, "player2", result[2].PlayerId)
	assert.Equal(t, int32(200), result[2].Score)
	assert.Equal(t, "player4", result[3].PlayerId)
	assert.Equal(t, int32(200), result[3].Score)
	assert.Equal(t, "player3", result[4].PlayerId)
	assert.Equal(t, int32(100), result[4].Score)

	// 测试获取超过实际数量的情况
	result = lb.GetTopN(10)
	assert.Len(t, result, 5)
}

func TestLeaderboard_GetPlayerRankRange(t *testing.T) {
	lb := NewLeaderboard()

	// 添加测试数据
	lb.UpdateScore("player1", 300, 1620000000)
	lb.UpdateScore("player2", 200, 1620000000)
	lb.UpdateScore("player3", 100, 1620000000)
	lb.UpdateScore("player4", 200, 1630000000)
	lb.UpdateScore("player5", 200, 1610000000)

	// 测试玩家周边排名
	result := lb.GetPlayerRankRange("player5", 3)
	assert.Len(t, result, 3)
	assert.Equal(t, "player1", result[0].PlayerId)
	assert.Equal(t, "player5", result[1].PlayerId)
	assert.Equal(t, "player2", result[2].PlayerId)

	// 测试边界情况
	result = lb.GetPlayerRankRange("player5", 10)
	assert.Len(t, result, 5)
	assert.Equal(t, "player1", result[0].PlayerId)
	assert.Equal(t, "player5", result[1].PlayerId)
	assert.Equal(t, "player2", result[2].PlayerId)
	assert.Equal(t, "player4", result[3].PlayerId)
	assert.Equal(t, "player3", result[4].PlayerId)

	// 测试边界情况（玩家不存在）
	result = lb.GetPlayerRankRange("nonexistent", 3)
	assert.Nil(t, result)
}
