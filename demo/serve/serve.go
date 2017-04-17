package main

import (
	"log"

	"github.com/gunsluo/middleware/cors"
	"github.com/gunsluo/middleware/logger"
	"github.com/gunsluo/rrd"

	"gopkg.in/kataras/iris.v5"
)

var (
	tool *rrd.Tool
)

type ChartData struct {
	DSName    string  `json:"dsName"`
	Mark      string  `json:"mark"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
}

type ListResponse struct {
	Total int         `json:"total"` //总数
	Data  interface{} `json:"data"`  //
}

type Response struct {
	Code int         `json:"code"`           //响应状态码
	Msg  string      `json:"msg"`            //响应消息
	Body interface{} `json:"body,omitempty"` //响应体
}

func main() {

	initRRD()
	framwork := iris.New()
	framwork.Use(logger.New())
	framwork.Use(newCorsMiddleware())

	api := framwork.Party("/api/v1")
	{
		api.Post("/push", Push)
		api.Get("/chart/d", Chart)
	}

	framwork.Listen("0.0.0.0:8800")
}

func initRRD() {

	tool = rrd.NewTool().Cfg(&rrd.Config{
		RRDDir: "./rrd",
		DSS: []rrd.DSConfig{
			{
				Name: "demo",
				Type: rrd.DSTypes.Gauge,
				Step: 60,
				RRAS: []rrd.RRAConfig{
					{
						Type:  rrd.RRATypes.Average,
						Steps: 1,
						Rows:  720,
					},

					{
						Type:  rrd.RRATypes.Average,
						Steps: 5,
						Rows:  576,
					},
				},
			},
		},
	})
}

func newCorsMiddleware() *cors.Cors {

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"OPTIONS", "HEAD", "GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Origin", "X-Requested-With", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           3600,
	})

	return crs
}

func Push(ctx *iris.Context) {

	d := new(ChartData)
	if err := ctx.ReadJSON(d); err != nil {
		log.Printf("%s %s %s %s\n", err.Error(), ctx.Request.Body(), ctx.MethodString(), ctx.PathString())
		ctx.JSON(iris.StatusOK, &Response{Code: 500, Msg: err.Error()})
		return
	}

	items := []*rrd.Item{
		&rrd.Item{
			Value:     rrd.JsonFloat(d.Value),
			Timestamp: d.Timestamp,
		},
	}
	err := tool.Write(d.DSName, d.Mark, items)
	if err != nil {
		log.Printf("%s %s %s %s\n", err.Error(), ctx.Request.Body(), ctx.MethodString(), ctx.PathString())
		ctx.JSON(iris.StatusOK, &Response{Code: 500, Msg: err.Error()})
		return
	}

	ctx.JSON(iris.StatusOK, &Response{Code: 200, Body: true})
}

// Chart chart/d?ds=?&type=?&mark=?&start=?&end=?&step=?
func Chart(ctx *iris.Context) {

	ds := ctx.URLParam("ds")
	typ := ctx.URLParam("type")
	mark := ctx.URLParam("mark")

	start, err := ctx.URLParamInt64("start")
	if err != nil {
		log.Printf("%s %s %s\n", err.Error(), ctx.MethodString(), ctx.PathString())
		ctx.JSON(iris.StatusOK, &Response{Code: 500, Msg: err.Error()})
		return
	}

	end, err := ctx.URLParamInt64("end")
	if err != nil {
		log.Printf("%s %s %s\n", err.Error(), ctx.MethodString(), ctx.PathString())
		ctx.JSON(iris.StatusOK, &Response{Code: 500, Msg: err.Error()})
		return
	}

	step, err := ctx.URLParamInt("step")
	if err != nil {
		log.Printf("%s %s %s\n", err.Error(), ctx.MethodString(), ctx.PathString())
		ctx.JSON(iris.StatusOK, &Response{Code: 500, Msg: err.Error()})
		return
	}

	itemsRet, err := tool.Fetch(typ, ds, mark, rrd.Unix(start), rrd.Unix(end), step)
	if err != nil {
		log.Printf("%s %s %s\n", err.Error(), ctx.MethodString(), ctx.PathString())
		ctx.JSON(iris.StatusOK, &Response{Code: 500, Msg: err.Error()})
		return
	}

	ctx.JSON(iris.StatusOK, &Response{Code: 200, Body: &ListResponse{Total: len(itemsRet), Data: itemsRet}})
}
