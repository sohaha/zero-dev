package parse

type errType uint

const (
	ErrException errType = iota + 1
	ErrModuleAlreadyExists
	ErrNotMigration

	errCount
)

var errDescriptions = [...]string{
	ErrException:           "异常",
	ErrModuleAlreadyExists: "模型名称已存在",
	ErrNotMigration:        "不支持表迁移",
}

var _ = [1]int{}[len(errDescriptions)-int(errCount)]

func (e errType) Error() string {
	if int(e) > len(errDescriptions) {
		return "未知错误"
	}

	return errDescriptions[e]
}
