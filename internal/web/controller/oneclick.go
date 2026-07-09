package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/mhsanaei/3x-ui/v3/internal/web/service"
	"github.com/mhsanaei/3x-ui/v3/internal/web/session"
)

// oneClickReality batch-generates VLESS + TCP + REALITY + Vision inbounds in a
// single call — the open-source "one-click config" feature. One request produces
// N fully-formed REALITY nodes, each with its own keypair, UUID and subId.
//
// The JSON body is optional; an empty body creates 10 nodes with default settings.
func (a *InboundController) oneClickReality(c *gin.Context) {
	var req service.OneClickRealityRequest
	// Body is optional — ignore bind errors (e.g. empty body) and use defaults.
	_ = c.ShouldBindJSON(&req)

	user := session.GetLoginUser(c)
	result, needRestart, err := a.inboundService.BatchCreateRealityVision(user.Id, req)
	if err != nil {
		jsonMsg(c, I18nWeb(c, "somethingWentWrong"), err)
		return
	}

	jsonObj(c, result, nil)
	if needRestart {
		a.xrayService.SetToNeedRestart()
	}
	a.broadcastInboundsUpdate(user.Id)
	notifyClientsChanged()
}
