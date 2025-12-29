package wirex

import (
	"os"
	"teaching_manage/dao"
	"teaching_manage/pkg/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGormDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./data/teaching_manage.db"), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}

	// Auto migrate the Student and Teacher models
	err = db.AutoMigrate(&dao.Student{}, &dao.Teacher{}, &dao.Order{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitLogger() logger.Logger {
	// 实现按文件大小切割日志
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/teaching_manage.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
	})

	// 配置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 创建核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), w),
		zap.DebugLevel,
	)

	l := zap.New(core, zap.AddCaller())
	return logger.NewZapLogger(l)
}
