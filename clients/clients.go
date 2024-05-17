package clients

type Config struct {
	Gateway string
}

var Clients Config

func Init() {
	Clients = Config{
		Gateway: "http://localhost:5555",
	}
}
