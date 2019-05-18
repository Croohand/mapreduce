package server

import "github.com/Croohand/mapreduce/common/responses"

func getMrConfig() *responses.MrConfig {
	return &responses.MrConfig{MaxRowLength: 1 << 15, BlockSize: 1 << 25, ReplicationFactor: 3}
}
