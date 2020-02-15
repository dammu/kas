package main

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ystyle/hcc/core"
	"github.com/ystyle/hcc/services"
	"golang.org/x/net/websocket"
	"net/http"
	"time"
)

func WS(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		wm := core.GetWsManager()
		wsClient := &core.WsClient{
			WsConn: ws,
			WsSend: make(chan core.Message, 10),
		}
		wm.Add(wsClient)
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func main() {
	log.EnableColor()
	log.SetLevel(log.DEBUG)

	wm := core.GetWsManager()
	wm.RegisterService("download", services.Submit)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	box := rice.MustFindBox("public")
	assetHandler := http.FileServer(box.HTTPBox())
	e.GET("/", echo.WrapHandler(assetHandler))
	e.GET("/asset/*", echo.WrapHandler(assetHandler))
	e.GET("/ws", WS)

	timer := time.NewTimer(time.Second * 5)
	go func() {
		for {
			select {
			case <-timer.C:
				timer.Reset(time.Second * 60)
				clients := len(wm.GetClients())
				fmt.Println("连接数为: ", clients)
			}
		}

	}()
	e.Logger.Fatal(e.Start(":1323"))
}