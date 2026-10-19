// Package annotate has a writer, a language model in the end, write names
// and prose about the subjects of a system's blueprints from the facts the
// survey established about each, and keeps what it writes only where it
// names nothing those facts do not. What it keeps is no fact: the
// blueprints show it beside the facts, as a model's reading of them.
package annotate

import (
	"context"
	"slices"

	"github.com/vlumi/copal/survey/blueprint"
)

// Writer writes a name and prose about a subject from its facts.
type Writer interface {
	Write(ctx context.Context, f blueprint.Facts) (name, text string, err error)
	// Identity is the writer's name, such as its model's, and the version
	// of the prompt it writes from.
	Identity() (writer, prompt string)
}

// Rejection is an annotation a writer wrote that names what its facts do
// not, and why it was not kept.
type Rejection struct {
	Annotation blueprint.Annotation
	Reason     string
}

// Paragraphs annotates the paragraphs of every program of a system, in
// order: an earlier annotation is kept as long as its facts are what they
// were and it is the writer's and prompt's, and for the rest the writer
// writes anew. What names what its facts do not is rejected, not kept.
func Paragraphs(ctx context.Context, s *blueprint.System, earlier []blueprint.Annotation, w Writer) ([]blueprint.Annotation, []Rejection, error) {
	writer, prompt := w.Identity()
	var kept []blueprint.Annotation
	var rejected []Rejection
	for _, p := range s.Programs {
		for _, para := range p.Paragraphs {
			facts, err := s.ParagraphFacts(p.Name, para.Name)
			if err != nil {
				return nil, nil, err
			}
			if i := slices.IndexFunc(earlier, func(a blueprint.Annotation) bool {
				return a.Kind == facts.Kind && a.Program == facts.Program && a.Subject == facts.Subject &&
					a.Hash == facts.Hash && a.Writer == writer && a.Prompt == prompt
			}); i >= 0 {
				kept = append(kept, earlier[i])
				continue
			}
			name, text, err := w.Write(ctx, facts)
			if err != nil {
				return nil, nil, err
			}
			a := blueprint.Annotation{Kind: facts.Kind, Program: facts.Program, Subject: facts.Subject, Hash: facts.Hash,
				Name: name, Text: text, Writer: writer, Prompt: prompt}
			if err := facts.Check(a); err != nil {
				rejected = append(rejected, Rejection{Annotation: a, Reason: err.Error()})
				continue
			}
			kept = append(kept, a)
		}
	}
	return kept, rejected, nil
}
