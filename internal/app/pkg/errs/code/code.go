//go:generate stringer -type=ErrorCode --linecomment

package code

type ErrorCode int

const (
	Unexpected            ErrorCode = 1000 + iota // 意外错误
	OutOfRange                                    // 数组索引越界
	SignFailed                                    // 签名失败
	AssertFailed                                  // 断言失败
	ParseFailed                                   // 解析失败
	RequestFailed                                 // 请求失败
	ResponseError                                 // 响应错误
	GetAddressFailed                              // 获取收货地址失败
	SelectAddressFailed                           // 选择收货地址错误
	NoValidAddress                                // 当前没有可用的收货地址
	NoValidProduct                                // 当前购物车中没有可购商品
	NoReserveTime                                 // 当前没有可用的运力
	NoReserveTimeAndRetry                         // 当前没有可用的运力, 请稍后再试
	ReserveTimeIsDisabled                         // 您选择的送达时间已经失效, 请重新选择
)

func (i ErrorCode) Int() int {
	return int(i)
}

func (i ErrorCode) Uint() uint {
	return uint(i)
}
