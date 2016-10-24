package master

//	"fmt"

type Getter interface {
	Get(string) ([]byte, error)
}

var (
	getter Getter
)

func RegisterGetter(g Getter) {
	getter = g
}

func Get(k string) ([]byte, error) {
	//	if !validKey(k) {
	//		return nil, fmt.Errorf("key is valid")
	//	}
	return service.Get(k)
}

/*
//type ReaderAndSeekerAndCloser interface {
//	Read(p []byte) (n int, err error)
//	// Seek sets the offset for the next Read or Write on file to offset, interpreted
//	// according to whence: 0 means relative to the origin of the file, 1 means
//	// relative to the current offset, and 2 means relative to the end.
//	// It returns the new offset and an error, if any.
//	// The behavior of Seek on a file opened with O_APPEND is not specified.
//	Seek(offset int64, whence int) (ret int64, err error)
//	Close()
//}

//	streamGetter ReaderAndSeekerAndCloser
//func RegisterSteamGetter(g ReaderAndSeekerAndCloser) {
//	streamGetter = g
//}
*/
