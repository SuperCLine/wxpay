package wxpay

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	payClient = NewPayClient(nil)

	ErrCode_Success = int64(0)
	ErrCode_ParseQuery = int64(1)
	ErrCode_Param = int64(2)
	ErrCode_PayData = int64(3)

	ErrCode_Product = int64(101)
)

type DefaultPayService struct {

}

func NewDefaultPayService() *DefaultPayService  {

	return &DefaultPayService{}
}

func (srv *DefaultPayService) Start()  {

	// login
	http.HandleFunc("/Login", srv.handlerLogin)

	// pay
	http.HandleFunc("/UnifiedOrder", srv.handlerUnifiedOrder)
	http.HandleFunc("/PayResult", srv.handlerPayResult)
	http.HandleFunc("/OrderQuery", srv.handlerOrderQuery)
	http.HandleFunc("/Micropay", srv.handlerMicropay)
	http.HandleFunc("/CloseOrder", srv.handlerCloseOrder)
	http.HandleFunc("/Refund", srv.handlerRefund)
	http.HandleFunc("/Reverse", srv.handlerReverse)
	http.HandleFunc("/RefundQuery", srv.handlerRefundQuery)
	http.HandleFunc("/DownloadBill", srv.handlerDownloadBill)

	log.Println("server running ...")

	err := http.ListenAndServeTLS(":80", "./keystore/server.crt", "./keystore/server.key", nil)
	//err := http.ListenAndServe(":80", nil)

	log.Println("start err ..."+err.Error())
}

func (srv *DefaultPayService) handlerError(w http.ResponseWriter, err int64, msg string) func() {

	data := NewPayData()
	return func() {

		data.Set("errcode", err)
		data.Set("errmsg", msg)
		io.WriteString(w, data.ToJson())

		log.Println(fmt.Sprintf("errcode: %d, errmsg: %s", err, msg))
	}
}

func (srv *DefaultPayService) handlerErrorXML(w http.ResponseWriter, err string, msg string) func() {

	data := NewPayData()
	return func() {

		data.Set("return_code", err)
		data.Set("return_msg", msg)
		io.WriteString(w, string(data.ToXml()))

		log.Println(fmt.Sprintf("return_code: %s, return_msg: %s", err, msg))
	}
}

func (srv *DefaultPayService) handlerLogin(w http.ResponseWriter, r *http.Request) {

	reqData, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		srv.handlerError(w, ErrCode_ParseQuery, err.Error())()
		return
	}

	jscode := reqData.Get("code")
	if jscode == "" {
		srv.handlerError(w, ErrCode_Param, "error param")()
		return
	}

	rspData, err := payClient.Login(jscode)
	if err != nil {
		srv.handlerError(w, ErrCode_PayData, err.Error())()
		return
	}

	log.Println("-----------------handlerLogin succeed-------------------")
	log.Println("openid: "+rspData.Get("openid"))
	log.Println("session_key: "+rspData.Get("session_key"))

	io.WriteString(w, rspData.ToJson())

	payListener.HandleLogicLogin(rspData)
}

func (srv *DefaultPayService) handlerUnifiedOrder(w http.ResponseWriter, r *http.Request) {

	reqData, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		srv.handlerError(w, ErrCode_ParseQuery, err.Error())()
		return
	}

	openId := reqData.Get("openid")
	productId := reqData.Get("productid")
	billIp := reqData.Get("ip")
	if openId == "" || productId == "" || billIp == "" {
		srv.handlerError(w, ErrCode_Param, "error param")()
		return
	}

	product := payListener.HandleLogicProduct(productId)
	if product == nil {
		srv.handlerError(w, ErrCode_Product, "no goods info")()
		return
	}

	tradeNo := NonceStr()

	apiData := NewPayData()
	apiData.Set("openid", openId)
	apiData.Set("product_id", productId)
	apiData.Set("body", product.Get("body"))
	apiData.Set("total_fee", product.Get("total_fee"))
	apiData.Set("detail", product.Get("detail"))
	apiData.Set("fee_type", product.Get("fee_type"))
	apiData.Set("trade_type", TradeType_JsApi)
	apiData.Set("out_trade_no", tradeNo)
	apiData.Set("spbill_create_ip", billIp)
	apiData.Set("notify_url", payConfig.NotifyUrl())
	apiData.Set("time_start", FormatTime(time.Now()))
	apiData.Set("time_expire", FormatTime(time.Now().Add(time.Minute * 10)))
	apiData.Set("nonce_str", NonceStr())

	rapiData, err := ApiUnifiedOrder(payClient, apiData)
	if err != nil {
		srv.handlerError(w, ErrCode_PayData, err.Error())()
		return
	}

	timeStamp := TimeStamp()
	nonceStr := NonceStr()
	packageStr := "prepay_id=" + rapiData.Get("prepay_id")

	rspData := NewPayData()
	rspData.Set("appId", payConfig.AppId())
	rspData.Set("timeStamp", timeStamp)
	rspData.Set("nonceStr", nonceStr)
	rspData.Set("package", packageStr)
	rspData.Set("signType", SignType_MD5)

	paySign := rspData.MakeSign(payConfig.ApiKey(), SignType_MD5)

	rspData.Set("errcode", ErrCode_Success)
	rspData.Set("errmsg", RCSuccess)
	rspData.Set("paySign", paySign)
	rspData.Set("tradeNo", tradeNo)

	io.WriteString(w, rspData.ToJson())

	log.Println("-----------------handlerUnifiedOrder succeed-------------------")
}

