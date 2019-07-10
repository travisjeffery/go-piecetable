package writegood

import (
	"bytes"
	"fmt"
)

type Document struct {
	// Content loaded from disk
	Original []byte
	// Content for user edits
	Added  []byte
	Pieces []*Piece
}

func (d *Document) Insert(offset int, b []byte) {
	added := &Piece{
		Start:  len(d.Added),
		Length: len(b),
		Type:   Added,
	}
	d.Added = append(d.Added, b...)

	if d.Pieces == nil {
		d.insert(0, added)
		return
	}

	curr := 0
	for i, p := range d.Pieces {
		if offset == curr {
			// Handle when change is made at start of piece. We need to add a new
			// piece at the front of the pieces.
			d.insert(i, added)
			return
		} else if curr+p.Length > offset {
			// Handle when change is made in middle of piece. We split the existing
			// piece to two pieces, so we update the existing piece to have a shorter
			// length, add a new piece for the added buffer in the middle, and add the
			// second piece for the splitted original piece.
			length := p.Length
			p.Length -= (p.Length + curr) - offset
			d.insert(i+1, added)
			d.insert(i+2, &Piece{
				Start:  p.Start + p.Length,
				Length: length - p.Length,
				Type:   p.Type,
			})
			return
		}
		curr += p.Length
	}

	// Handle when change is made at end. We need to add a piece at the end of the pieces.
	d.insert(len(d.Pieces), added)
}

func (d *Document) insert(i int, p *Piece) {
	d.Pieces = append(d.Pieces, &Piece{})
	copy(d.Pieces[i+1:], d.Pieces[i:])
	d.Pieces[i] = p
}

func (d *Document) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	for _, p := range d.Pieces {
		if p.Type == Original {
			buf.Write(d.Original[p.Start : p.Start+p.Length])
		} else if p.Type == Added {
			buf.Write(d.Added[p.Start : p.Start+p.Length])
		} else {
			return nil, fmt.Errorf("unknown piece type: %d", p.Type)
		}
	}
	return buf.Bytes(), nil
}

type Piece struct {
	// Which index to start reading from
	Start int
	// How many characters to read from that buffer
	Length int
	// Type of the piece
	Type PieceType
}

type PieceType int

const (
	Original PieceType = iota
	Added
)
