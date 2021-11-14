package main

import (
	"context"
	"gopokemon/config"
	"gopokemon/db"
	"gopokemon/log"
	"gopokemon/router"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

// @title gopokemon service API
// @version 1.0
// @description gopokemon service API
// @BasePath /
func main() {
	var err error
	rand.Seed(time.Now().UnixNano())
	log.Run()

	//if err = constant.CreateUploadFolder(); err != nil {
	//	fmt.Println("Failed create folder", err.Error())
	//	os.Exit(1)
	//}

	dbpool := db.Initialize()
	systemCleaning := make(chan struct{}, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go scheduler(systemCleaning, &wg)

	var shutdownCallback func() = func() {
		wg.Add(1)
		close(systemCleaning)

		log.System.Warn().Msg("Cleaning resources!")
		dbpool.Close()
		// INFO: bersihin semua sebelum shutdown
		wg.Done()
	}

	r := router.Init()
	go func() {
		if config.Environment == config.PRODUCTION {
			r.Server.RegisterOnShutdown(shutdownCallback)
			if err = r.Start(":" + config.ListenTo.Port); err != nil && err != http.ErrServerClosed {
				r.Logger.Fatal("Shutting down the server")
			}
		} else {
			r.TLSServer.RegisterOnShutdown(shutdownCallback)
			if err = r.StartTLS(":"+config.ListenTo.Port, config.CertificateFilePath, config.CertificateKeyFilePath); err != nil && err != http.ErrServerClosed {
				r.Logger.Fatal("Shutting down the serverTLS")
			}

			//r.Server.RegisterOnShutdown(shutdownCallback)
			//if err = r.Start(":" + config.ListenTo.Port); err != nil && err != http.ErrServerClosed {
			//	r.Logger.Fatal("Shutting down the server")
			//}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.Shutdown(ctx); err != nil {
		r.Logger.Fatal(err)
	}

	wg.Wait()
	log.System.Warn().Msg("Main System Shutdown!")
	log.CloseAll()
}

func scheduler(systemCleaning chan struct{}, wg *sync.WaitGroup) {
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	midnight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.Local)
	everyMidnight := time.NewTimer(midnight.Sub(now))
	everyMinute := time.NewTicker(time.Minute)

DailyLoop:
	for {
		select {
		case <-everyMidnight.C:
			everyMidnight.Reset(24 * time.Hour)
			log.ChangeDay()
			runtime.GC()
		case <-everyMinute.C:
		case <-systemCleaning:
			if !everyMidnight.Stop() {
				<-everyMidnight.C
			}
			break DailyLoop
		}
	}
	log.System.Warn().Msg("Scheduler shutdown")
	wg.Done()
}