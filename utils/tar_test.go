package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestTarGzDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Containerfile"), []byte("FROM scratch\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "greenboot"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "greenboot", "check.sh"), []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".git", "HEAD"), []byte("ref"), 0o644); err != nil {
		t.Fatal(err)
	}

	packed, err := TarGzDir(dir)
	if err != nil {
		t.Fatalf("TarGzDir: %v", err)
	}
	gz, err := gzip.NewReader(bytes.NewReader(packed))
	if err != nil {
		t.Fatal(err)
	}
	tr := tar.NewReader(gz)
	found := map[string]int64{}
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		found[hdr.Name] = hdr.Mode
	}
	if _, ok := found["Containerfile"]; !ok {
		t.Fatalf("Containerfile missing from archive: %v", found)
	}
	if mode, ok := found["greenboot/check.sh"]; !ok || mode&0o111 == 0 {
		t.Fatalf("greenboot/check.sh missing or not executable: %v", found)
	}
	for name := range found {
		if len(name) >= 4 && name[:4] == ".git" {
			t.Fatalf(".git should be excluded, found %q", name)
		}
	}
}
