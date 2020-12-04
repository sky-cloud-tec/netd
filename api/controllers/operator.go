package controllers

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/sky-cloud-tec/netd/cli"
	"github.com/sky-cloud-tec/proto/v1/common"
	"github.com/sky-cloud-tec/proto/v1/jrpc"
	"github.com/songtianyi/rrframework/logs"
	"strings"
)

func OperatorHotfix(c *gin.Context) {
	var req jrpc.OperatorHotfixRequest
	if err := c.ShouldBind(&req); err != nil {
		errResponse(c, common.Retcode_BAD_REQUEST, err)
		return
	}
	logs.Info("[hostfix]", req)
	t := strings.Join([]string{req.Vendor, req.Type, req.Version}, ".")
	op := cli.OperatorManagerInstance.Get(t)
	if op == nil {
		errResponse(c, common.Retcode_BAD_REQUEST, fmt.Errorf("no operator match %s", t))
		return
	}
	op.SetErrPatterns(fixRegex(&req, req.Errs, op.GetErrPatterns()))
	op.SetPrompts(req.Mode, fixRegex(&req, req.Prompts, op.GetPrompts(req.Mode)))
	c.JSON(http.StatusOK, &jrpc.IResponse{Code: common.Retcode_OK, Msg: "OK"})
}

func OperatorDump(c *gin.Context) {
	var req jrpc.OperatorHotfixRequest
	if err := c.ShouldBind(&req); err != nil {
		errResponse(c, common.Retcode_BAD_REQUEST, err)
		return
	}
	t := strings.Join([]string{req.Vendor, req.Type, req.Version}, ".")
	op := cli.OperatorManagerInstance.Get(t)
	if op == nil {
		errResponse(c, common.Retcode_BAD_REQUEST, fmt.Errorf("no operator match %s", t))
		return
	}
	for i, p := range op.GetPrompts(req.Mode) {
		logs.Info("[HOTFIX]", "prompt ->", i, p)
	}
	for i, p := range op.GetErrPatterns() {
		logs.Info("[HOTFIX]", "err pattern ->", i, p)
	}
	c.JSON(http.StatusOK, &jrpc.IResponse{Code: common.Retcode_OK, Msg: "OK"})
}

func fixRegex(req *jrpc.OperatorHotfixRequest, x []string, regs []*regexp.Regexp) []*regexp.Regexp {
	if x == nil {
		return regs
	}
	if regs == nil || req.FixType == jrpc.FixType_REPLACE_ALL {
		// int or reset old
		regs = make([]*regexp.Regexp, 0)
	}
	switch req.FixType {
	case jrpc.FixType_APPEND, jrpc.FixType_REPLACE_ALL:
		for _, v := range x {
			regs = append(regs, regexp.MustCompile(v))
		}
	default:
		// do nothing
		logs.Error("[HOTFIX]", fmt.Sprintf("fix type %s not implemented yet", req.FixType))
	}
	return regs
}

func errResponse(c *gin.Context, code common.Retcode, err error) {
	c.JSON(http.StatusOK, &jrpc.IResponse{
		Code: common.Retcode_DATABASE_OP_FAILED,
		Msg:  err.Error(),
	})

}
