package server

import "github.com/Croohand/mapreduce/common/responses"

func isAlive() *responses.SlaveStatus {
	return &responses.SlaveStatus{true, "slave"}
}
