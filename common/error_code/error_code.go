package error_code

const (
	Success         = 0
	Fail            = 1
	ServerBusy      = 2
	ParamError      = 3
)

func GetErrorString(error_code int) string {
	switch error_code {
	case Success: return "success"
	case Fail: return "fail"
	case ServerBusy: return "server busy"
	case ParamError: return "param error"
	default:
		return "unknown error code"
	}
}