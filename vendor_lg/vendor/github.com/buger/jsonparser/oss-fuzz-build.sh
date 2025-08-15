#!/bin/bash -eu

git clone https://github.com/dvyukov/go-fuzz-corpus
zip corpus.zip go-fuzz-corpus/json/corpus/*

cp corpus.zip $OUT/fuzzparsestring_seed_corpus.zip
compile_go_fuzzer github.com/buger/jsonparser FuzzParseString fuzzparsestring

cp corpus.zip $OUT/fuzzeachkey_seed_corpus.zip
compile_go_fuzzer github.com/buger/jsonparser FuzzEachKey fuzzeachkey

cp corpus.zip $OUT/fuzzdelete_seed_corpus.zip
compile_go_fuzzer github.com/buger/jsonparser FuzzDelete fuzzdelete

cp corpus.zip $OUT/fuzzset_seed_corpus.zip
compile_go_fuzzer github.com/buger/jsonparser FuzzSet fuzzset

cp corpus.zip $OUT/fuzzobjecteach_seed_corpus.zip
compile_go_fuzzer github.com/buger/jsonparser FuzzObjectEach fuzzobjecteach

cp corpus.zip $OUT/fuzzparsefloat_seed_corpus.zip
compile_go_fuzzer github.com/buger/jsonparser FuzzParseFloat fuzzparsefloat

cp corpus.zip $OUT/fuzzparseint_seed_corpus.zip
compile_go_fuzzer github.com/buger/jsonparser FuzzParseInt fuzzparseint

cp corpus.zip $OUT/fuzzparsebool_seed_corpus.zip
compile_go_fuzzer github.com/buger/jsonparser FuzzParseBool fuzzparsebool

cp corpus.zip $OUT/fuzztokenstart_seed_corpus.zip
compile_go_fuzzer github.com/buger/jsonparser FuzzTokenStart fuzztokenstart

cp corpus.zip $OUT/fuzzgetstring_seed_corpus.zip
compile_go_fuzzer github.com/buger/jsonparser FuzzGetString fuzzgetstring

cp corpus.zip $OUT/fuzzgetfloat_seed_corpus.zip
compile_go_fuzzer github.com/buger/jsonparser FuzzGetFloat fuzzgetfloat

cp corpus.zip $OUT/fuzzgetint_seed_corpus.zip
compile_go_fuzzer github.com/buger/jsonparser FuzzGetInt fuzzgetint

cp corpus.zip $OUT/fuzzgetboolean_seed_corpus.zip
compile_go_fuzzer github.com/buger/jsonparser FuzzGetBoolean fuzzgetboolean

cp corpus.zip $OUT/fuzzgetunsafestring_seed_corpus.zip
compile_go_fuzzer github.com/buger/jsonparser FuzzGetUnsafeString fuzzgetunsafestring

