package main

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/curt-hash/mkvbot/pkg/makemkvcon"
	"github.com/curt-hash/mkvbot/pkg/makemkvcon/defs"
	"github.com/curt-hash/mkvbot/pkg/moviedb"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	userInputPageName = "userInputPage"
	logsPageName      = "logsPage"

	progressBarFullChar  = '█'
	progressBarEmptyChar = '░'
)

type textUserInterface struct {
	*tview.Application

	interruptChan chan struct{}

	// Top
	statusBox *statusBox

	// Left
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

func newTextUserInterface() *textUserInterface {
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

	titleInfoBox := tview.NewTextView().SetWrap(false)
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

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(statusBox.box, 7, 0, false).
		AddItem(
			tview.NewFlex().
				AddItem(
					tview.NewFlex().
						SetDirection(tview.FlexRow).
						AddItem(driveInfoBox, 4, 0, false).
						AddItem(discInfoBox, 0, 20, false).
						AddItem(movieMetadataBox, 5, 0, false).
						AddItem(titleInfoBox, 0, 40, false),
					0,
					40,
					false,
				).
				AddItem(pages, 0, 60, false),
			0,
			80,
			false,
		)

	app.SetRoot(flex, true).SetFocus(flex)

	return &textUserInterface{
		Application: app,

		interruptChan: interruptChan,

		statusBox: statusBox,

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

func (t *textUserInterface) setDiscInfo(info makemkvcon.Info) {
	t.QueueUpdateDraw(func() {
		w := t.discInfoBox.BatchWriter()
		defer w.Close()

		w.Clear()
		for _, item := range info {
			fmt.Fprintln(w, item)
		}
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

		t.pages.SwitchToPage(userInputPageName)
		t.SetFocus(t.pages)
	})

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
					t.SetFocus(yearInput)
				} else {
					close(continueChan)
				}
			})

		t.pages.SwitchToPage(userInputPageName)
		t.SetFocus(t.pages)
	})

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

func (t *textUserInterface) getBestTitle(ctx context.Context, choices []*makemkvcon.Title) (*makemkvcon.Title, error) {
	continueChan := make(chan struct{})

	t.QueueUpdateDraw(func() {
		t.userInputIntroText.SetText("Heuristics identified multiple best titles. ¯\\_(ツ)_/¯\n\nChoose a title to backup and then hit Continue.")
		t.userInputForm.Clear(true)

		var (
			buf          strings.Builder
			choiceLabels []string
		)
		for i, title := range choices {
			fmt.Fprintf(&buf, "Title %d\n\n", title.Index)

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
			for _, attr := range attrs {
				v, err := title.GetAttr(attr)
				if err != nil {
					v = err.Error()
				}
				fmt.Fprintf(&buf, "%s: %s\n", attr.String(), v)
			}

			fmt.Fprintf(&buf, "\nStreams (%d):\n", len(title.Streams))
			for i, stream := range title.Streams {
				treeInfo, _ := stream.GetAttr(defs.TreeInfo)

				fmt.Fprintf(&buf, " %d. ", i)
				typ, _ := stream.GetAttr(defs.Type)
				switch typ {
				case "Video":
					size, _ := stream.GetAttr(defs.VideoSize)
					bitrate, _ := stream.GetAttr(defs.Bitrate)
					fmt.Fprintf(&buf, "Video (%s, %s @ %s)", treeInfo, size, bitrate)
				case "Audio":
					fmt.Fprintf(&buf, "Audio (%s)", treeInfo)
				case "Subtitles":
					lang, _ := stream.GetAttr(defs.LangName)
					fmt.Fprintf(&buf, "Subtitles (%s)", lang)
				}
				fmt.Fprintln(&buf)
			}

			label := strconv.Itoa(i + 1)
			h := len(attrs) + len(title.Streams) + 5
			t.userInputForm.AddTextView(label, buf.String(), 0, h, false, false)
			choiceLabels = append(choiceLabels, label)
		}

		t.userInputForm.AddDropDown("Choice", choiceLabels, -1, nil)
		t.userInputForm.AddButton("Continue", func() {
			close(continueChan)
		})

		t.pages.SwitchToPage(userInputPageName)
		t.SetFocus(t.pages)
	})

	select {
	case <-continueChan:
		t.QueueUpdateDraw(func() {
			t.pages.SwitchToPage(logsPageName)
		})
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	index, _ := t.userInputForm.GetFormItemByLabel("Choice").(*tview.DropDown).GetCurrentOption()
	if index < 0 || index >= len(choices) {
		return nil, fmt.Errorf("invalid choice")
	}

	return choices[index], nil
}

func (t *textUserInterface) setTitleInfo(title *makemkvcon.Title) {
	t.QueueUpdateDraw(func() {
		w := t.titleInfoBox.BatchWriter()
		defer w.Close()
		w.Clear()

		if title == nil {
			return
		}

		fmt.Fprintf(w, "Title %d\n\n", title.Index)

		for _, attr := range title.Info {
			fmt.Fprintf(w, "%s: %s\n", defs.Attr(attr.ID), attr.Value)
		}

		fmt.Fprintf(w, "\nStreams (%d):\n", len(title.Streams))
		for i, stream := range title.Streams {
			treeInfo, _ := stream.GetAttr(defs.TreeInfo)

			fmt.Fprintf(w, " %d. ", i)
			typ, _ := stream.GetAttr(defs.Type)
			switch typ {
			case "Video":
				size, _ := stream.GetAttr(defs.VideoSize)
				bitrate, _ := stream.GetAttr(defs.Bitrate)
				fmt.Fprintf(w, "Video (%s, %s @ %s)", treeInfo, size, bitrate)
			case "Audio":
				fmt.Fprintf(w, "Audio (%s)", treeInfo)
			case "Subtitles":
				lang, _ := stream.GetAttr(defs.LangName)
				fmt.Fprintf(w, "Subtitles (%s)", lang)
			default:
				fmt.Fprintf(w, "%s", typ)
			}
			fmt.Fprintln(w)
		}
	})
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
		w.Write(buf.Bytes())

		fmt.Fprintf(w, " %3d%%", int(b.progress*100))
	}
}
