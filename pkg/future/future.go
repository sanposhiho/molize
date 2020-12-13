package future

type Future struct {
	status FutureStatus
}

type FutureStatus int

const (
	FUTURE_STATUS_PENDING  FutureStatus = iota
	FUTURE_STATUS_RUNNING  FutureStatus = iota
	FUTURE_STATUS_FINISHED FutureStatus = iota
)
