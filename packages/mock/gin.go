package mock

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/textproto"
)

func GinMockSet(key string, value any) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set(key, value)
		context.Next()
	}
}

func CreateMockFileHeader(filename string, fileContent []byte) (*multipart.FileHeader, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	defer w.Close()

	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	if _, err := part.Write(fileContent); err != nil {
		return nil, err
	}

	// Now we have the complete form, we can parse it
	r := multipart.NewReader(&b, w.Boundary())
	p, err := r.NextPart()
	if err != nil {
		return nil, err
	}

	// Create a mock FileHeader from the part
	fileHeader := &multipart.FileHeader{
		Filename: filename,
		Header:   textproto.MIMEHeader{},
		Size:     int64(len(fileContent)),
	}

	fileHeader.Header.Set("Content-Disposition", p.Header.Get("Content-Disposition"))
	fileHeader.Header.Set("Content-Type", p.Header.Get("Content-Type"))

	return fileHeader, nil
}
