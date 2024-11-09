package socket

import (
	"log"
	"net"
	"os"
)

var socketPath = "/tmp/backup-service-socket.sock"

type ConnFunc func(net.Conn)

func SocketStart(connFuncs ...ConnFunc) {
	os.Remove(socketPath)
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal("Ошибка создания сокета:", err)
	}
	defer listener.Close()
	log.Println("Сокет запущен на", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Ошибка принятия соединения:", err)
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			for _, cf := range connFuncs {
				cf(c)
			}
		}(conn)
	}
}

func TriggerSocket() {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Printf("Не удалось подключится к сокету: %v", err)
		return
	}
	defer conn.Close()
}
