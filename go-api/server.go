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

const cAddr = "0x7bEd003f8F0f3F78F4c999DC852CB413472295BF"
const netUrl = "HTTP://127.0.0.1:7545"
const wsUrl = "ws://127.0.0.1:7545"

var (
	userErc20Adapter = outputadapter.NewUserERC20Adapter(netUrl, cAddr)
	userRepo         = outputadapter.NewUserFirestoreRepo()
	userServ         = user.NewUserService(userRepo, userErc20Adapter)
	userController   = inputadapter.NewUserController(userServ)

	holdErc20Adapter = outputadapter.NewHoldERC20Adapter(netUrl, wsUrl, cAddr)
	holdRepo         = outputadapter.NewHoldFirestoreRepo()
	holdServ         = hold.NewHoldService(holdRepo, userServ, holdErc20Adapter)
	holdController   = inputadapter.NewHoldController(holdServ)

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
