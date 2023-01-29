package log

import "os"

// 全局默认提供一个 Log 对外句柄，可以直接使用 API 系列调用
// SwinxLog 创建全局 log
var SwinxLog = NewLogger(os.Stderr, "", BitDefault)

// 因为 SwinxLog 对象 对所有输出方法做了一层包裹，所以在打印调用函数的时候，比正常的 logger 对象多一层调用
// 一般的 Logger 对象 calldDepth=2, SwinxLog 的 calldDepth=3
func init() {
	SwinxLog.calldDepth = 3
}

// Flags 获取 SwinxLog 标记位
func Flags() int {
	return SwinxLog.Flags()
}

//  ResetFlags 设置 SwinxLog 标记位
func ResetFlags(flag int) {
	SwinxLog.ResetFlags(flag)
}

// AddFlag 添加 flag 标记
func AddFlag(flag int) {
	SwinxLog.AddFlag(flag)
}

// SetPrefix 设置 SwinxLog 日志头前缀
func SetPrefix(prefix string) {
	SwinxLog.SetPrefix(prefix)
}

// SetLogFile 设置 SwinxLog 绑定的日志文件
func SetLogFile(fileDir string, fileName string) {
	SwinxLog.SetLogFile(fileDir, fileName)
}

// CloseDebug 设置关闭 debug
func CloseDebug() {
	SwinxLog.CloseDebug()
}

// OpenDebug 设置打开 debug
func OpenDebug() {
	SwinxLog.OpenDebug()
}

func Debugf(format string, v ...interface{}) {
	SwinxLog.Debugf(format, v...)
}

func Debug(v ...interface{}) {
	SwinxLog.Debug(v...)
}

func Infof(format string, v ...interface{}) {
	SwinxLog.Infof(format, v...)
}

func Info(v ...interface{}) {
	SwinxLog.Info(v...)
}

func Warnf(format string, v ...interface{}) {
	SwinxLog.Warnf(format, v...)
}

func Warn(v ...interface{}) {
	SwinxLog.Warn(v...)
}

func Errorf(format string, v ...interface{}) {
	SwinxLog.Errorf(format, v...)
}

func Error(v ...interface{}) {
	SwinxLog.Error(v...)
}

func Fatalf(format string, v ...interface{}) {
	SwinxLog.Fatalf(format, v...)
}

func Fatal(v ...interface{}) {
	SwinxLog.Fatal(v...)
}

func Panicf(format string, v ...interface{}) {
	SwinxLog.Panicf(format, v...)
}

func Panic(v ...interface{}) {
	SwinxLog.Panic(v...)
}

// ====> Stack  <====
func Stack(v ...interface{}) {
	SwinxLog.Stack(v...)
}
