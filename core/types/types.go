package types

type TranslationStatus int
type MangaCompletionStatus int
type MangaType int
type MangaFormat int
type MangaRating int

type MangaPage struct {
	ImageURL   string
	ImageIndex int
}

type ChapterPages []MangaPage

const (
	TS_UNKNOWN   TranslationStatus = 0
	TS_FINISHED  TranslationStatus = 1
	TS_CONTINUES TranslationStatus = 2
	TS_FREEZED   TranslationStatus = 3
)

const (
	MS_UNKNOWN   MangaCompletionStatus = 0
	MS_FINISHED  MangaCompletionStatus = 1
	MS_CONTINUES MangaCompletionStatus = 2
)

const (
	MT_UNKNOWN   MangaType = 0
	MT_MANGA_JP  MangaType = 1
	MT_MANGA_KOR MangaType = 2
	MT_MANGA_CN  MangaType = 3
	MT_MANGA_RU  MangaType = 4
	MT_COMIX     MangaType = 5
	MT_OEL_MANGA MangaType = 6
)

const (
	MF_UNKNOWN     MangaFormat = 0
	MF_SINGLE      MangaFormat = 1
	MF_WEB         MangaFormat = 2
	MF_IN_COLOR    MangaFormat = 3
	MF_DOJINSHI    MangaFormat = 4
	MF_4_COMA      MangaFormat = 5
	MF_COMPILATION MangaFormat = 6
)

const (
	MR_UNKNOWN MangaRating = 0
	MR_PG13    MangaRating = 1
	MR_NC17_R  MangaRating = 2
	MR_NC21    MangaRating = 3
)

type MangaChapter struct {
	Index         int
	VolumeNumber  string
	ChapterNumber string
	IsSpecial     bool
	Title         string
	URL           string
	Pages         ChapterPages
}

type MangaInfo struct {
	URL         string
	Year        int
	Volumes     int
	Translation TranslationStatus
	MangaStatus MangaCompletionStatus
	Type        MangaType
	Rating      MangaRating
	IsSinge     bool
	HasAgeLimit bool

	Title            string
	TitleAlternative string
	Description      string
	Publisher        string
	Hostname         string

	Format      []MangaFormat
	Covers      []string
	Genres      []string
	Authors     []string
	Artists     []string
	Translators []string

	Chapters []MangaChapter
}

type AbstractPlugin interface {
	SetUrl(url string)
	GetChapters()
	GetMetaData()
	GetData() MangaInfo
}
