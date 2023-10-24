package clients

import (
	kuiperapi "github.com/c12s/kuiper/pkg/api"
	magnetarapi "github.com/c12s/magnetar/pkg/api"
	oortapi "github.com/c12s/oort/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var Magnetar magnetarapi.MagnetarClient
var Kuiper kuiperapi.KuiperClient
var OortAdministrator oortapi.OortAdministratorClient

func Init() {
	Magnetar = newMagnetar("localhost:5000")
	Kuiper = newKuiper("localhost:9001")
	OortAdministrator = newOortAdministrator("localhost:8000")
}

func newMagnetar(address string) magnetarapi.MagnetarClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	return magnetarapi.NewMagnetarClient(conn)
}

func newKuiper(address string) kuiperapi.KuiperClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	return kuiperapi.NewKuiperClient(conn)
}

func newOortAdministrator(address string) oortapi.OortAdministratorClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	return oortapi.NewOortAdministratorClient(conn)
}
