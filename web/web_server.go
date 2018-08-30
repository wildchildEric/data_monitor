package web

import (
	"data_monitor/dtuserver"
	"strconv"
	"time"

	"github.com/iris-contrib/middleware/basicauth"
	"github.com/iris-contrib/middleware/cors"
	"github.com/iris-contrib/middleware/logger"
	"github.com/kataras/go-template/html"
	"github.com/kataras/iris"
)

/*StartServe ...*/
func StartServe(port int) {
	iris.Config.IsDevelopment = true
	//Templates Config
	htmlConfig := html.Config{Layout: "layouts/base.html"}
	htmlTempl := html.New(htmlConfig)
	iris.
		UseTemplate(htmlTempl).
		Directory("./web/templates", ".html")

	//Static files config
	iris.Static("/public", "./web/public/", 1)

	// Root level middleware
	iris.Use(logger.New())
	iris.Use(cors.Default())
	authConfig := basicauth.Config{
		Users: map[string]string{
			"admin":            "admin",
			"mySecondusername": "mySecondpassword"},
		Realm:      "Authorization Required", // if you don't set it it's "Authorization Required"
		ContextKey: "mycustomkey",            // if you don't set it it's "user"
		Expires:    time.Duration(30) * time.Minute,
	}
	auth := basicauth.New(authConfig)
	iris.Use(auth)

	iris.Get("/", index)
	iris.Listen(":" + strconv.Itoa(port))
}

func index(ctx *iris.Context) {
	dtus := dtuserver.MDtu.Dtus()
	data := map[string]interface{}{"dtus": dtus}
	ctx.MustRender("index.html", data)
}
