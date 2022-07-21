package inputadapter

import (
	"cleanGo/api/usecase/user"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type userStruct struct {
	Id      int64  `json:"id"`
	Address string `json:"address"`
}

type UserController interface {
	GetUsers(res http.ResponseWriter, req *http.Request)
	AddUser(res http.ResponseWriter, req *http.Request)
	GetUserBalance(res http.ResponseWriter, req *http.Request)
}

type userCtrl struct {
	serv user.UserService
}

func NewUserController(serv user.UserService) UserController {
	return &userCtrl{serv: serv}
}

func (uCtrl *userCtrl) GetUsers(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")

	users, err := uCtrl.serv.FindAllUsers()

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(ServiceError{Message: "Error getting users."})
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(users)
}

func (uCtrl *userCtrl) AddUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var newUser userStruct
	err := json.NewDecoder(req.Body).Decode(&newUser)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(ServiceError{Message: err.Error()})
		return
	}
	var usecaseUser user.User
	usecaseUser.Address = newUser.Address

	newUsecaseUser, err1 := uCtrl.serv.CreateUser(&usecaseUser)
	if err1 != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(ServiceError{Message: err1.Error()})
		return
	}
	var outUser userStruct
	outUser.Id = newUsecaseUser.Id
	outUser.Address = newUsecaseUser.Address

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(outUser)
}

func (uCtrl *userCtrl) GetUserBalance(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	// todo change the param to id instead of addr
	//userId, err := strconv.ParseInt(vars["id"], 10, 64)
	//if err != nil {
	//	res.WriteHeader(http.StatusInternalServerError)
	//	json.NewEncoder(res).Encode(ServiceError{Message: "Error getting user Id."})
	//	return
	//}

	addr := vars["addr"]

	balance, err := uCtrl.serv.FindUserBalance(addr)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(ServiceError{Message: "Error getting balance."})
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(balance)
}
