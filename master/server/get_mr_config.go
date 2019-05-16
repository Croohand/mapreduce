package server

import "github.com/Croohand/mapreduce/common/responses"

func getMrConfig() *responses.MrConfig {
	return &responses.MrConfig{MaxRowLength: 1 << 14, BlockSize: 1 << 20, ReplicationFactor: 3}
}
