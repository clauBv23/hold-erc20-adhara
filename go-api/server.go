package main

import (
	"cleanGo/api/inputadapter"
	http_layer "cleanGo/api/inputinfra"
	"cleanGo/api/outputadapter"
	"cleanGo/api/usecase/hold"
	"cleanGo/api/usecase/user"
	"fmt"
	"net/http"
)

var (
	userErc20Adapter = outputadapter.NewERC20Adapter("HTTP://127.0.0.1:7545", "0x60aBBa473620FF9361A30C047e5E885038d6AA01")
	userRepo         = outputadapter.NewUserFirestoreRepo()
	userServ         = user.NewUserService(userRepo, userErc20Adapter)
	userController   = inputadapter.NewUserController(userServ)

	holdRepo       = outputadapter.NewHoldFirestoreRepo()
	holdServ       = hold.NewHoldService(holdRepo, userServ)
	holdController = inputadapter.NewHoldController(holdServ)

	httpRouter = http_layer.NewMuxRouter()
)

func main() {
	port := ":8000"

	httpRouter.GET("/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(response, "Up and running  ")
	})
	httpRouter.GET("/holds", holdController.GetHolds)
	httpRouter.GET("/user/{address}/holds/", holdController.GetHoldsByUser)
	httpRouter.POST("/bet", holdController.AddHold)

	httpRouter.GET("/users", userController.GetUsers)
	httpRouter.GET("/balance/{addr}", userController.GetUserBalance)
	httpRouter.POST("/reg", userController.AddUser)

	httpRouter.SERVE(port)

}