func (srv *DefaultPayService) handlerPayResult(w http.ResponseWriter, r *http.Request) {

	log.Println("-----------------handlerPayResult succeed-------------------")

	resultData := NewPayData()
	err := resultData.FromXml(r.Body)
	if err != nil {
		srv.handlerErrorXML(w, RCFail, RCFail)()
		return
	}
	defer r.Body.Close()

	if !resultData.IsSet("transaction_id") ||
		resultData.Get("transaction_id") == "" {
		srv.handlerErrorXML(w, RCFail, "error transaction_id")()
		return
	}

	apiData := NewPayData()
	apiData.Set("transaction_id", resultData.Get("transaction_id"))
	rapiData, err := ApiOrderQuery(payClient, apiData)
	if err != nil {
		srv.handlerErrorXML(w, RCFail, "error order query")()
	} else {
		if rapiData.Get("return_code") == RCSuccess &&
			rapiData.Get("result_code") == RCSuccess {
			srv.handlerErrorXML(w, RCSuccess, "OK")()

			// TODO:
			//payListener.HandleLogicPay()
		} else {
			srv.handlerErrorXML(w, RCFail, "error order query")()
		}
	}
}

func (srv *DefaultPayService) handlerOrderQuery(w http.ResponseWriter, r *http.Request) {

	reqData, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		srv.handlerError(w, ErrCode_ParseQuery, err.Error())()
		return
	}

	tradeNo := reqData.Get("tradeno")
	if tradeNo == "" {
		srv.handlerError(w, ErrCode_Param, "error param")()
		return
	}

	apiData := NewPayData()
	apiData.Set("out_trade_no", tradeNo)

	rapiData, err := ApiOrderQuery(payClient, apiData)
	if err != nil {
		srv.handlerError(w, ErrCode_PayData, err.Error())()
		return
	}

	rspData := NewPayData()
	if rapiData.Get("return_code") == RCSuccess &&
		rapiData.Get("result_code") == RCSuccess {
		rspData.Set("errcode", ErrCode_Success)
	} else {
		rspData.Set("errcode", ErrCode_PayData)
	}
	rspData.Set("errmsg", rapiData.Get("return_msg"))

	io.WriteString(w, rspData.ToJson())

	log.Println("-----------------handlerOrderQuery succeed-------------------")
}

func (srv *DefaultPayService) handlerMicropay(w http.ResponseWriter, r *http.Request) {

	// TODO:
}

func (srv *DefaultPayService) handlerCloseOrder(w http.ResponseWriter, r *http.Request) {

	// TODO:
}

func (srv *DefaultPayService) handlerRefund(w http.ResponseWriter, r *http.Request) {

	// TODO:
}

func (srv *DefaultPayService) handlerReverse(w http.ResponseWriter, r *http.Request) {

	// TODO:
}

func (srv *DefaultPayService) handlerRefundQuery(w http.ResponseWriter, r *http.Request) {

	// TODO:
}

func (srv *DefaultPayService) handlerDownloadBill(w http.ResponseWriter, r *http.Request) {

	// TODO:
}
