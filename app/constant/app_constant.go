package constant

type ResponseStatus int

const (
	Success ResponseStatus = iota + 1
	InvalidRequest
	Unauthorized
	DataNotFound
	Conflict
	UnknownError
)

func (r ResponseStatus) GetResponseStatus() string {
	return [...]string{"SUCCESS", "INVALID_REQUEST", "UNAUTHORIZED", "DATA_NOT_FOUND", "CONFLICT", "UNKNOWN_ERROR"}[r-1]
}

func (r ResponseStatus) GetResponseMessage() string {
	return [...]string{"Success", "Invalid Request", "Unauthorized", "Data Not Found", "Conflict", "Unknown Error"}[r-1]
}
