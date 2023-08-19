package drafts

import _ "embed"

//go:embed js/query.js
var queryjs string

//go:embed js/trash.js
var trashjs string

//go:embed js/replace.js
var replacejs string

//go:embed js/load.js
var loadjs string
