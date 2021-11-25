package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/cmxjs/file-server/config"
	"github.com/cmxjs/file-server/models"
	"github.com/gin-gonic/gin"
)

type File struct {
	Path string `uri:"path"`
}

func Upload(c *gin.Context) {
	f, err := c.FormFile("f")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": "fail", "err": err.Error()})
		return
	}

	path, pathOk := c.GetPostForm("path")
	if !pathOk {
		c.JSON(http.StatusBadRequest, gin.H{"result": "fail", "err": "Failed to get path from form"})
		return
	}

	dst := filepath.Join(config.Cache, path, f.Filename)
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": "fail", "err": err.Error()})
		return
	}

	if err := c.SaveUploadedFile(f, dst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": "fail", "err": err.Error()})
		return
	}

	// update ttl
	ttl, ttlOk := c.GetPostForm("ttl")
	if ttlOk {
		ttl = strings.Trim(ttl, " ")

		ttlInt, err := strconv.ParseInt(ttl, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"result": "fail", "err": err.Error()})
			return
		}

		if err := models.UpdateTTL(dst, ttlInt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"result": "fail", "err": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"result": "success", "path": dst})
}

func TTL(c *gin.Context) {
	ttls, err := models.GetAllTTL()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Failed to get ttl info from db, err: %v\n", err.Error()))
		return
	}
	c.JSON(http.StatusOK, ttls)
}

func Download(c *gin.Context) {
	var file File
	if err := c.ShouldBindUri(&file); err != nil {
		log.Println(err)
		c.String(http.StatusNotAcceptable, fmt.Sprintln("Failed to get download file from url"))
		return
	}

	absolute_path := filepath.Join(config.Cache, file.Path)
	log.Println(absolute_path)

	if _, err := os.Stat(absolute_path); err != nil {
		c.String(http.StatusNotAcceptable, fmt.Sprintln("error, download file not exists"))
		return
	}

	_, filename := filepath.Split(absolute_path)
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(absolute_path)
}

func Preview(c *gin.Context, t *template.Template) {
	var file File
	if err := c.ShouldBindUri(&file); err != nil {
		log.Println(err)
		c.String(http.StatusNotAcceptable, fmt.Sprintln("Failed to get preview file from url"))
		return
	}

	absolute_path := filepath.Join(config.Cache, file.Path)
	log.Println(absolute_path)

	if _, err := os.Stat(absolute_path); err != nil {
		c.String(http.StatusNotAcceptable, fmt.Sprintln("error, preview file not exists"))
		return
	}

	r := regexp.MustCompile(`\.(\w+)$`)
	if r.FindStringSubmatch(absolute_path) == nil { // used when no suffix exists
		PreviewLog(c, absolute_path, t)
		return
	}

	switch r.FindStringSubmatch(absolute_path)[1] {
	case "docx":
		PreviewDocx(c, absolute_path)
	case "pdf":
		PreviewDocx(c, absolute_path)
	case "html":
		PreviewHtml(c, absolute_path)
	case "js":
		PreviewCode(c, absolute_path, "json")
	case "json":
		PreviewCode(c, absolute_path, "json")
	case "py":
		PreviewCode(c, absolute_path, "py")
	case "md":
		PreviewCode(c, absolute_path, "md")
	default:
		PreviewLog(c, absolute_path, t)
	}
}
