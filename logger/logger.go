package logger

import (
	"os"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	DebugLevelEmoji  = "🐛"
	InfoLevelEmoji   = ""
	WarnLevelEmoji   = "⚠"
	ErrorLevelEmoji  = "✖"
	DPanicLevelEmoji = "🚨"
	PanicLevelEmoji  = "🆘"
	FatalLevelEmoji  = "💀"
)

var Log *zap.Logger

const skipCaller = 1

func init() {
	Log, _ = zap.NewProduction(zap.AddCallerSkip(skipCaller))
}

type ZapConfig struct {
	RollFileConfig lumberjack.Logger `mapstructure:"roll_file_config"`
}

func Init(cfg *ZapConfig, debug bool) {
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel: func(l zapcore.Level, pae zapcore.PrimitiveArrayEncoder) {
			emoji := ""
			switch l {
			case zapcore.DebugLevel:
				emoji = DebugLevelEmoji
			case zapcore.InfoLevel:
				emoji = InfoLevelEmoji
			case zapcore.WarnLevel:
				emoji = WarnLevelEmoji
			}
			pae.AppendString("[" + emoji + "" + l.CapitalString() + "]")
		},
		EncodeTime:     zapcore.TimeEncoderOfLayout("[2006-01-02 15:04:05]"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	consoleCore := zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)
	fileCore := zapcore.NewCore(encoder, zapcore.AddSync(&cfg.RollFileConfig), zapcore.DebugLevel)

	core := zapcore.NewTee(consoleCore, fileCore)

	_log := zap.New(core, zap.AddCallerSkip(skipCaller), zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCaller())
	if debug {
		_log = _log.WithOptions(zap.Development())
	}
	Log = _log
	_log.Info("logger initialized")
}

func Sugar() *zap.SugaredLogger {
	return Log.Sugar()
}

func joinSpace(args ...interface{}) []interface{} {
	var newArgs []interface{}
	for _, arg := range args {
		newArgs = append(newArgs, arg)
		newArgs = append(newArgs, " ")
	}
	return newArgs
}

func Debug(args ...interface{})  { Sugar().Debug(args...) }
func Info(args ...interface{})   { Sugar().Info(joinSpace(args...)...) }
func Warn(args ...interface{})   { Sugar().Warn(joinSpace(args...)...) }
func Error(args ...interface{})  { Sugar().Error(joinSpace(args...)...) }
func DPanic(args ...interface{}) { Sugar().DPanic(joinSpace(args...)...) }
func Panic(args ...interface{})  { Sugar().Panic(joinSpace(args...)...) }
func Fatal(args ...interface{})  { Sugar().Fatal(joinSpace(args...)...) }

func Debugf(template string, args ...interface{}) { Sugar().Debugf(template, args...) }
func Infof(template string, args ...interface{})  { Sugar().Infof(template, args...) }
func Warnf(template string, args ...interface{})  { Sugar().Warnf(template, args...) }
func Errorf(template string, args ...interface{}) { Sugar().Errorf(template, args...) }
func Panicf(template string, args ...interface{}) { Sugar().Panicf(template, args...) }
func Fatalf(template string, args ...interface{}) { Sugar().Fatalf(template, args...) }

func Debugw(msg string, keysAndValues ...interface{}) { Sugar().Debugw(msg, keysAndValues...) }
func Infow(msg string, keysAndValues ...interface{})  { Sugar().Infow(msg, keysAndValues...) }
func Warnw(msg string, keysAndValues ...interface{})  { Sugar().Warnw(msg, keysAndValues...) }
func Errorw(msg string, keysAndValues ...interface{}) { Sugar().Errorw(msg, keysAndValues...) }

func OK(args ...interface{})            { Sugar().Info(append([]interface{}{"✅"}, args...)...) }
func OKf(t string, args ...interface{}) { Sugar().Infof("✅ "+t, args...) }
func OKw(msg string, kv ...interface{}) { Sugar().Infow("✅ "+msg, kv...) }

func Fail(args ...interface{})            { Sugar().Error(append([]interface{}{"❌"}, args...)...) }
func Failf(t string, args ...interface{}) { Sugar().Errorf("❌ "+t, args...) }
func Failw(msg string, kv ...interface{}) { Sugar().Errorw("❌ "+msg, kv...) }

