package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func hashFilename(filename string) string {
	filename = strings.TrimSpace(filename)                                      // Usunięcie spacji na początku i końcu
	base := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename)) // Pobranie nazwy pliku bez rozszerzenia
	ext := strings.ToLower(filepath.Ext(filename))                              // Pobranie rozszerzenia i zamiana na małe litery

	// Obsługiwane rozszerzenia
	validExtensions := map[string]string{
		".mp4":  ".mp4",
		".jpg":  ".jpg",
		".jpeg": ".jpg", // Konwersja .jpeg → .jpg
		".JPG":  ".jpg",
	}

	// Sprawdzenie obsługiwanego rozszerzenia
	newExt, ok := validExtensions[ext]
	if !ok {
		fmt.Printf("Ignoring file %s (unsupported extension: %s)\n", filename, ext)
		return filename // Jeśli nieobsługiwane, zwracamy oryginalną nazwę
	}

	// Hashowanie nazwy pliku (bez rozszerzenia)
	hasher := sha1.New()
	hasher.Write([]byte(base))
	hashed := hex.EncodeToString(hasher.Sum(nil))[:10] // Skrócony hash (10 znaków)

	return hashed + newExt
}

func renameVideosInFolder(rootFolder string) error {
	return filepath.Walk(rootFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing file %s: %v", path, err)
		}

		// Pomijamy katalogi
		if info.IsDir() {
			return nil
		}

		// Pobierz rozszerzenie pliku i zamień na małe litery
		ext := strings.ToLower(filepath.Ext(info.Name()))

		// Obsługiwane rozszerzenia
		if ext == ".jpg" || ext == ".mp4" {
			dir := filepath.Dir(path) // Pobierz katalog, w którym jest plik
			newName := hashFilename(info.Name())
			newPath := filepath.Join(dir, newName)

			// Sprawdzamy, czy plik o nowej nazwie już istnieje
			if _, err := os.Stat(newPath); err == nil {
				fmt.Printf("File %s already exists in %s, skipping...\n", newName, dir)
				return nil
			}

			err := os.Rename(path, newPath)
			if err != nil {
				fmt.Printf("Error renaming %s: %v\n", path, err)
			} else {
				fmt.Printf("Renamed %s -> %s\n", path, newPath)
			}
		}
		return nil
	})
}

func main() {
	rootFolder := "./media-section" // Ścieżka do katalogu głównego
	if err := renameVideosInFolder(rootFolder); err != nil {
		fmt.Println("Error:", err)
	}
}
