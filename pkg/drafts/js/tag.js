let d = Draft.find(input[0]);
input[1].forEach((t) => d.addTag(t));
d.update();
