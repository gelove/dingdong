package user_dto

type Result struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type Data struct {
	DoingRefundNum      int     `json:"doing_refund_num"`
	NoCommentOrderPoint int     `json:"no_comment_order_point"`
	NameNotice          string  `json:"name_notice"`
	NoPayOrderNum       int     `json:"no_pay_order_num"`
	DoingOrderNum       int     `json:"doing_order_num"`
	UserVip             VIP     `json:"user_vip"`
	UserSign            Sign    `json:"user_sign"`
	NotOnionTip         int     `json:"not_onion_tip"`
	NoDrawCouponMoney   string  `json:"no_draw_coupon_money"`
	PointNum            int     `json:"point_num"`
	Balance             Balance `json:"balance"`
	UserInfo            Info    `json:"user_info"`
	CouponNum           int     `json:"coupon_num"`
	NoCommentOrderNum   int     `json:"no_comment_order_num"`
}

type VIP struct {
	IsRenew                  int    `json:"is_renew"`
	VipSaveMoneyDescription  string `json:"vip_save_money_description"`
	VipDescription           string `json:"vip_description"`
	VipStatus                int    `json:"vip_status"`
	VipNotice                string `json:"vip_notice"`
	VipExpireTimeDescription string `json:"vip_expire_time_description"`
	VipUrl                   string `json:"vip_url"`
}

type Sign struct {
	IsTodaySign bool   `json:"is_today_sign"`
	SignSeries  int    `json:"sign_series"`
	SignText    string `json:"sign_text"`
}

type Balance struct {
	SetFingerPayPassword int    `json:"set_finger_pay_password"`
	Balance              string `json:"balance"`
	SetPayPassword       int    `json:"set_pay_password"`
}

type Info struct {
	Birthday       string `json:"birthday"`
	ShowInviteCode bool   `json:"show_invite_code"`
	NameInCheck    string `json:"name_in_check"`
	InviteCodeUrl  string `json:"invite_code_url"`
	Sex            int    `json:"sex"`
	Mobile         string `json:"mobile"`
	Avatar         string `json:"avatar"`
	ImUid          int    `json:"im_uid"`
	BindStatus     int    `json:"bind_status"`
	NameStatus     int    `json:"name_status"`
	NewRegister    bool   `json:"new_register"`
	ImSecret       string `json:"im_secret"`
	Name           string `json:"name"`
	ID             string `json:"id"`
	Introduction   string `json:"introduction"`
}
