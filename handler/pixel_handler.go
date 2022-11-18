package handler

import (
	cfg "biscotti/config"
	"github.com/nu7hatch/gouuid"
	"github.com/valyala/fasthttp"
	"net/url"
)

func setAdTokenCookie(key string) *fasthttp.Cookie {
	cookie := &fasthttp.Cookie{}
	cookie.SetKey(cfg.GetAppConfig().CookieKey)
	cookie.SetValue(key)
	cookie.SetMaxAge(10 * 365 * 24 * 3600)
	cookie.SetHTTPOnly(true)
	cookie.SetSecure(true)
	cookie.SetDomain(cfg.GetAppConfig().CookieDomain)
	cookie.SetSameSite(fasthttp.CookieSameSiteNoneMode)
	return cookie
}

func PixelHandler(ctx *fasthttp.RequestCtx) {
	uri := ctx.Request.URI()
	requestUri, _ := url.Parse(uri.String())

	networkID := string(ctx.QueryArgs().Peek("id"))
	redirectUrl, networkIsValid := cfg.GetAppConfig().NetworkToUrls[networkID]
	redirectUri, _ := url.Parse(redirectUrl)

	if !networkIsValid {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	adToken := string(ctx.Request.Header.Cookie(cfg.GetAppConfig().CookieKey))
	if adToken == "" {
		u, _ := uuid.NewV4()
		cookie := setAdTokenCookie(u.String())
		ctx.Response.Header.SetCookie(cookie)
		adToken = u.String()
	}
	q := redirectUri.Query()
	q.Set(cfg.GetAppConfig().PixelRedirectKey, adToken)
	for k, v := range requestUri.Query() {
		if k != "id" {
			for _, vv := range v {
				q.Add(k, vv)
			}
		}
	}
	redirectUri.RawQuery = q.Encode()
	ctx.Redirect(redirectUri.String(), fasthttp.StatusFound)
}
