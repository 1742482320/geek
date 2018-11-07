package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// StartHTTP StartHTTP
func StartHTTP(conf *Config) error {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.BasicAuth(func(user string, pass string, ctx echo.Context) (bool, error) {

		check := user + ":" + pass

		if check == conf.Admin {
			return true, nil
		}

		for i := range conf.HTTP.BasicAuth {
			ctx.Logger().Info(check, conf.HTTP.BasicAuth[i])
			if check == conf.HTTP.BasicAuth[i] {
				return true, nil
			}
		}

		return false, nil
	}))

	e.POST("/api/sync", func(c echo.Context) error {
		if atomic.LoadInt32(&downloading) == 1 {
			return c.String(http.StatusOK, "已运行同步，请稍等片刻")
		}
		go func() {
			if err := do(false); err != nil {
				fmt.Println(err)
			}
		}()
		return c.String(http.StatusOK, "starting")
	})
	// Routes
	e.Static("/", conf.DataDir)

	// Start server
	return e.Start(conf.HTTP.Listen)
}
