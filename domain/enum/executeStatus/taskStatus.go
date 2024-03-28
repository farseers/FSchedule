package executeStatus

type Enum int

const (
	None    Enum = iota //  未开始
	Working             //  执行中（已下发给客户端）
	Success             //  成功
	Fail                //  失败
)

func (e Enum) String() string {
	switch e {
	case None:
		return "None"
	case Working:
		return "Working"
	case Success:
		return "Success"
	case Fail:
		return "Fail"
	}
	return "None"
}
