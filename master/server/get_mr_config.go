package server

import "github.com/Croohand/mapreduce/common/responses"

func getMrConfig() *responses.MrConfig {
	return &responses.MrConfig{MaxRowLength: 1 << 18, BlockSize: 1 << 22, ReplicationFactor: 3}
}
