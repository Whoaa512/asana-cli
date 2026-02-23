package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/whoaa512/asana-cli/internal/config"
)

func TestListAttachments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/123/attachments" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if got := r.URL.Query().Get("limit"); got != "10" {
			t.Errorf("limit = %q, want %q", got, "10")
		}
		if got := r.URL.Query().Get("offset"); got != "abc" {
			t.Errorf("offset = %q, want %q", got, "abc")
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"gid":           "a1",
					"name":          "spec.pdf",
					"resource_type": "attachment",
				},
			},
			"next_page": map[string]any{"offset": "next"},
		})
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	client := newTestHTTPClient(server.URL)
	result, err := client.ListAttachments(context.Background(), "123", 10, "abc")
	if err != nil {
		t.Fatalf("ListAttachments() error = %v", err)
	}

	if len(result.Data) != 1 {
		t.Fatalf("len(result.Data) = %d, want 1", len(result.Data))
	}
	if result.Data[0].GID != "a1" {
		t.Errorf("result.Data[0].GID = %q, want %q", result.Data[0].GID, "a1")
	}
	if result.NextPage == nil || result.NextPage.Offset != "next" {
		t.Errorf("result.NextPage.Offset = %v, want %q", result.NextPage, "next")
	}
}

func TestGetAttachment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/attachments/a1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"gid":           "a1",
				"name":          "spec.pdf",
				"resource_type": "attachment",
				"size":          123,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	client := newTestHTTPClient(server.URL)
	attachment, err := client.GetAttachment(context.Background(), "a1")
	if err != nil {
		t.Fatalf("GetAttachment() error = %v", err)
	}

	if attachment.GID != "a1" {
		t.Errorf("attachment.GID = %q, want %q", attachment.GID, "a1")
	}
	if attachment.Size != 123 {
		t.Errorf("attachment.Size = %d, want %d", attachment.Size, 123)
	}
}

func TestUploadAttachment(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(filePath, []byte("hello attachment"), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tasks/123/attachments" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("unexpected auth header: %s", auth)
		}
		if contentType := r.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "multipart/form-data; boundary=") {
			t.Errorf("unexpected content type: %s", contentType)
		}

		if err := r.ParseMultipartForm(1024 * 1024); err != nil {
			t.Fatalf("ParseMultipartForm() error = %v", err)
		}

		if got := r.FormValue("name"); got != "override.txt" {
			t.Errorf("name = %q, want %q", got, "override.txt")
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("FormFile() error = %v", err)
		}
		defer func() { _ = file.Close() }()

		if header.Filename != "override.txt" {
			t.Errorf("uploaded filename = %q, want %q", header.Filename, "override.txt")
		}

		content, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("ReadAll(file) error = %v", err)
		}
		if string(content) != "hello attachment" {
			t.Errorf("uploaded content = %q, want %q", string(content), "hello attachment")
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"gid":           "a1",
				"name":          "override.txt",
				"resource_type": "attachment",
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	client := newTestHTTPClient(server.URL)
	attachment, err := client.UploadAttachment(context.Background(), "123", filePath, "override.txt")
	if err != nil {
		t.Fatalf("UploadAttachment() error = %v", err)
	}

	if attachment.GID != "a1" {
		t.Errorf("attachment.GID = %q, want %q", attachment.GID, "a1")
	}
}

func TestUploadAttachmentUsesBasenameWhenNameMissing(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "report.csv")
	if err := os.WriteFile(filePath, []byte("id,name\n1,test"), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1024 * 1024); err != nil {
			t.Fatalf("ParseMultipartForm() error = %v", err)
		}

		if got := r.FormValue("name"); got != "report.csv" {
			t.Errorf("name = %q, want %q", got, "report.csv")
		}

		_, header, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("FormFile() error = %v", err)
		}
		if header.Filename != "report.csv" {
			t.Errorf("uploaded filename = %q, want %q", header.Filename, "report.csv")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"gid":"a2","name":"report.csv","resource_type":"attachment"}}`))
	}))
	defer server.Close()

	client := newTestHTTPClient(server.URL)
	attachment, err := client.UploadAttachment(context.Background(), "123", filePath, "")
	if err != nil {
		t.Fatalf("UploadAttachment() error = %v", err)
	}

	if attachment.Name != "report.csv" {
		t.Errorf("attachment.Name = %q, want %q", attachment.Name, "report.csv")
	}
}

func TestUploadAttachmentMissingFile(t *testing.T) {
	var requestCount int32
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&requestCount, 1)
	}))
	defer server.Close()

	client := newTestHTTPClient(server.URL)
	_, err := client.UploadAttachment(context.Background(), "123", filepath.Join(t.TempDir(), "missing.txt"), "")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if atomic.LoadInt32(&requestCount) != 0 {
		t.Fatalf("expected no requests, got %d", requestCount)
	}
}

func TestDeleteAttachment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/attachments/a1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{}}`))
	}))
	defer server.Close()

	client := newTestHTTPClient(server.URL)
	if err := client.DeleteAttachment(context.Background(), "a1"); err != nil {
		t.Fatalf("DeleteAttachment() error = %v", err)
	}
}

func TestUploadAttachmentRetriesOnRateLimit(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "retry.txt")
	if err := os.WriteFile(filePath, []byte("retry content"), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"errors":[{"message":"rate limited"}]}`))
			return
		}

		if r.Method != http.MethodPost {
			t.Errorf("unexpected method on retry: %s", r.Method)
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Errorf("unexpected content type on retry: %s", r.Header.Get("Content-Type"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"gid":"a1","name":"retry.txt","resource_type":"attachment"}}`))
	}))
	defer server.Close()

	client := newTestHTTPClient(server.URL)
	attachment, err := client.UploadAttachment(context.Background(), "123", filePath, "")
	if err != nil {
		t.Fatalf("UploadAttachment() error = %v", err)
	}

	if got := atomic.LoadInt32(&attempts); got != 2 {
		t.Errorf("attempts = %d, want 2", got)
	}
	if attachment.GID != "a1" {
		t.Errorf("attachment.GID = %q, want %q", attachment.GID, "a1")
	}
}

func newTestHTTPClient(baseURL string) *HTTPClient {
	cfg := &config.Config{
		AccessToken: "test-token",
		Timeout:     5 * time.Second,
	}
	return NewHTTPClient(cfg, WithBaseURL(baseURL))
}
