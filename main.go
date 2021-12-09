package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"todosvc"
	//"github.com/gorilla/mux"
)


func main(){
	httpAddr := flag.String("http.addr",":8000","http listen address")
	flag.Parse()
	s := todosvc.Service.NewInmemService()
	h := todosvc.Service.MakeHTTPHandler(s)
	errs := make(chan error)
	go func(){
		c:= make(chan os.Signal,1)
		signal.Notify(c, syscall.SIGINT,syscall.SIGTERM)
		errs<-fmt.Errorf("%s",<-c)
	}()
	go func ()  {
		log.Fatal(http.ListenAndServe(*httpAddr,h))
		
	}()

	log.Fatal("exit" ,<-errs)
	

}
