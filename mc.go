package mcClient

type MC interface {
	GetString(key string) (val string,err error)
	SetString(key string, val string) (err error)
	Get(key string, val interface{}) (err error)
	Set(key string, val interface{}) (err error)
	Delete(key string) (err error)
}
