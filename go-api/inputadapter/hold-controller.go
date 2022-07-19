package inputadapter

import (
	"cleanGo/api/usecase/hold"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type holdStuct struct {
	Id     int64  `json:"id"`
	Amount int64  `json:"amount"`
	User   string `json:"user"`
}

type HoldController interface {
	GetHolds(res http.ResponseWriter, req *http.Request)
	AddHold(res http.ResponseWriter, req *http.Request)
	GetHoldsByUser(res http.ResponseWriter, req *http.Request)
}

type holdCtrl struct {
	serv hold.HoldService
}

func NewHoldController(serv hold.HoldService) HoldController {
	return &holdCtrl{serv: serv}
}

func (hCtrl *holdCtrl) GetHolds(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	holds, err := hCtrl.serv.FindAllHolds()

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(ServiceError{Message: "Error getting holds."})
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(holds)
}

func (hCtrl *holdCtrl) AddHold(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var newHold holdStuct
	err := json.NewDecoder(req.Body).Decode(&newHold)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(ServiceError{Message: err.Error()})
		return
	}
	var usecaseHold hold.Hold
	usecaseHold.Amount = newHold.Amount
	usecaseHold.User = newHold.User

	newUsecaseHold, err1 := hCtrl.serv.CreateHold(&usecaseHold)
	if err1 != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(ServiceError{Message: err1.Error()})
		return
	}
	var outHold holdStuct
	outHold.Id = newUsecaseHold.Id
	outHold.Amount = newUsecaseHold.Amount
	outHold.User = newUsecaseHold.User

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(outHold)
}

func (hCtrl *holdCtrl) GetHoldsByUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	userAddr := vars["address"]

	holds, err := hCtrl.serv.FindHoldsFromUser(userAddr)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(ServiceError{Message: "Error getting holds."})
		return
	}

	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(holds)
}

type ServiceError struct {
	Message string `json:"message"`
}
