/*
Copyright © 2020 si9ma <si9ma@si9ma.com>
*/
package cmd

import (
	"fmt"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/si9ma/simpleLB/config"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

var (
	port                    int
	autoTLS                 bool
	defaultAutoCertCacheDir = "./cache"
	autoCertCacheDir        string
)

const (
	defaultPort = 8080
)

var lbCmd = &cobra.Command{
	Use:   "lb",
	Short: "start a load balancer",
	Long:  "start a load balancer",
	Run: func(cmd *cobra.Command, args []string) {
		checkConfig()
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(lbCmd)

	lbCmd.Flags().IntVarP(&port, "port", "p", defaultPort, "listen port")
	lbCmd.Flags().BoolVar(&autoTLS, "autotls", false, "auto tls with Let's Encrypt")
	lbCmd.Flags().StringVar(&autoCertCacheDir, "autoCertCacheDir", defaultAutoCertCacheDir, "auto cert cache dir")
}

func checkConfig() {
	if len(lbConfig.LB) <= 0 {
		fmt.Fprintf(os.Stderr, "can't find any domain in the config")
		os.Exit(1)
	}
}

func startServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// forward all request to loadBalancer
	r.NoRoute(loadBalancer)

	var err error
	if autoTLS {
		domains := getDomains(lbConfig.LB)
		// service https with Let's Encrypt
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domains...),
			Cache:      autocert.DirCache(autoCertCacheDir),
		}
		err = autotls.RunWithManager(r, &m)
	} else {
		err = r.Run("0.0.0.0:" + strconv.Itoa(port))
	}

	if err != nil {
		log.Fatalf("run server fail:%s", err)
	}
}

func loadBalancer(c *gin.Context) {
	// if there's a port inside the host, parse pure host
	reqHost, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		reqHost = c.Request.Host
	}
	machineList, ok := lbConfig.LB[reqHost]
	if !ok || len(machineList.GetAvailableMachine()) <= 0 {
		log.Printf("can't find any downstream for %s\n", reqHost)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	// get a random machine
	availableMachineList := machineList.GetAvailableMachine()
	randN := rand.Intn(len(availableMachineList))
	machine := availableMachineList[randN]
	hostUrl, err := url.Parse("http://" + machine.Host)
	if err != nil {
		log.Printf("parse host:%s failed:%s\n", machine.Host, err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// forward request to the downstream
	proxy := httputil.NewSingleHostReverseProxy(hostUrl)
	proxy.ServeHTTP(c.Writer, c.Request)
}

// get domains from lbConfig
func getDomains(strMap map[string]config.MachineList) []string {
	keys := make([]string, len(strMap))

	i := 0
	for k := range strMap {
		keys[i] = k
		i++
	}

	return keys
}
