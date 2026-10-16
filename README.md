# Copal

*Reading the code that still runs the books.*

Copal reads a legacy COBOL batch system and draws its blueprints: the layout of every record, the structure of every program, which fields each step reads and writes, the business rules buried in conditions, and the job network that ties the programs together through fixed-format files. All of it is derived deterministically from the source.

It is developed against **the First Bank of Pangaea**, a small but working COBOL batch bank in [pangaea/](pangaea/): accounts, transactions, a daily batch and a month-end run, compiled and executed with GnuCOBOL.

Copal does not translate COBOL into anything. It produces the understanding a translation needs, and later the tests a rewrite has to pass.

**Status:** nothing built yet. [docs/brief.md](docs/brief.md) has the plan and [docs/roadmap.md](docs/roadmap.md) the order of work.

## The names

Copal is young amber: resin that has hardened but not quite fossilized, and the traditional base of varnish for old wood. Said aloud, it is not far from COBOL.

Pangaea is the single supercontinent of the age of the dinosaurs, before it broke into the continents of today. A monolith, waiting to be split.
