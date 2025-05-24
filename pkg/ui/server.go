package ui

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"gitlab.com/Hamed1984/provider/pkg/db"
)

type Resp struct {
	Msg  string      `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

func StartServer() {

	dbs := make(map[string]*db.StatsDB)
	dbs["hamed"] = db.NewStatsDB()
	dbs["ali"] = db.NewStatsDB()
	dbs["ahmad"] = db.NewStatsDB()
	dbs["erfan"] = db.NewStatsDB()
	dbs["ilia"] = db.NewStatsDB()

	e := echo.New()
	e.Use(middleware.Logger())

	e.POST("/v1/status/update/:provider/:id", func(c echo.Context) error {
		p := c.Param("provider")
		idStr := c.Param("id")
		id, _ := strconv.ParseUint(idStr, 10, 64)
		stat := new(db.OrderStatus)
		if err := c.Bind(stat); err != nil {
			log.Error(err.Error())
			resp := &Resp{
				Msg: "internal server error",
			}
			return c.JSON(http.StatusInternalServerError, resp)
		}
		pDb := dbs[p]
		if pDb == nil {
			resp := &Resp{
				Msg: "provider not found",
			}
			return c.JSON(http.StatusNotFound, resp)
		}
		pDb.UpdateStatus(id, stat)
		resp := &Resp{
			Msg: "success",
		}
		return c.JSON(http.StatusOK, resp)
	})

	e.GET("/v1/status/:provider/:id", func(c echo.Context) error {
		p := c.Param("provider")
		idStr := c.Param("id")
		id, _ := strconv.ParseUint(idStr, 10, 64)
		pDb := dbs[p]
		if pDb == nil {
			resp := &Resp{
				Msg: "provider not found",
			}
			return c.JSON(http.StatusNotFound, resp)
		}
		res := pDb.GetOrderStatuses(id)
		resp := &Resp{
			Msg:  "success",
			Data: res,
		}
		return c.JSON(http.StatusOK, resp)

	})

	server := &http.Server{
		Addr: ":9090",
	}

	go func() {
		err := e.StartServer(server)
		if err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

}
