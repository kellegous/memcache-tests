package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

// Proxy ...
type Proxy struct {
	feAddr, beAddr string
	lck            sync.RWMutex
	l              net.Listener
	cns            []net.Conn
}

// Start ...
func (p *Proxy) Start() error {
	p.lck.Lock()
	defer p.lck.Unlock()
	if p.l != nil {
		return nil
	}

	l, err := net.Listen("tcp", p.feAddr)
	if err != nil {
		return err
	}

	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				log.Println(err)
				return
			}
			go serveProxy(c, p)
		}
	}()

	p.l = l
	p.cns = nil
	log.Printf("proxy started %s -> %s", p.feAddr, p.beAddr)
	return nil
}

// IsStarted ...
func (p *Proxy) IsStarted() bool {
	p.lck.RLock()
	defer p.lck.RUnlock()
	return p.l != nil
}

// Drain ...
func (p *Proxy) Drain() {
	p.lck.Lock()
	defer p.lck.Unlock()

	p.l.Close()

	for _, cn := range p.cns {
		cn.Close()
	}

	p.l = nil
	p.cns = nil
	log.Println("proxy drained")
}

func (p *Proxy) register(c net.Conn) {
	p.lck.Lock()
	defer p.lck.Unlock()
	p.cns = append(p.cns, c)
}

func (p *Proxy) unregister(c net.Conn) {
	p.lck.Lock()
	defer p.lck.Unlock()
	var cns []net.Conn
	for _, cn := range p.cns {
		if cn != c {
			cns = append(cns, cn)
		}
	}
	p.cns = cns
}

func startMemcache(port int) (*os.Process, error) {
	c := exec.Command("memcached",
		fmt.Sprintf("--port=%d", port))
	if err := c.Start(); err != nil {
		return nil, err
	}
	return c.Process, nil
}

func sendJSON(w http.ResponseWriter,
	status int,
	data interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Panic(err)
	}
}

func serveProxy(c net.Conn, p *Proxy) {
	p.register(c)
	defer p.unregister(c)

	log.Printf("fr=%s to=%s", c.RemoteAddr(), p.beAddr)
	pc, err := net.Dial("tcp", p.beAddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer pc.Close()

	go func() {
		defer pc.Close()
		defer c.Close()

		if _, err := io.Copy(c, pc); err != nil {
			return
		}
	}()

	if _, err := io.Copy(pc, c); err != nil {
		return
	}
}

func serveHTTP(p *Proxy, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if s := r.FormValue("state"); s == "" {
			http.Error(w,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
		} else if s == "on" || s == "true" {
			if err := p.Start(); err != nil {
				log.Panic(err)
			}
			sendJSON(w, http.StatusOK, struct {
				State bool `json:"state"`
			}{
				true,
			})
		} else if s == "off" || s == "false" {
			p.Drain()
			sendJSON(w, http.StatusOK, struct {
				State bool `json:"state"`
			}{
				false,
			})
		} else {
			http.Error(w, http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
		}
	case "GET":
		sendJSON(w, http.StatusOK, struct {
			State bool `json:"state"`
		}{
			p.IsStarted(),
		})
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed),
			http.StatusMethodNotAllowed)
	}
}

func main() {
	flagFeAddr := flag.String(
		"fe",
		":11211",
		"frontend address")
	flagBePort := flag.Int(
		"be",
		11212,
		"backend port")
	flagHTTPAddr := flag.String(
		"http",
		":8080",
		"http address")

	flag.Parse()

	_, err := startMemcache(*flagBePort)
	if err != nil {
		log.Panic(err)
	}

	prx := Proxy{
		feAddr: *flagFeAddr,
		beAddr: fmt.Sprintf("localhost:%d", *flagBePort),
	}

	if err := prx.Start(); err != nil {
		log.Panic(err)
	}

	log.Panic(http.ListenAndServe(*flagHTTPAddr,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			serveHTTP(&prx, w, r)
		})))
}
