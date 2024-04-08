package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jmrflora/blogx/views"
	"github.com/jmrflora/blogx/views/paginas"
	"github.com/labstack/echo/v4"
)

func HandleUpload(c echo.Context) error {
	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination

	// Destination directory
	uploadDir := "../internal/assets/markdowns"

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		println("ola")
		// Create the directory if it doesn't exist
		if err := os.Mkdir(uploadDir, 0777); err != nil {
			return err
		}
	}
	dstPath := filepath.Join(uploadDir, file.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully</p>", file.Filename))
}

func HandleIndex(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	cmp := paginas.Index()

	return views.Renderizar(cmp, c)
}
