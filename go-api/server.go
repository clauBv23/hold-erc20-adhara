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
	userRepo       = outputadapter.NewUserFirestoreRepo()
	userServ       = user.NewUserService(userRepo)
	userController = inputadapter.NewUserController(userServ)

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
	httpRouter.GET("/balance/{id}", userController.GetUserBalance)
	httpRouter.POST("/reg", userController.AddUser)

	httpRouter.SERVE(port)

}
