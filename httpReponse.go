package main

import (
	"encoding/json"
	"github.com/josuehennemann/logger"
	"net/http"
	"strings"
)

type ApiResponseCode struct {
	code string
	text string
}

func NewApiResponseCode(c, t string) *ApiResponseCode {
	return &ApiResponseCode{code: c, text: t}
}
func (this *ApiResponseCode) Replace(m map[string]string) (t string) {
	//se não tem nada para fazer então sai fora
	t = this.text
	if m == nil {
		return
	}
	// as devidas trocas
	for k, v := range m {
		t = strings.ReplaceAll(t, "{{_"+k+"_}}", v)
	}
	return

}
func (this *ApiResponseCode) String() string {
	return "[" + this.code + " - " + this.text + "]"
}
func (this *ApiResponseCode) Code() string {
	return this.code
}
func (this *ApiResponseCode) Error() string {
	return this.text
}
func (this *ApiResponseCode) Response() (string, string) {
	return this.code, this.text
}

var _genericSuccess = NewApiResponseCode("1001", "Success")
var _invalidUserPassword = NewApiResponseCode("1002", "Usuário ou senha inválido")
var _invalidPassword = NewApiResponseCode("1003", "Senha inválida")
var _invalidName = NewApiResponseCode("1004", "Nome inválido")
var _invalidField = NewApiResponseCode("1005", "O campo {{_Field_}} é inválido")
var _errInvalidToken = NewApiResponseCode("1006", "Token inválido")
var _systemShutdown = NewApiResponseCode("1007", "System shutdown")
var _internalError = NewApiResponseCode("1008", "Ocorreu um erro interno. Tente novamente mais tarde ou entre em contato com nosso suporte")
var _dontHavePermission = NewApiResponseCode("1009", "Você não tem permissão para esta ação")

func responseInternalError(w http.ResponseWriter) {
	responseGenericError(w, _internalError, nil)
}

func responseDontHavePermission(w http.ResponseWriter) {
	responseGenericError(w, _dontHavePermission, nil)
}

func responseGenericError(w http.ResponseWriter, e *ApiResponseCode, replace map[string]string) {
	t := &StdResponse{Success: false, ErrorText: e.Replace(replace), Code: e.Code()}
	t.Response(w)
}

func responseGenericSuccess(w http.ResponseWriter, b interface{}) {
	t := &StdResponse{Success: true, Body: b, Code: "1000"}
	t.Response(w)
}

type StdResponse struct {
	Success   bool        `json:"success"`
	ErrorText string      `json:"error,omitempty"`
	Body      interface{} `json:"body"`
	Code      string      `json:"code"`
}

func (t *StdResponse) Response(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write(t.encode())
}

//faz o encode da resposta
func (t *StdResponse) encode() []byte {
	b, e := json.Marshal(t)
	if e != nil {
		Logger.Printf(logger.ERROR, "Falha gerar resposta da requisição [%s]", e)
		t.Success = false
		t.ErrorText = "Internal error (666)"
		t.Body = nil
		b, _ = json.Marshal(t)
	}

	return b
}
