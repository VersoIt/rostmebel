package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/rostmebel/backend/internal/domain/apperror"
	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

const (
	defaultImageVariantQuality = 82
	minImageVariantWidth       = 64
	maxImageVariantWidth       = 2560
	imageVariantCacheDirName   = ".variants"
)

var imageVariantLocks sync.Map

func (h *ProductHandler) GetImageVariant(w http.ResponseWriter, r *http.Request) {
	src := strings.TrimSpace(r.URL.Query().Get("src"))
	if src == "" {
		respondWithError(w, invalidQuery("src", src))
		return
	}

	width, quality, err := parseImageVariantParams(r)
	if err != nil {
		respondWithError(w, err)
		return
	}

	sourcePath, err := resolveUploadedImagePath(src)
	if err != nil {
		respondWithError(w, err)
		return
	}

	stat, err := os.Stat(sourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			respondWithError(w, apperror.New(apperror.CodeProjectNotFound, "Image not found", map[string]any{"src": src}))
			return
		}
		respondWithError(w, err)
		return
	}

	if stat.IsDir() {
		respondWithError(w, invalidQuery("src", src))
		return
	}

	config, format, err := decodeImageConfig(sourcePath)
	if err != nil {
		respondWithError(w, err)
		return
	}

	if config.Width <= width {
		serveStaticImageFile(w, r, sourcePath, detectImageContentType(sourcePath, format))
		return
	}

	cachePath, contentType := imageVariantCachePath(sourcePath, stat, width, quality, format)
	if _, cacheErr := os.Stat(cachePath); cacheErr == nil {
		serveStaticImageFile(w, r, cachePath, contentType)
		return
	}

	lock := imageVariantLock(cachePath)
	lock.Lock()
	defer lock.Unlock()

	if _, cacheErr := os.Stat(cachePath); cacheErr == nil {
		serveStaticImageFile(w, r, cachePath, contentType)
		return
	}

	if err := generateImageVariant(sourcePath, cachePath, width, quality, format); err != nil {
		respondWithError(w, err)
		return
	}

	if _, err := os.Stat(cachePath); err != nil {
		respondWithError(w, err)
		return
	}

	serveStaticImageFile(w, r, cachePath, contentType)
}

func parseImageVariantParams(r *http.Request) (int, int, error) {
	width := 1280
	if rawWidth := strings.TrimSpace(r.URL.Query().Get("w")); rawWidth != "" {
		parsedWidth, err := strconv.Atoi(rawWidth)
		if err != nil || parsedWidth < minImageVariantWidth || parsedWidth > maxImageVariantWidth {
			return 0, 0, invalidQuery("w", rawWidth)
		}
		width = parsedWidth
	}

	quality := defaultImageVariantQuality
	if rawQuality := strings.TrimSpace(r.URL.Query().Get("q")); rawQuality != "" {
		parsedQuality, err := strconv.Atoi(rawQuality)
		if err != nil || parsedQuality < 40 || parsedQuality > 95 {
			return 0, 0, invalidQuery("q", rawQuality)
		}
		quality = parsedQuality
	}

	return width, quality, nil
}

func resolveUploadedImagePath(src string) (string, error) {
	cleanPath := strings.TrimSpace(src)
	if cleanPath == "" || !strings.HasPrefix(cleanPath, "/uploads/") {
		return "", invalidQuery("src", src)
	}

	relativePath := strings.TrimPrefix(path.Clean(cleanPath), "/uploads")
	relativePath = strings.TrimPrefix(relativePath, "/")
	if relativePath == "" || strings.HasPrefix(relativePath, "..") {
		return "", invalidQuery("src", src)
	}

	uploadsRoot, err := filepath.Abs("./uploads")
	if err != nil {
		return "", err
	}

	resolvedPath := filepath.Join(uploadsRoot, relativePath)
	resolvedAbsPath, err := filepath.Abs(resolvedPath)
	if err != nil {
		return "", err
	}

	if resolvedAbsPath != uploadsRoot && !strings.HasPrefix(resolvedAbsPath, uploadsRoot+string(filepath.Separator)) {
		return "", invalidQuery("src", src)
	}

	return resolvedAbsPath, nil
}

func decodeImageConfig(sourcePath string) (image.Config, string, error) {
	file, err := os.Open(sourcePath)
	if err != nil {
		return image.Config{}, "", err
	}
	defer file.Close()

	if strings.EqualFold(filepath.Ext(sourcePath), ".webp") {
		config, err := webp.DecodeConfig(file)
		if err != nil {
			return image.Config{}, "", apperror.New(apperror.CodeUploadInvalidType, "Unsupported image type", map[string]any{
				"path": sourcePath,
			})
		}
		return config, "webp", nil
	}

	config, format, err := image.DecodeConfig(file)
	if err != nil {
		return image.Config{}, "", apperror.New(apperror.CodeUploadInvalidType, "Unsupported image type", map[string]any{
			"path": sourcePath,
		})
	}

	return config, format, nil
}

