package connection

import (
	"net"
	"net/http"

	"github.com/gobwas/ws"
)

func UpgradeToWS(w http.ResponseWriter, r *http.Request) (net.Conn, error) {
	var upgrader = ws.HTTPUpgrader{}
	conn, _, _, err := upgrader.Upgrade(r, w)
	return conn, err
}
