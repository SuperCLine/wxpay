package wxpay

import "fmt"

const (
	wxBase_url 	= "https://api.mch.weixin.qq.com/"

	//正式
	wxURL_UnifiedOrder      = wxBase_url + "pay/unifiedorder"                //统一下单
	wxURL_OrderQuery        = wxBase_url + "pay/orderquery"                  //查询订单
	wxURL_Micropay          = wxBase_url + "pay/micropay"                    //提交付款码支付
	wxURL_CloseOrder        = wxBase_url + "pay/closeorder"                  //关闭订单
	wxURL_Refund            = wxBase_url + "secapi/pay/refund"               //申请退款
	wxURL_Reverse           = wxBase_url + "secapi/pay/reverse"              //撤销订单
	wxURL_RefundQuery       = wxBase_url + "pay/refundquery"                 //查询退款
	wxURL_DownloadBill      = wxBase_url + "pay/downloadbill"                //下载对账单

	//支付类型
	TradeType_JsApi  = "JSAPI"
	TradeType_App    = "APP"
	TradeType_H5     = "MWEB"
	TradeType_Native = "NATIVE"

	SignType_MD5         = "MD5"
	SignType_SHA1        = "SHA1"
	SignType_HMAC_SHA256 = "HMAC-SHA256"

	RCSuccess = "SUCCESS"
	RCFail    = "FAIL"
)

func ApiUnifiedOrder(payClient *PayClient, payData *PayData) (*PayData, error)  {

	if !payData.IsSet("out_trade_no") ||
		!payData.IsSet("body") ||
		!payData.IsSet("total_fee") ||
		!payData.IsSet("trade_type") ||
		!payData.IsSet("notify_url") ||
		!payData.IsSet("spbill_create_ip"){
		return nil, fmt.Errorf("need pay param")
	}

	if payData.Get("trade_type") == TradeType_JsApi && !payData.IsSet("openid") {
		return nil, fmt.Errorf("need openid")
	}

	if payData.Get("trade_type") == TradeType_Native && !payData.IsSet("product_id") {
		return  nil, fmt.Errorf("need product_id")
	}

	payData.Set("appid", payConfig.AppId())
	payData.Set("mch_id", payConfig.MchId())
	payData.Set("sign_type", SignType_MD5)
	payData.Set("sign", payData.MakeSign(payConfig.ApiKey(), SignType_MD5))

	rpayData, err := payClient.PostXML(wxURL_UnifiedOrder, payData)
	if err != nil {
		return nil, err
	}

	return rpayData, nil
}

func ApiOrderQuery(payClient *PayClient, payData *PayData) (*PayData, error)  {

	if !payData.IsSet("out_trade_no") &&
		!payData.IsSet("transaction_id") {
		return nil, fmt.Errorf("need pay param")
	}

	payData.Set("appid", payConfig.AppId())
	payData.Set("mch_id", payConfig.MchId())
	payData.Set("nonce_str", NonceStr())
	payData.Set("sign_type", SignType_HMAC_SHA256)
	payData.Set("sign", payData.MakeSign(payConfig.ApiKey(), SignType_HMAC_SHA256))

	rpayData, err := payClient.PostXML(wxURL_OrderQuery, payData)
	if err != nil {
		return nil, err
	}

	return rpayData, nil
}

func ApiMicropay(payClient *PayClient, payData *PayData) (*PayData, error)  {

	if !payData.IsSet("out_trade_no") ||
		!payData.IsSet("body") ||
		!payData.IsSet("total_fee") ||
		!payData.IsSet("auth_code") ||
		!payData.IsSet("spbill_create_ip"){
		return nil, fmt.Errorf("need pay param")
	}

	payData.Set("appid", payConfig.AppId())
	payData.Set("mch_id", payConfig.MchId())
	payData.Set("nonce_str", NonceStr())
	payData.Set("sign_type", SignType_HMAC_SHA256)
	payData.Set("sign", payData.MakeSign(payConfig.ApiKey(), SignType_HMAC_SHA256))

	rpayData, err := payClient.PostXML(wxURL_Micropay, payData)
	if err != nil {
		return nil, err
	}

	return rpayData, nil
}

