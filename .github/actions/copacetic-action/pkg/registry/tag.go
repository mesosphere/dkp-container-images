package registry

import (
	"fmt"
	"regexp"
	"strconv"
)

// TagResolver resolves tag in series of tags.
type TagResolver interface {
	First() string
	Next(existingTags []string) string
	Latest(existingTags []string) string
}

// TODO(mh): Rename to numbered tag resolver.
func NewTagResolver(baseTag, suffix string) *tagResolver {
	r := regexp.MustCompile(
		fmt.Sprintf(
			"^%s-%s\\.(?P<version>\\d+)$",
			regexp.QuoteMeta(baseTag), regexp.QuoteMeta(suffix),
		),
	)

	return &tagResolver{
		baseTag: baseTag,
		suffix:  suffix,
		matcher: r,
	}
}

var _ TagResolver = &tagResolver{}

type tagResolver struct {
	matcher *regexp.Regexp
	baseTag string
	suffix  string
}

func (t *tagResolver) version(version int) string {
	return fmt.Sprintf("%s-%s.%d", t.baseTag, t.suffix, version)
}

func (t *tagResolver) Next(existingTags []string) string {
	nextVersion := 0

	for _, existingTag := range existingTags {
		if submatch := t.matcher.FindStringSubmatch(existingTag); submatch != nil {
			version, _ := strconv.Atoi(submatch[1])
			if version >= nextVersion {
				nextVersion = version + 1
			}
		}
	}

	return t.version(nextVersion)
}

func (t *tagResolver) First() string {
	return t.version(0)
}

func (t *tagResolver) Latest(existingTags []string) string {
	latestVersion := -1

	for _, existingTag := range existingTags {
		if submatch := t.matcher.FindStringSubmatch(existingTag); submatch != nil {
			version, _ := strconv.Atoi(submatch[1])
			if version >= latestVersion {
				latestVersion = version
			}
		}
	}

	if latestVersion >= 0 {
		return t.version(latestVersion)
	}

	return ""
}

func ParseBaseTag(patchedTag string) string {
	m := regexp.MustCompile(`^(.*)-[\S]+\.\d+`)
	if matches := m.FindStringSubmatch(patchedTag); matches != nil {
		return matches[1]
	}
	return ""
}
