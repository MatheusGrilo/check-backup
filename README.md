# Verificar Backups SQL
Aplicativo feito para verificar Backups gerados pelo MySQL 5.7.9, seja de arquivos SQL ou compactados (RAR/ZIP) baseado nas regras que informamos (ter ou não a tabela whatsapp_)

## Compilar
Versão do GO utilizada: 
> go version go1.23.1 windows/amd64

Precisa realizar a instalação do Fyne2 também:
> [https://fyne.io/fyne/v2](https://fyne.io/fyne/v2)

Comandos utilizados:
```bash
# Usar go mod tidy para atualizar arquivo go.sum
go mod tidy

# Compilar com o ome Backup-Check.exe
go build -o Backup-Check.exe -ldflags="-H windowsgui" main.go

# pnpm
go build -o SeuApp.exe -ldflags="-H windowsgui" 
```

## O que o aplicativo faz (Fluxo)?

1. Abre o aplicativo e pede para selecionar um arquivo (SQL, RAR, ZIP)
2. Verifica quantas tabelas existem pelo CREATE TABLE \`nome_da_tabela\`
3. Armazena as informações (quantidade de tabelas)
4. Mostra em tela as informações
5. Diz se o backup está completo ou não baseado na existência das tabelas \`whatsapp_\`