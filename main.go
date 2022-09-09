package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	settings := &Settings{
		ServerAddress: "0.0.0.0:9000",
		Description:   "Simple Go Server",
		MaxPlayers:    2137,
	}

	listener, err := net.Listen("tcp", settings.ServerAddress)
	if err != nil {
		log.Fatalln(err)
	}

	world := NewWorld(settings, listener)

	go func() {
		shutdownSignalsChannel := make(chan os.Signal, 1)
		signal.Notify(shutdownSignalsChannel, syscall.SIGINT, syscall.SIGTERM)

		<-shutdownSignalsChannel

		world.Shutdown()
	}()

	log.Println("server listening")

	for {
		connection, err := listener.Accept()
		if err != nil {
			if netOpError, ok := err.(*net.OpError); ok {
				if netOpError.Err.Error() == "use of closed network connection" {
					break
				}
			}

			log.Fatalln(err)
		}

		go handleConnection(world, connection)
	}

	log.Println("exiting")
}

func handleConnection(world *World, connection net.Conn) {
	player := NewPlayer(world, connection)
	world.RegisterPlayer(player)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("%v\n", r)
			player.Disconnect()
		}

		world.UnregisterPlayer(player)
	}()

	for {
		packetSize, err := ReadPacketSize(connection)
		if err != nil {
			if err == io.EOF {
				player.Disconnect()
				return
			}

			log.Printf("%v\n", err)
			player.Disconnect()
			return
		}

		if packetSize > MaxPacketSize {
			log.Println("invalid packet size")
			player.Disconnect()
			return
		}

		packetData := make([]byte, packetSize)
		_, err = connection.Read(packetData)
		if err != nil {
			log.Printf("%v\n", err)
			player.Disconnect()
			return
		}

		HandlePacket(player, packetData)
	}
}
