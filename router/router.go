package router

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cmxjs/file-server/config"
	"github.com/cmxjs/file-server/controller"
	"github.com/gin-gonic/gin"
)

func loadTemplate() (*template.Template, error) {
	absDir := func() string {
		absPath, err := filepath.Abs(os.Args[0])
		if err != nil {
			panic(err)
		}
		return filepath.Dir(absPath)
	}()

	templates := filepath.Join(absDir, "templates")
	files, err := os.ReadDir(templates)
	if err != nil {
		return nil, err
	}

	t := template.New("").Funcs(template.FuncMap{"add": func(a int, b int) int { return a + b }})
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".tmpl") {
			continue
		}
		h, err := os.ReadFile(filepath.Join(templates, file.Name()))
		if err != nil {
			return nil, err
		}
		t, err = t.New(file.Name()).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, err
}

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.MaxMultipartMemory = 64 << 20
	myTemplate, err := loadTemplate()
	if err != nil {
		log.Fatalf("Failed to load templates. err: %v", err.Error())
	}
	r.SetHTMLTemplate(myTemplate)

	r.StaticFS("/files", http.Dir(config.Cache))
	r.POST("/upload", controller.Upload)
	r.GET("/ttl", controller.TTL)
	r.GET("/download/*path", controller.Download)
	r.GET("/preview/*path", func(c *gin.Context) {
		controller.Preview(c, myTemplate)
	})

	return r
}
