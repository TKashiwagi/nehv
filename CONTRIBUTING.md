# Project Constraints

- The executable file built by this project must be named `configure`.
- Users must be able to launch the CLI using `./configure` in the top directory.
- Other binary names or aliases are prohibited.
- Refactor the code every time it grows by 100 lines.
- All documentation, comments, and code must be written in English.

These rules must be strictly followed in all phases of development, distribution, and operation.

# Running Tests

System tests are only executed in Linux/WSL environments and are skipped in other environments.

To run the tests, use the following command:

```sh
go test -v ./cmd/test/...
```

This will execute the relevant tests only in WSL environments and skip them in Windows environments.
If you have any other concerns about tests or functionality, please let us know.