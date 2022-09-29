package main

import (
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	settings := &Settings{
		ServerAddress:         "0.0.0.0:25565",
		Description:           "Simple Go Server",
		MaxPlayers:            2137,
		OnlineMode:            false,
		CompressionThreshold:  -1,
		IsDebug:               true,
		ViewDistance:          10,
		SimulationDistance:    10,
		KeepAliveSendInterval: 5,
		PlayerTimeout:         15,
	}

	world, err := NewWorld(settings)
	if err != nil {
		log.Fatalln(err)
	}

	startSigintListener(world)

	log.Println("server listening")

	err = world.Server().AcceptLoop(func(connection net.Conn, ip string) {
		player := NewPlayer(world, ip)
		packetHandler := NewPlayerPacketHandler(player, world, connection, ip)
		player.AssignPacketHandler(packetHandler)

		log.Printf("player connected from %s\n", ip)

		packetHandler.ReadLoop()

		log.Println("player disconnected")
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("exiting")
}

func startSigintListener(world *World) {
	go func() {
		shutdownSignalsChannel := make(chan os.Signal, 1)
		signal.Notify(shutdownSignalsChannel, syscall.SIGINT, syscall.SIGTERM)

		<-shutdownSignalsChannel

		world.Server().Shutdown()
	}()
}
