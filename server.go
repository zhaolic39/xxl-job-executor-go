package xxl

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ExecutorServer interface {
	// Run 运行服务
	Run() error
}

type Standalone struct {
	ExecutorIp   string
	ExecutorPort string
	Exec         *Executor
}

func (e *Standalone) Run() error {
	// 创建路由器
	mux := http.NewServeMux()
	// 设置路由规则
	mux.HandleFunc("/run", e.runTask)
	mux.HandleFunc("/kill", e.killTask)
	//mux.HandleFunc("/log", e.taskLog)
	mux.HandleFunc("/beat", e.beat)
	mux.HandleFunc("/idleBeat", e.idleBeat)

	// 创建服务器
	address := e.ExecutorIp + ":" + e.ExecutorPort
	server := &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	// 监听端口并提供服务
	log.Println("Starting server at ")
	go server.ListenAndServe()
	//quit := make(chan os.Signal)
	//signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	//<-quit
	//e.registryRemove()
	return nil
}

func (e *Standalone) runTask(writer http.ResponseWriter, request *http.Request) {
	req, _ := ioutil.ReadAll(request.Body)
	param := RunReq{}
	_ = json.Unmarshal(req, &param)
	s := (*e.Exec).(Executor)
	ret := s.RunTask(param)
	msg, _ := json.Marshal(ret)
	writer.Write(msg)
}

func (e *Standalone) killTask(writer http.ResponseWriter, request *http.Request) {
	req, _ := ioutil.ReadAll(request.Body)
	param := killReq{}
	_ = json.Unmarshal(req, &param)
	s := (*e.Exec).(Executor)
	ret := s.KillTask(param)
	msg, _ := json.Marshal(ret)
	writer.Write(msg)
}

func (e *Standalone) beat(writer http.ResponseWriter, request *http.Request) {
	s := (*e.Exec).(Executor)
	ret := s.Beat()
	msg, _ := json.Marshal(ret)
	writer.Write(msg)
}

func (e *Standalone) idleBeat(writer http.ResponseWriter, request *http.Request) {
	req, _ := ioutil.ReadAll(request.Body)
	param := idleBeatReq{}
	_ = json.Unmarshal(req, &param)
	s := (*e.Exec).(Executor)
	ret := s.IdleBeat(param)

	msg, _ := json.Marshal(ret)
	writer.Write(msg)
}
