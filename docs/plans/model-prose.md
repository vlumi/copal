# Plan: names and prose by a model

The survey's facts now hold on three specimens written by different hands, one of them in Japanese. What they cannot say is what a paragraph is for, in words a reader knows, or a rule as a sentence: `2100-AUTHORIZE` and a decision table of twelve columns are facts, not an explanation. A language model can write that explanation, and only that, on top of the facts the analysis established and the source they come from.

The deliverable is prose a reader can trust as far as it goes: every sentence about a fact the survey states, marked as a model's, kept apart from the facts, and checked by the survey before it is kept.

Part 1 is built. The rest waits, [milestone 16](../roadmap.md#16-names-and-prose-by-a-model), until the deterministic work before it is done: until then no part of the repository calls a model or needs an API key.

## Ground rules

- The model writes names and prose, and nothing else. It is given facts the analysis established and the source lines they come from, and it adds no fact: no link, rule, type, state or relation exists because a model wrote it.
- What it writes is kept apart from the facts, in `annotations/<specimen>/`, checked in, and the blueprints and the explorer show it beside the facts, marked as a model's.
- Reproducible without the model. Each annotation records a hash of the facts and the source it was written from, the model, and the prompt's version. The blueprints are rendered from the checked-in annotations, so `make ci` never calls a model and the blueprints stay golden. When the facts change, an annotation whose hash no longer matches is stale: the blueprint says so instead of showing it, and `make annotate` writes it anew.
- Checked before it is kept. A data item, paragraph, file or condition name a sentence uses must be one of those the model was given; an annotation that names anything else is rejected and reported, not stored.
- One prompt for every specimen. Nothing in the code or the prompts knows Pangaea, Laurasia or Gondwana, and annotating any of them is the same command. Gondwana's prose is English, as its documentation is, about facts whose names are Japanese.
- The history shows the prose as it shows the facts: a change of model, prompt or facts is a diff in `annotations/`.

## What it writes

- **Paragraphs and sections:** what each is for, in a sentence, and a name that says it, beside the name in the source.
- **Decision rules:** each rule of a decision table as one sentence, its conditions and its actions in the reader's words.
- **Programs:** what the program does, from its paragraphs, its files and its header comment.
- **Entities, fields and states:** a business name and a sentence for each, and a state machine as a short account of a record's life.

Not code, not a translation, and not a guess at an intent the facts do not support: where the facts are not enough, the annotation says so rather than inventing.

## Decisions

- **Model:** Claude Opus 5.5, `claude-opus-5-5`, through the Anthropic API, with the key from `ANTHROPIC_API_KEY`, called only by `make annotate`. One request per program and kind of annotation, about 150 requests for the three specimens in all.
- **Client:** the standard library's `net/http` against the Messages API, a request and a response type, so that the module keeps no dependencies; the Go SDK only if the client outgrows that.
- **Where it lives:** an internal package of the survey, `annotate`, which reads the blueprint model and writes the annotation files; the blueprints and the explorer read those files. It is not bedrock's: it is not language.

## Parts

1. **Annotations without a model.** The annotation files and their hashes, staleness, the check of the names a sentence uses, and their place in the blueprints and the explorer, all built and tested with a writer that needs no model and writes from the facts alone, so that the whole pipeline is golden before a model is called.

   As built: `blueprint` gives a paragraph's facts as text with its hash, checks an annotation's names against them, reads and writes `annotations/<specimen>/`, a file a program, and shows the prose on the program's page, and in the explorer on the paragraph's, marked as a model's, or that it is stale. `annotate.Paragraphs` keeps what is current and has a writer write the rest, rejecting what names anything its facts do not; it is tested with a writer that needs no model. The golden tests and `make blueprints` read the annotations; none exist yet, so no blueprint changed. A name the check takes as the code's is one in backquotes or in capitals with a hyphen, as COBOL's are; the prompt of part 2 asks for backquotes.
2. **The model and the paragraphs.** The client, the prompt, and `make annotate`, on the paragraphs and sections of the three specimens; the first annotations checked in and read against the source.
3. **Decision rules in prose.** The tables already have each rule as a sentence written from its conditions and statements, milestone 20; the model's prose says what a rule is for, in the business's words, and stands beside that sentence rather than in place of it.
4. **Programs, entities, fields and states.**
5. **Reading it against the code.** A sample of each kind read by hand against the source, what is wrong written down, the prompts adjusted, and the prose of each specimen checked for words that belong to another.

Each part one pull request or more.