func imageVariantCachePath(sourcePath string, stat os.FileInfo, width, quality int, format string) (string, string) {
	cacheRoot := filepath.Join(filepath.Dir(sourcePath), imageVariantCacheDirName)
	extension := ".jpg"
	contentType := "image/jpeg"

	switch strings.ToLower(format) {
	case "png", "gif":
		extension = ".png"
		contentType = "image/png"
	}

	hashInput := fmt.Sprintf("%s|%d|%d|%d|%d|%s", sourcePath, stat.Size(), stat.ModTime().UnixNano(), width, quality, format)
	sum := sha1.Sum([]byte(hashInput))
	fileName := hex.EncodeToString(sum[:]) + extension

	return filepath.Join(cacheRoot, fileName), contentType
}

func imageVariantLock(cachePath string) *sync.Mutex {
	lock, _ := imageVariantLocks.LoadOrStore(cachePath, &sync.Mutex{})
	return lock.(*sync.Mutex)
}

func generateImageVariant(sourcePath, cachePath string, width, quality int, format string) error {
	sourceImage, err := decodeSourceImage(sourcePath, format)
	if err != nil {
		return err
	}

	sourceBounds := sourceImage.Bounds()
	if sourceBounds.Dx() <= width {
		return copyFile(cachePath, sourcePath)
	}

	targetHeight := int(math.Round(float64(sourceBounds.Dy()) * (float64(width) / float64(sourceBounds.Dx()))))
	if targetHeight < 1 {
		targetHeight = 1
	}

	resized := image.NewRGBA(image.Rect(0, 0, width, targetHeight))
	xdraw.CatmullRom.Scale(resized, resized.Bounds(), sourceImage, sourceBounds, xdraw.Over, nil)

	if err := os.MkdirAll(filepath.Dir(cachePath), 0o755); err != nil {
		return err
	}

	tempPath := cachePath + ".tmp"
	outputFile, err := os.Create(tempPath)
	if err != nil {
		return err
	}

	encodeErr := encodeResizedImage(outputFile, resized, format, quality)
	closeErr := outputFile.Close()
	if encodeErr != nil {
		_ = os.Remove(tempPath)
		return encodeErr
	}
	if closeErr != nil {
		_ = os.Remove(tempPath)
		return closeErr
	}

	if err := os.Rename(tempPath, cachePath); err != nil {
		_ = os.Remove(tempPath)
		return err
	}

	return nil
}

func decodeSourceImage(sourcePath string, format string) (image.Image, error) {
	file, err := os.Open(sourcePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	switch strings.ToLower(format) {
	case "webp":
		return webp.Decode(file)
	default:
		img, _, err := image.Decode(file)
		return img, err
	}
}

func encodeResizedImage(w io.Writer, img image.Image, format string, quality int) error {
	switch strings.ToLower(format) {
	case "png", "gif":
		encoder := png.Encoder{CompressionLevel: png.BestSpeed}
		return encoder.Encode(w, img)
	default:
		return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
	}
}

func serveStaticImageFile(w http.ResponseWriter, r *http.Request, filePath string, contentType string) {
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.Header().Set("Cache-Control", "public, max-age=2592000")
	http.ServeFile(w, r, filePath)
}

func detectImageContentType(filePath, format string) string {
	switch strings.ToLower(format) {
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	case "webp":
		return "image/webp"
	default:
		switch strings.ToLower(filepath.Ext(filePath)) {
		case ".png":
			return "image/png"
		case ".gif":
			return "image/gif"
		case ".webp":
			return "image/webp"
		default:
			return "image/jpeg"
		}
	}
}

func copyFile(destinationPath, sourcePath string) error {
	if err := os.MkdirAll(filepath.Dir(destinationPath), 0o755); err != nil {
		return err
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	tempPath := destinationPath + ".tmp"
	outputFile, err := os.Create(tempPath)
	if err != nil {
		return err
	}

	_, copyErr := io.Copy(outputFile, sourceFile)
	closeErr := outputFile.Close()
	if copyErr != nil {
		_ = os.Remove(tempPath)
		return copyErr
	}
	if closeErr != nil {
		_ = os.Remove(tempPath)
		return closeErr
	}

	if err := os.Rename(tempPath, destinationPath); err != nil {
		_ = os.Remove(tempPath)
		return err
	}

	return nil
}
