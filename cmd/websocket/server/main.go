package main

import (
	"log"
	"net/http"
	"syscall"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/nht1206/go-study/websocket"
)

var epoller *websocket.Epoll

func main() {
	// Increase the limit of open file descriptors
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	// Initialize the epoll instance
	var err error
	epoller, err = websocket.NewEPoll()
	if err != nil {
		panic(err)
	}

	// Handle WebSocket upgrades
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			return
		}

		// Add the connection to the epoll instance
		if err := epoller.Add(conn); err != nil {
			log.Printf("Failed to add connection: %v", err)
			conn.Close()
		}
	})

	go Start()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func Start() {
	for {
		// Wait for events from the epoll instance
		connections, err := epoller.Wait()
		if err != nil {
			log.Printf("Failed to epoll wait: %v", err)
			continue
		}

		// Process incoming WebSocket messages from connections
		for _, conn := range connections {
			if conn == nil {
				break
			}
			msg, _, err := wsutil.ReadClientData(conn)
			if err != nil {
				// Remove the connection from the epoll instance if there's an error
				if err := epoller.Remove(conn); err != nil {
					log.Printf("Failed to remove: %v", err)
				}
				conn.Close()
			} else {
				log.Printf("msg: %s", string(msg))
			}
		}
	}
}
