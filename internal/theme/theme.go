package theme

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/davidborzek/spofi/pkg/rofi"
)

//go:embed theme.rasi
var themeResource embed.FS

// LoadTheme loads the embedded rofi theme and
// places it into a temporary directory (/tmp/spofi/theme.rasi), so rofi
// can use it.
func LoadTheme() {
	buf, err := themeResource.ReadFile("theme.rasi")
	if err != nil {
		log.Fatalln(err)
	}

	dir := filepath.Join(os.TempDir(), "spofi")

	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalln(err)
	}

	path := filepath.Join(dir, "theme.rasi")
	if err := os.WriteFile(path, buf, 0644); err != nil {
		log.Fatalln(err)
	}

	rofi.SetCustomTheme(path)
}
