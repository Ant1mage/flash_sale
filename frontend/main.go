package main

import (
	"context"
	"flash-sale/dao"
	"flash-sale/frontend/web/controllers"
	"flash-sale/helper"
	"flash-sale/services"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"time"
)

func main() {
	app := iris.New()

	app.Logger().SetLevel("debug")

	template := iris.HTML("./frontend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)

	app.StaticWeb("/public", "./frontend/web/public")

	app.StaticWeb("/html", "./frontend/web/htmlProductShow")

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		_ = ctx.View("shared/error.html")
	})
	sess := sessions.New(sessions.Config{
		Cookie:  "flashSaleCookie",
		Expires: 600 * time.Minute,
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userRepository := repositories.NewUserDao(helper.InstanceDB())
	userService := services.NewUserService(userRepository)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx, sess)
	userPro.Handle(new(controllers.UserController))

	product := repositories.NewProductDao(helper.InstanceDB())
	productService := services.NewProductService(product)

	proProduct := app.Party("/product")
	pro := mvc.New(proProduct)
	order := repositories.NewOrderDao(helper.InstanceDB())
	orderService := services.NewOrderService(order)
	//proProduct.Use(middleware.AuthConProduct)
	pro.Register(productService, orderService)
	pro.Handle(new(controllers.ProductController))

	_ = app.Run(
		iris.Addr("0.0.0.0:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
