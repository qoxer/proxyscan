package main
import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)
type Proxy struct {
	IP   string
	Port string
}
func main() {
	sitesFile, err := os.Open("./base/sites.txt")
	if err != nil {
		log.Fatal("Failed to open sites file:", err)
	}
	defer sitesFile.Close()
	sitesBytes, err := ioutil.ReadAll(sitesFile)
	if err != nil {
		log.Fatal("Failed to read sites file:", err)
	}
	sites := strings.Split(string(sitesBytes), "\n")
	var proxies []Proxy
	for _, site := range sites {
		if site == "" {
			continue
		}
		fmt.Println("Parsing proxies from:", site)
		resp, err := http.Get(site)
		if err != nil {
			log.Println("Failed to load website:", err)
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Failed to read HTML:", err)
			continue
		}
		re := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):(\d+)`)
		matches := re.FindAllStringSubmatch(string(body), -1)
		for _, match := range matches {
			ip := match[1]
			port := match[2]
			proxy := Proxy{IP: ip, Port: port}
			proxies = append(proxies, proxy)
		}
	}
	file, err := os.Create("./files/proxies.txt")
	if err != nil {
		log.Fatal("Failed to create file:", err)
	}
	defer file.Close()
	for _, proxy := range proxies {
		proxyStr := fmt.Sprintf("%s:%s\n", proxy.IP, proxy.Port)
		_, err := file.WriteString(proxyStr)
		if err != nil {
			log.Println("Failed to write to file:", err)
		}
	}
	fmt.Println("Proxies have been successfully parsed and saved to proxies.txt")
}