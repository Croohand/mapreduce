package server

import "github.com/Croohand/mapreduce/common/responses"

func isAlive() *responses.CoordinatorStatus {
	return &responses.CoordinatorStatus{true, "coordinator", state}
}
