let input = JSON.parse(draft.content);
let ds = Draft.query(...input);
if (ds.length > 0)
  context.addSuccessParameter("uuids", ds.map((d) => d.uuid).join(","));
