package main

import (
	"log"
    "sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type User struct {
    Name string
    Conn *websocket.Conn
}

type Room struct {
    Name string
    Users map[*User]bool
    Messages chan []byte
    mu sync.Mutex
}

// type Message struct {
//     Data string
//     user User
// }

var upgrader = websocket.Upgrader{}

var room = Room{
    Name: "chat",
    Users: make(map[*User]bool),
    Messages: make(chan []byte),
}

func handleChat(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Fatalln(err.Error())
        return
    }

    user := &User{
        Name: "andrey",
        Conn: conn,
    }

    room.mu.Lock()
    room.Users[user] = true
    room.mu.Unlock()

    go func() {
        defer func(conn *websocket.Conn) {
			room.mu.Lock()
			delete(room.Users, user) 
			room.mu.Unlock()
			conn.Close()
			log.Printf("Пользователь %s отключился\n", user.Name)
		}(conn)

        for {
            t, message, err := conn.ReadMessage()
            if err != nil {
                log.Fatal(err.Error())
                break
            }
            log.Printf("Message: %s", string(message))
            log.Println(room.Users)
            for user := range room.Users {
                err = user.Conn.WriteMessage(t, message)
                if err != nil {
                    log.Fatal(err.Error())
                    break
                }
            }
        }
     }()
}

func home(c *gin.Context) {
    c.HTML(200, "index.html", gin.H{"title": "main page"})
}

func main() {
    router := gin.Default()

    router.LoadHTMLGlob("front/html/*")

    router.GET("/", home)
    router.GET("/chat", handleChat)

    router.Run()
}