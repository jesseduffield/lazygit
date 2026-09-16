# Zsh completion for the `e2e` and `e2e-cli` recipes in lazygit's justfile.
#
# These recipes take integration-test names (e.g. submodule/reset). This makes
# `just e2e <Tab>` complete them from pkg/integration/tests/. To enable it, add
# the following to your ~/.zshrc, *after* the line that runs `compinit`:
#
#     source /path/to/lazygit/scripts/just_e2e_completion.zsh
#
# It is a no-op when `just` isn't installed, and only kicks in inside a project
# that has a justfile and a pkg/integration/tests/ directory, so it is harmless
# to source unconditionally.

(( $+commands[just] )) || return 0

# just's own completion is clap-dynamic and has no hook for completing a
# recipe's arguments, so we wrap it: handle the e2e recipes ourselves and
# delegate everything else (recipe names, flags, ...) to just's completer.
source <(JUST_COMPLETE=zsh just)   # defines _clap_dynamic_completer_just

_just_lazygit_e2e() {
    if (( CURRENT > 2 )); then
        case ${words[2]} in
        e2e | e2e-cli)
            # Find the justfile's directory, then complete the integration
            # tests under pkg/integration/tests/ relative to it.
            local dir=$PWD testdir=
            while [[ $dir != / ]]; do
                if [[ -e $dir/justfile || -e $dir/.justfile || -e $dir/Justfile ]]; then
                    testdir=$dir/pkg/integration/tests
                    break
                fi
                dir=${dir:h}
            done
            if [[ -d $testdir ]]; then
                # A test's name is its path under pkg/integration/tests/ without
                # the .go extension, e.g. submodule/reset. Build that list, then
                # let _multi_parts complete it one "/"-separated segment at a
                # time, so an empty <Tab> offers only categories.
                local -a tests
                tests=($testdir/**/*.go(.N:r))  # strip the .go extension
                tests=(${tests#$testdir/})      # make relative to the tests dir
                tests=(${(M)tests:#*/*})         # keep category/name (drop top-level helpers)
                tests=(${tests:#shared/*})       # drop the cross-directory shared package
                tests=(${tests:#*/shared})       # drop per-category shared.go helpers
                local expl
                _wanted tests expl 'integration test' _multi_parts / tests
                return
            fi
            ;;
        esac
    fi

    _clap_dynamic_completer_just "$@"
}

compdef _just_lazygit_e2e just   # bind last so this wins over the default
