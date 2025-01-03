package textwidget

import (
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/metafates/bento"
	"github.com/rivo/uniseg"
)

var _ bento.Widget = (*Span)(nil)

type Span struct {
	Style   lipgloss.Style
	Content string
}

func NewSpan(v string) *Span {
	return &Span{
		Style:   lipgloss.NewStyle(),
		Content: v,
	}
}

func (s *Span) Width() int {
	return uniseg.StringWidth(s.Content)
}

func (s *Span) Render(area bento.Rect, buffer *bento.Buffer) {
	area = area.Intersection(buffer.Area)
	if area.IsEmpty() {
		return
	}

	x, y := area.X, area.Y

	for i, grapheme := range s.StyledGraphemes(lipgloss.NewStyle()) {
		symbolWidth := grapheme.Width

		nextX := x + symbolWidth
		if nextX > area.Right() {
			break
		}

		switch {
		case i == 0:
			// the first grapheme is always set on the cell
			buffer.CellAt(bento.Position{X: x, Y: y}).SetSymbol(grapheme.Symbol).SetStyle(grapheme.Style)
		case x == area.X:
			// there is one or more zero-width graphemes in the first cell, so the first cell
			// must be appended to.
			buffer.CellAt(bento.Position{X: x, Y: y}).AppendSymbol(grapheme.Symbol).SetStyle(grapheme.Style)
		case symbolWidth == 0:
			// append zero-width graphemes to the previous cell
			buffer.CellAt(bento.Position{X: x - 1, Y: y}).AppendSymbol(grapheme.Symbol).SetStyle(grapheme.Style)
		default:
			// just a normal grapheme (not first, not zero-width, not overflowing the area)
			buffer.CellAt(bento.Position{X: x, Y: y}).SetSymbol(grapheme.Symbol).SetStyle(grapheme.Style)
		}

		for xHidden := x + 1; xHidden < nextX; xHidden++ {
			buffer.CellAt(bento.Position{X: xHidden, Y: y}).Reset()
		}

		x = nextX
	}
}

func (s *Span) StyledGraphemes(style lipgloss.Style) []StyledGrapheme {
	// TODO: patch styles
	// style = style.Patch(s.Style)

	var result []StyledGrapheme

	graphemes := uniseg.NewGraphemes(s.Content)

graphemes:
	for graphemes.Next() {
		for _, r := range graphemes.Runes() {
			if unicode.IsControl(r) {
				continue graphemes
			}
		}

		result = append(result, StyledGrapheme{
			Style:  style,
			Symbol: graphemes.Str(),
			Width:  graphemes.Width(),
		})
	}

	return result
}

func renderSpans(spans []Span, area bento.Rect, buffer *bento.Buffer, spanSkipWidth int) {
	for _, s := range spansAfterWidth(spans, spanSkipWidth) {
		area = area.IndentX(s.Offset)
		if area.IsEmpty() {
			break
		}

		s.Span.Render(area, buffer)

		area = area.IndentX(s.Width)
	}
}

type _SpanAfterWidth struct {
	Span   Span
	Width  int
	Offset int
}

func spansAfterWidth(spans []Span, skipWidth int) []_SpanAfterWidth {
	result := make([]_SpanAfterWidth, 0, len(spans))

	for _, s := range spans {
		spanWidth := s.Width()

		if skipWidth >= spanWidth {
			skipWidth = max(0, skipWidth-spanWidth)
			continue
		}

		availableWidth := max(0, spanWidth-skipWidth)
		skipWidth = 0

		if spanWidth <= availableWidth {
			result = append(result, _SpanAfterWidth{
				Span:   s,
				Width:  spanWidth,
				Offset: 0,
			})

			continue
		}

		// Span is only partially visible. As the end is truncated by the area width, only
		// truncate the start of the span.
		content, actualWidth := unicodeTruncateStart(s.Content, availableWidth)

		firstGraphemeOffset := max(0, availableWidth-actualWidth)

		result = append(result, _SpanAfterWidth{
			Span: Span{
				Style:   s.Style,
				Content: content,
			},
			Width:  actualWidth,
			Offset: firstGraphemeOffset,
		})
	}

	return result
}

func unicodeTruncateStart(s string, maxWidth int) (string, int) {
	graphemes := uniseg.NewGraphemes(uniseg.ReverseString(s))

	var (
		currentWidth int
		newWidth     int
	)

	byteIndex := len(s)

	for graphemes.Next() {
		from, to := graphemes.Positions()

		index := from
		width := to - from

		currentWidth += width

		if currentWidth <= maxWidth {
			byteIndex = index
			newWidth = currentWidth
		} else {
			break
		}
	}

	result := s[byteIndex:]

	return result, newWidth
}

type StyledGrapheme struct {
	Style  lipgloss.Style
	Symbol string
	Width  int
}
