let d = Draft.find(input[0]);
let res = {
  uuid: d.uuid,
  content: d.content,
  isFlagged: d.isFlagged,
  isArchived: d.isArchived,
  isTrashed: d.isTrashed,
};
context.addSuccessParameter("result", JSON.stringify(res));
