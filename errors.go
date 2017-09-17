package slaves

import "errors"

var (
	errworkIsNil      = errors.New("error: work is nil value")
	errAlreadyRunning = errors.New("error: pool is already running")
	errFuncNil        = errors.New("error: specified function is nil")
)
