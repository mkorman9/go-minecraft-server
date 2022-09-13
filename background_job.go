package main

import (
	"log"
	"time"
)

type BackgroundJob struct {
	world *World
}

func NewBackgroundJob(world *World) *BackgroundJob {
	return &BackgroundJob{
		world: world,
	}
}

func (bj *BackgroundJob) Start() {
	bj.startKeepAliveDaemon()
}

func (bj *BackgroundJob) startKeepAliveDaemon() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic in keep-alive job: %v, restarting\n", r)
				bj.startKeepAliveDaemon()
			}
		}()

		for {
			bj.world.BroadcastKeepAlive()
			time.Sleep(KeepAliveSendInterval)
		}
	}()
}
