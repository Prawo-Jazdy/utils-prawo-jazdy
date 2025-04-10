package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Question struct {
	ID                    int               `json:"id"`
	QuestionNumber        int               `json:"questionNumber"`
	Question              map[string]string `json:"question"`
	Answers               []Answer          `json:"answers"`
	CorrectAnswer         string            `json:"correctAnswer"`
	Media                 string            `json:"media"`
	Categories            []string          `json:"categories"`
	Type                  string            `json:"type"`
	Points                int               `json:"points"`
	LegalBasis            []LegalBasis      `json:"legalBasis"`
	AdditionalDescription string            `json:"additionalDescription"`
}

type Answer struct {
	ID   string            `json:"id"`
	Text map[string]string `json:"text"`
}

type LegalBasis struct {
	Name     string   `json:"name"`
	Articles []string `json:"articles"`
}

type Data []Question

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

func processFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	var questions Data
	if err := json.Unmarshal(data, &questions); err != nil {
		return fmt.Errorf("error parsing JSON in file %s: %v", filePath, err)
	}

	// Aktualizacja tylko pola "media"
	for i, q := range questions {
		if q.Media != "" {
			questions[i].Media = hashFilename(q.Media)
		}
	}

	updatedData, err := json.MarshalIndent(questions, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON in file %s: %v", filePath, err)
	}

	// Nadpisanie oryginalnego pliku z nowymi danymi
	if err := os.WriteFile(filePath, updatedData, 0644); err != nil {
		return fmt.Errorf("error writing file %s: %v", filePath, err)
	}

	fmt.Println("Updated file:", filePath)
	return nil
}

func processDirectory(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			fmt.Println("Processing:", path)
			if err := processFile(path); err != nil {
				fmt.Println("Error:", err)
			}
		}
		return nil
	})
}

func main() {
	folderPath := "./kategorie-json" // Zmień na ścieżkę do katalogu z plikami JSON
	if err := processDirectory(folderPath); err != nil {
		fmt.Println("Error processing directory:", err)
	}
}

// Pierwszy: KW_D16_354org.mp4
