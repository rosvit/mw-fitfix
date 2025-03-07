package fitfix_test

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/profile/filedef"
	"github.com/muktihari/fit/profile/typedef"
	"github.com/stretchr/testify/require"

	. "github.com/rosvit/mw-fitfix/internal/fitfix"
)

const (
	testFile    = "./testdata/test.fit"
	invalidFile = "./testdata/invalid.fit"
)

// noinspection GoUnhandledErrorResult
func Test_Process(t *testing.T) {
	t.Parallel()

	t.Run("success - default/empty options", func(t *testing.T) {
		t.Parallel()
		f := loadFile(t, testFile)
		defer f.Close()
		var w bytes.Buffer
		err := Process(f, &w, Options{})
		require.NoError(t, err)
		a := decodeActivity(t, &w)
		require.NotNil(t, a)
		// unmodified manufacturer, laps
		require.Equal(t, typedef.ManufacturerGarmin, a.FileId.Manufacturer)
		require.Len(t, a.Laps, 2)
		// fixed
		require.Len(t, a.Events, 2)
		require.Less(t, a.Events[0].Timestamp, a.Events[1].Timestamp)
		localTime := a.Activity.Timestamp.Add(1 * time.Hour)
		require.Equal(t, localTime, a.Activity.LocalTimestamp)
		require.Empty(t, a.FieldDescriptions)
		require.Len(t, a.Sessions, 1)
		start := a.Records[0].Timestamp
		end := a.Records[len(a.Records)-1].Timestamp
		sess := a.Sessions[0]
		require.Equal(t, start, sess.StartTime)
		require.Equal(t, end, sess.Timestamp)
	})

	t.Run("success - all options enabled", func(t *testing.T) {
		t.Parallel()
		f := loadFile(t, testFile)
		defer f.Close()
		var w bytes.Buffer
		err := Process(f, &w, Options{Device: true, Laps: true})
		require.NoError(t, err)
		a := decodeActivity(t, &w)
		require.NotNil(t, a)
		require.Equal(t, typedef.ManufacturerMywhoosh, a.FileId.Manufacturer)
		require.Len(t, a.Laps, 1)
		start := a.Records[0].Timestamp
		end := a.Records[len(a.Records)-1].Timestamp
		lap := a.Laps[0]
		require.Equal(t, start, lap.StartTime)
		require.Equal(t, end, lap.Timestamp)
		require.Equal(t, typedef.SportCycling, lap.Sport)
		require.Equal(t, typedef.SubSportVirtualActivity, lap.SubSport)
	})

	t.Run("error - decoding invalid FIT file", func(t *testing.T) {
		t.Parallel()
		f := loadFile(t, invalidFile)
		defer f.Close()
		var w bytes.Buffer
		err := Process(f, &w, Options{})
		require.ErrorContains(t, err, decoder.ErrNotFITFile.Error())
	})
}

func loadFile(t *testing.T, path string) *os.File {
	t.Helper()
	f, err := os.Open(path)
	require.NoError(t, err)
	return f
}

func decodeActivity(t *testing.T, r io.Reader) *filedef.Activity {
	t.Helper()
	dec, err := decoder.New(r).Decode()
	require.NoError(t, err)
	return filedef.NewActivity(dec.Messages...)
}
