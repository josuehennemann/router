package main

import (
	"github.com/josuehennemann/logger"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type MyHandler struct {
	Handler       func(w http.ResponseWriter, r *http.Request, userInfo *User)
	Authenticated bool
	Regexp        *regexp.Regexp
	OnlyAdmin     bool // referente ao usuário no cliente
	Type          string
}

const (
	REQUEST_NORMAL = "multipart/form-data"
	REQUEST_JSON   = "application/json"
)

func loadHandlers() {
	listHandlersHttp = map[string]*MyHandler{
		//métodos que não precisam de autenticação
		"/sign-in":                             &MyHandler{Handler: httpLogin, Authenticated: false, Type: REQUEST_JSON},
		"/recovery-password":                   &MyHandler{Handler: httpRecoveryPassword, Authenticated: false, Type: REQUEST_NORMAL},
		"/recovery-password/[hash]/[checksum]": &MyHandler{Handler: httpRecoveryPasswordCheck, Authenticated: false, Regexp: regexp.MustCompile("/recovery-password/[a-zA-Z0-9]+/[a-zA-Z0-9]+$"), Type: REQUEST_NORMAL},
		"/user/get":                            &MyHandler{Handler: httpUserGetInfo, Authenticated: true, Type: REQUEST_JSON},
		"/user/save":                           &MyHandler{Handler: httpUserSaveInfo, Authenticated: true, Type: REQUEST_JSON, OnlyAdmin: true},
	}
}

//struct personalizada, que herda a ResponseWrite nativa. Feito isso para conseguir pegar o Código http da request.
//Pois por padrão o Go não deixa esse código "exposto"
type MiddleEarth struct {
	http.ResponseWriter
	StatusCode int
}

// Escreve as 3 funções que a interface precisa ter
func (this *MiddleEarth) Header() http.Header {
	return this.ResponseWriter.Header()
}

func (this *MiddleEarth) Write(s []byte) (int, error) {
	return this.ResponseWriter.Write(s)
}
func (this *MiddleEarth) WriteHeader(statusCode int) {
	this.ResponseWriter.WriteHeader(statusCode)
	this.StatusCode = statusCode // aqui escreve no header o código http e tbm armazena em uma posição da struct
}

//Aqui é go repassa a requisição
func handlerHttp(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(sw http.ResponseWriter, r *http.Request) {
		w := &MiddleEarth{ResponseWriter: sw}
		defer PrepareClose(w, r)
		var user *User
		initTime := time.Now()
		defer func() {

			//caso tenha dado 404 ou for uma requisição "pública", então seta uma popula a variavel, para evitar panic
			if user == nil {
				user = &User{Code: "00"}
			}

			//grava o log de acesso no seguinte formato: url;cod_http;ip_acesso;tempo;user
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			Access.Printf(logger.ACCESS, "%s;%s;%d;%s;%0.6f;us_%s", r.Host, r.URL.Path, w.StatusCode, ip, time.Since(initTime).Seconds(), user.Code)
		}()

		uriNormalize := strings.ToLower(r.URL.Path)
		urlHandler, passUrl := listHandlersHttp[uriNormalize]

		//se nao achou nada no mapa, entao valida as que tem regexp
		if !passUrl {
			for _, v := range listHandlersHttp {
				//se é uma url com regexp então executa a regexp
				if v.Regexp != nil {
					passUrl = v.Regexp.MatchString(r.URL.Path)
				}

				//se achou alguma entao para
				if passUrl {
					urlHandler = v
					break
				}

			}
		}
		//se a url não bateu com nenhuma route, então da 404
		if !passUrl {
			NotFoundHandler(w, r)
			return
		}

		//se precisa ser authenticado, entao procura o token do usuário
		if urlHandler.Authenticated {
			user = getUserFromToken(r)

			//caso não consiga achar o usuário, retorna erro
			if user == nil {
				responseGenericError(w, _errInvalidToken, nil)
				return
			}

			//se é apenas para admin e não é um admin entao já retorna erro
			if urlHandler.OnlyAdmin && !user.Admin {
				responseDontHavePermission(w)
				return
			}

		}

		//executa a função do handler
		urlHandler.Handler(w, r, user)
	})
}

func PrepareClose(w http.ResponseWriter, r *http.Request) {
	r.Close = true
	r.Body.Close()
	return
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func startHttpServer() {
	// HTTP LISTEN
	loadHandlers()

	if tr, ok := http.DefaultTransport.(*http.Transport); ok {
		tr.DisableKeepAlives = true
		tr.MaxIdleConnsPerHost = 1
		tr.CloseIdleConnections()
	}
	Logger.Printf(logger.INFO, "Iniciando serviço [%s] ...", config.HttpAddress)

	//serviço na porta definida
	server := &http.Server{Addr: config.HttpAddress, ReadTimeout: 10 * time.Second, WriteTimeout: 20 * time.Second, Handler: handlerHttp(http.DefaultServeMux)}
	err := server.ListenAndServe()
	checkErrorAndKillMe(err)
}

type User struct {
	Code     string
	Name     string
	Email    string
	Password string
	Status   string
	Admin    bool
}

func getUserFromToken(r *http.Request) *User {
	return new(User)
}
