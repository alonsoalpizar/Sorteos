package notifier

import (
	"bytes"
	"embed"
	"html/template"
	"os"
	"path/filepath"
)

//go:embed templates/*.html
var embeddedTemplates embed.FS

// TemplateLoader carga plantillas de email desde archivos o embebidas
type TemplateLoader struct {
	templatesDir string
	useEmbedded  bool
	cache        map[string]*template.Template
}

// NewTemplateLoader crea un nuevo cargador de plantillas
func NewTemplateLoader(templatesDir string) *TemplateLoader {
	// Si no se especifica directorio o no existe, usar plantillas embebidas
	useEmbedded := templatesDir == "" || !dirExists(templatesDir)

	return &TemplateLoader{
		templatesDir: templatesDir,
		useEmbedded:  useEmbedded,
		cache:        make(map[string]*template.Template),
	}
}

// LoadTemplate carga una plantilla por nombre
func (tl *TemplateLoader) LoadTemplate(name string) (*template.Template, error) {
	// Verificar si está en caché
	if tmpl, ok := tl.cache[name]; ok {
		return tmpl, nil
	}

	var tmpl *template.Template
	var err error

	if tl.useEmbedded {
		// Cargar desde plantillas embebidas
		tmpl, err = template.ParseFS(embeddedTemplates, "templates/"+name)
	} else {
		// Cargar desde filesystem
		path := filepath.Join(tl.templatesDir, name)
		tmpl, err = template.ParseFiles(path)
	}

	if err != nil {
		return nil, err
	}

	// Guardar en caché
	tl.cache[name] = tmpl

	return tmpl, nil
}

// RenderTemplate renderiza una plantilla con los datos proporcionados
func (tl *TemplateLoader) RenderTemplate(name string, data interface{}) (string, error) {
	tmpl, err := tl.LoadTemplate(name)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ClearCache limpia el caché de plantillas
func (tl *TemplateLoader) ClearCache() {
	tl.cache = make(map[string]*template.Template)
}

// dirExists verifica si un directorio existe
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// VerificationEmailData datos para email de verificación
type VerificationEmailData struct {
	FirstName       string
	Code            string
	FrontendURL     string
	VerificationURL string // Opcional: si tienes link directo
}

// WelcomeEmailData datos para email de bienvenida
type WelcomeEmailData struct {
	FirstName   string
	FrontendURL string
}

// PasswordResetEmailData datos para email de reset de contraseña
type PasswordResetEmailData struct {
	FirstName   string
	ResetURL    string
	FrontendURL string
}

// PurchaseConfirmationData datos para email de confirmación de compra
type PurchaseConfirmationData struct {
	FirstName    string
	RaffleTitle  string
	RaffleID     string
	Numbers      []string
	TotalAmount  string
	DrawDate     string
	Prize        string
	FrontendURL  string
}

// WinnerNotificationData datos para email de ganador
type WinnerNotificationData struct {
	FirstName    string
	RaffleTitle  string
	WinnerNumber string
	Prize        string
	Instructions string
	FrontendURL  string
}
