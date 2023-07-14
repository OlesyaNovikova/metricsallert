package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// gzipsWriter реализует интерфейс http.ResponseWriter
type gzipWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newGzipWriter(w http.ResponseWriter) *gzipWriter {
	return &gzipWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *gzipWriter) Header() http.Header {
	return c.w.Header()
}

func (c *gzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		contentType := c.Header().Get("Content-Type")
		if headerCheck(contentType, "application/json") || headerCheck(contentType, "text/html") {
			c.w.Header().Set("Content-Encoding", "gzip")
		}
	}
	c.w.WriteHeader(statusCode)
}

func (c *gzipWriter) Write(p []byte) (int, error) {
	if headerCheck(c.Header().Get("Content-Encoding"), "gzip") {
		return c.zw.Write(p)
	}
	return c.w.Write(p)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *gzipWriter) Close() error {
	return c.zw.Close()
}

// gzipReader реализует интерфейс io.ReadCloser
type gzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newGzipReader(r io.ReadCloser) (*gzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c gzipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *gzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func WithGzip(h http.HandlerFunc) http.HandlerFunc {
	compFn := func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		if headerCheck(r.Header.Get("Content-Encoding"), "gzip") {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}
		// по умолчанию устанавливаем оригинальный http.ResponseWriter
		ow := w
		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		if headerCheck(r.Header.Get("Accept-Encoding"), "gzip") {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := newGzipWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			defer cw.Close()
		}

		// передаём управление хендлеру
		h.ServeHTTP(ow, r)

	}
	return http.HandlerFunc(compFn)
}

func headerCheck(str, par string) bool {
	options := strings.Split(str, ",")
	for _, option := range options {
		option = strings.TrimSpace(option)
		if option == par {
			return true
		}
	}
	return false
}
