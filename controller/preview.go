package controller

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cmxjs/file-server/config"
	"github.com/cmxjs/file-server/models"
	"github.com/gin-gonic/gin"
)

func PreviewCode(c *gin.Context, path string, fileType string) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Failed to read file content. err: %v", err.Error()))
		return
	}

	relPath, _ := filepath.Rel(config.Cache, path)

	c.HTML(http.StatusOK, "code.tmpl", gin.H{
		"FileName":    filepath.Base(path),
		"FilePath":    relPath,
		"Title":       fmt.Sprintf("%s preview", fileType),
		"CodeType":    fileType,
		"CodeContent": string(content),
	})
}

func PreviewDocx(c *gin.Context, path string) {
	m, _ := regexp.Compile(`\.docx$`)
	_file := filepath.Base(path)
	if m.FindString(_file) != "" {
		path = docx2pdf(path)
	}

	relPath, _ := filepath.Rel(config.Cache, path)
	// c.Redirect(http.StatusMovedPermanently, filepath.Join("/files", relPath))  //301 永久移动
	// c.Redirect(http.StatusPermanentRedirect, filepath.Join("/files", relPath)) //308 永久重定向
	c.Redirect(http.StatusTemporaryRedirect, filepath.Join("/files", relPath)) //307 临时重定向
}

func PreviewHtml(c *gin.Context, path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Failed to read file content. err: %v", err.Error()))
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
}

func docx2pdf(word_path string) string {
	outdir, file := filepath.Split(word_path)

	r, _ := regexp.Compile(`\.docx$`)
	file = r.ReplaceAllString(file, ".pdf")
	if _, err := os.Stat(filepath.Join(outdir, file)); err == nil {
		return filepath.Join(outdir, file)
	}

	c := exec.Command("libreoffice", "--headless", "--convert-to", "pdf:writer_pdf_Export", word_path, "--outdir", outdir)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Run()
	return filepath.Join(outdir, file)
}

type Log struct {
	Color string
	Info  string
}

func myCompile(e string) *regexp.Regexp { r, _ := regexp.Compile(e); return r }

var colors = map[string]*regexp.Regexp{
	"colorBlack":  myCompile(`\033\[1;30m`),
	"colorRed":    myCompile(`\033\[1;31m`),
	"colorGreen":  myCompile(`\033\[1;32m`),
	"colorOrange": myCompile(`\033\[1;33m`),
	"colorBlue":   myCompile(`\033\[1;34m`),
	"colorIndigo": myCompile(`\033\[1;35m`),
}

func color_log(log_content string) []Log {
	non_empty, _ := regexp.Compile(`\S+`)
	remove_color, _ := regexp.Compile(`\033\[0;35;0m`)
	end_color, _ := regexp.Compile(`\033\[0m`)

	startTestCase, _ := regexp.Compile(`Start a Testcase: (.+) \(\d+-\d+-\d+ \d+\:`)
	endTestCase, _ := regexp.Compile(`End all steps`)
	var mainTestCase bool = true // set false to include source testcase

	var new_content []Log
	var found_color_flag int = 0
	lines := strings.Split(log_content, "\n")
	for _, v := range lines {
		if non_empty.FindString(v) == "" {
			new_content = append(new_content, Log{
				Color: "colorBlack",
				Info:  v,
			})
			continue
		}

		if startTestCase.FindStringSubmatch(v) != nil {
			if mainTestCase {
				mainTestCase = false
			} else {
				new_content = append(new_content, Log{
					Color: "StartCase",
					Info:  strings.Trim(startTestCase.FindStringSubmatch(v)[1], ""),
				})
			}
		}
		if endTestCase.FindStringSubmatch(v) != nil {
			new_content = append(new_content, Log{
				Color: "EndCase",
				Info:  "",
			})
		}

		v = end_color.ReplaceAllString(v, "")
		v = remove_color.ReplaceAllString(v, "")

		found_color_flag = 0
		for color, r := range colors {
			if r.FindString(v) != "" {
				new_content = append(new_content, Log{
					Color: color,
					Info:  r.ReplaceAllString(v, ""),
				})
				found_color_flag = 1
				break
			}
		}

		if found_color_flag == 0 {
			new_content = append(new_content, Log{
				Color: "colorBlack",
				Info:  v,
			})
		}
	}

	return new_content
}

func PreviewLog(c *gin.Context, path string, t *template.Template) {
	value, err := models.GetLog(path)
	if err == nil {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, value)
		return
	}

	content, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Failed to read file content. err: %v", err.Error()))
		return
	}

	log_content := color_log(string(content))

	b := bytes.Buffer{}
	relPath, _ := filepath.Rel(config.Cache, path)
	if err = t.ExecuteTemplate(&b, "log.tmpl", gin.H{
		"Log_content": log_content,
		"Line_width":  len(strconv.Itoa(len(log_content))) * 10,
		"Log_file":    relPath,
	}); err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Failed to render html template log.tmpl. err: %v", err.Error()))
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, b.String())

	models.SetLog(path, b.String(), time.Hour*time.Duration(6))
}
