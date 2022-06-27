package object

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
	"xorm.io/xorm"
)

const (
	DataSourceCloseError = "database engine close error"
)

// global
var adapter *Adapter = nil

func InitAdapter() {
	adapter = NewAdapter(DatabaseSetting.driverName, DatabaseSetting.dataSourceName, DatabaseSetting.dbName)
	err := adapter.Engine.Sync(new(Short), new(Ticket))
	l := hlog.DefaultLogger()
	if err != nil {
		l.Error(InitDataBaseError)
	}
}

type Adapter struct {
	driverName     string
	dataSourceName string
	dbName         string
	Engine         *xorm.Engine
}

// finalizer is the destructor for Adapter.
func finalizer(a *Adapter) {
	err := a.Engine.Close()
	if err != nil {
		panic(err)
	}
}

// NewAdapter is the constructor for Adapter.
func NewAdapter(driverName string, dataSourceName string, dbName string) *Adapter {
	a := &Adapter{}
	a.driverName = driverName
	a.dataSourceName = dataSourceName
	a.dbName = dbName

	// Open the DB, create it if not existed.
	a.open()

	// Call the destructor when the object is released.
	runtime.SetFinalizer(a, finalizer)

	return a
}

func (a *Adapter) open() {
	engine, err := xorm.NewEngine(a.driverName, a.dataSourceName+a.dbName)
	if err != nil {
		panic(err)
	}

	a.Engine = engine
}

func (a *Adapter) close() {
	err := a.Engine.Close()
	if err != nil {
		panic(DataSourceCloseError)
	}
	a.Engine = nil
}
