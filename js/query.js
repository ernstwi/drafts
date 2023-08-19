let ds = Draft.query(...input);
let res = ds.map((d) => ({
  uuid: d.uuid,
  content: d.content,
  isFlagged: d.isFlagged,
  isArchived: d.isArchived,
  isTrashed: d.isTrashed,
}));
if (ds.length > 0) context.addSuccessParameter("result", JSON.stringify(res));
