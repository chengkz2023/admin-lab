package system

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	Mysql           = "mysql"
	InitSuccess     = "\n[%v] --> 初始化数据成功!\n"
	InitDataExist   = "\n[%v] --> %v 的初始化数据已存在!\n"
	InitDataFailed  = "\n[%v] --> %v 初始化数据失败! \nerr: %+v\n"
	InitDataSuccess = "\n[%v] --> %v 初始化数据成功!\n"
)

const (
	InitOrderSystem   = 10
	InitOrderInternal = 1000
	InitOrderExternal = 100000
)

var (
	ErrMissingDBContext        = errors.New("missing db in context")
	ErrMissingDependentContext = errors.New("missing dependent value in context")
	ErrDBTypeMismatch          = errors.New("db type mismatch")
)

type SubInitializer interface {
	InitializerName() string
	MigrateTable(ctx context.Context) (next context.Context, err error)
	InitializeData(ctx context.Context) (next context.Context, err error)
	TableCreated(ctx context.Context) bool
	DataInserted(ctx context.Context) bool
}

type TypedDBInitHandler interface {
	EnsureDB(ctx context.Context, conf *request.InitDB) (context.Context, error)
	WriteConfig(ctx context.Context) error
	InitTables(ctx context.Context, inits initSlice) error
	InitData(ctx context.Context, inits initSlice) error
}

type orderedInitializer struct {
	order int
	SubInitializer
}

type initSlice []*orderedInitializer

var (
	initializers initSlice
	cache        map[string]*orderedInitializer
)

func RegisterInit(order int, i SubInitializer) {
	if initializers == nil {
		initializers = initSlice{}
	}
	if cache == nil {
		cache = map[string]*orderedInitializer{}
	}
	name := i.InitializerName()
	if _, existed := cache[name]; existed {
		panic(fmt.Sprintf("Name conflict on %s", name))
	}
	ni := orderedInitializer{order, i}
	initializers = append(initializers, &ni)
	cache[name] = &ni
}

func systemInitializers() initSlice {
	systemInits := make(initSlice, 0, len(initializers))
	for _, init := range initializers {
		if init.order < InitOrderInternal {
			systemInits = append(systemInits, init)
		}
	}
	sort.Sort(&systemInits)
	return systemInits
}

type InitDBService struct{}

func (initDBService *InitDBService) InitDB(conf request.InitDB) error {
	ctx := context.TODO()
	ctx = context.WithValue(ctx, "adminPassword", conf.AdminPassword)
	if len(initializers) == 0 {
		return errors.New("no initializers registered")
	}
	sort.Sort(&initializers)

	initHandler := NewMysqlInitHandler()
	ctx = context.WithValue(ctx, "dbtype", Mysql)

	next, err := initHandler.EnsureDB(ctx, &conf)
	if err != nil {
		return err
	}

	db := next.Value("db").(*gorm.DB)
	global.GVA_DB = db

	if err = initHandler.InitTables(next, initializers); err != nil {
		return err
	}
	if err = initHandler.InitData(next, initializers); err != nil {
		return err
	}
	if err = initHandler.WriteConfig(next); err != nil {
		return err
	}

	initializers = initSlice{}
	cache = map[string]*orderedInitializer{}
	return nil
}

func (initDBService *InitDBService) EnsureSystemInitDataOnBoot(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	inits := systemInitializers()
	if len(inits) == 0 {
		return errors.New("no system initializers registered")
	}

	ctx := context.WithValue(context.Background(), "db", db)
	ctx = context.WithValue(ctx, "dbtype", Mysql)

	if err := createTables(ctx, inits); err != nil {
		return err
	}

	for _, init := range inits {
		if init.DataInserted(ctx) {
			continue
		}
		next, err := init.InitializeData(ctx)
		if err != nil {
			return err
		}
		ctx = next
		global.GVA_LOG.Info("seeded system initializer", zap.String("name", init.InitializerName()))
	}

	return nil
}

func createDatabase(dsn string, driver string, createSQL string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSQL)
	return err
}

func createTables(ctx context.Context, inits initSlice) error {
	next, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, init := range inits {
		if init.TableCreated(next) {
			continue
		}
		n, err := init.MigrateTable(next)
		if err != nil {
			return err
		}
		next = n
	}
	return nil
}

func (a initSlice) Len() int {
	return len(a)
}

func (a initSlice) Less(i, j int) bool {
	return a[i].order < a[j].order
}

func (a initSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
