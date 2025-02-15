package makemkv_test

import (
	"path/filepath"
	"testing"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/curt-hash/mkvbot/pkg/makemkv"
)

func TestParseDriveScan(t *testing.T) {
	for _, tc := range []struct {
		line     string
		expected *makemkv.DriveScan
	}{
		{
			`DRV:0,12,345,6789,"BD Brand Model SN123","A Disc","/dev/sr0"`,
			&makemkv.DriveScan{
				Pos: lexer.Position{
					Filename: "",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
				Index:      0,
				Visible:    12,
				Enabled:    345,
				Flags:      6789,
				DriveName:  "BD Brand Model SN123",
				DiscTitle:  "A Disc",
				VolumeName: "/dev/sr0",
			},
		},
		{
			`DRV:1,256,999,0,"","",""`,
			&makemkv.DriveScan{
				Pos: lexer.Position{
					Filename: "",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
				Index:      1,
				Visible:    256,
				Enabled:    999,
				Flags:      0,
				DriveName:  "",
				DiscTitle:  "",
				VolumeName: "",
			},
		},
	} {
		t.Run(tc.line, func(t *testing.T) {
			ol, err := makemkv.ParseLine(tc.line)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, ol.DriveScan)
		})
	}
}

func TestParseMessage(t *testing.T) {
	for _, tc := range []struct {
		line     string
		expected *makemkv.Message
	}{
		{
			`MSG:1234,56,1,"Foo ""bar"" (baz).","%1","Foo ""bar"" (baz)."`,
			&makemkv.Message{
				Pos: lexer.Position{
					Filename: "",
					Offset:   4,
					Line:     1,
					Column:   5,
				},
				Code:      1234,
				Flags:     56,
				NumParams: 1,
				Message:   `Foo "bar" (baz).`,
				Format:    "%1",
				Params:    []makemkv.Str{`Foo "bar" (baz).`},
			},
		},
	} {
		t.Run(tc.line, func(t *testing.T) {
			ol, err := makemkv.ParseLine(tc.line)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, ol.Message)
		})
	}
}

func TestParseDiscInfo(t *testing.T) {
	for _, tc := range []struct {
		line     string
		expected *makemkv.DiscInfo
	}{
		{
			`CINFO:31,6119,"<b>Source information</b><br>"`,
			&makemkv.DiscInfo{
				Pos: lexer.Position{
					Filename: "",
					Offset:   6,
					Line:     1,
					Column:   7,
				},
				Attribute: &makemkv.Attribute{
					Pos: lexer.Position{
						Filename: "",
						Offset:   6,
						Line:     1,
						Column:   7,
					},
					ID:    31,
					Code:  6119,
					Value: "<b>Source information</b><br>",
				},
			},
		},
	} {
		t.Run(tc.line, func(t *testing.T) {
			ol, err := makemkv.ParseLine(tc.line)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, ol.DiscInfo)
		})
	}
}

func TestParseTitleInfo(t *testing.T) {
	for _, tc := range []struct {
		line     string
		expected *makemkv.TitleInfo
	}{
		{
			`TINFO:0,12,345,"15 chapter(s) , 3.6 GB"`,
			&makemkv.TitleInfo{
				Pos: lexer.Position{
					Filename: "",
					Offset:   6,
					Line:     1,
					Column:   7,
				},
				TitleIndex: 0,
				Attribute: &makemkv.Attribute{
					Pos: lexer.Position{
						Filename: "",
						Offset:   8,
						Line:     1,
						Column:   9,
					},
					ID:    12,
					Code:  345,
					Value: "15 chapter(s) , 3.6 GB",
				},
			},
		},
	} {
		t.Run(tc.line, func(t *testing.T) {
			ol, err := makemkv.ParseLine(tc.line)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, ol.TitleInfo)
		})
	}
}

