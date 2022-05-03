package errors

type BadRequestError struct {
	Message string
}

func (err BadRequestError) Error() string {
	return err.Message
}

type InvalidQueryError struct {
	Query string
}

func (err InvalidQueryError) Error() string {
	return "Invalid query: " + err.Query
}

type ProcessQueryFailedError struct {
	Message string
}

func (err ProcessQueryFailedError) Error() string {
	return err.Message
}

type SubscribeTimeoutError struct{}

func (err SubscribeTimeoutError) Error() string {
	return "subscribe timeout"
}

type MaxRetryCountExceededError struct{}

func (err MaxRetryCountExceededError) Error() string {
	return "max retry count exceeded"
}
