package domain

type ThumbnailType string

const (
	TypeTiny    ThumbnailType = "tiny"
	TypeSmall   ThumbnailType = "small"
	TypeMedium  ThumbnailType = "medium"
	TypeLarge   ThumbnailType = "large"
	TypeSLarge  ThumbnailType = "slarge"
	TypeXLarge  ThumbnailType = "xlarge"
	TypeSquare  ThumbnailType = "square"
	TypeWide    ThumbnailType = "wide"
	TypeTall    ThumbnailType = "tall"
	TypePreview ThumbnailType = "preview"
)

type ThumbnailSize struct {
	Width  int
	Height int
	Label  string
	Type   ThumbnailType
}

var ThumbnailSizes = []ThumbnailSize{
	{Width: 100, Height: 100, Label: "100x100", Type: TypeTiny},
	{Width: 150, Height: 150, Label: "150x150", Type: TypeSmall},
	{Width: 300, Height: 300, Label: "300x300", Type: TypeMedium},
	{Width: 600, Height: 400, Label: "600x400", Type: TypeLarge},
	{Width: 400, Height: 600, Label: "400x600", Type: TypeLarge},
	{Width: 800, Height: 600, Label: "800x600", Type: TypeSLarge},
	{Width: 600, Height: 800, Label: "600x800", Type: TypeSLarge},
	{Width: 1024, Height: 768, Label: "1024x768", Type: TypeXLarge},
	{Width: 768, Height: 1024, Label: "768x1024", Type: TypeXLarge},
	{Width: 200, Height: 200, Label: "square", Type: TypeSquare},
	{Width: 400, Height: 200, Label: "wide", Type: TypeWide},
	{Width: 200, Height: 400, Label: "tall", Type: TypeTall},
	{Width: 100, Height: 100, Label: "preview", Type: TypePreview},
}
