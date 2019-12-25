package util_test

import (
	"fmt"
	"testing"

	"github.com/lughong/gin-api-demo/app/util"
)

var (
	path = "testdata"
	file = "test.txt"
)

func TestPathExists(t *testing.T) {
	exists, err := util.PathExists(path)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Logf("%s was not exists.", path)
	} else {
		t.Logf("%s was exists.", path)
	}
}

func TestCreateDir(t *testing.T) {
	ok, err := util.CreateDir(path)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Logf("Create %s was failed.", path)
	} else {
		t.Logf("Create %s success.", path)
	}
}

func TestCreateFile(t *testing.T) {
	fileName := fmt.Sprintf("%s/%s", path, file)
	ok, err := util.CreateFile(fileName)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Logf("Create %s was failed.", fileName)
	} else {
		t.Logf("Create %s success.", fileName)
	}
}

func BenchmarkCreateFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fileName := fmt.Sprintf("%s/%s%d", path, file, i)
		_, _ = util.CreateFile(fileName)
	}
}

func BenchmarkCreateFileConsuming(b *testing.B) {
	b.StopTimer()

	fileName := fmt.Sprintf("%s/%s", path, file)
	_, err := util.CreateFile(fileName)
	if err != nil {
		b.Fatal(err)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fileName := fmt.Sprintf("%s/%s%d", path, file, i)
		_, _ = util.CreateFile(fileName)
	}
}
