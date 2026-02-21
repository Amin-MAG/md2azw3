package handler

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/leotaku/mobi"
	"golang.org/x/text/language"

	ravandlog "github.com/Amin-MAG/md2azw3/pkg/log"
	"github.com/labstack/echo/v4"
)

// ConvertHandler handles markdown to AZW3 conversion requests.
type ConvertHandler struct {
	logger *ravandlog.Logger
}

// NewConvertHandler creates a new ConvertHandler.
func NewConvertHandler(logger *ravandlog.Logger) *ConvertHandler {
	return &ConvertHandler{logger: logger}
}

// Convert handles POST /convert.
// Accepts multipart form with:
//   - "markdown": the .md file (required)
//   - "cover": cover image file (optional)
//   - "title": book title (optional)
//   - "author": author name (optional)
func (h *ConvertHandler) Convert(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse the markdown file
	mdFile, err := c.FormFile("markdown")
	if err != nil {
		h.logger.WithError(err).Warn(ctx, "missing markdown file in request")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "markdown file is required",
		})
	}

	// Read the markdown content
	mdContent, err := readUploadedFile(mdFile)
	if err != nil {
		h.logger.WithError(err).Error(ctx, "failed to read markdown file")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to read markdown file",
		})
	}

	// Convert markdown to HTML
	htmlContent := mdToHTML(mdContent)

	// Build the book
	title := c.FormValue("title")
	if title == "" {
		title = replaceExt(filepath.Base(mdFile.Filename), "")
	}

	book := mobi.Book{
		Title:       title,
		CreatedDate: time.Now(),
		Language:    language.English,
		UniqueID:    rand.Uint32(),
		Chapters: []mobi.Chapter{
			{
				Title:  title,
				Chunks: mobi.Chunks(htmlContent),
			},
		},
	}

	if author := c.FormValue("author"); author != "" {
		book.Authors = []string{author}
	}

	// Handle optional cover image
	coverFile, coverErr := c.FormFile("cover")
	if coverErr == nil && coverFile != nil {
		coverImg, err := decodeUploadedImage(coverFile)
		if err != nil {
			h.logger.WithError(err).Error(ctx, "failed to decode cover image")
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to process cover image",
			})
		}
		book.CoverImage = coverImg
	}

	// Realize and write the AZW3 file
	h.logger.Info(ctx, "generating azw3")
	db := book.Realize()

	tmpDir, err := os.MkdirTemp("", "md2azw3-*")
	if err != nil {
		h.logger.WithError(err).Error(ctx, "failed to create temp directory")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	}
	defer os.RemoveAll(tmpDir)

	mdFilename := filepath.Base(mdFile.Filename)
	outputFilename := replaceExt(mdFilename, ".azw3")
	outputPath := filepath.Join(tmpDir, outputFilename)

	f, err := os.Create(outputPath)
	if err != nil {
		h.logger.WithError(err).Error(ctx, "failed to create output file")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	}

	if err = db.Write(f); err != nil {
		f.Close()
		h.logger.WithError(err).Error(ctx, "failed to write azw3")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to generate azw3",
		})
	}
	f.Close()

	h.logger.Info(ctx, "conversion successful, returning file")
	return c.Attachment(outputPath, outputFilename)
}

func mdToHTML(md []byte) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	opts := html.RendererOptions{Flags: html.CommonFlags | html.HrefTargetBlank}
	renderer := html.NewRenderer(opts)
	return string(markdown.ToHTML(md, p, renderer))
}

func readUploadedFile(fh *multipart.FileHeader) ([]byte, error) {
	src, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("open uploaded file: %w", err)
	}
	defer src.Close()
	return io.ReadAll(src)
}

func decodeUploadedImage(fh *multipart.FileHeader) (image.Image, error) {
	src, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("open uploaded file: %w", err)
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}
	return img, nil
}

func replaceExt(filename, newExt string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return filename + newExt
	}
	return filename[:len(filename)-len(ext)] + newExt
}
