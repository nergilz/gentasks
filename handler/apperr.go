package handler

type WorkerError struct {
	Err error
	Msg string
}

func (we WorkerError) Error() string {
	return we.Msg
}
