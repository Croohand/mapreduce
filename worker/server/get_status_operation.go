package server

import (
	"errors"

	"github.com/Croohand/mapreduce/common/responses"
)

func getStatusOperation(txId string) (*responses.OperationStatus, error) {
	opStatus, has := opStatuses[txId]
	if !has {
		return nil, errors.New("No operation with id " + txId)
	}
	return opStatus, nil
}
