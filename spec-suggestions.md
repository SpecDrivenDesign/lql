### X. Compiled Bytecode Signing and Verification

**Overview:**  
Implementations **MUST** support an optional process whereby DSL expressions can be compiled into a binary token stream (bytecode) that is cryptographically signed using RSA. This mechanism guarantees both the integrity and authenticity of the compiled bytecode.

**Requirements:**

- **Digital Signature Generation:**  
  When the compilation process is invoked with the appropriate flag (e.g. `-signed`), the compiler **MUST** compute an RSA signature over the serialized token stream.
    - The signature **MUST** be generated using an RSA private key provided via a PEM‑formatted file.
    - The resulting output bytecode **MUST** include a header magic value, a length field, the token stream, and the RSA signature appended at the end.

- **Verification During Execution:**  
  When executing compiled bytecode marked as signed, the engine **MUST** verify the signature using a corresponding RSA public key (also PEM‑formatted).
    - If the signature verification fails, execution **MUST** abort with an appropriate runtime error.

- **Interoperability:**  
  The details of the bytecode signing process (including the header, length encoding, and signature format) **MUST** be documented so that different implementations can interoperate when exchanging signed bytecode.

---

### X+1. Bytecode Serialization Format

**Overview:**  
In addition to textual DSL evaluation, the implementation **MAY** provide a compiler that serializes DSL expressions into a compact binary format (bytecode). This serialization is intended for later evaluation and (optionally) for secure distribution.

**Requirements:**

- **Header and Length Encoding:**  
  The serialized format **MUST** begin with a fixed header magic string (e.g. `"STOK"`) followed by a 4‑byte little‑endian encoded length field that indicates the size of the token stream.

- **Token Encoding:**  
  Each token **MUST** be encoded using a single byte to denote its type (according to a fixed mapping) followed by, when applicable, a length-prefixed literal string.

- **Extensibility:**  
  While the exact binary format is implementation‑defined, the specification **MAY** require that any two conforming implementations document their bytecode format in sufficient detail to allow for cross‑platform exchange if desired.

---

### X+2. Syntax Highlighting and Colorization

**Overview:**  
For improved readability in command‑line interfaces and development tools, implementations **MAY** provide syntax highlighting of DSL expressions.

**Requirements:**

- **Configurable Color Palettes:**  
  The DSL tool **MAY** support multiple themes (e.g. `mild`, `vivid`, `dracula`, and `solarized`). Users **MAY** select the desired theme via a command‑line flag (e.g. `-theme`).

- **Highlighting Targets:**  
  The colorization **MAY** be applied to:
    - Punctuation (brackets, commas, colons)
    - Literals (strings, numbers, booleans, null)
    - Operators
    - Identifiers and context references
    - Library and function names

- **Optional Enablement:**  
  Color output **MUST** be enabled only when the output device supports ANSI escapes or when an environment variable (e.g. `ENABLE_COLORS`) is set appropriately.

---

### X+3. Performance Benchmarking in the Test Runner

**Overview:**  
To assist developers in understanding and optimizing DSL expression performance, the test runner **MAY** include an optional benchmarking mode.

**Requirements:**

- **Iteration and Timing:**  
  When enabled (e.g. via a `--benchmark` flag), the test runner **MAY** evaluate DSL expressions multiple times (for example, 1000 iterations) and report:
    - The total elapsed time for the iterations.
    - The calculated operations per second.

- **Applicability:**  
  Benchmarking **MAY** be performed only on expressions that are free of errors and where the top‑level AST node is a function call or an expression with significant computation.

- **Output Reporting:**  
  Benchmark results **MAY** be displayed alongside the normal test result summary in either plain text or YAML output.

---

### X+4. Rich Test Runner Output and Reporting

**Overview:**  
Beyond simple pass/fail reporting, the DSL test runner **MUST** provide detailed output regarding each test case’s execution, including context, expected and actual results, and error details.

**Requirements:**

- **Output Formats:**  
  Implementations **MUST** support at least two output formats: plain text and YAML. The format can be selected using a command‑line flag (e.g. `--output`).

- **Detailed Test Information:**  
  For each test case, the report **MUST** include:
    - Test ID and description.
    - The DSL expression as interpreted by the parser (canonical form).
    - The provided context.
    - The expected result or error (if any).
    - The actual result or error message, including positional error information (line, column, and a snippet with a caret pointer).
    - Benchmark data (if benchmarking is enabled).

- **Colorized Output:**  
  When supported, the output **MAY** use ANSI color codes to visually distinguish passed, failed, and skipped tests.

- **Summary Statistics:**  
  A final summary **MUST** include counts of passed, failed, and skipped tests, as well as the total number of tests executed.

---

### X+5. PEM‑Based RSA Key Loading

**Overview:**  
As part of the security enhancements for bytecode signing and verification, the DSL implementation **MUST** load RSA keys from PEM‑formatted files.

**Requirements:**

- **Private Key Loading:**  
  The compiler **MUST** support a flag (e.g. `-private`) that specifies the path to an RSA private key in PEM format.
    - The key file **MUST** be validated to ensure it is a properly formatted PEM file containing an RSA private key.

- **Public Key Loading:**  
  Similarly, when executing signed bytecode, the engine **MUST** support a flag (e.g. `-public`) to specify the RSA public key file (in PEM format) used to verify the signature.

- **Error Handling:**  
  If the provided key file is missing, unreadable, or does not contain a valid RSA key, the DSL implementation **MUST** produce a clear runtime error message.

---

### X+6. Interactive REPL and Context Input

**Overview:**  
In addition to batch evaluation via test cases or compiled bytecode, the DSL implementation **MUST** provide an interactive Read–Eval–Print Loop (REPL) mode.

**Requirements:**

- **Expression Input:**  
  The REPL **MUST** allow users to supply a DSL expression (via a command‑line flag such as `-expr`) that remains active for multiple evaluations.

- **Context Data:**  
  The REPL **MAY** accept context input interactively from standard input.
    - The context **MUST** be provided in either JSON or YAML format.
    - If piped input is detected, the REPL **MUST** automatically parse the context accordingly.

- **Error Feedback:**  
  The REPL **MUST** display any parsing or evaluation errors in real time, using the same error reporting format (and optional colorization) as the test runner.

- **Exit Conditions:**  
  The REPL **MUST** allow the user to exit by providing an empty line or a specific exit command.

---

### X+7. Enhanced Lexer Functionality and Comment Handling

**Overview:**  
The DSL lexer implementation **MUST** include enhancements for robust tokenization, including support for comments and comprehensive escape sequence handling.

**Requirements:**

- **Whitespace and Comment Skipping:**  
  The lexer **MUST** ignore spaces, tabs, and newlines between tokens.
    - Lines beginning with a `#` character (ignoring preceding whitespace) **MUST** be treated as comments and skipped entirely.

- **Escape Sequence Handling:**  
  String literals **MUST** support standard JSON‑like escape sequences (e.g. `\n`, `\t`, `\\`, `\"`, `\'`) as well as Unicode escapes (e.g. `\uXXXX`), with proper validation.
    - Malformed escape sequences **MUST** trigger a lexical error, reporting the correct line and column.

- **Accurate Positional Tracking:**  
  The lexer **MUST** track line and column positions precisely so that any lexical errors can be reported with exact source locations.

- **Robust Numeric Parsing:**  
  The lexer **MUST** handle numeric literals with optional leading `+` or `-` signs, decimal points, and scientific notation (e.g. `1e10`), rejecting malformed numbers (e.g. `12..3`) with an appropriate lexical error.
