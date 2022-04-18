package errors

type BadRequestError struct {
	Message string
}

func (err BadRequestError) Error() string {
	return err.Message
}

type ProcessQueryFailedError struct {
	Message string
}

func (err ProcessQueryFailedError) Error() string {
	return err.Message
}

type SubscribeTimeoutError struct {
	Message string
}

func (err SubscribeTimeoutError) Error() string {
	return err.Message
}
