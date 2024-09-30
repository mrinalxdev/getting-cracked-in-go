package main

import (
	"fmt"
    "net/http"
    "sync"
)

type LoadBalancer struct {
    servers []string 
    current int 
    mu sync.Mutex
}

func NewLoadBalancer(servers []string) *LoadBalancer {
    return &LoadBalancer{
        servers: servers,
        current: 0,
    }
}

// returns the next available servers 
func (lb *LoadBalancer) GetNextServer() string {
    lb.mu.Lock()
    defer lb.mu.Unlock()

    server := lb.servers[lb.current]
    lb.current = (lb.current + 1) % len(lb.servers)
    return server
}

func (lb *LoadBalancer) ServeHttp(w http.ResponseWriter, r *http.Request){
    server := lb.GetNextServer()
    fmt.Println("Redirecting to %s\n", server)

    http.Redirect(w, r, server, http.StatusFound)
}

func main() {
    servers := []string {
        "http://server1.com",
        "http://server2.com",
        "http://server3.com", 
    }

    lb := NewLoadBalancer(servers)

    http.Handle("/", lb)

    fmt.Println("Load Balancer listening on port 8080")
    http.ListenAndServe(":8080", nil)
}

