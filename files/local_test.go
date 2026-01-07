package files

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalStorage_SaveGetDelete(t *testing.T) {
	dir, err := ioutil.TempDir("", "localstoragetest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	ls, err := NewLocalStorage(dir, "")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	content := []byte("hello, world")
	meta := Metadata{Name: "greeting.txt", ContentType: "text/plain", OwnerID: "u1"}

	id, err := ls.Save(ctx, bytes.NewReader(content), meta)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	// verify file exists on disk
	if _, err := os.Stat(filepath.Join(dir, id)); err != nil {
		t.Fatalf("file not exist: %v", err)
	}

	rc, gotMeta, err := ls.Get(ctx, id)
	if err != nil {
		t.Fatalf("get error: %v", err)
	}
	defer rc.Close()

	b, _ := ioutil.ReadAll(rc)
	if string(b) != string(content) {
		t.Fatalf("content mismatch: %s", string(b))
	}
	if gotMeta.Name != "greeting.txt" || gotMeta.OwnerID != "u1" {
		t.Fatalf("meta mismatch: %+v", gotMeta)
	}

	// stat
	m2, err := ls.Stat(ctx, id)
	if err != nil {
		t.Fatalf("stat error: %v", err)
	}
	if m2.ID != id {
		t.Fatalf("stat id mismatch: %s", m2.ID)
	}

	// delete
	if err := ls.Delete(ctx, id); err != nil {
		t.Fatalf("delete error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, id)); !os.IsNotExist(err) {
		t.Fatalf("file still exists after delete")
	}
}
