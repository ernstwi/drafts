package drafts

type CreateOptions struct {
	Tags    []string
	Folder  Folder
	Flagged bool
}

type QueryOptions struct {
	Tags             []string
	OmitTags         []string
	Sort             Sort
	SortDescending   bool
	SortFlaggedToTop bool
}
