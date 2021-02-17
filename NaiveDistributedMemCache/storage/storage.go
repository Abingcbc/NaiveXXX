package storage

type Storage interface {
	Put(key string, value Value) bool
	Get(key string) (Value, bool)
}

type Entry struct {
	key   string
	value Value
}

type Value struct {
	Bytes []byte
}

func (value *Value) Len() int64 {
	return int64(len(value.Bytes))
}
