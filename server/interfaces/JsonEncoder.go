package interfaces

//JsonEncoder is an interface for encoding json http responses
type JsonEncoder interface {
	Encode(v interface{}) ([]byte, error)
	EncodeResponse(status int, v interface{}) error
}