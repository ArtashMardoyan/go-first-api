package post

type Status string

const (
	StatusPublished   Status = "published"
	StatusUnpublished Status = "unpublished"
)

func (s Status) IsValid() bool {
	return s == StatusPublished || s == StatusUnpublished
}
