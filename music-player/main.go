package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Song struct {
	Title  string
	Artist string
	Path   string
}

type Player struct {
	Songs      []Song
	Current    int
	Mode       string // "sequence", "random", "repeat"
	mu         sync.Mutex
	isPlaying  bool
	shuffleIdx []int
}

func NewPlayer() *Player {
	return &Player{
		Mode:      "sequence",
		isPlaying: false,
	}
}

func (p *Player) SetSongs(songs []Song) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Songs = songs
	p.Current = 0
	p.Mode = "sequence"
	p.shuffleIdx = nil
}

func (p *Player) Next() Song {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.Songs) == 0 {
		return Song{}
	}

	switch p.Mode {
	case "random":
		if p.shuffleIdx == nil {
			p.initShuffle()
		}
		p.Current = p.shuffleIdx[0]
		p.shuffleIdx = p.shuffleIdx[1:]
		if len(p.shuffleIdx) == 0 {
			p.initShuffle()
		}
	case "repeat":
		// Do nothing, just play current song again
	case "sequence":
		p.Current = (p.Current + 1) % len(p.Songs)
	}

	return p.Songs[p.Current]
}

func (p *Player) initShuffle() {
	p.shuffleIdx = make([]int, len(p.Songs))
	for i := range p.shuffleIdx {
		p.shuffleIdx[i] = i
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(p.shuffleIdx), func(i, j int) {
		p.shuffleIdx[i], p.shuffleIdx[j] = p.shuffleIdx[j], p.shuffleIdx[i]
	})
}

func (p *Player) TogglePlay() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isPlaying = !p.isPlaying
	return p.isPlaying
}

func (p *Player) SetMode(mode string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Mode = mode
}

func main() {
	player := NewPlayer()

	// 添加一些示例歌曲
	player.SetSongs([]Song{
		{"Song 1", "Artist A", "/static/music/song1.mp3"},
		{"Song 2", "Artist B", "/static/music/song2.mp3"},
		{"Song 3", "Artist C", "/static/music/song3.mp3"},
		{"Song 4", "Artist D", "/static/music/song4.mp3"},
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Songs  []Song
			Player *Player
		}{
			Songs:  player.Songs,
			Player: player,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// 静态文件服务
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// API端点
	http.HandleFunc("/api/next", func(w http.ResponseWriter, r *http.Request) {
		nextSong := player.Next()
		if nextSong.Path == "" {
			http.Error(w, "No songs available", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"title":"` + nextSong.Title + `","artist":"` + nextSong.Artist + `","path":"` + nextSong.Path + `"}`))
	})

	http.HandleFunc("/api/toggle", func(w http.ResponseWriter, r *http.Request) {
		isPlaying := player.TogglePlay()
		w.Header().Set("Content-Type", "application/json")
		if isPlaying {
			w.Write([]byte(`{"status":"playing"}`))
		} else {
			w.Write([]byte(`{"status":"paused"}`))
		}
	})

	http.HandleFunc("/api/mode", func(w http.ResponseWriter, r *http.Request) {
		mode := r.URL.Query().Get("mode")
		player.SetMode(mode)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"mode":"` + mode + `"}`))
	})

	// 注意：在实际应用中，你需要提供真实的音乐文件
	// 这里只是模拟，实际应该设置一个真实的音乐文件服务

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
