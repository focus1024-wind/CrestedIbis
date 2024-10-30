package gb28181_server

import (
	"fmt"
	"github.com/ghettovoice/gosip/sip"
	"net/http"
)

func (config *GB28181Config) SipNotifyHandler(req sip.Request, tx sip.ServerTransaction) {
	fmt.Println(sip.NOTIFY)
	_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", ""))
}
