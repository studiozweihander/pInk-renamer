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

type Renomeacao struct {
	Antigo string
	Novo   string
}

func main() {
	for {
		fmt.Println("📚 pInk renamer")

		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		var formato int
		fmt.Print("\nEscolha o formato de numeração (2 = dezena, 3 = centena): ")
		fmt.Scan(&formato)

		if formato != 2 && formato != 3 {
			fmt.Println("❌ Formato inválido. Use 2 ou 3.")
			continue
		}

		files, err := os.ReadDir(dir)
		if err != nil {
			panic(err)
		}

		reNumero := regexp.MustCompile(`(\d+)(?:\D*$)?`)
		reMultiHifen := regexp.MustCompile(`-+`)
		var renomeacoes []Renomeacao

		for _, f := range files {
			if f.IsDir() {
				continue
			}

			oldName := f.Name()
			ext := filepath.Ext(oldName)
			if ext != ".cbr" && ext != ".cbz" {
				continue
			}

			base := strings.TrimSuffix(oldName, ext)
			numeroMatch := reNumero.FindStringSubmatch(base)
			if len(numeroMatch) == 0 {
				continue
			}

			numero, _ := strconv.Atoi(numeroMatch[1])
			numeroFmt := fmt.Sprintf("%0*d", formato, numero)

			base = strings.ReplaceAll(base, " ", "-")
			base = strings.ReplaceAll(base, ".", "-")
			base = strings.ReplaceAll(base, "_", "-")

			base = reNumero.ReplaceAllString(base, numeroFmt)

			base = reMultiHifen.ReplaceAllString(base, "-")

			if !strings.Contains(base, "-"+numeroFmt) {
				base = strings.TrimRight(base, "-") + "-" + numeroFmt
			}

			newName := strings.ToLower(base) + ext
			hasUpper := oldName != strings.ToLower(oldName)

			if !hasUpper && oldName == newName {
				fmt.Printf("🔷 Mantido: %s\n", oldName)
				return
			}

			if oldName != newName {
				renomeacoes = append(renomeacoes, Renomeacao{Antigo: oldName, Novo: newName})
				fmt.Printf("✅ %s → %s\n", oldName, newName)
			}
		}

		if len(renomeacoes) == 0 {
			fmt.Println("\n❌ Nenhum arquivo elegível para renomeação encontrado.")
			return
		}

		var confirma string
		fmt.Print("\nAssim ficará a renomeação. Deseja aplicar? (s/n): ")
		fmt.Scan(&confirma)

		if strings.ToLower(confirma) != "s" {
			fmt.Println("\n🔁 Operação cancelada. Reiniciando...\n")
			continue
		}

		start := time.Now()
		cpuCount := runtime.NumCPU()
		sem := make(chan struct{}, cpuCount)
		var wg sync.WaitGroup
		var renamed int64
		var mu sync.Mutex

		for _, r := range renomeacoes {
			wg.Add(1)
			sem <- struct{}{}

			go func(oldName, newName string) {
				defer wg.Done()
				defer func() { <-sem }()

				if err := os.Rename(oldName, newName); err != nil {
					fmt.Printf("❌ erro ao renomear %s: %v\n", oldName, err)
					return
				}

				mu.Lock()
				renamed++
				mu.Unlock()
			}(r.Antigo, r.Novo)
		}

		wg.Wait()
		elapsed := time.Since(start)
		fmt.Printf("\n📂 Foram renomeados %d arquivos em %.2f.\n", renamed, elapsed.Seconds(), cpuCount)
		break
	}
}
