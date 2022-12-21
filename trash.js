let input = JSON.parse(draft.content);
let d = Draft.find(input);
d.isTrashed = true;
d.update();
