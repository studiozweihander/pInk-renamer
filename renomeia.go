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

func logInfo(msg string, args ...interface{}) {
	fmt.Printf("\033[34m[INFO]\033[0m "+msg+"\n", args...)
}

func logSuccess(msg string, args ...interface{}) {
	fmt.Printf("\033[32m[SUCCESS]\033[0m "+msg+"\n", args...)
}

func logError(msg string, args ...interface{}) {
	fmt.Printf("\033[31m[ERROR]\033[0m "+msg+"\n", args...)
}

type Renomeacao struct {
	Antigo string
	Novo   string
}

func main() {
	for {
		fmt.Println(`
           █████            █████                                                                                   
          ░░███            ░░███                                                                                    
 ████████  ░███  ████████   ░███ █████    ████████   ██████  ████████    ██████   █████████████    ██████  ████████ 
░░███░░███ ░███ ░░███░░███  ░███░░███    ░░███░░███ ███░░███░░███░░███  ░░░░░███ ░░███░░███░░███  ███░░███░░███░░███
 ░███ ░███ ░███  ░███ ░███  ░██████░      ░███ ░░░ ░███████  ░███ ░███   ███████  ░███ ░███ ░███ ░███████  ░███ ░░░ 
 ░███ ░███ ░███  ░███ ░███  ░███░░███     ░███     ░███░░░   ░███ ░███  ███░░███  ░███ ░███ ░███ ░███░░░   ░███     
 ░███████  █████ ████ █████ ████ █████    █████    ░░██████  ████ █████░░████████ █████░███ █████░░██████  █████    
 ░███░░░  ░░░░░ ░░░░ ░░░░░ ░░░░ ░░░░░    ░░░░░      ░░░░░░  ░░░░ ░░░░░  ░░░░░░░░ ░░░░░ ░░░ ░░░░░  ░░░░░░  ░░░░░     
 ░███                                                                                                               
 █████                                                                                                              
░░░░░                                                                                                               
	`)

		var dir string
		if len(os.Args) > 1 {
			dir = os.Args[1]
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				logError(fmt.Sprintf("A pasta especificada não existe: %s", dir))
				return
			}
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				logError(fmt.Sprintf("Falha ao obter diretório atual: %v", err))
				return
			}
			dir = cwd
		}

		var formato int
		fmt.Print("\nEscolha o formato de numeração (2 = dezena, 3 = centena): ")
		fmt.Scan(&formato)

		if formato != 2 && formato != 3 {
			logError("Formato inválido. Use 2 ou 3.")
			continue
		}

		files, err := os.ReadDir(dir)
		if err != nil {
			logError(fmt.Sprintf("Falha ao ler diretório: %v", err))
			return
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
				logInfo(fmt.Sprintf("Mantido: %s", oldName))
				continue
			}

			if oldName != newName {
				renomeacoes = append(renomeacoes, Renomeacao{
					Antigo: filepath.Join(dir, oldName),
					Novo:   filepath.Join(dir, newName),
				})
				logSuccess(fmt.Sprintf("%s → %s", oldName, newName))
			}
		}

		if len(renomeacoes) == 0 {
			fmt.Println("")
			logError("Nenhum arquivo elegível para renomeação encontrado.")
			return
		}

		var confirma string
		fmt.Print("\nAssim ficará a renomeação. Deseja aplicar? (s/n): ")
		fmt.Scan(&confirma)

		if strings.ToLower(confirma) != "s" {
			logInfo("Operação cancelada. Reiniciando...\n")
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
					logError(fmt.Sprintf("Erro ao renomear %s: %v", oldName, err))
					return
				}

				mu.Lock()
				renamed++
				mu.Unlock()
			}(r.Antigo, r.Novo)
		}

		wg.Wait()
		elapsed := time.Since(start)
		logSuccess(fmt.Sprintf("Foram renomeados %d arquivos em %.2f segundos.", renamed, elapsed.Seconds()))
		break
	}
}

