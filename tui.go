package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/curt-hash/mkvbot/pkg/makemkv"
	"github.com/curt-hash/mkvbot/pkg/makemkv/defs"
	"github.com/curt-hash/mkvbot/pkg/moviedb"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	userInputPageName   = "userInputPage"
	chooseTitlePageName = "chooseTitlePage"
	logsPageName        = "logsPage"

	progressBarFullChar  = '█'
	progressBarEmptyChar = '░'
)

type textUserInterface struct {
	*beeper
	*tview.Application

	interruptChan chan struct{}

	// Top
	statusBox *statusBox

	// Left
	leftFlex         *tview.Flex
	driveInfoBox     *tview.TextView
	discInfoBox      *tview.TextView
	movieMetadataBox *tview.TextView
	titleInfoBox     *tview.TextView

	// Right
	userInputIntroText *tview.TextView
	userInputForm      *tview.Form
	logBox             *tview.TextView
	pages              *tview.Pages
}

func newTextUserInterface(beeper *beeper) *textUserInterface {
	app := tview.NewApplication()

	// Notify the main application about the Ctrl+C by closing the channel rather
	// than stopping the tview application abruptly (the default behavior).
	interruptChan := make(chan struct{})
	closeOnce := sync.OnceFunc(func() {
		close(interruptChan)
	})
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			closeOnce()
			return nil
		default:
			return event
		}
	})

	statusBox := newStatusBox()

	driveInfoBox := tview.NewTextView().SetWrap(false)
	driveInfoBox.SetBorder(true).SetTitle("Drive Information")

	discInfoBox := tview.NewTextView().SetWrap(false)
	discInfoBox.SetBorder(true).SetTitle("Disc Information")

	movieMetadataBox := tview.NewTextView().SetWrap(false)
	movieMetadataBox.SetBorder(true).SetTitle("Movie Metadata")

	titleInfoBox := tview.NewTextView().SetWrap(true)
	titleInfoBox.SetBorder(true).SetTitle("Title Information")

	userInputIntroText := tview.NewTextView().SetWrap(true)
	userInputForm := tview.NewForm()
	logBox := tview.NewTextView().
		SetWrap(true).
		SetScrollable(false).
		SetChangedFunc(func() { app.Draw() })
	pages := tview.NewPages().
		AddPage(
			userInputPageName,
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(userInputIntroText, 7, 0, false).
				AddItem(userInputForm, 0, 80, true),
			true,
			false,
		).
		AddPage(logsPageName, logBox, true, true)
	pages.SetBorder(true)

	leftFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(driveInfoBox, 4, 0, false).
		AddItem(discInfoBox, 10, 0, false).
		AddItem(movieMetadataBox, 5, 0, false).
		AddItem(titleInfoBox, 0, 40, false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(statusBox.box, 7, 0, false).
		AddItem(
			tview.NewFlex().
				AddItem(leftFlex, 0, 40, false).
				AddItem(pages, 0, 60, false),
			0,
			80,
			false,
		)

	app.SetRoot(flex, true).SetFocus(flex)

	return &textUserInterface{
		beeper: beeper,

		Application: app,

		interruptChan: interruptChan,

		statusBox: statusBox,

		leftFlex:         leftFlex,
		driveInfoBox:     driveInfoBox,
		discInfoBox:      discInfoBox,
		movieMetadataBox: movieMetadataBox,
		titleInfoBox:     titleInfoBox,

		userInputIntroText: userInputIntroText,
		userInputForm:      userInputForm,
		logBox:             logBox,
		pages:              pages,
	}
}

// waitForInterrupt returns after Ctrl+C.
func (t *textUserInterface) waitForInterrupt() {
	<-t.interruptChan
}

func (t *textUserInterface) run() error {
	return t.Run()
}

func (t *textUserInterface) setDriveInfo(name, volume string) {
	t.QueueUpdateDraw(func() {
		t.driveInfoBox.SetText(fmt.Sprintf("Name: %s\nVolume: %s", name, volume))
	})
}

func (t *textUserInterface) setStatus(format string, args ...any) {
	t.statusBox.setStatus(fmt.Sprintf(format, args...))
	t.updateStatusBox()
}

func (t *textUserInterface) setTask(format string, args ...any) {
	t.statusBox.setTask(fmt.Sprintf(format, args...))
	t.updateStatusBox()
}

func (t *textUserInterface) setSubtask(format string, args ...any) {
	t.statusBox.setSubtask(fmt.Sprintf(format, args...))
	t.updateStatusBox()
}

func (t *textUserInterface) setProgress(progress float64) {
	t.statusBox.progress = progress
	t.updateStatusBox()
}

func (t *textUserInterface) updateStatusBox() {
	t.QueueUpdateDraw(t.statusBox.update)
}

func (t *textUserInterface) setDiscInfo(info makemkv.Info) {
	t.QueueUpdateDraw(func() {
		w := t.discInfoBox.BatchWriter()
		defer w.Close()

		w.Clear()
		for _, item := range info {
			fmt.Fprintln(w, item)
		}

		t.leftFlex.ResizeItem(t.discInfoBox, len(info)+2, 0)
	})
}

