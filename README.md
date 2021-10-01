# PaaS IO Assignment

Report network IO statistics.

You are writing a PaaS, and you need a way to bill customers based
on network and filesystem usage.

Create a wrapper for network connections and files that can report IO
statistics. The wrapper must report:

- The total number of bytes read/written.
- The total number of read/write operations.

## Running the tests

To run the tests, run the command `go test` from within the assignment directory.

If the test suite contains benchmarks, you can run these with the `--bench` and `--benchmem`
flags:

    go test -v --bench . --benchmem

## Submission Instructions
Email your solution as a public Github repository link OR a .zip package to engineeringjobs@cirrusmd.com, include 'Staff Engineer Take-Home' in the subject.

## Questions

If you have questions about the instructions, please ask. We want you to be successful. If you have a question about how
to handle something that wasn't specifically addressed, make a decision and feel free to call it out in your own project
readme or submission email with your reasoning behind your decision. No right or wrong answers for these types of things.
