package clients

import magnetarapi "github.com/c12s/magnetar/pkg/api"

var Magnetar magnetarapi.MagnetarClient

func Init() {
	Magnetar = newMagnetar()
}
