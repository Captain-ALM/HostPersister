package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"golang.captainalm.com/HostPersister/hosts"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	buildVersion = "develop"
	buildDate    = ""
)

func main() {
	log.Printf("[Main] Starting up Host Persister #%s (%s)\n", buildVersion, buildDate)
	y := time.Now()

	//Hold main thread till safe shutdown exit:
	wg := &sync.WaitGroup{}
	wg.Add(1)

	//Load environment file:
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	//Load ENVs
	hostsPath := os.Getenv("HOSTS_FILE")
	sourcePath := os.Getenv("SOURCE_FILE")
	syncTime, err := strconv.ParseInt(os.Getenv("SYNC_TIME"), 10, 64)
	if err != nil {
		log.Println("[Main] Invalid SYNC_TIME; defaulting to one-shot execution.")
		syncTime = 0
	}
	overwriteMode := os.Getenv("HOSTS_OVERWRITE") == "1"

	//Load hosts and source files
	hostsFile, err := hosts.NewHostsFile(hostsPath)
	if err != nil {
		log.Fatalln("Failed to load HOSTS_FILE")
	}
	sourceFile, err := hosts.NewHostsFile(sourcePath)
	if err != nil {
		log.Fatalln("Failed to load SOURCE_FILE")
	}

	if syncTime < 1 {
		//One-shot execution:
		z := time.Now().Sub(y)
		log.Printf("[Main] Took '%s' to fully initialize modules\n", z.String())
		executePersistence(hostsFile, sourceFile, overwriteMode)
	} else {
		//Sync execution:
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		z := time.Now().Sub(y)
		log.Printf("[Main] Took '%s' to fully initialize modules\n", z.String())

		exec := true
		syncDur := time.Duration(syncTime) * time.Millisecond

		go func(exec *bool) {
			for *exec {
				executePersistence(hostsFile, sourceFile, overwriteMode)
				time.Sleep(syncDur)
			}
		}(&exec)

		go func(exec *bool) {
			<-sigs
			fmt.Printf("\n")

			*exec = false

			a := time.Now()
			log.Printf("[Main] Signalling program exit...\n")

			b := time.Now().Sub(a)
			log.Printf("[Main] Took '%s' to fully shutdown modules\n", b.String())
			wg.Done()
		}(&exec)

		wg.Wait()
	}
	log.Println("[Main] Goodbye")
}

func executePersistence(hostsFile *hosts.File, sourceFile *hosts.File, ovrw bool) {
	for _, entry := range sourceFile.Entries {
		for _, domain := range entry.Domains {
			if (!hostsFile.HasDomain(domain)) || ovrw {
				hostsFile.OverwriteDomainSingleton(domain, entry.IPAddress)
			}
		}
	}
	err := hostsFile.WriteHostsFile()
	if err != nil {
		log.Println("[Main] Error Writing Hosts File.")
	}
}
