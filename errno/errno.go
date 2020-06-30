package errno

var (
	// Common errors
	OK                        = &Errno{Code: 0, Result: "OK"}
	RequestBodyParseError     = &Errno{Code: 90001, Result: "Request body parse error"}
	InternalServerError       = &Errno{Code: 90002, Result: "Internal server error"}
	ActivityCodeRequiredError = &Errno{Code: 90003, Result: "Activity code required error"}
	CardRequiredError         = &Errno{Code: 90004, Result: "Card required error"}
	GroupCodeRequiredError    = &Errno{Code: 90005, Result: "Group code required error"}
)

type Errno struct {
	Code   int    `json:"code"`
	Result string `json:"result"`
}

func (err Errno) Error() string {
	return err.Result
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Result
	}
	switch typed := err.(type) {
	case *Errno:
		return typed.Code, typed.Result

	}
	return InternalServerError.Code, err.Error()

}
