package mem

import (
	"testing"
	"time"
)

func TestFileDataNameRace(t *testing.T) {
	t.Parallel()
	const someName = "someName"
	const someOtherName = "someOtherName"
	d := FileData{
		name: someName,
	}

	if d.Name() != someName {
		t.Errorf("Failed to read correct Name, was %v", d.Name())
	}

	ChangeFileName(&d, someOtherName)
	if d.Name() != someOtherName {
		t.Errorf("Failed to set Name, was %v", d.Name())
	}

	go func() {
		ChangeFileName(&d, someName)
	}()

	if d.Name() != someName && d.Name() != someOtherName {
		t.Errorf("Failed to read either Name, was %v", d.Name())
	}
}

func TestFileDataModTimeRace(t *testing.T) {
	t.Parallel()
	someTime := time.Now()
	someOtherTime := someTime.Add(1 * time.Minute)

	d := FileData{
		modtime: someTime,
	}

	s := FileInfo{
		FileData: &d,
	}

	if s.ModTime() != someTime {
		t.Errorf("Failed to read correct value, was %v", s.ModTime())
	}

	SetModTime(&d, someOtherTime)
	if s.ModTime() != someOtherTime {
		t.Errorf("Failed to set ModTime, was %v", s.ModTime())
	}

	go func() {
		SetModTime(&d, someTime)
	}()

	if s.ModTime() != someTime && s.ModTime() != someOtherTime {
		t.Errorf("Failed to read either modtime, was %v", s.ModTime())
	}
}

func TestFileDataModeRace(t *testing.T) {
	t.Parallel()
	const someMode = 0777
	const someOtherMode = 0660

	d := FileData{
		mode: someMode,
	}

	s := FileInfo{
		FileData: &d,
	}

	if s.Mode() != someMode {
		t.Errorf("Failed to read correct value, was %v", s.Mode())
	}

	SetMode(&d, someOtherMode)
	if s.Mode() != someOtherMode {
		t.Errorf("Failed to set Mode, was %v", s.Mode())
	}

	go func() {
		SetMode(&d, someMode)
	}()

	if s.Mode() != someMode && s.Mode() != someOtherMode {
		t.Errorf("Failed to read either mode, was %v", s.Mode())
	}
}

func TestFileDataIsDirRace(t *testing.T) {
	t.Parallel()

	d := FileData{
		dir: true,
	}

	s := FileInfo{
		FileData: &d,
	}

	if s.IsDir() != true {
		t.Errorf("Failed to read correct value, was %v", s.IsDir())
	}

	go func() {
		s.Lock()
		d.dir = false
		s.Unlock()
	}()

	//just logging the value to trigger a read:
	t.Logf("Value is %v", s.IsDir())
}

func TestFileDataSizeRace(t *testing.T) {
	t.Parallel()

	const someData = "Hello"
	const someOtherDataSize = "Hello World"

	d := FileData{
		data: []byte(someData),
		dir:  false,
	}

	s := FileInfo{
		FileData: &d,
	}

	if s.Size() != int64(len(someData)) {
		t.Errorf("Failed to read correct value, was %v", s.Size())
	}

	go func() {
		s.Lock()
		d.data = []byte(someOtherDataSize)
		s.Unlock()
	}()

	//just logging the value to trigger a read:
	t.Logf("Value is %v", s.Size())

	//Testing the Dir size case
	d.dir = true
	if s.Size() != int64(42) {
		t.Errorf("Failed to read correct value for dir, was %v", s.Size())
	}
}
