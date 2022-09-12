package main

import (
	"log"
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
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic in background job: %v, restarting\n", r)
				bj.Start()
			}
		}()

		// TODO
	}()
}