func TestParseStreamInfo(t *testing.T) {
	for _, tc := range []struct {
		line     string
		expected *makemkv.StreamInfo
	}{
		{
			`SINFO:0,1,30,0,"DD 3/2+1 English"`,
			&makemkv.StreamInfo{
				Pos: lexer.Position{
					Filename: "",
					Offset:   6,
					Line:     1,
					Column:   7,
				},
				TitleIndex:  0,
				StreamIndex: 1,
				Attribute: &makemkv.Attribute{
					Pos: lexer.Position{
						Filename: "",
						Offset:   10,
						Line:     1,
						Column:   11,
					},
					ID:    30,
					Code:  0,
					Value: "DD 3/2+1 English",
				},
			},
		},
	} {
		t.Run(tc.line, func(t *testing.T) {
			ol, err := makemkv.ParseLine(tc.line)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, ol.StreamInfo)
		})
	}
}

func TestParseCurrentTask(t *testing.T) {
	for _, tc := range []struct {
		line     string
		expected *makemkv.CurrentTask
	}{
		{
			`PRGT:5018,0,"Scanning CD-ROM devices"`,
			&makemkv.CurrentTask{
				Pos: lexer.Position{
					Filename: "",
					Offset:   5,
					Line:     1,
					Column:   6,
				},
				Task: &makemkv.Task{
					Pos: lexer.Position{
						Filename: "",
						Offset:   5,
						Line:     1,
						Column:   6,
					},
					ID:   5018,
					Code: 0,
					Name: "Scanning CD-ROM devices",
				},
			},
		},
	} {
		t.Run(tc.line, func(t *testing.T) {
			ol, err := makemkv.ParseLine(tc.line)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, ol.CurrentTask)
		})
	}
}

func TestParseCurrentSubtask(t *testing.T) {
	for _, tc := range []struct {
		line     string
		expected *makemkv.CurrentSubtask
	}{
		{
			`PRGC:5018,0,"Scanning CD-ROM devices"`,
			&makemkv.CurrentSubtask{
				Pos: lexer.Position{
					Filename: "",
					Offset:   5,
					Line:     1,
					Column:   6,
				},
				Task: &makemkv.Task{
					Pos: lexer.Position{
						Filename: "",
						Offset:   5,
						Line:     1,
						Column:   6,
					},
					ID:   5018,
					Code: 0,
					Name: "Scanning CD-ROM devices",
				},
			},
		},
	} {
		t.Run(tc.line, func(t *testing.T) {
			ol, err := makemkv.ParseLine(tc.line)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, ol.CurrentSubtask)
		})
	}
}

func TestParseProgress(t *testing.T) {
	for _, tc := range []struct {
		line     string
		expected *makemkv.Progress
	}{
		{
			`PRGV:3744,313,65536`,
			&makemkv.Progress{
				Pos: lexer.Position{
					Filename: "",
					Offset:   5,
					Line:     1,
					Column:   6,
				},
				SubtaskValue: 3744,
				TaskValue:    313,
				Max:          65536,
			},
		},
	} {
		t.Run(tc.line, func(t *testing.T) {
			ol, err := makemkv.ParseLine(tc.line)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, ol.Progress)
		})
	}
}

func TestParseTitleCount(t *testing.T) {
	for _, tc := range []struct {
		line     string
		expected *makemkv.TitleCount
	}{
		{
			`TCOUNT:123`,
			&makemkv.TitleCount{
				Pos: lexer.Position{
					Filename: "",
					Offset:   7,
					Line:     1,
					Column:   8,
				},
				Count: 123,
			},
		},
	} {
		t.Run(tc.line, func(t *testing.T) {
			ol, err := makemkv.ParseLine(tc.line)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, ol.TitleCount)
		})
	}
}

func TestParseFile(t *testing.T) {
	matches, err := filepath.Glob("testdata/logs/*.txt")
	require.NoError(t, err)
	require.NotEmpty(t, matches)
	for _, match := range matches {
		t.Run(match, func(t *testing.T) {
			output, err := makemkv.ParseFile(match)
			assert.NoError(t, err)
			assert.NotEmpty(t, output.Lines)
		})
	}
}
