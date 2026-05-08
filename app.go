package main

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx  context.Context
	lang string
}

// translations holds all UI strings per language
var translations = map[string]map[string]string{
	"it": {
		"urlEmpty":    "L'URL non può essere vuoto.",
		"startEmpty":  "Il tempo di inizio non può essere vuoto.",
		"endEmpty":    "Il tempo di fine non può essere vuoto.",
		"cancelled":   "Operazione annullata.",
		"errorPrefix": "Errore durante il download:\n",
		"savedIn":     "Video salvato in: ",
	},
	"en": {
		"urlEmpty":    "URL cannot be empty.",
		"startEmpty":  "Start time cannot be blank.",
		"endEmpty":    "End time cannot be blank.",
		"cancelled":   "Operation cancelled.",
		"errorPrefix": "Error while downloading:\n",
		"savedIn":     "Video saved in: ",
	},
}

// DownloadResult is the result returned to the frontend
type DownloadResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// SetLanguage sets the active language for backend messages
func (a *App) SetLanguage(lang string) {
	if _, ok := translations[lang]; ok {
		a.lang = lang
	}
}

// t is a helper that returns the translated string for a key
func (a *App) t(key string) string {
	if msgs, ok := translations[a.lang]; ok {
		if val, ok := msgs[key]; ok {
			return val
		}
	}
	return key
}

// CheckYtDlp checks if yt-dlp is installed
func (a *App) CheckYtDlp() bool {
	_, err := exec.LookPath("yt-dlp")
	return err == nil
}

// DownloadSegment downloads a YouTube video segment using yt-dlp
func (a *App) DownloadSegment(url, start, end, output string) DownloadResult {

	if strings.TrimSpace(url) == "" {
		return DownloadResult{Success: false, Message: a.t("urlEmpty")}
	}

	if strings.TrimSpace(start) == "" {
		return DownloadResult{Success: false, Message: a.t("startEmpty")}
	}

	if strings.TrimSpace(end) == "" {
		return DownloadResult{Success: false, Message: a.t("endEmpty")}
	}

	if strings.TrimSpace(output) == "" {
		output = "video.mp4"
	}

	// Ensure output has .mp4 extension
	if filepath.Ext(output) == "" {
		output = output + ".mp4"
	}

	// Opens the native dialog to choose where to save
	path, err := wailsRuntime.SaveFileDialog(a.ctx, wailsRuntime.SaveDialogOptions{
		DefaultFilename: output,
		Filters: []wailsRuntime.FileFilter{
			{DisplayName: "Video (*.mp4)", Pattern: "*.mp4"},
		},
	})

	if err != nil || path == "" {
		return DownloadResult{Success: false, Message: a.t("cancelled")}
	}

	section := fmt.Sprintf("*%s-%s", start, end)

	cmd := exec.CommandContext(a.ctx,
		"yt-dlp",
		url,
		"--download-sections", section,
		"-o", path,
		"--merge-output-format", "mp4",
	)

	combinedOutput, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := strings.TrimSpace(string(combinedOutput))
		if errMsg == "" {
			errMsg = err.Error()
		}
		return DownloadResult{
			Success: false,
			Message: a.t("errorPrefix") + errMsg,
		}
	}

	return DownloadResult{
		Success: true,
		Message: a.t("savedIn") + path,
	}
}
