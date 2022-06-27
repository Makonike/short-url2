package controller

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"net/http"
	"short-url/object"
)

func ShortUrl(h *server.Hertz) {
	hl := hlog.DefaultLogger()
	h.GET("/short", func(c context.Context, ctx *app.RequestContext) {
		sUrl := object.ToShort(ctx.Query("url"))
		ctx.JSON(http.StatusOK, utils.H{"shortUrl": sUrl})
	})
	h.GET("/a/:url", func(c context.Context, ctx *app.RequestContext) {
		url := ctx.Param("url")
		// get the long url
		longUrl := object.GetLong(url)
		// set 302 code then redirect
		hl.Info(longUrl)
		ctx.Redirect(http.StatusFound, []byte(longUrl))
	})
}
