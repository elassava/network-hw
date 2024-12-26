package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// Bağlı kullanıcıları takip etmek için global değişkenler
var (
	clients      = make(map[string]net.Conn)
	clientsMutex sync.RWMutex
)

func main() {
	// Start the server on port 9000

	dstream, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer dstream.Close()

	fmt.Printf("Server is running at %s\n", dstream.Addr().String())

	var wg sync.WaitGroup

	// 10 saniyede bir bağlı kullanıcıları listele
	go func() {
		for {
			time.Sleep(10 * time.Second)
			listConnectedClients()
		}
	}()

	for {
		// Accept new client connections
		conn, err := dstream.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		wg.Add(1)
		go func(c net.Conn) {
			defer wg.Done()
			handleClient(c)
		}(conn)
	}
}

// Bağlı kullanıcıları listeleyen fonksiyon
func listConnectedClients() {
	clientsMutex.RLock()
	defer clientsMutex.RUnlock()

	fmt.Println("\nBağlı kullanıcılar:")
	for username := range clients {
		fmt.Printf("- %s\n", username)
	}
	fmt.Printf("Toplam kullanıcı sayısı: %d\n", len(clients))
}

// Handle each client connection
func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading username:", err)
		return
	}
	username = username[:len(username)-1]

	// Sadece bağlantı bilgisi göster
	fmt.Printf("User connected: %s (%s)\n", username, conn.RemoteAddr().String())

	clientsMutex.Lock()
	clients[username] = conn
	clientsMutex.Unlock()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()

		// DM kontrolü
		if message[0] == '@' {
			parts := strings.SplitN(message[1:], " ", 2)
			if len(parts) == 2 {
				targetUser := parts[0]
				content := parts[1]

				clientsMutex.RLock()
				targetConn, exists := clients[targetUser]
				clientsMutex.RUnlock()

				if exists {
					// DM'i sadece hedef kullanıcıya gönder
					dmMsg := fmt.Sprintf("[DM from %s]: %s\n", username, content)
					targetConn.Write([]byte(dmMsg))
					// Gönderene sessiz onay
					conn.Write([]byte(fmt.Sprintf("✓ Sent to %s\n", targetUser)))
				} else {
					conn.Write([]byte(fmt.Sprintf("❌ User %s not found\n", targetUser)))
				}
			}
		}
		// Normal mesajları işleme almıyoruz (sadece DM'ler çalışacak)
	}

	// Kullanıcı çıkış yaptığında
	fmt.Printf("User %s disconnected.\n", username)
	clientsMutex.Lock()
	delete(clients, username)
	clientsMutex.Unlock()
}
