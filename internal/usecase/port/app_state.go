package port

type AppStateRepository interface {
	Get(key string) (string, error)
	Set(key, value string) error
}
