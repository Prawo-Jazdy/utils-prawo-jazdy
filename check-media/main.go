package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Question struct {
	Media string `json:"media"`
}

type Data []Question

// Sprawdza, czy plik istnieje
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Przenosi plik do katalogu "exists/"
func moveFile(sourcePath, destDir string) error {
	// Tworzenie katalogu docelowego, jeśli nie istnieje
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %v", destDir, err)
	}

	// Nowa ścieżka docelowa
	destPath := filepath.Join(destDir, filepath.Base(sourcePath))

	// Przeniesienie pliku
	if err := os.Rename(sourcePath, destPath); err != nil {
		return fmt.Errorf("error moving file %s to %s: %v", sourcePath, destPath, err)
	}

	return nil
}

// Szuka pliku w podfolderach mediaDir
func findFileInSubdirs(mediaDir, filename string) (string, bool) {
	var foundPath string
	err := filepath.Walk(mediaDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing %s: %v\n", path, err)
			return err
		}
		if !info.IsDir() && info.Name() == filename {
			foundPath = path
			return filepath.SkipDir // Zatrzymujemy przeszukiwanie
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory %s: %v\n", mediaDir, err)
		return "", false
	}
	return foundPath, foundPath != ""
}

// Szuka pliku w katalogu "exists/"
func findFileInExists(existsDir, filename string) bool {
	existsPath := filepath.Join(existsDir, filename)
	return fileExists(existsPath)
}

// Przetwarza pojedynczy plik JSON
func processFile(jsonPath, mediaDir, existsDir string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", jsonPath, err)
	}

	var questions Data
	if err := json.Unmarshal(data, &questions); err != nil {
		return fmt.Errorf("error parsing JSON in file %s: %v", jsonPath, err)
	}

	// Iteracja po obiektach w JSON
	for _, q := range questions {
		if q.Media == "" {
			continue
		}

		// Najpierw sprawdzamy, czy plik już został przeniesiony
		if findFileInExists(existsDir, q.Media) {
			continue // Jeśli plik już istnieje w "exists/", pomijamy
		}

		// Szukamy pliku w podfolderach
		if mediaPath, found := findFileInSubdirs(mediaDir, q.Media); found {
			moveFile(mediaPath, existsDir) // Przenosi plik
		} else {
			fmt.Printf("Missing: %s\n", q.Media) // Wyświetla brakujące pliki, jeśli nie są w "exists/"
		}
	}
	return nil
}

// Przechodzi przez wszystkie pliki JSON w katalogu
func processDirectory(jsonDir, mediaDir, existsDir string) error {
	return filepath.Walk(jsonDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			fmt.Println("Processing JSON:", path)
			if err := processFile(path, mediaDir, existsDir); err != nil {
				fmt.Println("Error:", err)
			}
		}
		return nil
	})
}

func main() {
	jsonFolder := "./kategorie-json" // Katalog z plikami JSON
	mediaFolder := "./media-section" // Katalog z folderami zawierającymi multimedia
	existsFolder := "./exists"       // Katalog dla znalezionych plików

	if err := processDirectory(jsonFolder, mediaFolder, existsFolder); err != nil {
		fmt.Println("Error processing directory:", err)
	}
}
