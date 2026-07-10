package tools

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"turtle/core/lgr"

	"github.com/gin-gonic/gin"
)

func ProxyMiddleware(targetURL string) gin.HandlerFunc {

	return func(c *gin.Context) {
		// Create a new HTTP request using the target URL and original request method

		root := "http://127.0.0.1:5000"

		final_path := root + targetURL

		lgr.Info("Heading for: ", final_path)

		targetRequest, err := http.NewRequest(c.Request.Method, final_path, c.Request.Body)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Copy headers from the original request to the target request
		for key, value := range c.Request.Header {
			targetRequest.Header[key] = value
		}

		// Copy form data from the original request to the target request
		if err := c.Request.ParseMultipartForm(32 << 20); err == nil {
			for key, values := range c.Request.MultipartForm.Value {
				for _, value := range values {
					targetRequest.Form.Add(key, value)
				}
			}
		}

		// Copy cookies from the original request to the target request
		for _, cookie := range c.Request.Cookies() {
			targetRequest.AddCookie(cookie)
		}

		// Send the target request and handle the response
		targetClient := http.DefaultClient
		targetResponse, err := targetClient.Do(targetRequest)
		if err != nil {
			c.AbortWithError(http.StatusBadGateway, err)
			return
		}
		defer targetResponse.Body.Close()

		// Copy headers and status code from the target response to the original response
		for key, value := range targetResponse.Header {
			c.Writer.Header()[key] = value
		}
		c.Writer.WriteHeader(targetResponse.StatusCode)

		// Copy the target response body to the original response body
		_, err = io.Copy(c.Writer, targetResponse.Body)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}

func ProxyMiddleware2() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a new HTTP request using the target URL and original request method
		root := "http://127.0.0.1:5000"
		finalPath := root + c.Request.RequestURI

		lgr.Info("Requesting:", finalPath)

		// Clone the request body if it exists
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}

		// Create a new target request with the cloned body
		targetRequest, err := http.NewRequest(c.Request.Method, finalPath, bytes.NewReader(bodyBytes))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Copy headers from the original request to the target request
		for key, value := range c.Request.Header {
			targetRequest.Header[key] = value
		}

		// Reassign the body to the original request so Gin can read it again
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Parse form data if present
		if err := c.Request.ParseMultipartForm(32 << 20); err == nil {
			targetRequest.MultipartForm = c.Request.MultipartForm
			targetRequest.Form = c.Request.Form
		}

		// Copy cookies from the original request to the target request
		for _, cookie := range c.Request.Cookies() {
			targetRequest.AddCookie(cookie)
		}

		// Send the target request and handle the response
		targetClient := http.DefaultClient
		targetResponse, err := targetClient.Do(targetRequest)
		if err != nil {
			c.AbortWithError(http.StatusBadGateway, err)
			return
		}
		defer targetResponse.Body.Close()

		// Copy headers and status code from the target response to the original response
		for key, value := range targetResponse.Header {
			c.Writer.Header()[key] = value
		}
		c.Writer.WriteHeader(targetResponse.StatusCode)

		// Copy the target response body to the original response body
		_, err = io.Copy(c.Writer, targetResponse.Body)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}

func ProxyUseBoth(fn func(*gin.Context)) gin.HandlerFunc {
	proxed := ProxyWrapper()
	return func(c *gin.Context) {
		fn(c)
		proxed(c)
	}
}

func ProxySwapWrapper(fn func(*gin.Context)) gin.HandlerFunc {
	if false {
		return ProxyWrapper()
	} else {
		return fn
	}
}

func PreventProxy(fn func(func(*gin.Context), *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(ProxyWrapper(), c)
	}
}

func ProxyWrapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		root := "http://127.0.0.1:5000"
		finalPath := root + c.Request.RequestURI

		// Create a new HTTP request using the target URL and original request method
		targetRequest, err := http.NewRequest(c.Request.Method, finalPath, c.Request.Body)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Copy headers from the original request to the target request
		for key, value := range c.Request.Header {
			targetRequest.Header[key] = value
		}

		// Copy cookies from the original request to the target request
		for _, cookie := range c.Request.Cookies() {
			targetRequest.AddCookie(cookie)
		}

		// Create a custom HTTP client to prevent automatic redirects
		targetClient := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Prevent following redirects automatically
				// Return http.ErrUseLastResponse to return the response to the caller
				return http.ErrUseLastResponse
			},
		}

		// Send the target request and handle the response
		targetResponse, err := targetClient.Do(targetRequest)
		if err != nil {
			c.AbortWithError(http.StatusBadGateway, err)
			return
		}
		defer targetResponse.Body.Close()

		// If the response is a redirect, handle it
		if targetResponse.StatusCode >= 300 && targetResponse.StatusCode < 400 {
			// Copy redirect location header
			location := targetResponse.Header.Get("Location")
			if location != "" {
				// Optionally, you can modify the location header here
				c.Redirect(targetResponse.StatusCode, location)
				return
			}
		}

		// Copy headers and status code from the target response to the original response
		for key, value := range targetResponse.Header {
			c.Writer.Header()[key] = value
		}
		c.Writer.WriteHeader(targetResponse.StatusCode)

		// Copy the target response body to the original response body
		_, err = io.Copy(c.Writer, targetResponse.Body)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}

func cloneMultipartForm(originalForm *multipart.Form) (*multipart.Form, error) {
	clonedForm := &multipart.Form{
		Value: make(map[string][]string),
		File:  make(map[string][]*multipart.FileHeader),
	}

	// Clone regular form values
	for key, values := range originalForm.Value {
		clonedForm.Value[key] = append([]string{}, values...) // Create a new slice and copy values
	}

	for key, fileHeaders := range originalForm.File {
		clonedHeaders := make([]*multipart.FileHeader, len(fileHeaders))
		for i, fileHeader := range fileHeaders {
			// Creating a new FileHeader (shallow copy)
			clonedHeaders[i] = &multipart.FileHeader{
				Filename: fileHeader.Filename,
				Header:   cloneMIMEHeader(fileHeader.Header),
				Size:     fileHeader.Size,
			}
		}
		clonedForm.File[key] = clonedHeaders
	}

	return clonedForm, nil
}

// This function clones a textproto.MIMEHeader, which is used in multipart.FileHeader.
func cloneMIMEHeader(header textproto.MIMEHeader) textproto.MIMEHeader {
	clonedHeader := make(textproto.MIMEHeader)
	for key, values := range header {
		clonedHeader[key] = append([]string{}, values...)
	}
	return clonedHeader
}
