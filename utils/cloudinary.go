package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type CloudinaryResponse struct {
	SecureURL string `json:"secure_url"`
	PublicID  string `json:"public_id"`
}

const (
	CloudinaryAPIKey    = "239243922124375"
	CloudinaryAPISecret = "yj3z35qj8JxNBFo84sBA6yyDl_5o"
	CloudinaryCloudName = "dxwnbbcpo"
)

// UploadFormFileToCloudinary uploads a Gin uploaded file to Cloudinary using unsigned preset
func UploadFormFileToCloudinary(c *gin.Context, fileHeader *multipart.FileHeader) (string, error) {
	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer file.Close()

	// Prepare multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Attach file
	part, err := writer.CreateFormFile("file", filepath.Base(fileHeader.Filename))
	if err != nil {
		log.Println(err)
		return "", err
	}
	if _, err = io.Copy(part, file); err != nil {
		log.Println(err)
		return "", err
	}

	// Add unsigned preset field
	_ = writer.WriteField("upload_preset", "practise") // 👈 your unsigned preset

	writer.Close()

	// Send request to Cloudinary
	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", CloudinaryCloudName),
		body)
	if err != nil {
		log.Println(err)
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Println(err)
		return "", fmt.Errorf("cloudinary upload failed: %s", string(respBody))
	}

	// Parse secure_url from Cloudinary response
	type uploadResp struct {
		SecureURL string `json:"secure_url"`
	}
	var result uploadResp
	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Println(err)
		return "", err
	}

	log.Println("Cloudinary upload successful:", result.SecureURL)
	return result.SecureURL, nil
}
