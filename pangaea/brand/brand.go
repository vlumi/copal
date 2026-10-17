// Package brand holds the First Bank of Pangaea's marks, for the console
// and the documentation to show: the seal, the globe from its center as an
// icon, and the bank's unofficial mascot.
package brand

import "embed"

// Files are the marks: seal.png, the round seal with the bank's name and
// motto; icon.png, the globe and its Dimetrodon alone, for small sizes;
// and rex.png, a Tyrannosaurus in a suit with a briefcase, who is from a
// later era than the bank's name admits and appears where something has
// gone wrong.
//
//go:embed *.png
var Files embed.FS
