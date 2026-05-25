package post

type Status string

const (
	StatusPublished   Status = "published"
	StatusUnpublished Status = "unpublished"
)

var validStatuses = map[Status]bool{
	StatusPublished:   true,
	StatusUnpublished: true,
}

func (s Status) IsValid() bool {
	return validStatuses[s]
}
