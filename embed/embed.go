package embed

import (
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/nabomhalang/halangcordgo/config"
)

var (
	log *config.Logger = config.NewLogger("embed")
)

type Embed struct {
	*discordgo.MessageEmbed
}

func NewEmbed() *Embed {
	return &Embed{&discordgo.MessageEmbed{
		Footer: &discordgo.MessageEmbedFooter{
			Text: "ⓒ 2023 nabomhalang. All rights reserved.",
		},
	}}
}

// SetTitle sets the title of the embed.
func (e *Embed) SetTitle(name string) *Embed {
	e.Title = name
	return e
}

// SetDescription sets the description of the embed.
func (e *Embed) SetDescription(description string) *Embed {
	if len(description) > 2028 {
		description = description[:2028]
	}

	e.Description = description
	return e
}

// AddField adds a field to the embed.
// name is the title of the field.
// value is the content of the field.
func (e *Embed) AddField(name, value string, inline bool) *Embed {
	if len(value) > 1024 {
		value = value[:1024]
	}

	if len(name) > 1024 {
		name = name[:1024]
	}

	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})

	return e
}

// SetFooter sets the footer of the embed.
// args[0] = text
// args[1] = iconURL
// args[2] = proxyURL
func (e *Embed) SetFooter(args ...string) *Embed {
	iconURL := ""
	text := ""
	proxyURL := ""

	switch {
	case len(args) > 2:
		proxyURL = args[2]
		fallthrough
	case len(args) > 1:
		iconURL = args[1]
		fallthrough
	case len(args) > 0:
		text = args[0]
	case len(args) == 0:
		return e
	}

	e.Footer = &discordgo.MessageEmbedFooter{
		IconURL:      iconURL,
		Text:         text,
		ProxyIconURL: proxyURL,
	}

	return e
}

// SetImage sets the image of the embed.
// args[0] = URL
// args[1] = proxyURL
func (e *Embed) SetImage(args ...string) *Embed {
	var URL string
	var proxyURL string

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		URL = args[0]
	}
	if len(args) > 1 {
		proxyURL = args[1]
	}
	e.Image = &discordgo.MessageEmbedImage{
		URL:      URL,
		ProxyURL: proxyURL,
	}
	return e
}

// SetThumbnail sets the thumbnail of the embed.
// args[0] = URL
// args[1] = proxyURL
func (e *Embed) SetThumbnail(args ...string) *Embed {
	var URL string
	var proxyURL string

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		URL = args[0]
	}
	if len(args) > 1 {
		proxyURL = args[1]
	}
	e.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:      URL,
		ProxyURL: proxyURL,
	}
	return e
}

// SetAuthor sets the author of the embed.
// args[0] = name
// args[1] = iconURL
// args[2] = URL
// args[3] = proxyURL
func (e *Embed) SetAuthor(args ...string) *Embed {
	var (
		name     string
		iconURL  string
		URL      string
		proxyURL string
	)

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		name = args[0]
	}
	if len(args) > 1 {
		iconURL = args[1]
	}
	if len(args) > 2 {
		URL = args[2]
	}
	if len(args) > 3 {
		proxyURL = args[3]
	}

	e.Author = &discordgo.MessageEmbedAuthor{
		Name:         name,
		IconURL:      iconURL,
		URL:          URL,
		ProxyIconURL: proxyURL,
	}

	return e
}

// SetURL sets the URL of the embed.
// url is the URL of the embed.
func (e *Embed) SetURL(url string) *Embed {
	e.URL = url
	return e
}

// SetColor sets the color of the embed.
// color can be an int representing the color value or a string.
// Available string values are: "red", "green", "blue", "yellow", "orange",
// "purple", "pink", "white", "black", and "random".
func (e *Embed) SetColor(color interface{}) *Embed {
	switch color.(type) {
	case int:
		e.Color = color.(int)
	case string:
		switch color.(string) {
		case "red":
			e.Color = 0xFF0000
		case "green":
			e.Color = 0x00FF00
		case "blue":
			e.Color = 0x0000FF
		case "yellow":
			e.Color = 0xFFFF00
		case "orange":
			e.Color = 0xFFA500
		case "purple":
			e.Color = 0x800080
		case "pink":
			e.Color = 0xFFC0CB
		case "white":
			e.Color = 0xFFFFFF
		case "black":
			e.Color = 0x000000
		case "random":
			e.Color = rand.Intn(0xFFFFFF)
		}
	default:
		e.Color = 0x000000
	}

	return e
}

// InlineAllFields sets all fields to be inline.
func (e *Embed) InlineAllFields() *Embed {
	for _, v := range e.Fields {
		v.Inline = true
	}
	return e
}

// Truncate truncates the embed to the maximum allowed length.
func (e *Embed) Truncate() *Embed {
	e.TruncateDescription()
	e.TruncateFields()
	e.TruncateFooter()
	e.TruncateTitle()
	return e
}

func (e *Embed) TruncateFields() *Embed {
	if len(e.Fields) > 25 {
		e.Fields = e.Fields[:EmbedLimitField]
	}

	for _, v := range e.Fields {

		if len(v.Name) > EmbedLimitFieldName {
			v.Name = v.Name[:EmbedLimitFieldName]
		}

		if len(v.Value) > EmbedLimitFieldValue {
			v.Value = v.Value[:EmbedLimitFieldValue]
		}

	}
	return e
}

func (e *Embed) TruncateDescription() *Embed {
	if len(e.Description) > EmbedLimitDescription {
		e.Description = e.Description[:EmbedLimitDescription]
	}
	return e
}

func (e *Embed) TruncateTitle() *Embed {
	if len(e.Title) > EmbedLimitTitle {
		e.Title = e.Title[:EmbedLimitTitle]
	}
	return e
}

func (e *Embed) TruncateFooter() *Embed {
	if e.Footer != nil && len(e.Footer.Text) > EmbedLimitFooter {
		e.Footer.Text = e.Footer.Text[:EmbedLimitFooter]
	}
	return e
}
