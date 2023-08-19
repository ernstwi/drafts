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

//go:embed js/archive.js
var archivejs string

//go:embed js/get.js
var getjs string
