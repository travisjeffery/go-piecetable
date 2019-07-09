package writegood_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	writegood "github.com/travisjeffery/writegood"
)

func TestInsert(t *testing.T) {
	d := &writegood.Document{}
	req := require.New(t)

	// initial case, empty doc
	d.Insert(0, []byte("hello"))
	t.Logf("added: %s\n", d.Added)
	req.Equal([]byte("hello"), d.Added)
	req.Equal([]*writegood.Piece{{
		Start:  0,
		Length: len("hello"),
		Type:   writegood.Added,
	}}, d.Pieces)
	b, _ := d.Bytes()
	req.Equal([]byte("hello"), b)
	t.Logf("added: %s\n", d.Added)

	// add to start
	d.Insert(0, []byte("world"))
	req.Equal([]*writegood.Piece{{
		Start:  len("hello"),
		Length: len("world"),
		Type:   writegood.Added,
	}, {
		Start:  0,
		Length: len("hello"),
		Type:   writegood.Added,
	}}, d.Pieces)
	b, _ = d.Bytes()
	req.Equal([]byte("worldhello"), b)

	// add to middle
	d.Insert(len("world"), []byte("bye"))
	t.Logf("added: %s\n", d.Added)
	req.Equal([]*writegood.Piece{{
		Start:  len("hello"),
		Length: len("world"),
		Type:   writegood.Added,
	}, {
		Start:  len("helloworld"),
		Length: len("bye"),
		Type:   writegood.Added,
	}, {
		Start:  0,
		Length: len("hello"),
		Type:   writegood.Added,
	}}, d.Pieces)
	b, _ = d.Bytes()
	req.Equal([]byte("worldbyehello"), b)

	d.Insert(len("wor"), []byte("yes"))
	t.Logf("added: %s\n", d.Added)
	req.Equal([]*writegood.Piece{{
		Start:  len("hello"),
		Length: len("wor"),
		Type:   writegood.Added,
	}, {
		Start:  len("helloworldbye"),
		Length: len("yes"),
		Type:   writegood.Added,
	}, {
		Start:  len("hellowor"),
		Length: len("ld"),
		Type:   writegood.Added,
	}, {
		Start:  len("helloworld"),
		Length: len("bye"),
		Type:   writegood.Added,
	}, {
		Start:  0,
		Length: len("hello"),
		Type:   writegood.Added,
	}}, d.Pieces)
	b, _ = d.Bytes()
	req.Equal([]byte("woryesldbyehello"), b)

	// add to end
	d.Insert(len("woryesldbyehello"), []byte("grief"))
	t.Logf("added: %s\n", d.Added)
	req.Equal([]*writegood.Piece{{
		Start:  len("hello"),
		Length: len("wor"),
		Type:   writegood.Added,
	}, {
		Start:  len("helloworldbye"),
		Length: len("yes"),
		Type:   writegood.Added,
	}, {
		Start:  len("hellowor"),
		Length: len("ld"),
		Type:   writegood.Added,
	}, {
		Start:  len("helloworld"),
		Length: len("bye"),
		Type:   writegood.Added,
	}, {
		Start:  0,
		Length: len("hello"),
		Type:   writegood.Added,
	}, {
		Start:  len("woryesldbyehello"),
		Length: len("grief"),
		Type:   writegood.Added,
	}}, d.Pieces)
	b, _ = d.Bytes()
	req.Equal([]byte("woryesldbyehellogrief"), b)
}
