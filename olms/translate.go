package olms

import (
	"encoding/json"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle
var i18nMessageFile *i18n.MessageFile

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	languages, err := filepath.Glob(joinPath(dir(Self), "languages/*.json"))
	if err != nil {
		return
	}
	for _, language := range languages {
		i18nMessageFile, _ = bundle.LoadMessageFile(language)
	}
}

func localize(c *gin.Context) map[string]string {
	lang, _ := c.Cookie("lang")
	localizer := i18n.NewLocalizer(bundle, lang, c.GetHeader("Accept-Language"))
	if i18nMessageFile != nil {
		translate := make(map[string]string, len(i18nMessageFile.Messages))
		var tag language.Tag
		for _, Message := range i18nMessageFile.Messages {
			var message string
			message, tag, _ = localizer.LocalizeWithTag(&i18n.LocalizeConfig{MessageID: Message.ID})
			translate[Message.ID] = message
		}
		c.Header("Content-Language", tag.String())
		return translate
	}
	return nil
}
