let d = Draft.find(input[0]);
if (d) {
  let res = {
    uuid: d.uuid,
    content: d.content,
    tags: d.tags,
    isFlagged: d.isFlagged,
    isArchived: d.isArchived,
    isTrashed: d.isTrashed,
  };
  context.addSuccessParameter("result", JSON.stringify(res));
  context.addSuccessParameter("app", "kitty.app");
}
