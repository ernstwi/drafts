let { program, input, app } = JSON.parse(draft.content);
let header = `let input = ${JSON.stringify(input)};`;
eval(
  [
    `let input = ${JSON.stringify(input)};`,
    `let app = ${JSON.stringify(app)};`,
    program,
  ].join("\n"),
);
