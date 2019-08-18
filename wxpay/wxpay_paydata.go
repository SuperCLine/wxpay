package wxpay

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"hash"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var BufPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 4<<10)) // 4KB
	},
}

type PayData struct {

	data map[string]interface{}
}

func NewPayData() *PayData  {

	return &PayData{
		data:make(map[string]interface{}),
	}
}

func (pd *PayData) IsSet(key string) bool  {

	_, ok := pd.data[key]
	return ok
}

func (pd *PayData) Set(key string, val interface{})  {

	vKind := reflect.ValueOf(val).Kind()
	switch vKind {
	case reflect.String:
		pd.data[key] = val.(string)
	case reflect.Int:
		pd.data[key] = strconv.Itoa(val.(int))
	case reflect.Int64:
		pd.data[key] = strconv.FormatInt(val.(int64), 10)
	case reflect.Float32:
		pd.data[key] = strconv.FormatFloat(float64(val.(float32)), 'f', -1, 32)
	case reflect.Float64:
		pd.data[key] = strconv.FormatFloat(val.(float64), 'f', -1, 64)
	case reflect.Ptr:
		pd.data[key] = val
	case reflect.Struct:
		pd.data[key] = val
	case reflect.Map:
		pd.data[key] = val
	case reflect.Slice:
		pd.data[key] = val
	default:
		pd.data[key] = ""
	}
}

func (pd *PayData) Get(key string) string  {

	val, ok := pd.data[key]
	if !ok {
		return ""
	}

	_, oks := val.(string)
	if oks {
		return val.(string)
	} else {
		b, err := json.Marshal(val)
		if err != nil {
			return ""
		}
		str := string(b)
		if str == "null" {
			return ""
		}
		return str
	}
}

func (pd *PayData) ToXml() []byte  {

	buf := BufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer BufPool.Put(buf)

	buf.WriteString("<xml>")
	for key := range pd.data {

		buf.WriteString("<")
		buf.WriteString(key)
		buf.WriteString("><![CDATA[")
		buf.WriteString(pd.Get(key))
		buf.WriteString("]]></")
		buf.WriteString(key)
		buf.WriteString(">")
	}
	buf.WriteString("</xml>")

	return buf.Bytes()
}

func (pd *PayData) FromXml(r io.Reader) error  {

	buf := BufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer BufPool.Put(buf)

	_, err := buf.ReadFrom(r)
	if err != nil {
		return err
	}

	return xml.Unmarshal(buf.Bytes(), pd)
}

func (pd *PayData) ToUrl() string  {

	url := ""
	for key := range pd.data {

		v := pd.Get(key)
		if key != "sign" && v != "" {

			url += key + "=" + v + "&"
		}
	}

	return strings.Trim(url, "&")
}

func (pd *PayData) ToJson() string  {

	b, err := json.Marshal(pd.data)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (pd *PayData) FromJson(r io.Reader) error  {

	buf := BufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer BufPool.Put(buf)

	_, err := buf.ReadFrom(r)
	if err != nil {
		return err
	}

	var jsonTemplate interface{}
	err = json.Unmarshal(buf.Bytes(), &jsonTemplate)
	if err != nil {
		return err
	}

	pd.data = jsonTemplate.(map[string]interface{})

	return nil
}

func (pd *PayData) FromJsonStr(s string) error  {

	var jsonTemplate interface{}
	err := json.Unmarshal([]byte(s), &jsonTemplate)
	if err != nil {
		return err
	}

	pd.data = jsonTemplate.(map[string]interface{})

	return nil
}

func (pd *PayData) MakeSign(apiKey, signType string) string {

	var h hash.Hash
	if signType == SignType_HMAC_SHA256 {
		h = hmac.New(sha256.New, []byte(apiKey))
	} else {
		h = md5.New()
	}

	keys := make([]string, 0, len(pd.data))
	for k := range pd.data {
		if k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	bufw := bufio.NewWriterSize(h, 128)
	for _, k := range keys {
		v := pd.Get(k)
		if v == "" {
			continue
		}
		bufw.WriteString(k)
		bufw.WriteByte('=')
		bufw.WriteString(v)
		bufw.WriteByte('&')
	}
	bufw.WriteString("key=")
	bufw.WriteString(apiKey)
	bufw.Flush()

	signature := make([]byte, hex.EncodedLen(h.Size()))
	hex.Encode(signature, h.Sum(nil))

	return string(bytes.ToUpper(signature))
}

func (pd *PayData) CheckSign(apiKey, signType string) error {

	if !pd.IsSet("sign") || pd.Get("sign") == "" {
		return fmt.Errorf("签名不正确")
	}

	return_sign := pd.Get("sign")
	cal_sign := pd.MakeSign(apiKey, signType)

	if return_sign != cal_sign {
		return fmt.Errorf("签名验证错误")
	} else {
		return nil
	}
}


type xmlMapEntry struct {

	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func (pd *PayData) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	var value string
	for k, v := range pd.data {

		vKind := reflect.ValueOf(v).Kind()

		switch vKind {
		case reflect.String:
			value = v.(string)
		case reflect.Int:
			value = strconv.Itoa(v.(int))
		case reflect.Int64:
			value = strconv.FormatInt(v.(int64), 10)
		case reflect.Float32:
			value = strconv.FormatFloat(float64(v.(float32)), 'f', -1, 32)
		case reflect.Float64:
			value = strconv.FormatFloat(v.(float64), 'f', -1, 64)
		default:
			value = ""
		}
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: value})
	}

	return e.EncodeToken(start.End())
}

func (pd *PayData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	for {
		var e xmlMapEntry
		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		pd.Set(e.XMLName.Local, e.Value)
	}

	return nil
}