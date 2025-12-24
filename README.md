```
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
```

## Visão Geral

pInk renamer é um script que automatiza o processo de renomeação dos arquivos disponibilizados no projeto pInk. 

## Funcionalidades

- Extração automática de números
- Formatação em dezenas (01, 02) ou centenas (001, 002)
- Preview antes de aplicar
- Processamento paralelo

## Como Usar
### Execução Padrão

```
go run renomeia.go
```
Obs: Essa execução rodará no diretório em que o script se encontra.

### Execução com Diretório Específico

```
go run renomeia.go ./meu_diretorio
```

## Formatação dos Números

### 1. Formato Dezena (2 dígitos)
```
Batman-01.cbr
Batman-02.cbr
...
Batman-10.cbr
```

### 2. Formato Centena (3 dígitos)
```
Batman-001.cbr
Batman-002.cbr
...
Batman-100.cbr
```

## Alterações Aplicadas

| Antes | Depois | Regra |
|-------|--------|-------|
| `Batman 01.cbr` | `batman-01.cbr` | Espaços → hífens |
| `Batman.01.cbz` | `batman-01.cbz` | Pontos → hífens |
| `BATMAN_01.cbr` | `batman-01.cbr` | Minúsculas |

## Exemplo Completo
```
$ go run renomeia.go ./quadrinhos

Escolha o formato de numeração (2 = dezena, 3 = centena): 2

[SUCCESS] Batman 01.cbr → batman-01.cbr
[SUCCESS] Batman.02.cbz → batman-02.cbz
[SUCCESS] BATMAN_03.cbr → batman-03.cbr

Assim ficará a renomeação. Deseja aplicar? (s/n): s

[SUCCESS] Foram renomeados 3 arquivos em 0.01 segundos.
```

## Limitações e Observações
1. **Detecção de números**: Extrai o último grupo numérico no nome
2. **Case-insensitive**: Converte tudo para minúsculas
3. **Hífens múltiplos**: `batman--01.cbr` → `batman-01.cbr`
4. **Arquivos sem números**: São ignorados
5. **Espaços**: Todos os espaços são transformados em hífens
