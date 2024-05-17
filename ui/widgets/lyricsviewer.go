package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	fynelyrics "github.com/dweymouth/fyne-lyrics"
	"github.com/dweymouth/supersonic/backend/mediaprovider"
)

type LyricsViewer struct {
	widget.BaseWidget

	noLyricsLabel widget.Label
	viewer        *fynelyrics.LyricsViewer
	lyrics        *mediaprovider.Lyrics
	nextLyricLine int

	container *fyne.Container
	isEmpty   bool
}

func NewLyricsViewer() *LyricsViewer {
	l := &LyricsViewer{noLyricsLabel: widget.Label{
		Text: "Lyrics not available",
	}, isEmpty: true}
	l.ExtendBaseWidget(l)
	l.container = container.NewStack(&l.noLyricsLabel)
	return l
}

func (l *LyricsViewer) SetLyrics(lyrics *mediaprovider.Lyrics) {
	l.lyrics = lyrics
	l.nextLyricLine = 0
	if lyrics == nil || len(lyrics.Lines) == 0 {
		if !l.isEmpty {
			l.container.Objects[0] = &l.noLyricsLabel
			l.isEmpty = true
			l.Refresh()
		}
		return
	}

	if l.viewer == nil {
		l.viewer = fynelyrics.NewLyricsViewer()
		l.viewer.ActiveLyricPosition = fynelyrics.ActiveLyricPositionTopThird
	}
	lines := make([]string, len(lyrics.Lines))
	for i, line := range lyrics.Lines {
		lines[i] = line.Text
	}
	l.viewer.SetLyrics(lines, lyrics.Synced)
	if l.isEmpty {
		l.container.Objects[0] = l.viewer
		l.isEmpty = false
		l.Refresh()
	}
}

func (l *LyricsViewer) UpdatePlayPos(timeSecs float64) {
	if l.lyrics == nil || !l.lyrics.Synced {
		return
	}
	// advance if needed
	if l.lyrics.Lines[l.nextLyricLine].Start <= timeSecs {
		l.viewer.NextLine()
		if l.nextLyricLine < len(l.lyrics.Lines) {
			l.nextLyricLine++
		}
	}
}

func (l *LyricsViewer) OnSeeked(timeSecs float64) {
	if l.lyrics == nil || !l.lyrics.Synced {
		return
	}
	curLine := 0
	for i, l := range l.lyrics.Lines {
		if l.Start < timeSecs {
			curLine = i + 1
			break
		}
	}
	l.viewer.SetCurrentLine(curLine)
}

func (l *LyricsViewer) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(l.container)
}
