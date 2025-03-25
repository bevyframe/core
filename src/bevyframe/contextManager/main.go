package contextManager

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func Run(packageName string) {

	socketAddr := fmt.Sprintf("/opt/bevyframe/sockets/%s", packageName)
	socket, err := net.Listen("unix", socketAddr)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chmod(socketAddr, 0700)
	if err != nil {
		log.Fatal(err)
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
		_ = os.Remove(socketAddr)
		os.Exit(1)
	}()

	context := map[string]Variable{}

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
			rawCommand := string(buf[:n])
			rawCommand = strings.TrimSuffix(rawCommand, "\n")
			args := strings.Split(rawCommand, " ")
			if args[0] == "get" {
				variable, ok := GetVar(context, args[1], args[2])
				if !ok {
					_, err = conn.Write([]byte("null 0"))
				} else {
					length := strconv.Itoa(len(variable.VarData))
					out := fmt.Sprintf("%s %s\n", variable.VarType, length)
					_, err = conn.Write([]byte(out))
					buf = make([]byte, 4096)
					_, err = conn.Read(buf)
					_, err = conn.Write(variable.VarData)
				}
			} else if args[0] == "set" {
				uAddr := args[1]
				varType := args[2]
				varName := args[3]
				length, _ := strconv.Atoi(args[4])
				var varData []byte
				_, err := conn.Write([]byte("OK"))
				if err != nil {
					return
				}
				for {
					buf = make([]byte, 4096)
					n, err = conn.Read(buf)
					if err != nil {
						err = conn.Close()
						if err != nil {
							return
						}
					}
					varData = append(varData, buf[:n]...)
					if len(varData) >= length {
						break
					}
				}
				ok := SetVar(&context, uAddr, varName, varType, varData)
				if ok {
					_, err = conn.Write([]byte(strconv.Itoa(len(varData))))
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			if err != nil {
				log.Fatal(err)
			}

		}(conn)
	}
}
