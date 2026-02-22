package log

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value any
}

func NewField(key string, value any) Field {
	return Field{
		Key:   key,
		Value: value,
	}
}
