package ui

import (
	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/phuonganhniie/spt-terminal/config"
	"github.com/phuonganhniie/spt-terminal/entities"
)

var (
	cfg   = config.NewConfig()
	color = entities.NewColors()
)

var (
	TrackStyle         tcell.Style
	AlbumStyle         tcell.Style
	ArtistStyle        tcell.Style
	TimeStyle          tcell.Style
	GenreStyle         tcell.Style
	PlaylistNavStyle   tcell.Style
	NavStyle           tcell.Style
	ContextMenuStyle   tcell.Style
	NotSelectableStyle tcell.Style
)

var (
	borders = map[bool]map[string]rune{
		true: {
			"TopLeft":          '╭',
			"TopRight":         '╮',
			"BottomRight":      '╯',
			"BottomLeft":       '╰',
			"Vertical":         '│',
			"Horizontal":       '─',
			"TopLeftFocus":     '╭',
			"TopRightFocus":    '╮',
			"BottomRightFocus": '╯',
			"BottomLeftFocus":  '╰',
			"VerticalFocus":    '│',
			"HorizontalFocus":  '─',
		},
		false: {
			"TopLeft":          tview.Borders.TopLeft,
			"TopRight":         tview.Borders.TopRight,
			"BottomRight":      tview.Borders.BottomRight,
			"BottomLeft":       tview.Borders.BottomLeft,
			"Vertical":         tview.Borders.Vertical,
			"Horizontal":       tview.Borders.Horizontal,
			"TopLeftFocus":     tview.Borders.TopLeftFocus,
			"TopRightFocus":    tview.Borders.TopRightFocus,
			"BottomRightFocus": tview.Borders.BottomRightFocus,
			"BottomLeftFocus":  tview.Borders.BottomLeftFocus,
			"VerticalFocus":    tview.Borders.VerticalFocus,
			"HorizontalFocus":  tview.Borders.HorizontalFocus,
		},
	}
)

func SetStyles() {
	TrackStyle = color.Track.Style()
	AlbumStyle = color.Album.Style()
	ArtistStyle = color.Artist.Style()
	TimeStyle = color.Timestamp.Style()
	GenreStyle = color.Genre.Style()
	PlaylistNavStyle = color.PlaylistNav.Style()
	NavStyle = color.Nav.Style()
	ContextMenuStyle = color.ContextMenu.Style()
	NotSelectableStyle = color.Null.Style()
	tview.Styles.BorderColorFocus = color.BorderFocus.Foreground()
	tview.Styles.BorderColor = color.Border.Foreground()
}

func SetBorderRunes() {
	tview.Borders.TopLeft = borders[cfg.RoundedCorners]["TopLeft"]
	tview.Borders.TopRight = borders[cfg.RoundedCorners]["TopRight"]
	tview.Borders.BottomRight = borders[cfg.RoundedCorners]["BottomRight"]
	tview.Borders.BottomLeft = borders[cfg.RoundedCorners]["BottomLeft"]
	tview.Borders.Vertical = borders[cfg.RoundedCorners]["Vertical"]
	tview.Borders.Horizontal = borders[cfg.RoundedCorners]["Horizontal"]
	tview.Borders.TopLeftFocus = borders[cfg.RoundedCorners]["TopLeftFocus"]
	tview.Borders.TopRightFocus = borders[cfg.RoundedCorners]["TopRightFocus"]
	tview.Borders.BottomRightFocus = borders[cfg.RoundedCorners]["BottomRightFocus"]
	tview.Borders.BottomLeftFocus = borders[cfg.RoundedCorners]["BottomLeftFocus"]
	tview.Borders.VerticalFocus = borders[cfg.RoundedCorners]["VerticalFocus"]
	tview.Borders.HorizontalFocus = borders[cfg.RoundedCorners]["HorizontalFocus"]
}
