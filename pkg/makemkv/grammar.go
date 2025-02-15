package makemkv

import (
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

//nolint:govet
type (
	// Output represents the multi-line output of makemkvcon.
	//
	// See https://makemkv.com/developers/usage.txt.
	Output struct {
		Pos lexer.Position

		Lines []*Line `( @@ EOL )*`
	}

	// Line represents a single line of output from makemkvcon.
	Line struct {
		Pos lexer.Position

		DriveScan      *DriveScan      `  "DRV"    ":" @@`
		Message        *Message        `| "MSG"    ":" @@`
		DiscInfo       *DiscInfo       `| "CINFO"  ":" @@`
		TitleInfo      *TitleInfo      `| "TINFO"  ":" @@`
		StreamInfo     *StreamInfo     `| "SINFO"  ":" @@`
		CurrentTask    *CurrentTask    `| "PRGT"   ":" @@`
		CurrentSubtask *CurrentSubtask `| "PRGC"   ":" @@`
		Progress       *Progress       `| "PRGV"   ":" @@`
		TitleCount     *TitleCount     `| "TCOUNT" ":" @@`
	}

	// DriveScan represents a makemkvcon "DRV" output line, which describes a
	// disc drive.
	DriveScan struct {
		Pos lexer.Position

		Index      int `@Int`
		Visible    int `"," @Int`
		Enabled    int `"," @Int`
		Flags      int `"," @Int`
		DriveName  Str `"," @String`
		DiscTitle  Str `"," @String`
		VolumeName Str `"," @String`
	}

	// Message represents a makemkvcon "MSG" output line, which is an
	// informational logging line.
	Message struct {
		Pos lexer.Position

		Code      int   `@Int`
		Flags     int   `"," @Int`
		NumParams int   `"," @Int`
		Message   Str   `"," @String`
		Format    Str   `"," @String`
		Params    []Str `( "," @String )*`
	}

	// DiscInfoLine represents a makemkvcon "CINFO" output line, which provides
	// information about a disc.
	DiscInfo struct {
		Pos lexer.Position

		Attribute *Attribute `@@`
	}

	// TitleInfo represents a makemkvcon "TINFO" output line, which provides
	// information about a title.
	TitleInfo struct {
		Pos lexer.Position

		TitleIndex int        `@Int`
		Attribute  *Attribute `"," @@`
	}

	// StreamInfo represents an "SINFO" makemkvcon output line, which provides
	// information about a stream.
	StreamInfo struct {
		Pos lexer.Position

		TitleIndex  int        `@Int`
		StreamIndex int        `"," @Int`
		Attribute   *Attribute `"," @@`
	}

	// Attribute is the common representation of the "CINFO", "TINFO" and "SINFO"
	// makemkvcon output lines, which describe an attribute of a disc, title, or
	// stream.
	Attribute struct {
		Pos lexer.Position

		// ID is an integer that identifies the attribute.
		ID int `@Int`

		// Code is an integer that corresponds to Value, if Value is an enumeration.
		Code int `"," @Int`

		// Value is the value of the attribute identified by ID.
		Value Str `"," @String`
	}

	// CurrentTask represents a makemkvcon "PRGT" output line, which describes
	// the overall task being performed.
	CurrentTask struct {
		Pos lexer.Position

		Task *Task `@@`
	}

	// CurrentSubtask represents a makemkvcon "PRGC" output line, which describes
	// the current sub-task.
	CurrentSubtask struct {
		Pos lexer.Position

		Task *Task `@@`
	}

	// Task is the common representation for makemkvcon "PRGT" and "PRGC" output
	// lines, which describe the current task and subtask, respectively.
	Task struct {
		Pos lexer.Position

		ID   int `@Int`
		Code int `"," @Int`
		Name Str `"," @String`
	}

	// Progress represents a makemkvcon "PRGV" output line, which describes the
	// progress of a task and sub-task.
	Progress struct {
		Pos lexer.Position

		SubtaskValue int `@Int`
		TaskValue    int `"," @Int`

		// Max is a constant denominator used to calculate the progress percentage.
		Max int `"," @Int`
	}

	// TitleCount represents a "TCOUNT" makemkvcon output line, which describes
	// the number of titles found on a disc.
	TitleCount struct {
		Pos lexer.Position

		Count int `@Int`
	}
)

type Str string

var _ participle.Capture = (*Str)(nil)

func (s *Str) Capture(values []string) error {
	*s = Str(strings.ReplaceAll(strings.Trim(values[0], `"`), `""`, `"`))
	return nil
}

func (s *Str) String() string {
	return string(*s)
}

var (
	lexerOption = participle.Lexer(lexer.MustSimple([]lexer.SimpleRule{
		{Name: `Keyword`, Pattern: `\b(DRV|MSG|CINFO|TINFO|SINFO|PRGT|PRGC|PRGV|TCOUNT)\b`},
		{Name: `Int`, Pattern: `\d+`},
		{Name: `String`, Pattern: `"(""|[^"])*"`},
		{Name: `Punct`, Pattern: `[,:]`},
		{Name: `EOL`, Pattern: `\r?\n`},
		{Name: `whitespace`, Pattern: `\s+`},
	}))
	lineParser   = participle.MustBuild[Line](lexerOption)
	outputParser = participle.MustBuild[Output](lexerOption)
)
