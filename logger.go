package processor

import (
	"encoding/json"
	"log"
	"os"
	"runtime/debug"
)

type LoggerType int8

const (
	EMERGENCY     LoggerType = iota
	ALERT         LoggerType = iota
	CRITICAL      LoggerType = iota
	ERROR         LoggerType = iota
	WARNING       LoggerType = iota
	NOTICE        LoggerType = iota
	INFORMATIONAL LoggerType = iota
	DEBUG         LoggerType = iota
)

type LoggerConfig struct {
	Engine     *log.Logger
	LevelToLog LoggerType
}

type logMessage struct {
	Level      LoggerType
	Message    string
	StackTrace string
}

type logFunction struct {
	Level      LoggerType
	Message    string
	StackTrace string
	Caller     string
	Params     []interface{}
}

type logMethod struct {
	Level      LoggerType
	Message    string
	StackTrace string
	Object     string
	Caller     string
	Params     []interface{}
}

type logBool struct {
	Level      LoggerType
	Message    string
	StackTrace string
	Name       string
	Value      bool
}

type logGeneric struct {
	Level      LoggerType
	Message    string
	StackTrace string
	Data       []interface{}
}

type LoggerMessage struct {
	config     *LoggerConfig
	level      LoggerType // Log Level Type
	additional string     // Log Message
}

// Additional information for logging message.
func (msg *LoggerMessage) Message(message string) *LoggerMessage {
	msg.additional = message
	return msg
}

// No additional information is needed.
// Just log message with stack trace.
func (msg *LoggerMessage) None() {
	data := &logMessage{msg.level, msg.additional, string(debug.Stack())}
	msg.do(data)
}

// Log function with message and stack trace.
func (msg *LoggerMessage) Function(name string, params ...interface{}) {
	data := &logFunction{msg.level, msg.additional, string(debug.Stack()), name, params}
	msg.do(data)
}

// Log object method.
func (msg *LoggerMessage) Method(object, name string, params ...interface{}) {
	data := &logMethod{msg.level, msg.additional, string(debug.Stack()), object, name, params}
	msg.do(data)
}

// Log boolean.
func (msg *LoggerMessage) Bool(name string, value bool) {
	data := &logBool{msg.level, msg.additional, string(debug.Stack()), name, value}
	msg.do(data)
}

func (msg *LoggerMessage) Any(info ...interface{}) {
	data := &logGeneric{msg.level, msg.additional, string(debug.Stack()), info}
	msg.do(data)
}

func (msg *LoggerMessage) do(data interface{}) {
	raw, err := json.Marshal(data)
	if err != nil {
		msg.config.Engine.Fatalf("Logger can not create Function JSON message: %s", err)
	}
	message := string(raw)

	if msg.level > msg.config.LevelToLog {
		// Only log when level to log is less than or greater. Might be counter intuitive, but
		// that is just how it jelly rolls. Also other implementations appear to work this way as
		// well.
		return
	}

	if msg.level == EMERGENCY {
		// All hell has escaped and only the Hero can save us now!
		// GAME OVER.
		// Either use save or start from the beginning.
		msg.config.Engine.Fatalln(message)
		return // This is dead code. Fatalln will exit. Keep anyway for clarity of intent.
	}

	if msg.level < WARNING {
		// Not sure if Warning should included here.
		// Meant to print line in log.
		msg.config.Engine.Panicln(message)
		return // This is dead code. Panicln will panic, which will return execution to caller.
	}

	msg.config.Engine.Println(message)
}

type Logger struct {
	config *LoggerConfig // Used for logging.
}

func (l *Logger) Emergency() *LoggerMessage {
	return l.Level(EMERGENCY)
}

func (l *Logger) Alert() *LoggerMessage {
	return l.Level(ALERT)
}

func (l *Logger) Critical() *LoggerMessage {
	return l.Level(CRITICAL)
}

func (l *Logger) Error() *LoggerMessage {
	return l.Level(ERROR)
}

func (l *Logger) Warning() *LoggerMessage {
	return l.Level(WARNING)
}

func (l *Logger) Notice() *LoggerMessage {
	return l.Level(NOTICE)
}

func (l *Logger) Info() *LoggerMessage {
	return l.Level(INFORMATIONAL)
}

func (l *Logger) Debug() *LoggerMessage {
	return l.Level(DEBUG)
}

func (l *Logger) Level(level LoggerType) *LoggerMessage {
	message := new(LoggerMessage)
	message.config = l.config
	message.level = level
	return message
}

func NewLogger(config *LoggerConfig) *Logger {
	return &Logger{config}
}

func DefaultLogger() *Logger {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	config := &LoggerConfig{logger, ERROR}
	return &Logger{config}
}
