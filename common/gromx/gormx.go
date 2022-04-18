package gromx

import (
	"admonitor-server/common/errorx"
	"admonitor-server/common/global"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGromDB() *gorm.DB {
	dsn := DataSource()
	mysqlConfig := mysql.Config{
		DSN:                       dsn,  // DSN data source name
		DefaultStringSize:         191,  // string 类型字段的默认长度
		SkipInitializeWithVersion: true, // 根据版本自动配置
	}
	db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig())
	if err != nil {
		global.LOG.Error(errorx.New("连接mysql数据库失败", err))
		panic(err)
	}
	//配置数据连接参数！
	mySQL, err := db.DB()
	if err != nil {
		global.LOG.Error(errorx.New("配置数据连接参数失败", err))
		panic(err)
	}
	// defer mySQL.Close()
	//设置最大空闲连接
	num := viper.GetInt("mysql.max-idle-conns")
	if num == 0 {
		num = 10
	}
	mySQL.SetMaxIdleConns(num)
	//设置最大连接数
	num2 := viper.GetInt("mysql.max-open-conns")
	if num2 == 0 {
		num2 = 100
	}
	mySQL.SetMaxOpenConns(num2)
	//设置连接超时时间:1分钟
	mySQL.SetConnMaxLifetime(time.Minute)
	return db
}

func gormConfig() *gorm.Config {
	newLogger := logger.New(
		NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             400 * time.Millisecond, // 慢 SQL 阈值
			LogLevel:                  logger.Warn,            // 日志级别
			IgnoreRecordNotFoundError: true,                   // 是否忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,                  // 彩色打印
		},
	)
	switch viper.GetString("mysql.log-level") {
	case "silent", "Silent":
		newLogger.LogMode(logger.Silent)
	case "error", "Error":
		newLogger.LogMode(logger.Error)
	case "warn", "Warn":
		newLogger.LogMode(logger.Warn)
	case "info", "Info":
		newLogger.LogMode(logger.Info)
	default:
		newLogger.LogMode(logger.Info)
	}
	return &gorm.Config{
		Logger: newLogger,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		DryRun: false,
	}
}

func DataSource() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", viper.GetString("mysql.user"), viper.GetString("mysql.password"), viper.GetString("mysql.addr"), viper.GetString("mysql.db"), viper.GetString("mysql.params"))
}

func datasource2() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", viper.GetString("mysql.user"), viper.GetString("mysql.password2"), viper.GetString("mysql.addr"), viper.GetString("mysql.db"), viper.GetString("mysql.params"))
}

type writer struct {
	logger.Writer
}

// NewWriter writer 构造函数
func NewWriter(w logger.Writer) *writer {
	return &writer{Writer: w}
}

// Printf 格式化打印日志，实现 Writer 接口
func (w *writer) Printf(message string, data ...interface{}) {
	global.LOG.Info(fmt.Sprintf("【GORM日志】"+message+"\n", data...))
	if viper.GetString("model") == "dev" {
		w.Writer.Printf(message, data...)
	}
}
