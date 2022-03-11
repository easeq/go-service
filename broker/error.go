package broker

type Code uint32

const (
	// UNKNOWN indicates the error could not be determined
	UNKNOWN = 0
	// OK indicates operation doesn't contain any error
	OK Code = 1
	// ERR indicates the operation contains an error
	ERR Code = 2
	// WARN indicates the operation may have an error
	WARN Code = 3
)

// BrokerError returns a broker status error
type BrokerError func(msg string) error

var (
	Ok      = WithCode(OK)
	Err     = WithCode(ERR)
	Warn    = WithCode(WARN)
	Unknown = WithCode(UNKNOWN)
)

type BrokerStatus struct {
	Code    Code
	Message string
}

func NewBrokerStatus(code Code, msg string) *BrokerStatus {
	return &BrokerStatus{code, msg}
}

func (b *BrokerStatus) Error() string {
	return b.Message
}

func WithCode(code Code) BrokerError {
	return func(msg string) error {
		return NewBrokerStatus(code, msg)
	}
}
