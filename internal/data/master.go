package data

type MasterQ interface {
	New() MasterQ
	ClientQ() ClientQ
	ApplicationQ() ApplicationQ
	LocationQ() LocationQ
	ApplicationStatusQ() ApplicationStatusQ
	Transaction(func(data interface{}) error, interface{}) error
}

type CacheQ interface {
	SetCode(key string, value string) error
	GetCode(key string) (string, error)
	DelCode(key string) error
}
