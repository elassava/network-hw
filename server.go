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

	// 20 sn kullanıcı listesini duyur
	go func() {
		for {
			time.Sleep(30 * time.Second)
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

// Bağlı kullanıcıları listeleyen ve duyuran fonksiyon
func listConnectedClients() {
	clientsMutex.RLock()
	// Kullanıcı listesini oluştur
	userList := "\nOnline Users:\n"
	for username := range clients {
		userList += fmt.Sprintf("🟢 %s\n", username)
	}

	// Her kullanıcıya listeyi gönder
	for _, conn := range clients {
		conn.Write([]byte(userList))
	}
	clientsMutex.RUnlock()
}

// Handle each client connection
func handleClient(conn net.Conn) {
	defer conn.Close()

	// Username validation loop
	var username string
	for {
		reader := bufio.NewReader(conn)
		tempUsername, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading username:", err)
			return
		}
		tempUsername = strings.TrimSpace(tempUsername)

		// Check if username exists
		clientsMutex.RLock()
		_, exists := clients[tempUsername]
		clientsMutex.RUnlock()

		if exists {
			conn.Write([]byte("USERNAME_TAKEN\n"))
		} else {
			username = tempUsername
			conn.Write([]byte("USERNAME_ACCEPTED\n"))
			break
		}
	}

	// Sadece bağlantı bilgisi göster
	fmt.Printf("User connected: %s (%s)\n", username, conn.RemoteAddr().String())

	clientsMutex.Lock()
	clients[username] = conn
	// Yeni kullanıcı bağlandığında diğer kullanıcılara duyur
	connectMsg := fmt.Sprintf("\n🟢 %s joined the chat\n", username)
	for user, client := range clients {
		if user != username {
			client.Write([]byte(connectMsg))
		}
	}
	clientsMutex.Unlock()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()

		// Mesaj kontrolü
		if strings.HasPrefix(message, "[") && strings.Contains(message, "]") {
			endIndex := strings.Index(message, "]")
			if endIndex > 1 {
				targetUser := message[1:endIndex]
				content := strings.TrimSpace(message[endIndex+1:])

				if targetUser == "all" {
					// Broadcast mesajı - tüm kullanıcılara gönder
					broadcastMsg := fmt.Sprintf("\n[Broadcast from %s]: %s\n", username, content)
					clientsMutex.RLock()
					for recipient, recipientConn := range clients {
						if recipient != username { // Kendisine gönderme
							recipientConn.Write([]byte(broadcastMsg))
						}
					}
					clientsMutex.RUnlock()
					// Gönderene onay
					conn.Write([]byte("\n✓ Broadcast sent\n"))
				} else {
					// Normal DM işlemi
					clientsMutex.RLock()
					targetConn, exists := clients[targetUser]
					clientsMutex.RUnlock()

					if exists {
						dmMsg := fmt.Sprintf("\n[DM from %s]: %s\n", username, content)
						targetConn.Write([]byte(dmMsg))
						conn.Write([]byte(fmt.Sprintf("\n✓ Sent to %s\n", targetUser)))
					} else {
						conn.Write([]byte(fmt.Sprintf("\n❌ User %s not found\n", targetUser)))
					}
				}
			}
		}
		// Normal mesajları işleme almıyoruz (sadece DM'ler çalışacak)
	}

	// Kullanıcı çıkış yaptığında
	fmt.Printf("User %s disconnected.\n", username)
	clientsMutex.Lock()
	// Kullanıcı çıkış yaptığında diğer kullanıcılara duyur
	disconnectMsg := fmt.Sprintf("\n🔴 %s left the chat\n", username)
	for user, client := range clients {
		if user != username {
			client.Write([]byte(disconnectMsg))
		}
	}
	delete(clients, username)
	clientsMutex.Unlock()
}
