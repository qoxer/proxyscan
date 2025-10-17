package main
import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)
type Proxy struct {
	IP   string
	Port string
}
func main() {
	file, err := os.Open("./files/proxies.txt")
	if err != nil {
		log.Fatal("Failed to open file:", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var proxies []Proxy
	for scanner.Scan() {
		proxyStr := scanner.Text()
		parts := strings.Split(proxyStr, ":")
		if len(parts) != 2 {
			log.Println("Invalid proxy format:", proxyStr)
			continue
		}
		proxy := Proxy{
			IP:   parts[0],
			Port: parts[1],
		}
		proxies = append(proxies, proxy)
	}
	var wg sync.WaitGroup
	concurrentRequests := 10
	semaphore := make(chan struct{}, concurrentRequests)
	validFile, err := os.Create("./files/valid.txt")
	if err != nil {
		log.Fatal("Failed to create valid file:", err)
	}
	defer validFile.Close()
	for _, proxy := range proxies {
		wg.Add(1)
		go func(p Proxy) {
			defer wg.Done()
			semaphore <- struct{}{}
			valid := checkProxy(p.IP, p.Port)
			if valid {
				fmt.Printf("+ %s:%s\n", p.IP, p.Port)
				validFile.WriteString(p.IP + ":" + p.Port + "\n")
			} else {
				fmt.Printf("- %s:%s\n", p.IP, p.Port)
			}
			<-semaphore
		}(proxy)
	}
	wg.Wait()
	fmt.Println("Proxy checking completed. Valid proxies have been saved to valid.txt")
}
func checkProxy(ip, port string) bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	proxyURL, err := url.Parse(fmt.Sprintf("http://%s:%s", ip, port))
	if err != nil {
		log.Println("Failed to parse proxy URL:", err)
		return false
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client.Transport = transport
	req, err := http.NewRequest("GET", "http://google.com", nil)
	if err != nil {
		log.Println("Failed to create request:", err)
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}
