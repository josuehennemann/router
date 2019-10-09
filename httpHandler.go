package main

import (
	"net/http"
)

func httpLogin(w http.ResponseWriter, r *http.Request, auth *User) {
	responseGenericSuccess(w, "Aqui responde do jeito que precisar")
}
func httpRecoveryPassword(w http.ResponseWriter, r *http.Request, auth *User) {
	responseGenericSuccess(w, []string{"Aqui responde do jeito que precisar"})
}
func httpRecoveryPasswordCheck(w http.ResponseWriter, r *http.Request, auth *User) {
	responseGenericSuccess(w, map[string]string{"Aqui": "responde do jeito que precisar"})
}
func httpUserGetInfo(w http.ResponseWriter, r *http.Request, auth *User) {
	responseGenericSuccess(w, auth)
}
func httpUserSaveInfo(w http.ResponseWriter, r *http.Request, auth *User) {

	responseGenericSuccess(w, 123213123)
}
