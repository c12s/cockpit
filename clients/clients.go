package clients

type Config struct {
	Magnetar string
	Kuiper   string
	Oort     string
}

var Clients Config

func Init() {
	Clients = Config{
		Magnetar: "http://localhost:5555",
		Kuiper:   "http://localhost:9001",
		Oort:     "http://localhost:8000",
	}
}
