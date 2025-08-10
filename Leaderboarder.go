package leaderboard

import (
	"github.com/liyiheng/zset"
)

const EndTime = 2145916800 // 2038-01-01 00:00:00

type Leaderboarder interface {
	// UpdateScore 更新玩家分数
	UpdateScore(playerId string, score int32, timestamp int32)
	// GetPlayerRank 获取玩家当前排名
	GetPlayerRank(playerId string) RankInfo
	// GetTopN 获取排行榜前N名
	GetTopN(n int32) []RankInfo
	// GetPlayerRankRange 获取玩家周边排名
	GetPlayerRankRange(playerId string, rangeNum int32) []RankInfo
}

type RankInfo struct {
	PlayerId string
	Rank     int32
	Score    int32
}

type Leaderboard struct {
	//zset.SortedSet
	zset *zset.SortedSet[string]
}

func NewLeaderboard() *Leaderboard {
	return &Leaderboard{
		zset: zset.New[string](),
	}
}

func (lb *Leaderboard) UpdateScore(playerId string, score int32, timestamp int32) {
	newScore := float64(score) + float64(EndTime-timestamp)/1000000000
	lb.zset.Set(newScore, playerId)
}

func (lb *Leaderboard) GetPlayerRank(playerId string) RankInfo {
	rank, score := lb.zset.GetRank(playerId, true)
	if rank == -1 {
		return RankInfo{
			PlayerId: playerId,
			Rank:     0,
			Score:    0,
		}
	}
	return RankInfo{
		PlayerId: playerId,
		Rank:     int32(rank) + 1,
		Score:    int32(score),
	}
}

func (lb *Leaderboard) GetTopN(n int32) []RankInfo {
	var rankInfos []RankInfo
	rank := int32(1)
	lb.zset.RevRange(0, int64(n-1), func(score float64, id string) {
		rankInfos = append(rankInfos, RankInfo{
			PlayerId: id,
			Score:    int32(score),
			Rank:     rank,
		})
		rank++
	})
	return rankInfos
}

func (lb *Leaderboard) GetPlayerRankRange(playerId string, rangeNum int32) []RankInfo {
	myRank, _ := lb.zset.GetRank(playerId, true)
	if myRank == -1 {
		return nil
	}

	start := myRank - int64(rangeNum)/2
	end := start + int64(rangeNum) - 1

	rankLen := lb.zset.Length()
	if end >= rankLen {
		start -= end - rankLen + 1
		end = rankLen - 1
	}
	if start < 0 {
		end -= start
		start = 0
	}

	var rankInfos []RankInfo
	rank := int32(start) + 1
	lb.zset.RevRange(start, end, func(score float64, id string) {
		rankInfos = append(rankInfos, RankInfo{
			PlayerId: id,
			Score:    int32(score),
			Rank:     rank,
		})
		rank++
	})
	return rankInfos
}
