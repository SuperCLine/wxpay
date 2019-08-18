package wxpay

import (
	"math/rand"
	"strconv"
	"time"
)

var (
	BeijingLocation = time.FixedZone("Asia/Shanghai", 8*60*60)
)

func RandomString(n int) string {

	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	length := len(letter)

	for i := range b {
		b[i] = letter[rand.Intn(length)]
	}

	return string(b)
}

func NonceStr() string {

	return RandomString(32)
}

// FormatTime time.Time => "yyyyMMddHHmmss."
func FormatTime(t time.Time) string {

	return t.In(BeijingLocation).Format("20060102150405")
}

// ParseTime "yyyyMMddHHmmss" => time.Time.
func ParseTime(value string) (time.Time, error) {

	return time.ParseInLocation("20060102150405", value, BeijingLocation)
}

func TimeStamp() string  {

	return strconv.FormatInt(time.Now().Unix(), 10)
}