# Build in a stock Go builder container
FROM fuzzers/go-fuzz:1.2.0 as builder

ADD . /ferretdb
WORKDIR /ferretdb

RUN cd internal/util/testutil && GO_ENABLED=0 go-fuzz-build -libfuzzer -o fuzz.a .
RUN cd internal/util/testutil && clang -fsanitize=fuzzer fuzz.a -o fuzz.libfuzzer

FROM ubuntu:20.04
COPY --from=builder /ferretdb/internal/util/testutil/fuzz.libfuzzer /