package perror

type Type int

const (
	TypeUnspecified Type = iota
	TypeNotFound
	TypeAlreadyExists
	TypeNotExists
	TypeInvalidParam
	TypeTooManyRequests
	TypeNoRowsUpdated
	TypeInvalidStatus
)
