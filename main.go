package main

import (
	"flag"
	"fmt"
	"github.com/josuehennemann/logger"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	//faz o parse da variaveis que sao passadas para o binario

	flag.Parse()
	initService()
	//trava a execução da main
	select {}
}

//verifica se deu erro ou nao, em caso de erro, printa o erro e mata o script
func checkErrorAndKillMe(e error) {
	if e == nil {
		return
	}
	fmt.Println(e.Error())
	os.Exit(2)
}

func initService() {
	ServiceName = "Router"

	var err error

	//inicia o processo para subir o binario como serviço no linux
	//	procDaemon, err = daemon.Daemonize(*pidPath)
	checkErrorAndKillMe(err)

	//cria a goroutine que sabe tratar o desligamento do serviço
	go killMeSignal()

	_init()

	//inicia o servidor http
	go startHttpServer()
}

//inicializa as variaveis que podem ser utilizadas caso rode sem ser serviço. Ex: caso execute algum go test
func _init() {
	//carrega para a memoria o arquivo de inicialização
	err := initConfig()
	checkErrorAndKillMe(err)
	setOutput()
	//inicia o arquivo de log
	Logger, err = logger.New(config.LogPath+ServiceName+".log", logger.LEVEL_PRODUCTION, true)
	checkErrorAndKillMe(err)
	Logger.SetRemoveAfter(30) // define em 30 dias para remover os logs

	Access, err = logger.New(config.LogPath+ServiceName+"-Access.log", logger.LEVEL_PRODUCTION, true)
	checkErrorAndKillMe(err)
}

func setOutput() {
	if config.IsDev {
		return
	}

	dirLog := filepath.Dir(config.LogPath)
	if dirLog == "." {
		dirLog = ""
	} else {
		dirLog += "/"
	}

	//IMPORTANTE ISSO NAO FUNCIONA NO WINDOWS
	daemon, err := os.OpenFile(dirLog+ServiceName+"-Daemon.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	checkErrorAndKillMe(err)
	syscall.Dup2(int(daemon.Fd()), 1)
	syscall.Dup2(int(daemon.Fd()), 2)
}

func killMeSignal() {
	c := make(chan os.Signal)
	signal.Notify(c)
	for {
		s := <-c
		if s == syscall.SIGTERM || s == syscall.SIGINT {
			Logger.Printf(logger.INFO, "Sinal de desligamento recebido")

			Logger.Close()
			os.Exit(0)
		}
	}
	return
}
