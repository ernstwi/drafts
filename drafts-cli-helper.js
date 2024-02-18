let { program, input } = JSON.parse(draft.content);
let header = `let input = ${JSON.stringify(input)};`;
eval([header, program].join('\n'));
