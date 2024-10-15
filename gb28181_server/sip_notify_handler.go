package gb28181_server

import (
	"github.com/ghettovoice/gosip/sip"
	"net/http"
)

func (config *GB28181Config) SipNotifyHandler(req sip.Request, tx sip.ServerTransaction) {
	_ = tx.Respond(sip.NewResponseFromRequest("", req, http.StatusOK, "OK", ""))
}
