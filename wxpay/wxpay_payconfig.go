package wxpay

type PayConfig struct {

	appId  string
	appSecret string
	mchId  string
	apiKey string
	notifyUrl string
}

func (pcf *PayConfig) AppId() string {
	return pcf.appId
}
func (pcf *PayConfig) AppSecret() string  {
	return  pcf.appSecret
}
func (pcf *PayConfig) MchId() string {
	return pcf.mchId
}
func (pcf *PayConfig) ApiKey() string {
	return pcf.apiKey
}
func (pcf *PayConfig) NotifyUrl() string  {
	return pcf.notifyUrl
}

func NewPayConfig() *PayConfig  {

	return &PayConfig{

		appId:"",
		appSecret:"",
		mchId:"",
		apiKey:"",
		notifyUrl:"",
	}
}

var (

	payConfig = NewPayConfig()
)