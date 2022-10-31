package conf

const (
	RequestSuccess     = 200
	RequestRedirect    = 302
	RequestTokenExpire = 401
	RequestFailed      = 402
)

// 操作日志-结果类型
const OperateFail = 0
const OperateSuccess = 1

// 日志操作类型
const LogOperateTypeAdd = 1
const LogOperateTypeReceive = 2
const LogOperateTypeAppeal = 3
const LogOperateTypeNotify = 4
const LogOperateTypeCancel = 5
const LogOperateTypeEdit = 6
const LogOperateTypeConfirm = 7

const (
	DayLayout      = "02/01/2006"
	DateTimeLayout = "02/01/2006 15:04:05"
)
