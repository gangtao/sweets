package config

type KVStore interface {
	GetConfig(dataId string, group string) (string, error)
	PublishConfig(dataId string, group string, content string) error
	DeleteConfig(dataId string, group string) error
	ListenConfig(dataId string, group string, timeout int) (string, error)
}
