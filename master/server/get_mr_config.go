package server

import "github.com/Croohand/mapreduce/common/responses"

func getMrConfig() *responses.MrConfig {
	return &responses.MrConfig{BlockSize: 1 << 24, ReplicationFactor: 3, MinReplicationFactor: 2}
}
