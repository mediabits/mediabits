package wui

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

var public = flag.Bool("public", false, "open this mediabits server to the public - to use this you must use the -password option as well")
var port = flag.Uint("port", 0, "")
var password = flag.String("password", "", "require a password to access the web interface")
var passwordOverride = flag.Bool("unsafe-insecure-yes-i-want-to-get-hacked-because-im-an-idiot", false, "overrides the password requirement for public servers - DO NOT USE THIS UNLESS YOU WANT HACKERS TO STEAL YOUR DATA, PILLAGE YOUR HOUSE AND DESTROY ALL OF YOUR HOPES AND DREAMS")

func localAddresses() (*[]net.IP, error) {
	var localAddrs []net.IP

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		// Ignore down
		if i.Flags&net.FlagUp == 0 {
			continue
		}

		// Ignore loopback
		if i.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Get the addresses
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPNet:
				// TODO: IPv6 support
				ip := v.IP.To4()
				if ip == nil {
					continue
				}

				localAddrs = append(localAddrs, ip)
			}

		}
	}

	return &localAddrs, nil
}

func WuiMain() {
	// Seed the RNG - this is only used for getting a random port and does not need to be secure
	rand.Seed(time.Now().UTC().UnixNano())

	// Require password for public servers
	if *password == "" && *public && !*passwordOverride {
		fmt.Println("A password is required for public servers.")
		return
	}

	// Basic password check
	if *password != "" && len(*password) < 8 {
		fmt.Println("Passwords must be a minimum of eight characters long.")
		return
	}

	// Default port between 13000 and 14000
	if *port == 0 {
		tmp := uint(13000 + rand.Intn(1000))
		port = &tmp
	}

	// Listening on
	var authStr string

	if *password != "" {
		authStr = "mediabits:" + *password + "@"
	}

	fmt.Println("Web addresses:")
	if *public {
		localIPs, err := localAddresses()
		if err != nil {
			fmt.Printf("Failed to determine local address: %s", err.Error())
			return
		}
		for _, ip := range *localIPs {
			fmt.Printf("http://%s%s:%d/\n", authStr, ip, *port)
		}
	} else {
		fmt.Printf("http://%s%s:%d/\n", authStr, "127.0.0.1", *port)
	}

	// Construct laddr
	laddr := ":" + strconv.FormatUint(uint64(*port), 10)
	if !*public {
		laddr = "127.0.0.1" + laddr
	}

	// Listen TODO:ipv6
	listener, err := net.Listen("tcp4", laddr)
	if err != nil {
		fmt.Printf("Failed to listen: %s\n", err.Error())
		return
	}
	defer listener.Close()

	// Build the mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", authReq(handleIndex))
	mux.HandleFunc("/listfiles", authReq(handleListFiles))
	mux.HandleFunc("/movie", authReq(handleMovie))
	mux.HandleFunc("/tv", authReq(handleTV))

	http.Serve(listener, mux)
}
