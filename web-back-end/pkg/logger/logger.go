package logger

import (
	"os"

	"github.com/haoran-mc/wx_scan_login/web-back-end/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger // 只支持强类型的结构化日志记录，比 Sugared Logger 更快
	// SugaredLogger *zap.SugaredLogger // 支持结构化和 printf 风格的日志记录
)

var zapLogLevel map[string]zapcore.LevelEnabler = map[string]zapcore.LevelEnabler{
	"debug":  zap.DebugLevel,
	"info":   zap.InfoLevel,
	"warn":   zap.WarnLevel,
	"error":  zap.ErrorLevel,
	"dpanic": zap.DPanicLevel,
	"panic":  zap.PanicLevel,
	"fatal":  zap.FatalLevel,
}

func init() {
	core := zapcore.NewCore(
		getEncoder(),   // 日志编码器，如何写入日志
		getLogWriter(), // 指定日志文件写入位置
		getlogLevel(),  // 设置日志级别
	)

	Logger = zap.New(core, zap.AddCaller())
	// SugaredLogger = Logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 修改时间编码器
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 在日志中使用大写字母记录日志级别

	// NewConsoleEncoder 打印更符合人们观察的方式
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	return zapcore.AddSync(os.Stdout)
}

func getlogLevel() zapcore.LevelEnabler {
	return zapLogLevel[config.Conf.Logger.Level]
}
