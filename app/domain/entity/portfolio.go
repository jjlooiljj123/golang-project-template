package entity

type Album struct {
	ID    AlbumID
	Title string
}

type AlbumID string

func (pid *AlbumID) String() string {
	return string(*pid)
}
