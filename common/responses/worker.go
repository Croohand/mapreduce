package responses

type WorkerStatus struct {
	Alive bool
	Type  string
}

type BlockStatus struct {
	Exists bool
}
