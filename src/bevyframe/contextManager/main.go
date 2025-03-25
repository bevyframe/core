package contextManager

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func contextManager() {
	if os.Getuid() != 0 {
		log.Fatal("context_manager must be run as root")
	}

	socket, err := net.Listen("unix", "/opt/bevyframe/context.sock")
	if err != nil {
		return
	}
	err = os.Chmod("/opt/bevyframe/context.sock", 0777)
	if err != nil {
		panic(err)
	}
	defer func(socket net.Listener) {
		err := socket.Close()
		if err != nil {
			panic(err)
		}
	}(socket)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		_ = os.Remove("/opt/bevyframe/context.sock")
		os.Exit(1)
	}()

	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(conn net.Conn) {
			defer func(conn net.Conn) {
				err := conn.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(conn)

			buf := make([]byte, 4096)
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			_, err = conn.Write(buf[:n])
			if err != nil {
				log.Fatal(err)
			}

		}(conn)
	}
}