func (t *textUserInterface) getMovieTitleForSearch(ctx context.Context, q string) (string, error) {
	continueChan := make(chan struct{})

	t.QueueUpdateDraw(func() {
		t.userInputIntroText.SetText("The query below will be used to search the movie database for metadata. Correct it if necessary and then hit Continue.")
		t.userInputForm.
			Clear(true).
			AddInputField("Query", q, 0, nil, nil).
			AddButton("Continue", func() {
				close(continueChan)
			})
		t.userInputForm.SetFocus(1)

		t.pages.SwitchToPage(userInputPageName)
		t.SetFocus(t.pages)
	})

	t.beep()

	select {
	case <-continueChan:
		t.QueueUpdateDraw(func() {
			t.pages.SwitchToPage(logsPageName)
		})
	case <-ctx.Done():
		return "", ctx.Err()
	}

	return t.userInputForm.GetFormItemByLabel("Query").(*tview.InputField).GetText(), nil
}

func (t *textUserInterface) getMovieMetadata(ctx context.Context, md *moviedb.Metadata) (*moviedb.Metadata, error) {
	continueChan := make(chan struct{})

	t.QueueUpdateDraw(func() {
		var buf strings.Builder
		fmt.Fprintf(&buf, "Correct the movie metadata below if necessary and then hit Continue.")
		if md.Tag != "" {
			fmt.Fprintf(&buf, "\n\nhttps://www.imdb.com/title/%s/", strings.TrimPrefix(md.Tag, "imdb-"))
		}
		t.userInputIntroText.SetText(buf.String())
		t.userInputForm.
			Clear(true).
			AddInputField("Name", md.Name, 0, nil, nil).
			AddInputField("Year", strconv.Itoa(md.Year), 4, func(textToCheck string, lastChar rune) bool {
				return len(textToCheck) <= 4 && unicode.IsDigit(lastChar)
			}, nil).
			AddInputField("Tag", md.Tag, 0, nil, nil).
			AddButton("Continue", func() {
				yearInput := t.userInputForm.GetFormItemByLabel("Year").(*tview.InputField)
				var err error
				if md.Year, err = strconv.Atoi(yearInput.GetText()); err != nil {
					t.userInputForm.SetFocus(1)
				} else {
					close(continueChan)
				}
			})
		t.userInputForm.SetFocus(3)

		t.pages.SwitchToPage(userInputPageName)
		t.SetFocus(t.pages)
	})

	t.beep()

	select {
	case <-continueChan:
		t.QueueUpdateDraw(func() {
			t.pages.SwitchToPage(logsPageName)
		})
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	md.Name = t.userInputForm.GetFormItemByLabel("Name").(*tview.InputField).GetText()
	md.Tag = t.userInputForm.GetFormItemByLabel("Tag").(*tview.InputField).GetText()
	return md, nil
}

func (t *textUserInterface) setMovieMetadata(md *moviedb.Metadata) {
	var s string
	if md != nil {
		s = fmt.Sprintf("Name: %s\nYear: %d\nTag: %s", md.Name, md.Year, md.Tag)
	}

	t.QueueUpdateDraw(func() {
		t.movieMetadataBox.SetText(s)
	})
}

func (t *textUserInterface) getBestTitle(ctx context.Context, choices []*makemkv.Title) (*makemkv.Title, error) {
	continueChan := make(chan struct{})

	var index int
	table := tview.NewTable().
		SetSelectable(true, false).
		SetSelectionChangedFunc(func(r, _ int) {
			index = r - 1
			if index >= 0 && index < len(choices) {
				t.setTitleInfoFunc(choices[index])()
			}
		}).
		SetSelectedFunc(func(r, _ int) {
			index = r - 1
			close(continueChan)
		})

	attrs := []defs.Attr{
		defs.TreeInfo,
		defs.MetadataLanguageName,
		defs.ChapterCount,
		defs.Duration,
		defs.DiskSize,
		defs.AngleInfo,
		defs.SegmentsMap,
		defs.Comment,
	}

	header := []string{"Index"}
	for _, attr := range attrs {
		header = append(header, attr.String())
	}
	header = append(header, "StreamCount")
	for i, s := range header {
		table.SetCell(0, i, tview.NewTableCell(s).SetSelectable(false).SetExpansion(1))
	}

	for i, title := range choices {
		r := i + 1
		table.SetCellSimple(r, 0, strconv.Itoa(title.Index))
		for j, attr := range attrs {
			table.SetCellSimple(r, j+1, title.GetAttrDefault(attr, "-"))
		}
		table.SetCellSimple(r, len(attrs)+1, strconv.Itoa(len(title.Streams)))
	}

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewTextView().
				SetWrap(true).
				SetText("Heuristics identified multiple best titles. ¯\\_(ツ)_/¯\n\nInformation about the highlighted title is shown to the left. Use the arrow keys to change the highlighted title and scroll the table. Press Enter to choose the highlighted title."),
			5,
			0,
			false,
		).
		AddItem(table, 0, 100, true)

	t.setTitleInfo(choices[0])
	t.QueueUpdateDraw(func() {
		t.pages.AddAndSwitchToPage(chooseTitlePageName, flex, true)
		t.SetFocus(t.pages)
	})

	t.beep()

	select {
	case <-continueChan:
		t.QueueUpdateDraw(func() {
			t.pages.RemovePage(chooseTitlePageName)
			t.pages.SwitchToPage(logsPageName)
		})
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if index < 0 || index >= len(choices) {
		return nil, fmt.Errorf("invalid choice")
	}

	return choices[index], nil
}

func (t *textUserInterface) setTitleInfoFunc(title *makemkv.Title) func() {
	return func() {
		w := t.titleInfoBox.BatchWriter()
		defer w.Close()

		w.Clear()
		if title != nil {
			writeTitleInfo(w, title)
		}
	}
}

func (t *textUserInterface) setTitleInfo(title *makemkv.Title) {
	t.QueueUpdateDraw(t.setTitleInfoFunc(title))
}

type statusBox struct {
	box      *tview.TextView
	status   string
	task     string
	subtask  string
	progress float64
}

func newStatusBox() *statusBox {
	box := tview.NewTextView().SetWrap(false)
	box.SetBorder(true).SetTitle("Status")

	return &statusBox{
		box:      box,
		progress: -1,
	}
}

func (b *statusBox) setStatus(s string) {
	b.status = s
	b.setTask("")
}

func (b *statusBox) setTask(s string) {
	b.task = s
	b.progress = -1
	b.setSubtask("")
}

func (b *statusBox) setSubtask(s string) {
	b.subtask = s
}

func (b *statusBox) update() {
	w := b.box.BatchWriter()
	defer w.Close()

	w.Clear()

	fmt.Fprintln(w, b.status)
	if b.task != "" {
		fmt.Fprintln(w, b.task)
		if b.subtask != "" {
			fmt.Fprintln(w, b.subtask)
		}
	}

	if b.progress >= 0 {
		fmt.Fprintln(w)

		_, _, width, _ := b.box.GetInnerRect()
		progressBarChars := width - len(" 100%")
		fullChars := int(float64(progressBarChars) * b.progress)
		var buf bytes.Buffer
		buf.Grow(width)
		for i := range progressBarChars {
			if i < fullChars {
				buf.WriteRune(progressBarFullChar)
			} else {
				buf.WriteRune(progressBarEmptyChar)
			}
		}
		_, _ = w.Write(buf.Bytes())

		fmt.Fprintf(w, " %3d%%", int(b.progress*100))
	}
}

func writeTitleInfo(w io.Writer, title *makemkv.Title) {
	fmt.Fprintf(w, "Title %d\n\n", title.Index)

	for _, attr := range title.Info {
		if attr.ID == defs.PanelTitle {
			continue
		}

		fmt.Fprintf(w, "%s: %s\n", attr.ID, attr.Value)
	}
	fmt.Fprintf(w, "StreamCount: %d\n", len(title.Streams))

	audioStreamCount := 0
	audioLanguagesByType := make(map[string][]string)
	subtitlesStreamCount := 0
	subtitlesByLanguage := make(map[string]int)
	for _, stream := range title.Streams {
		treeInfo := stream.GetAttrDefault(defs.TreeInfo, "-")

		switch stream.Type() {
		case defs.TypeCodeVideo:
			size := stream.GetAttrDefault(defs.VideoSize, "unknown size")
			bitrate := stream.GetAttrDefault(defs.Bitrate, "unknown bit rate")
			if bitrate == "" {
				bitrate = "unknown bit rate"
			}
			fmt.Fprintf(w, "\nVideo: %s (%s @ %s)\n", treeInfo, size, bitrate)
		case defs.TypeCodeAudio:
			audioStreamCount++

			codec := stream.GetAttrDefault(defs.CodecLong, "Unknown Codec")
			layout := stream.GetAttrDefault(defs.AudioChannelLayoutName, "Unknown Layout")
			lang := stream.GetAttrDefault(defs.LangName, "Unknown")
			key := fmt.Sprintf("%s %s", codec, layout)
			audioLanguagesByType[key] = append(audioLanguagesByType[key], lang)
		case defs.TypeCodeSubtitles:
			subtitlesStreamCount++

			lang := stream.GetAttrDefault(defs.LangName, "Unknown")
			if n, ok := subtitlesByLanguage[lang]; ok {
				subtitlesByLanguage[lang] = n + 1
			} else {
				subtitlesByLanguage[lang] = 1
			}
		}
	}

	fmt.Fprintf(w, "\nAudio (%d streams):\n", audioStreamCount)
	for key, langs := range audioLanguagesByType {
		fmt.Fprintf(w, "  * %s (%s)\n", key, strings.Join(langs, ", "))
	}

	if len(subtitlesByLanguage) > 0 {
		fmt.Fprintf(w, "\nSubtitles (%d streams):\n", subtitlesStreamCount)
		for lang, count := range subtitlesByLanguage {
			fmt.Fprintf(w, "  * %s (%d)\n", lang, count)
		}
	}
}