func Pending(args ...interface{})            { Sugar().Info(append([]interface{}{"⏳"}, args...)...) }
func Pendingf(t string, args ...interface{}) { Sugar().Infof("⏳ "+t, args...) }
func Pendingw(msg string, kv ...interface{}) { Sugar().Infow("⏳ "+msg, kv...) }

func Start(args ...interface{})            { Sugar().Info(append([]interface{}{"🚀"}, args...)...) }
func Startf(t string, args ...interface{}) { Sugar().Infof("🚀 "+t, args...) }
func Startw(msg string, kv ...interface{}) { Sugar().Infow("🚀 "+msg, kv...) }

func Done(args ...interface{})            { Sugar().Info(append([]interface{}{"🏁"}, args...)...) }
func Donef(t string, args ...interface{}) { Sugar().Infof("🏁 "+t, args...) }
func Donew(msg string, kv ...interface{}) { Sugar().Infow("🏁 "+msg, kv...) }

func Attn(args ...interface{})            { Sugar().Info(append([]interface{}{"⚠️"}, args...)...) }
func Attnf(t string, args ...interface{}) { Sugar().Infof("⚠️ "+t, args...) }
func Attnw(msg string, kv ...interface{}) { Sugar().Infow("⚠️ "+msg, kv...) }

func Note(args ...interface{})            { Sugar().Info(append([]interface{}{"ℹ️"}, args...)...) }
func Notef(t string, args ...interface{}) { Sugar().Infof("ℹ️ "+t, args...) }
func Notew(msg string, kv ...interface{}) { Sugar().Infow("ℹ️ "+msg, kv...) }

type Stack struct {
	Caller string `json:"caller"`
	Line   int    `json:"line"`
	Func   string `json:"func"`
	Module string `json:"module"`
}

func LogStack(level zapcore.Level, args ...interface{}) {

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "<unknown>"
		line = 0
	}
	Sugar().Log(level, append([]interface{}{
		map[string]interface{}{
			"stacks": map[string]interface{}{
				"file": file,
				"line": line,
			},
		},
	}, args[:]...)...)
}

func DebugStack(args ...interface{}) {
	LogStack(zap.DebugLevel, args...)
}

func WarnStack(args ...interface{}) {
	LogStack(zap.WarnLevel, args...)
}

// LogStartupBanner 输出启动成功的 Unicode logo
func LogStartupBanner(appName, version, env, serverAddr string, startTime time.Time) {
	startTimeStr := startTime.Format("2006-01-02 15:04:05")

	// 计算启动耗时
	uptime := time.Since(startTime)

	// 使用醒目的分隔线
	separator := "═══════════════════════════════════════════════════════════════"

	// 构建 banner，每行单独输出以保持格式
	bannerLines := []string{
		"",
		separator,
		"",
		"     ███████╗██╗   ██╗██████╗ ███████╗██████╗ ██╗██╗",
		"     ██╔════╝██║   ██║██╔══██╗██╔════╝██╔══██╗██║██║",
		"     ███████╗██║   ██║██████╔╝█████╗  ██████╔╝██║██║",
		"     ╚════██║██║   ██║██╔═══╝ ██╔══╝  ██╔══██╗██║██║",
		"     ███████║╚██████╔╝██║     ███████╗██║  ██║██║██║",
		"     ╚══════╝ ╚═════╝ ╚═╝     ╚══════╝╚═╝  ╚═╝╚═╝╚═╝",
		"",
		"              ╔═══════════════════════════════╗",
		"              ║   🚀 服务启动成功！🚀         ║",
		"              ╚═══════════════════════════════╝",
		"",
		separator,
		"  📦 应用名称: " + appName,
		"  📌 版本信息: " + version,
		"  🌍 运行环境: " + env,
		"  ⏰ 启动时间: " + startTimeStr,
		"  ⚡ 启动耗时: " + uptime.String(),
		"  🌐 服务地址: " + serverAddr,
		separator,
		"",
	}

	// 逐行输出，确保在日志文件中醒目
	for _, line := range bannerLines {
		if line == "" {
			Log.Info("")
		} else {
			Log.Info(line)
		}
	}
}
