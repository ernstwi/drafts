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
	res := Get(uuid)
	assert.Equal(t, text, res)
}

func TestCreateTags(t *testing.T) {
	text := rand()
	tag := rand()
	uuid := Create(text, CreateOptions{Tags: []string{tag}})
	defer func() {
		Trash(uuid)
	}()
	res := Query("", FilterInbox, QueryOptions{Tags: []string{tag}})
	assert.Equal(t, uuid, res[0].UUID)
}

func TestCreateFolder(t *testing.T) {
	text := rand()
	uuid := Create(text, CreateOptions{Folder: FolderArchive})
	defer func() {
		Trash(uuid)
	}()
	res := Query(text, FilterArchive, QueryOptions{})
	empty := Query(text, FilterInbox, QueryOptions{})
	assert.Equal(t, uuid, res[0].UUID)
	assert.Equal(t, 0, len(empty))
}

func TestCreateFlagged(t *testing.T) {
	text := rand()
	uuid := Create(text, CreateOptions{Flagged: true})
	uuid_ := Create(text, CreateOptions{})
	defer func() {
		Trash(uuid)
		Trash(uuid_)
	}()
	res := Query(text, FilterFlagged, QueryOptions{})
	assert.Equal(t, 1, len(res))
	assert.Equal(t, uuid, res[0].UUID)
}

func TestCreateAction(t *testing.T) {
	text := rand()
	uuid := Create(text, CreateOptions{Action: "Indent"})
	defer func() {
		Trash(uuid)
	}()
	res := Get(uuid)
	assert.Equal(t, "	"+text, res)
}

// Skipped: Requires interaction in Drafts GUI
func TestCreateAllowEmpty(t *testing.T) {
	t.SkipNow()
	text := ""
	uuid := Create(text, CreateOptions{Action: "Indent", AllowEmpty: false})
	defer func() {
		Trash(uuid)
	}()
	res := Get(uuid)
	assert.Equal(t, "", res)
}

func TestPrepend(t *testing.T) {
	text := rand()
	prefix := rand()
	uuid := Create(text, CreateOptions{})
	defer func() {
		Trash(uuid)
	}()
	Prepend(uuid, prefix)
	res := Get(uuid)
	assert.Equal(t, prefix+"\n"+text, res)
}

func TestAppend(t *testing.T) {
	text := rand()
	suffix := rand()
	uuid := Create(text, CreateOptions{})
	defer func() {
		Trash(uuid)
	}()
	Append(uuid, suffix)
	res := Get(uuid)
	assert.Equal(t, text+"\n"+suffix, res)
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

	res := uuids(Query("", FilterInbox, QueryOptions{Tags: []string{"test"}}))
	assert.EqualSlice(t, []string{a, b, c}, res)

	res = uuids(Query("", FilterInbox, QueryOptions{Tags: []string{"test", "a"}}))
	assert.EqualSlice(t, []string{a}, res)

	res = uuids(Query("", FilterInbox, QueryOptions{
		Tags:     []string{"test"},
		OmitTags: []string{"a"},
	}))
	assert.EqualSlice(t, []string{b, c}, res)

	// TODO: Testing Sort requires draft modification

	res = uuids(Query("", FilterInbox, QueryOptions{
		Tags:           []string{"test"},
		SortDescending: true,
	}))
	assert.EqualSlice(t, []string{c, b, a}, res)

	res = uuids(Query("", FilterInbox, QueryOptions{
		Tags:             []string{"test"},
		SortFlaggedToTop: true,
	}))
	assert.EqualSlice(t, []string{b, a, c}, res)
}

// ---- Helpers ----------------------------------------------------------------

// Return a random string.
func rand() string {
	return fmt.Sprint(time.Now().UnixNano())
}

// Extract UUIDs from a slice of Drafts
func uuids(ds []Draft) []string {
	res := make([]string, len(ds))
	for i := range ds {
		res[i] = ds[i].UUID
	}
	return res
}
