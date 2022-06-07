//go:generate stringer -type=ErrorCode --linecomment

package code

type ErrorCode int

const OK ErrorCode = 0

// Error codes 应该只定义业务逻辑上的错误，而不应该定义系统级别的错误。
const (
	InternalError         ErrorCode = 1000 + iota // 内部错误
	SelectSessionFailed                           // 选择session文件错误
	GetUserDetailFailed                           // 获取用户详情失败
	GetAddressFailed                              // 获取收货地址失败
	GetFlowDetailFailed                           // 获取首页流水详情失败
	NoValidAddress                                // 当前没有可用的收货地址
	CheckAllFailed                                // 购物车全选失败
	GetCartFailed                                 // 获取购物车失败
	NoValidProduct                                // 当前购物车中没有可购商品
	GetReserveTimeFailed                          // 获取运力失败
	NoReserveTime                                 // 当前没有可用的运力
	NoReserveTimeAndRetry                         // 当前没有可用的运力, 请稍后再试
	ReserveTimeIsDisabled                         // 您选择的送达时间已经失效, 请重新选择
	CheckOrderFailed                              // 订单校验失败
	SubmitOrderFailed                             // 提交订单失败
	NotifyFailed                                  // 通知失败
)

func (i ErrorCode) Int() int {
	return int(i)
}

func (i ErrorCode) Uint() uint {
	return uint(i)
}
