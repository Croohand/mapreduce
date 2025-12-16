package server

import "github.com/Croohand/mapreduce/common/responses"

func isAlive() *responses.WorkerStatus {
	return &responses.WorkerStatus{true, "worker"}
}
