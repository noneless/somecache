package master

type Getter interface {
	Get(string) ([]byte, error)
}

var getter Getter

func RegisterGetter(g Getter) {
	getter = g
}

func Get(k string) ([]byte, error) {
	return service.Get(k)
}
