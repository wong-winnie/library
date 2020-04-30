package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"library/dao/common"
	"library/dao/config"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type LogMgr struct {
	programName string
	debug       bool
}

//TODO : 主进程不要太快结束, 不然看不到日志
func InitLog(cfg *config.LogCfg) *LogMgr {
	if cfg.ProgramName == "" {
		common.SimplePanic("ProgramName Is Empty")
	}
	SetLogger(cfg.ProgramName+"Info", NewLogger(cfg.ProgramName+"_info.log", cfg.Debug))
	SetLogger(cfg.ProgramName+"Error", NewLogger(cfg.ProgramName+"_err.log", cfg.Debug))
	SetLogger(cfg.ProgramName+"Panic", NewLogger(cfg.ProgramName+"_panic.log", cfg.Debug))
	return &LogMgr{
		programName: cfg.ProgramName,
		debug:       cfg.Debug,
	}
}

func (mgr *LogMgr) ZapCustom(level string) {
	SetLogger(mgr.programName+level, NewLogger(mgr.programName+"_"+level+".log", mgr.debug))
}

func (mgr *LogMgr) ZapSimpleLog(sign string, keysAndValues ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	keysAndValues = append(keysAndValues, file)
	keysAndValues = append(keysAndValues, line)
	GetLogger(mgr.programName+"Info").Infow(sign, keysAndValues...)
}

func (mgr *LogMgr) ZapErrorLog(sign string, keysAndValues ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	keysAndValues = append(keysAndValues, file)
	keysAndValues = append(keysAndValues, line)
	GetLogger(mgr.programName+"Error").Infow(sign, keysAndValues...)
}

func (mgr *LogMgr) ZapPanicLog(sign string, keysAndValues ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	keysAndValues = append(keysAndValues, file)
	keysAndValues = append(keysAndValues, line)
	GetLogger(mgr.programName+"Panic").Infow(sign, keysAndValues...)
}

func (mgr *LogMgr) ZapCustomLog(level string, sign string, keysAndValues ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	keysAndValues = append(keysAndValues, file)
	keysAndValues = append(keysAndValues, line)
	GetLogger(mgr.programName+level).Infow(sign, keysAndValues...)
}

//-----------------------------------------------------------------------------------------

var maploggers map[string]*zap.SugaredLogger

func init() {
	maploggers = make(map[string]*zap.SugaredLogger)
}

func GetWorkDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func GetLogger(logName string) *zap.SugaredLogger {
	if v, ok := maploggers[logName]; ok {
		return v
	}
	return nil
}

func SetLogger(logName string, logger *zap.SugaredLogger) {
	maploggers[logName] = logger
}

func NewLogger(fileName string, debug bool) *zap.SugaredLogger {

	tmp := GetLogger(fileName)
	if tmp != nil {
		return tmp
	}

	fileName = GetWorkDir() + "/logs/" + fileName
	fmt.Println("log path:", fileName)

	if tmp != nil {
		return tmp
	}

	//highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl >= zapcore.ErrorLevel
	//})
	//lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl < zapcore.ErrorLevel
	//})

	//consoleDebugging := zapcore.Lock(os.Stdout)

	//consoleErrors := zapcore.Lock(os.Stderr)

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    20,   // megabytes
		MaxBackups: 20,   //最多保留20个备份
		MaxAge:     7,    // days
		Compress:   true, //是否压缩备份文件
	})

	newDevelopmentEncoderConfig := zap.NewDevelopmentEncoderConfig()
	newDevelopmentEncoderConfig.TimeKey = "T"
	newDevelopmentEncoderConfig.CallerKey = "T"
	newDevelopmentEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if debug {
		core := zapcore.NewTee(
			//zapcore.NewCore(zapcore.NewConsoleEncoder(newDevelopmentEncoderConfig), consoleErrors, highPriority),
			//zapcore.NewCore(zapcore.NewConsoleEncoder(newDevelopmentEncoderConfig), consoleDebugging, highPriority),

			//zapcore.NewCore(zapcore.NewConsoleEncoder(newDevelopmentEncoderConfig), consoleDebugging, zapcore.InfoLevel),
			zapcore.NewCore(zapcore.NewJSONEncoder(newDevelopmentEncoderConfig), w, zapcore.InfoLevel),
		)
		logger := zap.New(core)
		return logger.Sugar()
	}

	newProductionEncoderConfig := zap.NewProductionEncoderConfig()
	newProductionEncoderConfig.TimeKey = "T"
	newProductionEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	newProductionEncoderConfig.CallerKey = "F"

	core := zapcore.NewTee(
		//zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), w, highPriority),
		zapcore.NewCore(zapcore.NewJSONEncoder(newProductionEncoderConfig), w, zapcore.ErrorLevel),
	)
	logger := zap.New(core, zap.AddCaller())
	return logger.Sugar()
}
