package main

import (
	"archive/zip"
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	_dialog "fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/nwaples/rardecode"
	"github.com/sqweek/dialog"
)

var iconData []byte

func loadIcon() []byte {
	return iconData
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Validador de Backup SQL")

	iconResource := fyne.NewStaticResource("icon", loadIcon())
	myWindow.SetIcon(iconResource)

	resultLabel := widget.NewLabel("")
	spinner := widget.NewProgressBarInfinite()
	spinner.Hide()

	// Botão de informação
	infoButton := widget.NewButton("INFO", func() {
		_dialog.ShowInformation("🦗 Grilo 🦗", "Criado por Matheus Ferreira", myWindow)
	})

	// Botão para selecionar arquivo
	selectButton := widget.NewButton("Selecionar backup", func() {
		filePath, err := dialog.File().
			Filter("SQL Files, ZIP, RAR", "sql", "zip", "rar").
			Filter("All Files", "*").
			Title("Selecionar arquivo do backup").
			Load()
		if err == nil {
			spinner.Show()
			resultLabel.SetText("Processando...")

			go func() {
				defer spinner.Hide()
				processFile(filePath, resultLabel)
			}()
		} else {
			resultLabel.SetText("Erro ao abrir o arquivo: " + err.Error())
		}
	})

	footer := container.NewGridWithColumns(2,
		selectButton,
		infoButton,
	)

	content := container.NewVBox(
		resultLabel,
		spinner,
		layout.NewSpacer(),
		footer,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 300))
	myWindow.SetFixedSize(true)

	myWindow.ShowAndRun()
}

func processFile(filePath string, resultLabel *widget.Label) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".sql":
		validateSQLFile(filePath, resultLabel)
	case ".zip":
		err := processZipFile(filePath, resultLabel)
		if err != nil {
			resultLabel.SetText("Erro ao processar o arquivo ZIP: " + err.Error())
		}
	case ".rar":
		err := processRarFile(filePath, resultLabel)
		if err != nil {
			resultLabel.SetText("Erro ao processar o arquivo RAR: " + err.Error())
		}
	default:
		resultLabel.SetText("Formato de arquivo não suportado.")
	}
}

func processZipFile(zipPath string, resultLabel *widget.Label) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		if strings.HasSuffix(file.Name, ".sql") {
			rc, err := file.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			reader := bufio.NewReader(rc)
			validateSQLFromReader(reader, resultLabel)
			return nil
		}
	}
	resultLabel.SetText("Nenhum arquivo SQL encontrado no ZIP.")
	return nil
}

func processRarFile(rarPath string, resultLabel *widget.Label) error {
	file, err := os.Open(rarPath)
	if err != nil {
		return err
	}
	defer file.Close()

	rarReader, err := rardecode.NewReader(file, "")
	if err != nil {
		return err
	}

	for {
		header, err := rarReader.Next()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		if strings.HasSuffix(header.Name, ".sql") {
			reader := bufio.NewReader(rarReader)
			validateSQLFromReader(reader, resultLabel)
			return nil
		}
	}
	resultLabel.SetText("Nenhum arquivo SQL encontrado no RAR.")
	return nil
}

func validateSQLFromReader(reader *bufio.Reader, resultLabel *widget.Label) {
	var tableCount int
	var lastTables []string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				resultLabel.SetText("Erro durante a leitura do arquivo: " + err.Error())
			}
			break
		}
		if strings.HasPrefix(line, "CREATE TABLE") {
			tableCount++
			tableName := extractTableName(line)
			lastTables = append(lastTables, tableName)
			if len(lastTables) > 5 {
				lastTables = lastTables[1:]
			}
		}
	}

	result := fmt.Sprintf("O arquivo tem %d tabelas\n\nÚltimas 5 tabelas criadas:\n", tableCount)

	for _, table := range lastTables {
		result += fmt.Sprintf("\n- %s", table)
	}

	if isBackupComplete(lastTables) {
		result += "\n\n\nBackup completo"
	} else {
		result += "\n\n\nBackup incompleto"
	}

	resultLabel.SetText(result)
}

func validateSQLFile(filepath string, resultLabel *widget.Label) {
	file, err := os.Open(filepath)
	if err != nil {
		resultLabel.SetText("Erro ao abrir o arquivo: " + err.Error())
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	validateSQLFromReader(reader, resultLabel)
}

func extractTableName(line string) string {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return ""
	}
	return strings.Trim(parts[2], "`")
}

func isBackupComplete(tables []string) bool {
	for _, table := range tables {
		if strings.HasPrefix(table, "whatsapp_") {
			return true
		}
	}
	return false
}
