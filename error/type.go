package e

type JsonError struct {
	Error string `json:"error"`
}

const (
	Internal StatusType = iota
	NotFound
	BadInput
	Conflict
	Forbidden
	Unauthorize
)

type StatusType int
