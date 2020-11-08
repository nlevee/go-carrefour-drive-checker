package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/nlevee/go-carrefour-drive-checker/internal/api"
	"github.com/nlevee/go-carrefour-drive-checker/pkg/carrefour"
)

func main() {
	driveID := flag.String("id", "", "The drive Id")
	postalCode := flag.String("cp", "", "The Postal Code")
	listenHost := flag.String("host", "0.0.0.0", "Start a server and listen on this host")
	listenPort := flag.String("port", "", "Start a server and listen on this port")
	proxies := flag.String("proxies", "", "use this proxy list randomly to connect to carrefour API (http://ip:port,http://ip:port,...)")
	flag.Parse()

	// config proxies
	if *proxies != "" {
		carrefour.Proxies = strings.Split(*proxies, ",")
	}

	// recherche du driveId si code postal
	if *driveID == "" && *postalCode != "" {
		storeIDs, _ := carrefour.GetStoreIDByPostalCode(*postalCode)
		if len(storeIDs) > 0 {
			driveID = &storeIDs[0]
		} else {
			log.Fatal("no stores found")
		}
	}

	if *listenPort != "" && *listenHost != "" {
		if *driveID != "" {
			go carrefour.NewDriveHandler(*driveID)
		}
		api.StartServer(*listenHost, *listenPort)
	} else {
		if *driveID == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
		carrefour.NewDriveHandler(*driveID)
	}
}
