package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("📚 pInk renamer")
	fmt.Println("---------------------------------------")
	fmt.Print("Escolha o formato de numeração (2 = dezena, 3 = centena): ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	numDigits := 2
	if input == "3" {
		numDigits = 3
	}

	fmt.Printf("→ Usando formatação com %d dígitos.\n\n", numDigits)

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("❌ Erro ao obter diretório atual: %v\n", err)
	}

	entries, err := os.ReadDir(currentDir)
	if err != nil {
		log.Fatalf("❌ Erro ao listar arquivos: %v\n", err)
	}

	trailingRe := regexp.MustCompile(`(?i)^(.*?)[ _-]?0*(\d+)$`)

	skippedSelf := false

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		oldName := entry.Name()

		if oldName == "renomeia.go" {
			if !skippedSelf {
				fmt.Printf("⚠️  Ignorando o script %s (por segurança).\n", oldName)
				skippedSelf = true
			}
			continue
		}

		ext := filepath.Ext(oldName)
		baseWithPossibleNum := strings.TrimSuffix(oldName, ext)

		var base string
		var numStr string
		var hasNum bool

		if matches := trailingRe.FindStringSubmatch(baseWithPossibleNum); len(matches) == 3 {
			base = strings.TrimSpace(matches[1])
			numStr = matches[2]
			hasNum = true
		} else {
			base, numStr, hasNum = extractLastNumber(baseWithPossibleNum)
		}

		if !hasNum {
			numStr = "1"
		}

		numInt, convErr := strconv.Atoi(numStr)
		if convErr != nil {
			fmt.Printf("❌ Não foi possível ler número em '%s' (interpretado: '%s'), pulando.\n", oldName, numStr)
			continue
		}

		formattedNum := fmt.Sprintf("%0*d", numDigits, numInt)

		cleanBase := strings.TrimRight(base, " _-")
		if cleanBase == "" {
			cleanBase = "item"
		}

		newName := fmt.Sprintf("%s-%s%s", cleanBase, formattedNum, ext)

		if oldName == newName {
			fmt.Printf("🔹 Mantido: %s\n", oldName)
			continue
		}

		if _, err := os.Stat(newName); err == nil {
			fmt.Printf("⚠️  Já existe: %s → %s (pulando para evitar colisão)\n", oldName, newName)
			continue
		}

		if err := os.Rename(oldName, newName); err != nil {
			fmt.Printf("❌ Falha ao renomear %s → %s: %v\n", oldName, newName, err)
			continue
		}

		fmt.Printf("✅ %s → %s\n", oldName, newName)
	}

	fmt.Println("\n✨ Renomeação concluída.")
}

func extractLastNumber(s string) (base string, num string, found bool) {
	runes := []rune(s)
	end := -1
	start := -1

	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] >= '0' && runes[i] <= '9' {
			if end == -1 {
				end = i + 1
			}
			start = i
		} else if end != -1 {
			break
		}
	}

	if start != -1 && end != -1 {
		num = string(runes[start:end])
		base = strings.TrimSpace(strings.TrimRight(string(runes[:start]), " _-"))
		return base, num, true
	}

	return s, "", false
}
