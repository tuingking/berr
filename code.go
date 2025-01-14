package berr

type ErrCode uint32

// preset error codes
const (
	// internal
	CodeInternal ErrCode = iota
	CodeJsonMarshal
	CodeJsonUnmarshal

	// business error
	CodeBusiness ErrCode = iota + 1000
	CodeValidation

	// db error
	CodeDB ErrCode = iota + 2000
	CodeSQLRead
	CodeSQLInsert
	CodeSQLUpdate
	CodeSQLDelete
	CodeSQLLastInsertID

	// external error
	CodeExternal ErrCode = iota + 3000
	CodeHttpRequest
	CodeKafkaPublish
	CodeKafkaConsume
)
