package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("📚 pInk renamer")
	fmt.Println("---------------------------------------")

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var formato int
	fmt.Print("Escolha o formato de numeração (2 = dezena, 3 = centena): ")
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

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		oldName := f.Name()
		ext := filepath.Ext(oldName)
		base := strings.TrimSuffix(oldName, ext)

		if ext != ".cbr" && ext != ".cbz" {
			continue
		}

		numeroMatch := reNumero.FindStringSubmatch(base)
		if len(numeroMatch) == 0 {
			continue
		}

		numero, _ := strconv.Atoi(numeroMatch[1])
		numeroFmt := fmt.Sprintf("%0*d", formato, numero)

		newBase := strings.ToLower(fmt.Sprintf("hellblazer-%s", numeroFmt))
		newName := newBase + ext

		hasUpper := oldName != strings.ToLower(oldName)

		if !hasUpper && oldName == newName {
			fmt.Printf("🔷 Mantido: %s\n", oldName)
			continue
		}

		err := os.Rename(oldName, newName)
		if err != nil {
			fmt.Printf("❌ erro ao renomear %s: %v\n", oldName, err)
			continue
		}

		fmt.Printf("✅ %s → %s\n", oldName, newName)
	}
}
