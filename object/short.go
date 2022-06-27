package object

import (
	"crypto/md5"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io"
	"math"
	"strings"
	"time"
)

const (
	InitDataBaseError = "init database for short error"
	Source            = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

type Short struct {
	Id       int64
	LongUrl  string `xorm:"varchar(255)"`
	ShortUrl string `xorm:"varchar(50)"`
	Md5      string `xorm:"varchar(50)"`
	Expire   time.Time
	Created  time.Time `xorm:"created"`
	Deleted  time.Time `xorm:"deleted"`
}

func ToShort(url string) string {
	hl := hlog.DefaultLogger()
	var short Short
	// 获取分号器id(短链id)
	id := GetTicket()
	// 将长链md5压缩，利于索引搜索
	h := md5.New()
	_, err := io.WriteString(h, url)
	if err != nil {
		hl.Error("write md5 string io error")
		return ""
	}
	short.LongUrl = url
	short.Md5 = fmt.Sprintf("%x", h.Sum(nil))
	short.Id = id
	// get short url
	short.ShortUrl = changeShort(id)
	// save information of short url
	_, err2 := adapter.Engine.Insert(&short)
	if err2 != nil {
		hl.Error("save information of short url error")
		return ""
	}
	return short.ShortUrl
}

// 将短链id转换为短链（10进制转62进制）
func changeShort(id int64) (url string) {
	var bytes []byte
	for id > 0 {
		bytes = append(bytes, Source[id%62])
		id = id / 62
	}
	reverse(bytes)
	return string(bytes)
}

// 反转字符串
func reverse(a []byte) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}

// 62进制转10进制
func decode(str string) int64 {
	var num int64
	n := len(str)
	for i := 0; i < n; i++ {
		pos := strings.IndexByte(Source, str[i])
		num += int64(math.Pow(62, float64(n-i-1)) * float64(pos))
	}
	return num
}

// GetLong 根据短链(62)获取长链
func GetLong(url string) (longUrl string) {
	var short Short
	_, err := adapter.Engine.Table("short").Where("id = ?", decode(url)).Get(&short)
	if err != nil {
		hlog.Info(err)
		return ""
	}
	return short.LongUrl
}
