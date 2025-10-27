package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Println("📚 pInk renamer")

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var formato int
	fmt.Println("Escolha o formato de numeração (2 = dezena, 3 = centena): ")
	fmt.Scan(&formato)

	if formato != 2 && formato != 3 {
		fmt.Println("❌ Formato inválido. Use 2 ou 3.")
		return
	}

	fmt.Printf("→ Usando formatação com %d dígitos.\n\n", formato)

	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	reNumero := regexp.MustCompile(`(\d+)(?:\D*$)?`)
	start := time.Now()

	cpuCount := runtime.NumCPU()
	sem := make(chan struct{}, cpuCount)
	var wg sync.WaitGroup
	var renamed int64
	var mu sync.Mutex

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		oldName := f.Name()
		ext := filepath.Ext(oldName)
		if ext != ".cbr" && ext != ".cbz" {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(oldName string) {
			defer wg.Done()
			defer func() { <-sem }()

			base := strings.TrimSuffix(oldName, filepath.Ext(oldName))
			numeroMatch := reNumero.FindStringSubmatch(base)
			if len(numeroMatch) == 0 {
				return
			}

			numero, _ := strconv.Atoi(numeroMatch[1])
			numeroFmt := fmt.Sprintf("%0*d", formato, numero)

			base = strings.ReplaceAll(base, " ", "-")
			base = reNumero.ReplaceAllString(base, numeroFmt)
			newName := strings.ToLower(base) + filepath.Ext(oldName)

			hasUpper := oldName != strings.ToLower(oldName)
			if !hasUpper && oldName == newName {
				fmt.Printf("🔷 Mantido: %s\n", oldName)
				return
			}

			if err := os.Rename(oldName, newName); err != nil {
				fmt.Printf("❌ Erro ao renomear %s: %v\n", oldName, err)
				return
			}

			mu.Lock()
			renamed++
			mu.Unlock()

			fmt.Printf("✅ %s → %s\n", oldName, newName)
		}(oldName)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("\n📂 Foram renomeados %d arquivos em %.2f segundos usando %d núcleos.\n", renamed, elapsed.Seconds(), cpuCount)
}
