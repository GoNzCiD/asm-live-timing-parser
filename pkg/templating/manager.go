package templating

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Notes from
// * https://github.com/sgulics/go-chi-example
// * https://gist.github.com/logrusorgru/abd846adb521a6fb39c7405f32fec0cf
// * https://github.com/asit-dhal/golang-template-layout/blob/master/src/templmanager/templatemanager.go

type TemplateManager struct {
	dir       string    // root directory
	layoutDir string    // template layout directory
	devel     bool      // reload every time
	loadedAt  time.Time // loaded at (last loading time)
	ext       string    // extension for template files
	templates map[string]*template.Template
}

func NewTemplateManager(dir string) (tmpl *TemplateManager, err error) {
	if dir, err = filepath.Abs(dir); err != nil {
		return
	}

	layoutDir := ""
	if layoutDir, err = filepath.Abs("templates/layouts"); err != nil {
		return
	}

	tmpl = new(TemplateManager)
	tmpl.dir = dir
	tmpl.layoutDir = layoutDir
	tmpl.ext = ".gohtml"
	tmpl.devel = false

	tmpl.templates = make(map[string]*template.Template)

	if err = tmpl.Load(); err != nil {
		tmpl = nil
	}

	return
}

func (t *TemplateManager) Load() (err error) {
	t.loadedAt = time.Now()

	layoutFiles, err := filepath.Glob(t.layoutDir + "/*.gohtml")
	if err != nil {
		return err
	}

	var walkFunc = func(path string, info os.FileInfo, err error) (_ error) {
		// handle walking error if any
		if err != nil {
			return err
		}

		// skip all except regular files
		if !info.Mode().IsRegular() {
			return
		}

		// filter by extension
		if filepath.Ext(path) != t.ext {
			return
		}

		// get relative path
		var rel string
		if rel, err = filepath.Rel(t.dir, path); err != nil {
			return err
		}

		// Ignore files in the layout directory
		if filepath.Dir(path) == t.layoutDir {
			return
		}

		// name of a template is its relative path
		// without extension
		rel = strings.TrimSuffix(rel, t.ext)

		var (
			nt = template.New(rel) //.Funcs(sprig.FuncMap()).Funcs(t.funcs).Funcs(viewHelpers())
			b  []byte
		)

		if b, err = os.ReadFile(path); err != nil {
			return err
		}
		tmpl, err := nt.ParseFiles(layoutFiles...)
		if err != nil {
			return err
		}

		tmpl, err = nt.Parse(string(b))
		if err != nil {
			return err
		}

		t.templates[tmpl.Name()] = tmpl

		return err
	}

	if err = filepath.Walk(t.dir, walkFunc); err != nil {
		return
	}

	return
}

// IsModified lookups directory for changes to
// reload (or not to reload) templates_oild if development
// pin is true.
func (t *TemplateManager) IsModified() (yep bool, err error) {

	var errStop = fmt.Errorf("stop")

	var walkFunc = func(path string, info os.FileInfo, err error) (_ error) {
		// handle walking error if any
		if err != nil {
			return err
		}

		// skip all except regular files
		if !info.Mode().IsRegular() {
			return
		}

		// filter by extension
		if filepath.Ext(path) != t.ext {
			return
		}

		if yep = info.ModTime().After(t.loadedAt); yep == true {
			return errStop
		}

		return
	}

	// clear the errStop
	if err = filepath.Walk(t.dir, walkFunc); err == errStop {
		err = nil
	}

	return
}

func (t *TemplateManager) Template(name string) (tmpl *template.Template, err error) {
	// if development
	if t.devel == true {
		// lookup directory for changes
		modified, err := t.IsModified()
		if err != nil {
			return nil, err
		}

		// reload
		if modified == true {
			if err = t.Load(); err != nil {
				return nil, err
			}
		}
	}

	tmpl, ok := t.templates[name]
	if !ok {
		return nil, fmt.Errorf("template not found")
	}
	return tmpl.Clone()
}

func (t *TemplateManager) Render(w io.Writer, name string, data interface{}) (err error) {

	tmpl, err := t.Template(name)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}
