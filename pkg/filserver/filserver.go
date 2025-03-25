package fileserver

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// FileServer represents a file server instance
type FileServer struct {
	rootDir string
	host    string
	port    string
}

// New creates a new FileServer instance
func New(rootDir string, host string, opts ...FileServerOpt) *FileServer {
	f := &FileServer{
		rootDir: rootDir,
		host:    host,
		port:    ":8000", // default port
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

// Run starts the file server
func (fs *FileServer) Run() error {
	if _, err := os.Stat(fs.rootDir); os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist", fs.rootDir)
	}

	http.HandleFunc("/", fs.fileHandler)

	address := fs.host + fs.port
	fmt.Printf("Serving %s on http://%s\n", fs.rootDir, address)
	return http.ListenAndServe(fs.port, nil)
}

// fileHandler handles file requests
func (fs *FileServer) fileHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(fs.rootDir, r.URL.Path)
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	if fileInfo.IsDir() {
		if r.URL.Query().Get("download") == "true" {
			fs.compressAndDownloadDir(w, path)
		} else {
			fs.dirList(w, path)
		}
	} else {
		http.ServeFile(w, r, path)
	}
}

// dirList displays the contents of a directory
func (fs *FileServer) dirList(w http.ResponseWriter, dirPath string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	relativePath, err := filepath.Rel(fs.rootDir, dirPath)
	if err != nil {
		http.Error(w, "Error creating relative path", http.StatusInternalServerError)
		return
	}
	if relativePath == "." {
		relativePath = ""
	}

	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; margin: 20px; }
			h1 { color: #333; }
			.nav-links { margin: 10px 0; }
			.nav-links a { margin-right: 10px; }
			.file-list { list-style: none; padding: 0; }
			.file-item { display: flex; align-items: center; padding: 5px 0; }
			.file-link { text-decoration: none; color: #0066cc; margin-right: 15px; }
			.download-btn { 
				padding: 3px 10px;
				background-color: #4CAF50;
				color: white;
				border: none;
				border-radius: 3px;
				text-decoration: none;
				font-size: 0.9em;
			}
			.download-btn:hover { background-color: #45a049; }
		</style>
	</head>
	<body>
	`)

	fmt.Fprintf(w, "<h1>Directory: /%s</h1>", relativePath)

	fmt.Fprintf(w, "<div class=\"nav-links\">")
	if dirPath != fs.rootDir {
		parentPath := "/" + filepath.ToSlash(filepath.Dir(relativePath))
		if parentPath == "/"+"." {
			parentPath = "/"
		}
		fmt.Fprintf(w, `<a href="%s">Back</a>`, parentPath)
	}
	fmt.Fprintf(w, `<a href="/">Root</a></div>`)

	fmt.Fprintf(w, "<ul class=\"file-list\">")
	for _, file := range files {
		name := file.Name()
		isDir := file.IsDir()
		if isDir {
			name += "/"
		}
		relativePath, _ := filepath.Rel(fs.rootDir, filepath.Join(dirPath, name))
		slashedPath := filepath.ToSlash(relativePath)

		fmt.Fprintf(w, "<li class=\"file-item\">")
		fmt.Fprintf(w, `<a class="file-link" href="/%s">%s</a>`, slashedPath, name)
		downloadParam := ""
		if isDir {
			downloadParam = "?download=true"
		}
		fmt.Fprintf(w, `<a class="download-btn" href="/%s%s" download>Download</a>`, slashedPath, downloadParam)
		fmt.Fprintf(w, "</li>")
	}
	fmt.Fprintf(w, "</ul></body></html>")
}

// compressAndDownloadDir compresses a directory and sends it as a zip file
func (fs *FileServer) compressAndDownloadDir(w http.ResponseWriter, dirPath string) {
	relativePath, err := filepath.Rel(fs.rootDir, dirPath)
	if err != nil {
		http.Error(w, "Error creating relative path", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", relativePath))

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if it's the root directory itself
		if path == dirPath {
			return nil
		}

		// Create a relative path for the zip file
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		// Create zip header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		http.Error(w, "Error creating zip file", http.StatusInternalServerError)
		return
	}
}
