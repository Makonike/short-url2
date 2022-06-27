package object

import "github.com/cloudwego/hertz/pkg/common/hlog"

type Ticket struct {
	Id   int64
	Stub rune
}

// GetTicket 获取分号发布器id
// 在主键自增的情况下，分布式场景获取唯一id
func GetTicket() (id int64) {
	var tik Ticket
	var err error
	hl := hlog.DefaultLogger()
	// replace into 刷新主键自增
	_, err = adapter.Engine.Exec("REPLACE INTO `ticket` (stub) VALUES (?);", 'a')
	if err != nil {
		hl.Error("replace into sql error")
		return 0
	}
	// 获取id
	_, err = adapter.Engine.Table("ticket").Where("stub = ?", 'a').Desc("id").Get(&tik)
	if err != nil {
		hl.Error("get ticket sql error")
		return 0
	}
	hl.Info(tik)
	hl.Info(tik.Id)
	return tik.Id
}
