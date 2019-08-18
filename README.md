# wxpay
This project is wexin mini game's pay SDK with golang.

# [官方文档](http://supercline.com/game/tool-sdk/wxpay-with-go.html) #

**go get github.com/SuperCLine/wxpay**

- `wxpay_payconfig.go` 支付配置模块
- `wxpay_util.go` 支付通用模块
- `wxpay_paydata.go` 支付通信数据模块
- `wxpay_payapi.go` 支付API模块
- `wxpay_payclient.go` 与微信后台通信模块
		
**以上模块构成go支付SDK**

- `wxpay_default_listener.go` 支付回调监听
- `wxpay_default_service.go` 基于支付SDK实现的默认支付服务
		
**开发者可以根据以上模块实现自己的支付服务**