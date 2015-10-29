package interrupt

import (
	"github.com/tsuka611/golang_sandbox/log"
	"os"
	"os/signal"
	"syscall"
)

func SignalTrap() <-chan int {
	exit := make(chan int)
	go listenSignal(exit)
	return exit
}

func listenSignal(exit chan<- int) {
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGHUP, syscall.SIGINFO, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	for {
		sig := <-chSig
		switch sig {
		case syscall.SIGHUP:
			log.TRACE.Printlnf("SIGHUP received. %v", sig)
			exit <- 129
			return
		case syscall.SIGINT:
			log.TRACE.Printlnf("SIGINT received. %v", sig)
			exit <- 130
			return
		case syscall.SIGTERM:
			log.TRACE.Printlnf("SIGTERM received. %v", sig)
			exit <- 143
			return
		case syscall.SIGQUIT, syscall.SIGKILL:
			log.TRACE.Printlnf("SIGQUIT/SIGKILL received. %v", sig)
			exit <- 1
			return
		default:
			log.ERROR.Printlnf("Other Signal Received. %v", sig)
			exit <- 1
			return
		}
	}
}
