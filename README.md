reppl
=====

The simplest possible REPeatable project PipeLine.

This is a very basic system for managing chains and pipelines of formulas, pushing updates through the whole picture.
This is not the final system that will be used by repeatr, though it may have some similar traits.

(It's also feels a bit like a [REPL](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop)
for repeatr formulas, since it evaluates and prints after each step, hence the name.)


Example use
-----------

```
reppl init                  ## makes a .reppl history file in the cwd
reppl name k v              ## adds a name mapping to the .reppl file: look up `"k"` and get `hash: "v"`
# human: writes some formulas, using names instead of hashes in the inputs.  commits to git.
reppl formula repin f.frm   ## looks up the names in the .reppl files, and adds in the hashes, and writes the file back
# human: may choose to git commit again, to save the pinned formula -- it's now runnable by regular `repeatr run` if desired
reppl eval f.frm            ## runs the formula, adds any named output results to the .reppl file, and also a map of `"$doc_hash: $results"`.
reppl eval f.frm            ## *doesn't* run the formula, because it can already find the doc_hash in the .reppl file!
reppl eval -f f.frm         ## forces running the formula regardless of remembered results
```

Big picture intention: for any project with a bunch of isolated steps, some of which feed their outputs into later steps, your build script can look something like this:

```
reppl eval subproj/basis.formula
reppl eval subproj/comp-green.formula
reppl eval subproj/green-test.formula
reppl eval mainproj.formula
reppl eval shrinkwrap.formula
```

and that is it.  Just list everything; yes, even the base image builds.

Every time you want a rebuild, you just always run all that.
You don't need to think about it, because the machine does: it intelligently no-ops itself and thus runs fast and clean.
At the same time, the magic of content-addressable formulas means if you make changes, they're guaranteed to be included correctly.

### Integrating with [you-name-it]

Integrate reppl with other input sources easily with the `reppl name` command.

For example, let's say you have a project where most of the changes are in git repo, and you're running the `reppl` tool inside the repo directory.
You can push the current git commit checked out into the reppl named wares by simply running:

```
reppl name my-project-checkout "$(git rev-parse HEAD)"
```

And voilà; any of your formulas with `"my-project-checkout"` as an input name will be automatically rebuilt with the up-to-date value when you next run `reppl eval`.


### What doesn't it do?

Reppl does not have an automatic dependency order resolver.  We feel that it represents unnecessarily complexity for the problem reppl is solving.  Just list your tasks in a valid order; the no-op for any task that's been run before is instantaneous.

Reppl is not meant to replace your build tool (especially your fancy build tool that does have an automatic dependency order resolver).
Think of Reppl as the tool you want to reach for when packaging releases, segmenting out builds of major separate components of software, and so on.
Use `make`, or `bazel`, or the `go` tool, or `rake`, or `ant`, or *whatever* your tool-of-choice is for managing language-specific build goals.
Reppl *may* be the One Ring to Rule Them All at an empire-management sort of level, but is not particularly targetted at local granular task management.



Why is this not the final system?
---------------------------------

This disregards all hard problems about global versioning, name reservation, signing, forward-locking, believably-historically-valid checking for older values, and so on.

In other words, it disregards a *lot*, and tackling some of those subjects is *not optional* for building a *public, distributable* package manager.

(Then again, half the so-called package managers in the wild today don't implement any of those systems, aside from some name reservation implemented in a skeezy web form barely strapped to a postgres database, so, who am I to say what's suitable.)



License
-------

Apache v2
