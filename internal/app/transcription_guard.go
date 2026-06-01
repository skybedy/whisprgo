package app

import (
	"regexp"
	"strings"
	"unicode"
)

var commonCzechWords = map[string]struct{}{
	"a": {}, "ale": {}, "ano": {}, "co": {}, "do": {}, "i": {}, "jak": {}, "jako": {}, "je": {}, "jsem": {},
	"jsi": {}, "jsme": {}, "jste": {}, "k": {}, "na": {}, "ne": {}, "o": {}, "od": {}, "po": {}, "pro": {},
	"se": {}, "si": {}, "s": {}, "tak": {}, "to": {}, "tohle": {}, "toho": {}, "to je": {}, "u": {}, "v": {},
	"ve": {}, "z": {}, "za": {}, "že": {}, "už": {}, "tady": {}, "ten": {}, "ta": {}, "tohoto": {},
}

var wordSplitRe = regexp.MustCompile(`[A-Za-zÀ-ž']+`)

func shouldRetryParakeetForCzech(lang, text string) bool {
	if strings.ToLower(strings.TrimSpace(lang)) != "cs" {
		return false
	}
	t := strings.TrimSpace(text)
	if t == "" {
		return false
	}
	if containsCyrillic(t) {
		return true
	}
	words := wordSplitRe.FindAllString(strings.ToLower(t), -1)
	if len(words) < 3 {
		return false
	}
	czHits := 0
	for _, w := range words {
		if _, ok := commonCzechWords[w]; ok {
			czHits++
		}
	}
	if czHits <= 1 && likelyEnglish(words) {
		return true
	}
	return false
}

func containsCyrillic(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Cyrillic) {
			return true
		}
	}
	return false
}

func likelyEnglish(words []string) bool {
	en := map[string]struct{}{
		"a": {}, "an": {}, "and": {}, "are": {}, "be": {}, "can": {}, "do": {}, "for": {}, "hello": {},
		"how": {}, "i": {}, "is": {}, "it": {}, "my": {}, "not": {}, "of": {}, "please": {}, "test": {},
		"don't": {},
		"that": {}, "the": {}, "this": {}, "to": {}, "what": {}, "you": {}, "your": {},
	}
	hits := 0
	for _, w := range words {
		if _, ok := en[w]; ok {
			hits++
		}
	}
	return hits >= 2
}
