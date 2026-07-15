package controller

import (
	"errors"
	"io"

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
	// An empty body means "use the defaults" (io.EOF from the JSON decoder).
	// A body that is present but malformed is a mistake, not a request for
	// defaults: swallowing that error made a truncated
	// {"protocol":"ss2022","count":2, silently create 10 REALITY nodes and
	// still report success. Fail loudly instead.
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		jsonMsg(c, I18nWeb(c, "somethingWentWrong"), err)
		return
	}

	user := session.GetLoginUser(c)
	if user == nil {
		jsonMsg(c, I18nWeb(c, "somethingWentWrong"), errors.New("no login user"))
		return
	}
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