func ApiCloseOrder(payClient *PayClient, payData *PayData) (*PayData, error) {

	if !payData.IsSet("out_trade_no") {
		return nil, fmt.Errorf("need pay param")
	}

	payData.Set("appid", payConfig.AppId())
	payData.Set("mch_id", payConfig.MchId())
	payData.Set("nonce_str", NonceStr())
	payData.Set("sign_type", SignType_HMAC_SHA256)
	payData.Set("sign", payData.MakeSign(payConfig.ApiKey(), SignType_HMAC_SHA256))

	rpayData, err := payClient.PostXML(wxURL_CloseOrder, payData)
	if err != nil {
		return nil, err
	}

	return rpayData, nil
}

func ApiRefund(payClient *PayClient, payData *PayData) (*PayData, error)  {

	if !payData.IsSet("out_trade_no") ||
		!payData.IsSet("out_refund_no") ||
		!payData.IsSet("total_fee") ||
		!payData.IsSet("refund_fee") ||
		!payData.IsSet("op_user_id") {
		return nil, fmt.Errorf("need pay param")
	}

	payData.Set("appid", payConfig.AppId())
	payData.Set("mch_id", payConfig.MchId())
	payData.Set("nonce_str", NonceStr())
	payData.Set("sign_type", SignType_HMAC_SHA256)
	payData.Set("sign", payData.MakeSign(payConfig.ApiKey(), SignType_HMAC_SHA256))

	rpayData, err := payClient.PostXML(wxURL_Refund, payData)
	if err != nil {
		return nil, err
	}

	return rpayData, nil
}

func ApiReverse(payClient *PayClient, payData *PayData) (*PayData, error)  {

	if !payData.IsSet("out_trade_no") {
		return nil, fmt.Errorf("need pay param")
	}

	payData.Set("appid", payConfig.AppId())
	payData.Set("mch_id", payConfig.MchId())
	payData.Set("nonce_str", NonceStr())
	payData.Set("sign_type", SignType_HMAC_SHA256)
	payData.Set("sign", payData.MakeSign(payConfig.ApiKey(), SignType_HMAC_SHA256))

	rpayData, err := payClient.PostXML(wxURL_Reverse, payData)
	if err != nil {
		return nil, err
	}

	return rpayData, nil
}

func ApiRefundQuery(payClient *PayClient, payData *PayData) (*PayData, error)  {

	if !payData.IsSet("out_trade_no") &&
		!payData.IsSet("out_refund_no") &&
		!payData.IsSet("transaction_id") &&
		!payData.IsSet("refund_id") {
		return nil, fmt.Errorf("need pay param")
	}

	payData.Set("appid", payConfig.AppId())
	payData.Set("mch_id", payConfig.MchId())
	payData.Set("nonce_str", NonceStr())
	payData.Set("sign_type", SignType_HMAC_SHA256)
	payData.Set("sign", payData.MakeSign(payConfig.ApiKey(), SignType_HMAC_SHA256))

	rpayData, err := payClient.PostXML(wxURL_RefundQuery, payData)
	if err != nil {
		return nil, err
	}

	return rpayData, nil
}

func ApiDownloadBill(payClient *PayClient, payData *PayData) (*PayData, error)  {

	if !payData.IsSet("bill_date") {
		return nil, fmt.Errorf("need pay param")
	}

	payData.Set("appid", payConfig.AppId())
	payData.Set("mch_id", payConfig.MchId())
	payData.Set("nonce_str", NonceStr())
	payData.Set("sign_type", SignType_HMAC_SHA256)
	payData.Set("sign", payData.MakeSign(payConfig.ApiKey(), SignType_HMAC_SHA256))

	rpayData, err := payClient.PostXML(wxURL_DownloadBill, payData)
	if err != nil {
		return nil, err
	}

	return rpayData, nil
}
