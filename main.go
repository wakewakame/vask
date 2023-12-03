package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"

	"vask/internal/db"
	"vask/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	// 引数解析
	var (
		bind = flag.String("b", "127.0.0.1", "bind to this address")
		port = flag.Int("p", 8080, "bind to this port")
	)
	flag.Parse()

	// http server 初期化
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetLevel(log.INFO)

	// db 読み込み
	db, err := db.Open("./projects.db")
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	// ハンドラ追加
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(301, "/static/index.html")
	})
	e.Static("/static", "assets")
	api := e.Group("/api")
	api.GET("/project", func(c echo.Context) error {
		projects, err := db.GetProjects()
		if err != nil {
			return c.JSON(500, "{}")
		}
		projectsJson, err := json.Marshal(projects)
		if err != nil {
			return c.JSON(500, "{}")
		}
		return c.JSONBlob(200, projectsJson)
	})
	api.POST("/project", func(c echo.Context) error {
		project := model.Project{}
		err := json.NewDecoder(c.Request().Body).Decode(&project)
		if err != nil {
			return c.JSON(400, "{}")
		}
		id, err := db.AddProject(project.Name)
		if err != nil {
			return c.JSON(500, "{}")
		}
		return c.JSON(200, map[string]interface{}{"status": "ok", "id": id})
	})
	api.GET("/project/:id", func(c echo.Context) error {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(400, "{}")
		}
		project, err := db.GetProject(id)
		if err != nil {
			return c.JSON(500, "{}")
		}
		projectJson, err := json.Marshal(project)
		if err != nil {
			return c.JSON(500, "{}")
		}
		return c.JSONBlob(200, projectJson)
	})
	api.PUT("/project/:id", func(c echo.Context) error {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(400, "{}")
		}
		project := model.Project{}
		err = json.NewDecoder(c.Request().Body).Decode(&project)
		if err != nil {
			return c.JSON(400, "{}")
		}
		err = db.SetProject(id, project.Name)
		if err != nil {
			return c.JSON(500, "{}")
		}
		return c.JSON(200, map[string]string{"status": "ok"})
	})
	api.DELETE("/project/:id", func(c echo.Context) error {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(400, "{}")
		}
		err = db.DeleteProject(id)
		if err != nil {
			return c.JSON(500, "{}")
		}
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// http server 起動
	e.Logger.Infof("start http://%s:%d", *bind, *port)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", *bind, *port)))
}
