package responses

type MasterStatus struct {
	Alive bool
	Type  string
}

type TransactionStatus struct {
	Alive bool
}

type MrConfig struct {
	BlockSize            int
	ReplicationFactor    int
	MinReplicationFactor int
}

type FileStatus struct {
	Exists bool
}
