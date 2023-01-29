package drafts

type CreateOptions struct {
	Tags       []string
	Folder     Folder
	Flagged    bool
	Action     string
	AllowEmpty bool
	// Omitted: RetParam
}

type QueryOptions struct {
	Tags             []string
	OmitTags         []string
	Sort             Sort
	SortDescending   bool
	SortFlaggedToTop bool
}
