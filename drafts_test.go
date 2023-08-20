package drafts

import (
	"fmt"
	"testing"
	"time"

	"github.com/ernstwi/drafts/assert"
)

func TestCreateDefault(t *testing.T) {
	text := rand()
	uuid := Create(text, CreateOptions{})
	defer func() {
		Trash(uuid)
	}()
	draft := Get(uuid)
	assert.DeepEqual(t, Draft{
		UUID:    uuid,
		Content: text,
		Tags:    []string{},
	}, draft)
}

func TestCreateFlagged(t *testing.T) {
	text := rand()
	uuid := Create(text, CreateOptions{Flagged: true})
	defer func() {
		Trash(uuid)
	}()
	draft := Get(uuid)
	assert.DeepEqual(t, Draft{
		UUID:      uuid,
		Content:   text,
		Tags:      []string{},
		IsFlagged: true,
	}, draft)
}

func TestCreateArchived(t *testing.T) {
	text := rand()
	uuid := Create(text, CreateOptions{Folder: FolderArchive})
	defer func() {
		Trash(uuid)
	}()
	draft := Get(uuid)
	assert.DeepEqual(t, Draft{
		UUID:       uuid,
		Content:    text,
		Tags:       []string{},
		IsArchived: true,
	}, draft)
}

func TestCreateTags(t *testing.T) {
	text := rand()
	tag := rand()
	uuid := Create(text, CreateOptions{Tags: []string{tag}})
	defer func() {
		Trash(uuid)
	}()
	draft := Get(uuid)
	assert.DeepEqual(t, Draft{
		UUID:    uuid,
		Content: text,
		Tags:    []string{tag},
	}, draft)
}

func TestPrepend(t *testing.T) {
	text := rand()
	prefix := rand()
	uuid := Create(text, CreateOptions{})
	defer func() {
		Trash(uuid)
	}()
	Prepend(uuid, prefix)
	content := Get(uuid).Content
	assert.Equal(t, prefix+"\n"+text, content)
}

func TestAppend(t *testing.T) {
	text := rand()
	suffix := rand()
	uuid := Create(text, CreateOptions{})
	defer func() {
		Trash(uuid)

	}()
	Append(uuid, suffix)
	content := Get(uuid).Content
	assert.Equal(t, text+"\n"+suffix, content)
}

func TestReplace(t *testing.T) {
	text := rand()
	replacement := rand()
	uuid := Create(text, CreateOptions{})
	defer func() {
		Trash(uuid)
	}()
	Replace(uuid, replacement)
	content := Get(uuid).Content
	assert.Equal(t, replacement, content)
}

func TestTrash(t *testing.T) {
	text := rand()
	uuid := Create(text, CreateOptions{})
	Trash(uuid)
	draft := Get(uuid)
	assert.DeepEqual(t, Draft{
		UUID:      uuid,
		Content:   text,
		Tags:      []string{},
		IsTrashed: true,
	}, draft)
}

func TestArchive(t *testing.T) {
	text := rand()
	uuid := Create(text, CreateOptions{})
	Archive(uuid)
	draft := Get(uuid)
	assert.DeepEqual(t, Draft{
		UUID:       uuid,
		Content:    text,
		Tags:       []string{},
		IsArchived: true,
	}, draft)
}

func TestQuery(t *testing.T) {
	a := Create("A", CreateOptions{Tags: []string{"test", "a"}})
	b := Create("B", CreateOptions{Tags: []string{"test", "b"}, Flagged: true})
	c := Create("C", CreateOptions{Tags: []string{"test", "c"}})

	defer func() {
		for _, uuid := range []string{a, b, c} {
			Trash(uuid)
		}
	}()

	uuids := getUUIDs(Query("", FilterInbox, QueryOptions{Tags: []string{"test"}}))
	assert.EqualSlice(t, []string{a, b, c}, uuids)

	uuids = getUUIDs(Query("", FilterInbox, QueryOptions{Tags: []string{"test", "a"}}))
	assert.EqualSlice(t, []string{a}, uuids)

	uuids = getUUIDs(Query("", FilterInbox, QueryOptions{
		Tags:     []string{"test"},
		OmitTags: []string{"a"},
	}))
	assert.EqualSlice(t, []string{b, c}, uuids)

	// TODO: Testing Sort requires draft modification

	uuids = getUUIDs(Query("", FilterInbox, QueryOptions{
		Tags:           []string{"test"},
		SortDescending: true,
	}))
	assert.EqualSlice(t, []string{c, b, a}, uuids)

	uuids = getUUIDs(Query("", FilterInbox, QueryOptions{
		Tags:             []string{"test"},
		SortFlaggedToTop: true,
	}))
	assert.EqualSlice(t, []string{b, a, c}, uuids)
}

func TestSelect(t *testing.T) {
	a := Create("a", CreateOptions{})
	b := Create("b", CreateOptions{})
	defer func() {
		Trash(a)
		Trash(b)
	}()
	b_ := Get(Active()).Content
	Select(a)
	a_ := Get(Active()).Content
	assert.Equal(t, "a", a_)
	assert.Equal(t, "b", b_)
}

func TestGetSpecialChars(t *testing.T) {
	t.Skip()
	// https://en.wikipedia.org/wiki/URL_encoding#Percent-encoding_reserved_characters
	chars := []string{"␣", "!", "\"", "#", "$", "%", "&", "'", "(", ")", "*", "+", ",", "/", ":", ";", "=", "?", "@", "[", "]"}
	for _, c := range chars {
		uuid := Create(c, CreateOptions{})
		defer func() {
			Trash(uuid)
		}()
		content := Get(uuid).Content
		assert.Equal(t, c, content)
	}
}

func TestTag(t *testing.T) {
	text := rand()
	tag := rand()
	uuid := Create(text, CreateOptions{})
	defer func() {
		Trash(uuid)
	}()
	Tag(uuid, tag)
	draft := Get(uuid)
	assert.DeepEqual(t, Draft{
		UUID:    uuid,
		Content: text,
		Tags:    []string{tag},
	}, draft)
}

// ---- Helpers ----------------------------------------------------------------

// Return a random string.
func rand() string {
	return fmt.Sprint(time.Now().UnixNano())
}

// Extract UUIDs from a slice of Drafts
func getUUIDs(ds []Draft) []string {
	uuids := make([]string, len(ds))
	for i := range ds {
		uuids[i] = ds[i].UUID
	}
	return uuids
}
