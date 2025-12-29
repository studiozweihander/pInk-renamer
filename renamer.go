package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const (
	asciiColor = "\033[38;5;212m"
	resetColor = "\033[0m"
	separator  = "────────────────────────────────────────────────────────────"
)

type RenamePlan struct {
	OldPath string
	NewPath string
}

type Config struct {
	WorkDir string
	Digits  int
}

func main() {
	printASCII()

	cfg, err := parseArgs()
	if err != nil {
		fmt.Println()
		logError("%v", err)
		os.Exit(1)
	}

	if err := run(cfg); err != nil {
		fmt.Println()
		logError("%v", err)
		os.Exit(1)
	}
}

func run(cfg Config) error {
	logInfo("Diretório de trabalho: %s", cfg.WorkDir)

	entries, err := os.ReadDir(cfg.WorkDir)
	if err != nil {
		return fmt.Errorf("falha ao ler diretório")
	}

	plans, err := buildRenamePlan(cfg.WorkDir, entries, cfg.Digits)
	if err != nil {
		return err
	}

	if len(plans) == 0 {
		return fmt.Errorf("nenhum arquivo elegível para renomeação encontrado")
	}

	fmt.Println()
	fmt.Println(separator)
	fmt.Println()

	logInfo("Preview dos arquivos que serão renomeados")
	fmt.Println()

	for _, p := range plans {
		fmt.Printf("  %s → %s\n",
			filepath.Base(p.OldPath),
			filepath.Base(p.NewPath),
		)
	}

	fmt.Println()
	if !confirmExecution() {
		logInfo("Operação cancelada pelo usuário")
		return nil
	}

	fmt.Println()
	fmt.Println(separator)
	fmt.Println()

	if err := validateCollisions(plans); err != nil {
		return err
	}

	jobs := make(chan RenamePlan, len(plans))
	var wg sync.WaitGroup

	workers := runtime.NumCPU()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go renameWorker(jobs, &wg)
	}

	for _, p := range plans {
		jobs <- p
	}
	close(jobs)

	wg.Wait()

	fmt.Println()
	logSuccess(
		"Processo finalizado com sucesso. Renomeados: %d arquivos.",
		len(plans),
	)

	return nil
}

func renameWorker(jobs <-chan RenamePlan, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		if err := os.Rename(job.OldPath, job.NewPath); err != nil {
			fmt.Println()
			logError("erro ao renomear %s", filepath.Base(job.OldPath))
			os.Exit(1)
		}

		fmt.Printf("%s → %s\n",
			filepath.Base(job.OldPath),
			filepath.Base(job.NewPath),
		)
	}
}

func buildRenamePlan(dir string, entries []os.DirEntry, digits int) ([]RenamePlan, error) {
	var plans []RenamePlan

	re := regexp.MustCompile(`(?i)^\s*([a-z0-9\s._-]+?)\s*#?\s*(\d+)`)

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		ext := filepath.Ext(name)
		base := strings.TrimSuffix(name, ext)

		matches := re.FindStringSubmatch(base)
		if len(matches) != 3 {
			continue
		}

		title := matches[1]
		number, _ := strconv.Atoi(matches[2])
		numFmt := fmt.Sprintf("%0*d", digits, number)

		title = strings.ToLower(title)
		title = strings.ReplaceAll(title, " ", "-")
		title = strings.ReplaceAll(title, "_", "-")
		title = strings.ReplaceAll(title, ".", "-")
		title = strings.Trim(title, "-")

		newName := fmt.Sprintf("%s-%s%s", title, numFmt, ext)

		if name == newName {
			continue
		}

		plans = append(plans, RenamePlan{
			OldPath: filepath.Join(dir, name),
			NewPath: filepath.Join(dir, newName),
		})
	}

	return plans, nil
}

func validateCollisions(plans []RenamePlan) error {
	seen := make(map[string]string)

	for _, p := range plans {
		if prev, exists := seen[p.NewPath]; exists {
			return fmt.Errorf(
				"dois arquivos resultariam no mesmo nome: %s e %s",
				filepath.Base(prev),
				filepath.Base(p.OldPath),
			)
		}
		seen[p.NewPath] = p.OldPath
	}

	return nil
}

func parseArgs() (Config, error) {
	var cfg Config

	args := os.Args[1:]

	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		abs, err := filepath.Abs(args[0])
		if err != nil {
			return cfg, fmt.Errorf("diretório inválido")
		}

		info, err := os.Stat(abs)
		if err != nil || !info.IsDir() {
			return cfg, fmt.Errorf("diretório não encontrado: %s", abs)
		}

		cfg.WorkDir = abs
	} else {
		dir, err := os.Getwd()
		if err != nil {
			return cfg, fmt.Errorf("falha ao obter diretório atual")
		}
		cfg.WorkDir = dir
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Escolha o formato de numeração (2 ou 3 dígitos): ")
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	d, err := strconv.Atoi(line)
	if err != nil || (d != 2 && d != 3) {
		return cfg, fmt.Errorf("formato inválido")
	}

	cfg.Digits = d
	return cfg, nil
}

func confirmExecution() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Deseja continuar? (s/N): ")

	resp, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	resp = strings.TrimSpace(strings.ToLower(resp))
	return resp == "s" || resp == "sim"
}

func printASCII() {
	fmt.Print(asciiColor)
	fmt.Print(`
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
	fmt.Print(resetColor)
}

func logInfo(msg string, args ...any) {
	fmt.Printf("\033[34m[INFO]\033[0m "+msg+"\n", args...)
}

func logSuccess(msg string, args ...any) {
	fmt.Printf("\033[32m[SUCCESS]\033[0m "+msg+"\n", args...)
}

func logError(msg string, args ...any) {
	fmt.Printf("\033[31m[ERROR]\033[0m "+msg+"\n", args...)
}
