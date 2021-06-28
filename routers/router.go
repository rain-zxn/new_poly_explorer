package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"new_poly_explorer/controllers"
)

func init() {
	ns := beego.NewNamespace("/api/v1",
		beego.NSRouter("/getCrossTx", &controllers.ExplorerController{}, "get:GetCrossTx"),
		beego.NSRouter("/getcrosstxlist/", &controllers.ExplorerController{}, "post:GetCrossTxList"),
		beego.NSRouter("/getexplorerinfo/", &controllers.ExplorerController{}, "post:GetExplorerInfo"),
		beego.NSRouter("/gettokentxlist/", &controllers.ExplorerController{}, "post:GetTokenTxList"),
		beego.NSRouter("/getaddresstxlist/", &controllers.ExplorerController{}, "post:GetAddressTxList"),
	)
	beego.AddNamespace(ns)
	beego.Router("/", &controllers.MainController{})

}
