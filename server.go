package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/sys/unix"
)

var (
	addr     = flag.String("addr", "localhost:6969", "http service address")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func main() {
	flag.Parse()
	r := gin.Default()
	r.GET("/ws", ws)
	r.LoadHTMLGlob("templates/*")
	r.Static("/dist", "./dist")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.Run(*addr)
}

func ws(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		log.Println(err)
		return
	}
	include := []string{"dist/", "./templates"}
	fd, err := unix.InotifyInit()
	if err != nil {
		log.Println("Cant watch files error: ", err)
	}
	for _, dir := range include {
		unix.InotifyAddWatch(fd, dir, unix.IN_MODIFY)
	}

	pollfds := make([]unix.PollFd, 1, 1)
	pollfds[0].Fd = int32(fd)
	pollfds[0].Events = unix.POLLIN

	lastPoll := time.Now()
	for {
		num, err := unix.Poll(pollfds, -1)
		if err != nil {
			log.Println("poll error: ", err)
		}
		if num > 0 {
			buf := make([]byte, 4096)
			unix.Read(fd, buf)
		}
		newPoll := time.Now()
		if newPoll.Sub(lastPoll) < 50*time.Millisecond {
			continue
		}
		ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := ws.WriteMessage(websocket.TextMessage, []byte{}); err != nil {
			log.Println("ws send update error:", err)
			return
		}
		lastPoll = newPoll
	}
}
