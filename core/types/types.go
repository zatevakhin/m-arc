package types

type TranslationStatus int
type MangaCompletionStatus int

const (
	TS_UNKNOWN     TranslationStatus = 0
	TS_FINISHED    TranslationStatus = 1
	TS_IN_PROGRESS TranslationStatus = 2
)

const (
	MS_UNKNOWN     MangaCompletionStatus = 0
	MS_FINISHED    MangaCompletionStatus = 1
	MS_IN_PROGRESS MangaCompletionStatus = 2
)

type MangaChapter struct {
	Index         int
	VolumeNumber  string
	ChapterNumber string
	IsSpecial     bool
	Name          string
	URL           string
}

type MangaInfo struct {
	URL         string
	Year        int
	Volumes     int
	Translation TranslationStatus
	MangaStatus MangaCompletionStatus
	IsSinge     bool
	HasAgeLimit bool

	Title            string
	TitleAlternative string
	Description      string

	Covers      []string
	CoversBig   []string
	Genres      []string
	Authors     []string
	Translators []string
	Rating      []string

	Chapters []MangaChapter
}

type AbstractPlugin interface {
	SetUrl(url string)
	GetChapters()
	GetMetaData()
}
