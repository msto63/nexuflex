// main.go
/**
 * Nexuflex Client - Main Application
 *
 * This file contains the entry point for the nexuflex client application,
 * which provides a text-based user interface (TUI) for accessing nexuflex services.
 *
 * @author msto63
 * @version 1.0.0
 * @date 2025-03-12
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/msto63/nexuflex/nexuflex-client/config"
	"github.com/msto63/nexuflex/nexuflex-client/core"
	"github.com/msto63/nexuflex/nexuflex-client/i18n"
	"github.com/msto63/nexuflex/nexuflex-client/ui"
)

func main() {
	// Define command line parameters
	configFile := flag.String("config", "", "Path to config file")
	serverAddr := flag.String("server", "", "Server address (IP or hostname)")
	serverPort := flag.Int("port", 0, "Server port")
	discoverMode := flag.Bool("discover", false, "Enable automatic server discovery")
	discoverTimeout := flag.Int("discover-timeout", 5, "Timeout for server discovery in seconds")
	debug := flag.Bool("debug", false, "Enable debug output")
	language := flag.String("lang", "", "Language code (e.g., 'en', 'de')")
	flag.Parse()

	// Configure debug logging
	if *debug {
		logFile := filepath.Join(os.TempDir(), "nexuflex-client.log")
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		log.SetOutput(f)
		log.Println("Nexuflex client started")
	} else {
		// Disable logging
		log.SetOutput(os.NewFile(0, os.DevNull))
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Command line parameters override configuration file
	if *serverAddr != "" {
		cfg.Server.Address = *serverAddr
	}
	if *serverPort != 0 {
		cfg.Server.Port = *serverPort
	}
	if *discoverMode {
		cfg.Server.AutoDiscover = true
	}
	if *discoverTimeout != 5 {
		cfg.Server.DiscoverTimeoutSeconds = *discoverTimeout
	}
	if *language != "" {
		cfg.UI.Language = *language
	}

	// Initialize language files
	if err := i18n.LoadLanguage(cfg.UI.Language); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading language files: %v\n", err)
		fmt.Fprintf(os.Stderr, "Using English as fallback language\n")
		// Try loading default language (English)
		if err := i18n.LoadLanguage("en"); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading default language: %v\n", err)
			os.Exit(1)
		}
	}

	// Create client
	client := core.NewClient(&cfg, log.Printf)

	// Create TUI
	tui := ui.NewTUI(client)

	// Start TUI
	if err := tui.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing user interface: %v\n", err)
		os.Exit(1)
	}

	// Automatic server discovery, if configured
	if cfg.Server.AutoDiscover {
		err := client.DiscoverServer(time.Duration(cfg.Server.DiscoverTimeoutSeconds) * time.Second)
		if err != nil {
			tui.ShowError(fmt.Sprintf(i18n.GetMessage("error.discovery"), err))
		}
	} else if cfg.Server.Address != "" && cfg.Server.Port != 0 {
		// Connect to configured server
		err := client.Connect(cfg.Server.Address, cfg.Server.Port, cfg.Server.UseTLS)
		if err != nil {
			tui.ShowError(fmt.Sprintf(i18n.GetMessage("error.connection"), err))
		}
	}

	// Close client when application exits
	defer client.Close()
}
