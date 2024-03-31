let { program, input, app } = JSON.parse(draft.content);
eval(
  [
    `let input = ${JSON.stringify(input)};`,
    `let app = ${JSON.stringify(app)};`,
    program,
  ].join("\n"),
);
