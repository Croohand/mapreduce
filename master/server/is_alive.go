package server

import "github.com/Croohand/mapreduce/common/responses"

func isAlive() *responses.MasterStatus {
	return &responses.MasterStatus{true, "master"}
}
