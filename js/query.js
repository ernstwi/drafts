let ds = Draft.query(...input);
let res = ds.map((d) => ({
  uuid: d.uuid,
  content: d.content,
}));
if (ds.length > 0) context.addSuccessParameter("result", JSON.stringify(res));
