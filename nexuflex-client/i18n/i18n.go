// i18n.go
/**
 * Nexuflex Client - Internationalization
 *
 * This file contains the implementation for language support and message loading.
 *
 * @author msto63
 * @version 1.0.0
 * @date 2025-03-12
 */

package i18n

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

// Current loaded language and messages
var (
	currentLanguage string
	messages        map[string]string
)

// LoadLanguage loads a language file based on the specified language code
func LoadLanguage(langCode string) error {
	// If no language code is provided, try to detect from environment
	if langCode == "" {
		langCode = detectLanguage()
	}

	// Initialize messages map
	messages = make(map[string]string)

	// Find language file paths
	langPaths := findLangFilePaths(langCode)
	if len(langPaths) == 0 {
		return fmt.Errorf("no language file found for code '%s'", langCode)
	}

	// Load each language file found
	for _, path := range langPaths {
		if err := loadLangFile(path); err != nil {
			return err
		}
	}

	// Set current language
	currentLanguage = langCode
	return nil
}

// GetMessage returns a localized message for the given key
func GetMessage(key string) string {
	if msg, ok := messages[key]; ok {
		return msg
	}
	// If key doesn't exist, return the key itself as fallback
	return key
}

// GetCurrentLanguage returns the currently loaded language code
func GetCurrentLanguage() string {
	return currentLanguage
}

// GetAvailableLanguages returns a list of available language codes
func GetAvailableLanguages() ([]string, error) {
	langCodes := make([]string, 0)

	// Check standard paths for language files
	paths := getStandardLangDirs()
	for _, dir := range paths {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue // Skip this directory if it can't be read
		}

		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".ini") {
				// Extract language code from filename (e.g., "en.ini" -> "en")
				langCode := strings.TrimSuffix(file.Name(), ".ini")
				if isValidLangCode(langCode) {
					langCodes = append(langCodes, langCode)
				}
			}
		}
	}

	if len(langCodes) == 0 {
		return nil, errors.New("no language files found")
	}

	return langCodes, nil
}

// Helper functions

// detectLanguage tries to detect the system language
func detectLanguage() string {
	// Try LANG environment variable first (UNIX-like systems)
	langEnv := os.Getenv("LANG")
	if langEnv != "" {
		// Extract language code (e.g., "en_US.UTF-8" -> "en")
		parts := strings.Split(langEnv, "_")
		if len(parts) > 0 && isValidLangCode(parts[0]) {
			return parts[0]
		}
	}

	// Try alternative environment variables
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANGUAGE"} {
		langEnv = os.Getenv(env)
		if langEnv != "" {
			parts := strings.Split(langEnv, "_")
			if len(parts) > 0 && isValidLangCode(parts[0]) {
				return parts[0]
			}
		}
	}

	// Fallback to English
	return "en"
}

// isValidLangCode checks if a language code is valid
func isValidLangCode(code string) bool {
	// Simple validation: 2-3 characters, all lowercase
	return len(code) >= 2 && len(code) <= 3 && code == strings.ToLower(code)
}

// getStandardLangDirs returns standard directories to look for language files
func getStandardLangDirs() []string {
	dirs := []string{
		"lang",    // Local directory
		"i18n",    // Local directory alternative
		"locales", // Local directory alternative
	}

	// Add user config directory
	if configDir, err := os.UserConfigDir(); err == nil {
		dirs = append(dirs, filepath.Join(configDir, "nexuflex", "lang"))
	}

	// Add executable directory
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		dirs = append(dirs, filepath.Join(exeDir, "lang"))
	}

	return dirs
}

// findLangFilePaths finds all language files for a given language code
func findLangFilePaths(langCode string) []string {
	paths := []string{}

	// Check standard directories
	for _, dir := range getStandardLangDirs() {
		langFile := filepath.Join(dir, langCode+".ini")
		if _, err := os.Stat(langFile); err == nil {
			paths = append(paths, langFile)
		}
	}

	return paths
}

// loadLangFile loads messages from a language file
func loadLangFile(path string) error {
	// Load INI file
	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}

	// Process all sections
	for _, section := range cfg.Sections() {
		sectionName := section.Name()

		// Skip default section with empty name
		if sectionName == "DEFAULT" {
			// Load keys from DEFAULT section directly into messages map
			for _, key := range section.Keys() {
				messages[key.Name()] = key.Value()
			}
			continue
		}

		// For other sections, prefix the keys with section name
		for _, key := range section.Keys() {
			messageKey := sectionName + "." + key.Name()
			messages[messageKey] = key.Value()
		}
	}

	return nil
}
