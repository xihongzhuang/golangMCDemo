package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	api_service "github.com/xihong-coding-exercise/internal/api-service"
)

func main() {
	iPort := 8080
	sPort := os.Getenv("PORT")
	if sPort != "" {
		i, err := strconv.Atoi(sPort)
		if err == nil {
			iPort = i
		}
	}
	//try to shut down server gracefully
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	instance := api_service.NewAPIServiceInstance(iPort)
	if err := instance.StartService(); err != nil {
		log.Fatal("failed to start service", err)
	}

	<-done
	log.Println("ctrl+C or sigTerm received, start shutdown server gracefully")
	instance.Shutdown(context.Background())
	log.Println("server shutdown gracefully")
}
