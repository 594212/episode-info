package main

import (
	"flag"
	"html/template"
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
	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int { return a + b },
	})
	r.LoadHTMLGlob("templates/*")

	r.Static("/dist", "./dist")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", data)
	})
	r.Run(*addr)

}

type Evaluation struct {
	Like        uint32
	Dislike     uint32
	Rating      float64
	WorldRating float64
}

type ATag struct {
	Content any
	Href    string
}
type Author struct {
	Name       string
	SecondName string
	Url        string
}
type Number struct {
	Total  int
	Number int
}

// type Release struct {
// 	Weekday     time.Weekday
// 	Hour        int
// 	Releasetime time.Time
// }

var now = time.Now()
var future = now.AddDate(0, 0, 7)
var targetTime = time.Date(future.Year(), future.Month(), future.Day(), 20, 0, 0, 0, future.Location())

var data = struct {
	Title       string
	Name        string
	Image       string
	Number      Number
	Author      ATag
	Evaluation  Evaluation
	ReleaseTime string
	Tags        []ATag
	Year        ATag
	Status      ATag
	Type        ATag
	Studio      ATag
	Voiceover   []ATag
}{
	ReleaseTime: targetTime.Format(time.RFC3339),
	Number: Number{
		Total:  128,
		Number: 78,
	},
	Title: "«Противостояние святого»",
	Name:  "«Противостояние святого»",
	Image: "https://amedia.lol/uploads/posts/2024-12/original.webp",
	Author: ATag{
		Content: Author{
			"Xian Ni",
			"Renegade Immortal",
			"renegade_immortal",
		},
		Href: "/renegade_immortal",
	},
	Evaluation: Evaluation{
		Like:        32747,
		Dislike:     558,
		Rating:      9.83,
		WorldRating: 7.23,
	},
	Tags: []ATag{
		{"экшен", "action"},
		{"приключение", "adventure"},
		{"фэнтези", "fantasy"},
		{"исторический", "history"},
		{"китайское", "china"},
		{"3D", "3D"},
	},
	Year: ATag{"2023", "/year/2023"},
	Status: ATag{
		"Онгоинги",
		"ongoing",
	},
	Type: ATag{
		"ONA", "ona",
	},
	Studio: ATag{
		"BUILD DREAM",
		"build_dream",
	},
	Voiceover: []ATag{
		{"AniStar", "anistar"},
		{"AnimeVost", "anivost"},
		{"Animy", "animy"},
		{"FSG N.N. Азия.Subtitles", "fsg_n.n._asia.subtitles"},
	},
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
		log.Println("Can't watch files error: ", err)
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
		}
		return
		lastPoll = newPoll
	}
}
