package library

const defaultWastepaperSliceSize = 10

// Library stores wastepaper for reading
type Library struct {
	wastepaper []Wastepaper
}

// NewLibrary creates new library
func NewLibrary() *Library {
	lib := &Library{}
	return lib
}

// GetWastepaper returns a wastepaper at a specified position in library
func (lib *Library) GetWastepaper(pos int) Wastepaper {
	if lib.wastepaper == nil {
		return nil
	}
	if pos > len(lib.wastepaper) || pos < 1 {
		return nil
	}
	if lib.wastepaper[pos-1].IsTaken() {
		return nil
	}
	lib.wastepaper[pos-1].Take()
	return lib.wastepaper[pos-1]
}

// PutWastepaper makes the wp Wastepaper available for reading
func (lib *Library) PutWastepaper(wp Wastepaper) {
	wp.Put()
}

// AddWastepaper add wastepaper to library
func (lib *Library) AddWastepaper(wp Wastepaper) {
	if wp == nil {
		return
	}
	if lib.wastepaper == nil {
		lib.wastepaper = make([]Wastepaper, 0, defaultWastepaperSliceSize)
	}
	lib.wastepaper = append(lib.wastepaper, wp)
}
