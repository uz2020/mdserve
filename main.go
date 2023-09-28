package main

import (
	"context"
	"embed"
	"flag"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

//go:embed dist
var content embed.FS

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ReadFile() {
}

func stream(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	wg.Add(1)
	go func() {
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read message:", err)
				break
			}
			if mt == websocket.CloseMessage {
				log.Println("client close")
				break
			}

			name := string(message)
			err = watcher.Add(name)
			if err != nil {
				log.Println("watch err", err)
				break
			}
			log.Println("watch for", mt, name)

			b, err := ioutil.ReadFile(name)
			if err != nil {
				log.Println("read:", err)
				break
			}
			err = c.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
		cancel()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		ticker := time.NewTicker(time.Minute * 1)
		for {
			select {
			case <-ticker.C:
				log.Println("tick")
			case event := <-watcher.Events:
				log.Println("modify", event.Name, event.String())

				b, err := ioutil.ReadFile(event.Name)
				if err != nil {
					log.Println("read:", err)
					return
				}
				err = c.WriteMessage(websocket.TextMessage, b)
				if err != nil {
					log.Println("write:", err)
					return
				}
			case <-ctx.Done():
				log.Println("cancel")
				wg.Done()
				return
			}
		}

	}()
	wg.Wait()
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	mux := http.NewServeMux()

	fSys, err := fs.Sub(content, "dist")
	if err != nil {
		panic(err)
	}
	fs := http.FileServer(http.FS(fSys))
	mux.Handle("/", fs)

	mux.HandleFunc("/stream", stream)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
