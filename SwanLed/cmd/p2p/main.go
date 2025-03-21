package p2p

import (
	"bytes"
	"encoding/gob"
	"github.com/briandowns/spinner"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var log = logrus.New()

type Peer struct {
	Conn *websocket.Conn
	Addr string
}

var (
	PeersList  = make(map[string]*Peer)
	PeerIpList = []string{}
	mu         sync.RWMutex
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections for testing
	},
}

func HandlePeerConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to upgrade connection: %v\n", err)
		return
	}
	defer func() {
		mu.Lock()
		delete(PeersList, conn.RemoteAddr().String())
		mu.Unlock()
		err := conn.Close()
		if err != nil {
			log.Printf("[ERROR] Failed to close connection: %v", err)
		}
		log.Printf("[INFO] Peer %s disconnected", conn.RemoteAddr().String())
	}()

	peerAddr := conn.RemoteAddr().String()

	mu.Lock()
	PeersList[peerAddr] = &Peer{
		Conn: conn,
		Addr: peerAddr,
	}
	mu.Unlock()

	log.Printf("[INFO] New peer connected: %s", peerAddr)

	// Set ping/pong handlers
	conn.SetPingHandler(func(message string) error {
		log.Printf("[INFO] Received ping from %s", peerAddr)
		return conn.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second))
	})

	conn.SetPongHandler(func(message string) error {
		log.Printf("[INFO] Received pong from %s", peerAddr)
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[ERROR] Unexpected close error from %s: %v", peerAddr, err)
			} else {
				log.Printf("[ERROR] Read error from %s: %v (Type: %T)", peerAddr, err, err)
			}
			break
		}

		if len(message) == 0 {
			log.Printf("[WARN] Empty message received from %s", peerAddr)
			continue
		}

		log.Printf("[MESSAGE] Received from %s: %s", peerAddr, string(message))
		go handleMessage(message, conn)
	}
}

func handleMessage(message []byte, conn *websocket.Conn) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recovered from panic in handleMessage: %v", r)
		}
	}()

	var msgType string
	err := Decode(message, &msgType)
	if err != nil {
		log.Printf("[ERROR] Failed to decode message: %v", err)
		return
	}

	log.Printf("[INFO] Handling message of type: %s", msgType)
}

func Broadcast(message []byte) {
	for _, peer := range PeersList {
		err := peer.Conn.WriteMessage(websocket.BinaryMessage, message)
		if err != nil {
			log.Printf("[ERROR] Failed to broadcast to %s: %v\n", peer.Addr, err)
		}
	}
}

func Encode(message interface{}) []byte {
	var buffer bytes.Buffer
	err := gob.NewEncoder(&buffer).Encode(message)
	if err != nil {
		log.Printf("[ERROR] Failed to encode message: %v", err)
		return nil
	}
	return buffer.Bytes()
}

func Decode(data []byte, v interface{}) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(v)
}

func startPeerServer(port string) {
	http.HandleFunc("/peer", HandlePeerConnection)
	log.Printf("[INFO] Starting server on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("[ERROR] Failed to start server on port %s: %v", port, err)
	}
}

func connectToPeer(peerAddr string) {
	if !strings.Contains(peerAddr, ":") {
		peerAddr += ":8080" // Default port if not specified
	}

	var retries = 3
	for retries > 0 {
		spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Customize spinner speed
		spinner.Start()
		log.Printf("[INFO] Attempting to connect to peer %s...", peerAddr)

		conn, _, err := websocket.DefaultDialer.Dial("ws://"+peerAddr+"/peer", nil)
		spinner.Stop()

		if err != nil {
			log.Printf("[ERROR] Failed to connect to peer %s: %v", peerAddr, err)
			retries--
			if retries == 0 {
				log.Printf("[ERROR] Giving up on connecting to peer %s after multiple attempts", peerAddr)
				return
			}
			time.Sleep(1 * time.Second) // Wait before retrying
			continue
		}

		// Add the peer to the PeersList
		mu.Lock()
		PeersList[peerAddr] = &Peer{
			Conn: conn,
			Addr: peerAddr,
		}
		mu.Unlock()

		log.Printf("[INFO] Connected to peer %s", peerAddr)
		
		go func() {
			ticker := time.NewTicker(10 * time.Minute) // Ping every 10 minutes
			defer ticker.Stop()
			defer func(conn *websocket.Conn) {
				err := conn.Close()
				if err != nil {
					log.Printf("[ERROR] Failed to close peer connection: %v", err)
				}
			}(conn)

			for {
				select {
				case <-ticker.C:
					err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second))
					if err != nil {
						log.Printf("[ERROR] Failed to send ping to %s: %v", peerAddr, err)
						return
					}
					log.Printf("[INFO] Sent ping to peer %s", peerAddr)
				}
			}
		}()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("[ERROR] Unexpected close error from %s: %v", peerAddr, err)
				} else {
					log.Printf("[ERROR] Read error from %s: %v (Type: %T)", peerAddr, err, err)
				}
				break
			}

			log.Printf("[MESSAGE] Received from %s: %s", peerAddr, string(message))
			go handleMessage(message, conn)
		}
		mu.Lock()
		delete(PeersList, peerAddr)
		mu.Unlock()
		log.Printf("[INFO] Peer %s disconnected", peerAddr)
	}
}

func addToPeer(peerAddr string) {
	mu.Lock()
	PeerIpList = append(PeerIpList, peerAddr)
	mu.Unlock()
	connectToPeer(peerAddr)
}

func splitPeers(peersStr string) []string {
	var peers []string
	for _, peer := range strings.Split(peersStr, ",") {
		peers = append(peers, peer)
	}
	return peers
}

func Start(port, peerAddr string) {

	if peerAddr == "" {
		go startPeerServer(port)
	} else {
		peerAddresses := splitPeers(peerAddr)
		for _, peer := range peerAddresses {
			addToPeer(peer)
		}
		go startPeerServer(port)
	}

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh

	mu.Lock()
	for _, peer := range PeersList {
		err := peer.Conn.Close()
		if err != nil {
			return
		}
	}
	mu.Unlock()

	log.Println("[INFO] Server stopped gracefully.")
}
