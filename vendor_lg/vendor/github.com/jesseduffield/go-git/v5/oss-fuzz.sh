#!/bin/bash -eu
# Copyright 2023 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
################################################################################


go mod download
go get github.com/AdamKorcz/go-118-fuzz-build/testing

if [ "$SANITIZER" != "coverage" ]; then
    sed -i '/func (s \*DecoderSuite) TestDecode(/,/^}/ s/^/\/\//' plumbing/format/config/decoder_test.go
    sed -n '35,$p' plumbing/format/packfile/common_test.go >> plumbing/format/packfile/delta_test.go
    sed -n '20,53p' plumbing/object/object_test.go >> plumbing/object/tree_test.go
    sed -i 's|func Test|// func Test|' plumbing/transport/common_test.go
fi

compile_native_go_fuzzer $(pwd)/internal/revision                       FuzzParser              fuzz_parser
compile_native_go_fuzzer $(pwd)/plumbing/format/config                  FuzzDecoder             fuzz_decoder_config
compile_native_go_fuzzer $(pwd)/plumbing/format/packfile                FuzzPatchDelta          fuzz_patch_delta
compile_native_go_fuzzer $(pwd)/plumbing/object                         FuzzParseSignedBytes    fuzz_parse_signed_bytes
compile_native_go_fuzzer $(pwd)/plumbing/object                         FuzzDecode              fuzz_decode
compile_native_go_fuzzer $(pwd)/plumbing/protocol/packp                 FuzzDecoder             fuzz_decoder_packp
compile_native_go_fuzzer $(pwd)/plumbing/transport                      FuzzNewEndpoint         fuzz_new_endpoint
