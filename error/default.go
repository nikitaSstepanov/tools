package e

var (
	InternalErr = New("Something going wrong...", Internal)
	BadInputErr = New("Bad input.", BadInput)
)
