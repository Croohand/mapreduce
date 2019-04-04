package responses

type SlaveStatus struct {
	Alive bool
	Type  string
}

type BlockStatus struct {
	Exists bool
}
