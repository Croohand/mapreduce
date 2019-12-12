package responses

import "github.com/Croohand/mapreduce/common/fsutil"

type MasterStatus struct {
	Alive bool
	Type  string
	State string
}

type TransactionStatus struct {
	Alive bool
}

type StartedTransaction struct {
	Id string
}

type MrConfig struct {
	MaxRowLength      int
	BlockSize         int
	ReplicationFactor int
}

type FileStatus struct {
	Exists bool
}

type PathBlocks []fsutil.BlockInfoEx

type ListedFiles []string
