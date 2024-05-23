package hub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	engine  *gin.Engine
	rootDir string
}

func NewServer() *Server {
	engine := gin.Default()
	engine.MaxMultipartMemory = 8 << 20

	engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	})

	server := &Server{engine: engine}

	engine.POST("/blobs/uploads", server.handleUploadBlob)
	engine.GET("/blobs/:digest", server.handleDownloadBlob)
	engine.POST("/manifest/:name/:tag/uploads", server.handleUploadManifest)
	engine.GET("/manifest/:name/:tag", server.handleDownloadManifest)
	engine.GET("/manifest", server.handleListPlugins)
	engine.GET("/manifest/:name", server.handleListTags)

	return server
}

func (s *Server) Run(port int, rootDir string) error {
	s.rootDir = rootDir
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.engine,
	}

	go func() {
		logrus.Infof("start server on port %d", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("failed to listen and serve: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	ch := <-sig
	logrus.Infof("server receive signal: %s", ch.String())
	return server.Shutdown(ctx)
}

func (s *Server) getBlobFilePath(digest *Digest) string {
	_ = os.MkdirAll(fmt.Sprintf("%s/blob/%s", s.rootDir, digest.Prefix()), os.ModePerm)
	return fmt.Sprintf("%s/blob/%s/%s", s.rootDir, digest.Prefix(), digest.Value())
}

func (s *Server) getManifestFilePath(name, tag string) string {
	_ = os.MkdirAll(fmt.Sprintf("%s/manifest/%s", s.rootDir, name), os.ModePerm)
	return fmt.Sprintf("%s/manifest/%s/%s", s.rootDir, name, tag)
}

func (s *Server) existsAndIsFile(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	if f.IsDir() {
		return false
	}
	return true
}

func (s *Server) handleListTags(ctx *gin.Context) {
	tags := make([]string, 0)

	file, err := os.Open(fmt.Sprintf("%s/manifest/%s", s.rootDir, ctx.Param("name")))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	defer file.Close()

	entries, err := file.Readdir(-1)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			tags = append(tags, entry.Name())
		}
	}
	ctx.JSON(http.StatusOK, tags)
}

func (s *Server) handleListPlugins(ctx *gin.Context) {
	plugins := make([]string, 0)
	file, err := os.Open(fmt.Sprintf("%s/manifest", s.rootDir))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	defer file.Close()

	entries, err := file.Readdir(-1)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			plugins = append(plugins, entry.Name())
		}
	}

	ctx.JSON(http.StatusOK, plugins)
}

func (s *Server) handleDownloadManifest(ctx *gin.Context) {
	name := ctx.Param("name")
	tag := ctx.Param("tag")

	filePath := s.getManifestFilePath(name, tag)
	logrus.Infof("manifest file path: %s", filePath)
	f, err := os.ReadFile(filePath)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	mf := new(Manifest)
	err = json.Unmarshal(f, mf)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if s.existsAndIsFile(filePath) {
		ctx.JSON(http.StatusOK, mf)
		return
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "manifest not exist"})
	}
}

func (s *Server) handleUploadManifest(ctx *gin.Context) {
	manifest := new(Manifest)
	err := ctx.ShouldBindJSON(manifest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if manifest.MediaType != ManifestListMediaType {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "not supported media type"})
		return
	}

	for _, m := range manifest.Manifests {
		if m.MediaType != ManifestMediaType {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "not supported media type"})
			return
		}

		digest, err := ParseDigest(m.Digest)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		filePath := s.getBlobFilePath(digest)
		if !s.existsAndIsFile(filePath) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "blob not exist"})
			return
		}
	}

	name := ctx.Param("name")
	tag := ctx.Param("tag")
	manifestJson, err := json.Marshal(manifest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = os.WriteFile(s.getManifestFilePath(name, tag), manifestJson, os.ModePerm)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (s *Server) handleDownloadBlob(ctx *gin.Context) {
	digest, err := ParseDigest(ctx.Param("digest"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	filePath := s.getBlobFilePath(digest)
	if s.existsAndIsFile(filePath) {
		ctx.File(filePath)
		return
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "blob not exist"})
	}
}

func (s *Server) handleUploadBlob(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	dst := fmt.Sprintf("%s/tmp/%s-%d", s.rootDir, file.Filename, time.Now().UnixMicro())
	err = ctx.SaveUploadedFile(file, dst)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	digest, err := GetDigest(dst)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = os.Rename(dst, s.getBlobFilePath(digest))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": digest.String()})
}
