The JSON files in this directory are machine-generated; please do not edit.

Translating lazygit happens at https://crowdin.com/project/lazygit/.

# Updating translations from Crowdin

We regularly need to pull changes from Crowdin and integrate them here. This is
done by downloading a zip file of the translations from Crowdin, unzipping it,
and calling `scripts/update_language_files.sh` with the unzipped directory as an
argument.

# Uploading the English file to Crowdin

The English version of all the texts is still maintained in
`pkg/i18n/english.go`; it needs to be uploaded to Crowdin regularly. To do this,
call `go run cmd/i18n/main.go`; this will create an unversioned file `en.json`
in the root of the repository. Upload this to
`https://crowdin.com/project/lazygit/sources/files` and delete it from the
working copy again.
