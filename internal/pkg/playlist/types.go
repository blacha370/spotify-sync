package playlist

type simpleTrack struct {
	Artists string
	Name    string
	Uri     string
}

type artist struct {
	Name string `json:"name"`
}

type track struct {
	Artists []artist `json:"artists"`
	Uri     string   `json:"uri"`
	Name    string   `json:"name"`
}

type trackContainer struct {
	Track track `json:"track"`
}

type Playlist struct {
	Next   string           `json:"next"`
	Buffer []trackContainer `json:"items"`
	Tracks map[string]simpleTrack
}

func (p *Playlist) cleanBuffer() {
	for _, trackContainer := range p.Buffer {
		t := simpleTrack{Name: trackContainer.Track.Name, Uri: trackContainer.Track.Uri}
		as := ""
		for _, a := range trackContainer.Track.Artists {
			as = as + a.Name + ";"
		}
		t.Artists = as
		p.Tracks[trackContainer.Track.Name+"-"+as] = t
	}
	p.Buffer = []trackContainer{}
}
