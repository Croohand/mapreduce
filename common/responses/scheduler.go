package responses

import (
	"fmt"
	"time"
)

type PreparedOperation struct {
	ReadTxId, WriteTxId string
}

type OperationStatus struct {
	MappersDone  int
	MappersAll   int
	ReducersDone int
	ReducersAll  int
	Error        string
	Started      time.Time
}

func (opStatus OperationStatus) Done() bool {
	return opStatus.MappersDone == opStatus.MappersAll && opStatus.ReducersDone == opStatus.ReducersAll
}

func (opStatus OperationStatus) Failed() bool {
	return len(opStatus.Error) > 0
}

func (opStatus OperationStatus) String() string {
	done := fmt.Sprintf("Mappers done: %d out of %d\nReducers done: %d out of %d\nExecution time: %.3f\n",
		opStatus.MappersDone,
		opStatus.MappersAll,
		opStatus.ReducersDone,
		opStatus.ReducersAll,
		time.Since(opStatus.Started).Seconds())
	if opStatus.Failed() {
		return fmt.Sprintf("Operation failed with error: %s\n", opStatus.Error) + done
	}
	if opStatus.Done() {
		return "Operation is done\n" + done
	}
	return "Operation is running...\n" + done
}
