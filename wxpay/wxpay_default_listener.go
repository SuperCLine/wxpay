package wxpay

type DefaultPayListener struct {

}

func NewDefaultPayListener() *DefaultPayListener {

	return &DefaultPayListener{}
}

func (listener *DefaultPayListener) HandleLogicLogin(pdata *PayData)  {

	// 1.如果逻辑服务与登录服务分开时，则通知逻辑服登录结果
	// 2.如果逻辑服务与登录服务在一起，则直接处理
}

func (listener *DefaultPayListener) HandleLogicProduct(id string) *PayData  {

	// 根据商品ID从商品配置文件中或从商品库表中读取

	pd := NewPayData()
	pd.Set("product_id", id)
	pd.Set("body", "b")
	pd.Set("detail", "d")
	pd.Set("fee_type", "CNY")
	pd.Set("total_fee", 100)

	return pd
}

func (listener *DefaultPayListener) HandleLogicPay(pdata *PayData)  {

	// 1.如果逻辑服务与登录服务分开时，则通知逻辑服支付结果
	// 2.如果逻辑服务与登录服务在一起，则直接处理(如写日志表、写订单表等)
}

// TODO: add more logic handler at here

var (
	payListener = NewDefaultPayListener()
)