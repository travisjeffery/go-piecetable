package piecetable_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/travisjeffery/piecetable"
)

func TestInsert_Empty(t *testing.T) {
	d := &piecetable.PieceTable{}
	req := require.New(t)

	// initial case, empty doc
	d.Insert(0, []byte("hello"))
	t.Logf("added: %s\n", d.Add)
	req.Equal([]byte("hello"), d.Add)
	req.Equal([]*piecetable.Piece{{
		Start:  0,
		Length: len("hello"),
		Type:   piecetable.Add,
	}}, d.Pieces)
	b, _ := d.Bytes()
	req.Equal([]byte("hello"), b)
	t.Logf("added: %s\n", d.Add)

	// add to start
	d.Insert(0, []byte("world"))
	req.Equal([]*piecetable.Piece{{
		Start:  len("hello"),
		Length: len("world"),
		Type:   piecetable.Add,
	}, {
		Start:  0,
		Length: len("hello"),
		Type:   piecetable.Add,
	}}, d.Pieces)
	b, _ = d.Bytes()
	req.Equal([]byte("worldhello"), b)

	// add to middle
	d.Insert(len("world"), []byte("bye"))
	t.Logf("added: %s\n", d.Add)
	req.Equal([]*piecetable.Piece{{
		Start:  len("hello"),
		Length: len("world"),
		Type:   piecetable.Add,
	}, {
		Start:  len("helloworld"),
		Length: len("bye"),
		Type:   piecetable.Add,
	}, {
		Start:  0,
		Length: len("hello"),
		Type:   piecetable.Add,
	}}, d.Pieces)
	b, _ = d.Bytes()
	req.Equal([]byte("worldbyehello"), b)

	d.Insert(len("wor"), []byte("yes"))
	t.Logf("added: %s\n", d.Add)
	req.Equal([]*piecetable.Piece{{
		Start:  len("hello"),
		Length: len("wor"),
		Type:   piecetable.Add,
	}, {
		Start:  len("helloworldbye"),
		Length: len("yes"),
		Type:   piecetable.Add,
	}, {
		Start:  len("hellowor"),
		Length: len("ld"),
		Type:   piecetable.Add,
	}, {
		Start:  len("helloworld"),
		Length: len("bye"),
		Type:   piecetable.Add,
	}, {
		Start:  0,
		Length: len("hello"),
		Type:   piecetable.Add,
	}}, d.Pieces)
	b, _ = d.Bytes()
	req.Equal([]byte("woryesldbyehello"), b)

	// add to end
	d.Insert(len("woryesldbyehello"), []byte("grief"))
	t.Logf("added: %s\n", d.Add)
	req.Equal([]*piecetable.Piece{{
		Start:  len("hello"),
		Length: len("wor"),
		Type:   piecetable.Add,
	}, {
		Start:  len("helloworldbye"),
		Length: len("yes"),
		Type:   piecetable.Add,
	}, {
		Start:  len("hellowor"),
		Length: len("ld"),
		Type:   piecetable.Add,
	}, {
		Start:  len("helloworld"),
		Length: len("bye"),
		Type:   piecetable.Add,
	}, {
		Start:  0,
		Length: len("hello"),
		Type:   piecetable.Add,
	}, {
		Start:  len("woryesldbyehello"),
		Length: len("grief"),
		Type:   piecetable.Add,
	}}, d.Pieces)
	b, _ = d.Bytes()
	req.Equal([]byte("woryesldbyehellogrief"), b)
}

func TestInsert_Existing(t *testing.T) {
	req := require.New(t)
	d := &piecetable.PieceTable{
		Original: []byte("helloworld"),
		Pieces: []*piecetable.Piece{{
			Start:  0,
			Length: len("helloworld"),
			Type:   piecetable.Original,
		}},
	}
	d.Insert(len("hello"), []byte("earth"))
	t.Logf("original: %s, added: %s\n", d.Original, d.Add)
	req.Equal([]*piecetable.Piece{{
		Start:  0,
		Length: len("hello"),
		Type:   piecetable.Original,
	}, {
		Start:  0,
		Length: len("earth"),
		Type:   piecetable.Add,
	}, {
		Start:  len("hello"),
		Length: len("world"),
		Type:   piecetable.Original,
	}}, d.Pieces)
	b, _ := d.Bytes()
	req.Equal([]byte("helloearthworld"), b)
}

func TestDelete(t *testing.T) {
	req := require.New(t)
	d := &piecetable.PieceTable{
		Original: []byte("helloworld"),
		Pieces: []*piecetable.Piece{{
			Start:  0,
			Length: len("helloworld"),
			Type:   piecetable.Original,
		}},
	}
	d.Insert(len("hello"), []byte("earth"))
	b, _ := d.Bytes()
	req.Equal([]byte("helloearthworld"), b)

	// delete whole earth piece
	d.Delete(len("hello"), len("helloearth"))
	req.Equal([]*piecetable.Piece{{
		Start:  0,
		Length: len("hello"),
		Type:   piecetable.Original,
	}, {
		Start:  len("hello"),
		Length: len("world"),
		Type:   piecetable.Original,
	}}, d.Pieces)
	b, _ = d.Bytes()
	req.Equal([]byte("helloworld"), b)

	// delete part of piece
	d.Delete(len("hel"), len("hello"))
	req.Equal([]*piecetable.Piece{{
		Start:  0,
		Length: len("hel"),
		Type:   piecetable.Original,
	}, {
		Start:  len("hello"),
		Length: len("world"),
		Type:   piecetable.Original,
	}}, d.Pieces)
	b, _ = d.Bytes()
	req.Equal([]byte("helworld"), b)

	// delete middle of piece
	d.Delete(len("helwor"), len("helworl"))
	req.Equal([]*piecetable.Piece{{
		Start:  0,
		Length: len("hel"),
		Type:   piecetable.Original,
	}, {
		Start:  len("hello"),
		Length: len("wor"),
		Type:   piecetable.Original,
	}, {
		Start:  len("helloworl"),
		Length: len("d"),
		Type:   piecetable.Original,
	}}, d.Pieces)
	b, _ = d.Bytes()
	req.Equal([]byte("helword"), b)
}
