package server

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/lukawaay/zitadel_test_app/internal/server/config"

	"github.com/gin-gonic/gin"
	"github.com/gwatts/gin-adapter"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
	sdk_oidc "github.com/zitadel/zitadel-go/v3/pkg/authentication/oidc"
)


//go:embed templates/*
var templateFS embed.FS

type oidcHandlerType = *sdk_oidc.UserInfoContext[*oidc.IDTokenClaims, *oidc.UserInfo]

func getAuth(logger *slog.Logger, ctx context.Context, cfg *config.Config, zt *zitadel.Zitadel) (*authentication.Authenticator[oidcHandlerType], error) {
	oidcAuth := sdk_oidc.DefaultAuthentication(cfg.ClientID, cfg.RedirectURI, cfg.Key)

	var err error
	var res *authentication.Authenticator[oidcHandlerType]

	for i := 1; i <= 20; i += 1 {
		res, err = authentication.New(ctx, zt, cfg.Key, oidcAuth)
		if err == nil {
			return res, nil
		}
		logger.Error(fmt.Sprintf("Failed to create authentication (attempt %d out of %d): %s", i, 20, err))
		time.Sleep(2 * time.Second)
	}

	return nil, err
}

func Start(cfg *config.Config, logger *slog.Logger) error {
	ctx := context.Background()

	var ztopt zitadel.Option
	if cfg.InstanceSecure {
		ztopt = zitadel.WithPort(cfg.InstancePort)
	} else {
		ztopt = zitadel.WithInsecure(fmt.Sprint(cfg.InstancePort))
	}
	zt := zitadel.New(cfg.InstanceDomain, ztopt)

	authN, err := getAuth(logger, ctx, cfg, zt)
	if err != nil {
		return err
	}

	r := gin.Default()

	tmpl, err := template.ParseFS(templateFS, "templates/*")
	if err != nil {
		return err
	}
	r.SetHTMLTemplate(tmpl)

	mw := authentication.Middleware(authN)

	r.Use(func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}

		logger.Error(fmt.Sprint(err))

		c.HTML(c.Writer.Status(), "error", gin.H {
			"error": err.Error(),
		})
	})

	r.Any(config.SDKAuthEndpoint + "/login", gin.WrapH(authN))
	r.Any(config.SDKAuthEndpoint + "/callback", gin.WrapH(authN))
	r.Any(config.SDKAuthEndpoint + "/logout", gin.WrapH(authN))
	
	r.GET(config.AppEndpoint, adapter.Wrap(mw.RequireAuthentication()), func(c *gin.Context) {
		authCtx := mw.Context(c.Request.Context())
		userInfo := authCtx.GetUserInfo()
		token := authCtx.GetTokens().AccessToken

		c.HTML(http.StatusOK, "app", gin.H {
			"user": userInfo,
			"token": token,
		})
	})

	return r.Run(fmt.Sprintf(":%d", cfg.Port))
}
