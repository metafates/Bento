package popupwidget

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/clearwidget"
)

var _ bento.Widget = (*Popup)(nil)

type Popup struct {
	Block                *blockwidget.Block
	Horizontal, Vertical bento.Flex
	Style                bento.Style
	Width, Height        bento.Constraint
	Content              bento.Widget
	Padding              bento.Padding
}

func New(content bento.Widget) Popup {
	return Popup{
		Block:      nil,
		Horizontal: bento.FlexCenter,
		Vertical:   bento.FlexCenter,
		Style:      bento.NewStyle(),
		Width:      bento.ConstraintPercentage(60),
		Height:     bento.ConstraintPercentage(20),
		Padding:    bento.NewPadding(),
		Content:    content,
	}
}

func (p Popup) WithBlock(block blockwidget.Block) Popup {
	p.Block = &block
	return p
}

func (p Popup) WithStyle(style bento.Style) Popup {
	p.Style = style
	return p
}

func (p Popup) WithHeight(height bento.Constraint) Popup {
	p.Height = height
	return p
}

func (p Popup) WithWidth(width bento.Constraint) Popup {
	p.Width = width
	return p
}

func (p Popup) WithPadding(padding bento.Padding) Popup {
	p.Padding = padding
	return p
}

func (p Popup) Top() Popup {
	p.Vertical = bento.FlexStart
	return p
}

func (p Popup) Middle() Popup {
	p.Vertical = bento.FlexCenter
	return p
}

func (p Popup) Bottom() Popup {
	p.Vertical = bento.FlexEnd
	return p
}

func (p Popup) Left() Popup {
	p.Horizontal = bento.FlexStart
	return p
}

func (p Popup) Center() Popup {
	p.Horizontal = bento.FlexCenter
	return p
}

func (p Popup) Right() Popup {
	p.Horizontal = bento.FlexEnd
	return p
}

func (p Popup) Render(area bento.Rect, buffer *bento.Buffer) {
	vertical := bento.NewLayout(p.Height).Vertical().WithPadding(p.Padding).WithFlex(p.Vertical)
	horizontal := bento.NewLayout(p.Width).Horizontal().WithPadding(p.Padding).WithFlex(p.Horizontal)

	area = vertical.Split(area).Unwrap()
	area = horizontal.Split(area).Unwrap()

	clearwidget.New().Render(area, buffer)
	buffer.SetStyle(area, p.Style)

	if p.Block != nil {
		p.Block.Render(area, buffer)
		area = p.Block.Inner(area)
	}

	if area.IsEmpty() {
		return
	}

	p.Content.Render(area, buffer)
}
