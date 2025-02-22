# Logical Query Language (LQL) – Public Interface Usage Guide

**Module Path:**  
`github.com/SpecDrivenDesign/lql`

LQL is a domain‑specific language for constructing logical queries over JSON‑like data. It is implemented as a Go module with several interrelated packages. This guide explains how to parse and evaluate DSL expressions, as well as how to compile and execute bytecode.

---

## 1. Parsing and Abstract Syntax Trees (AST)

**Package:**  
`github.com/SpecDrivenDesign/lql/pkg/parser`

**Purpose:**  
The parser converts a DSL expression (provided as a string) into an immutable Abstract Syntax Tree (AST). The AST nodes are defined in  
`github.com/SpecDrivenDesign/lql/pkg/ast/expressions`  
and all implement the following public interface methods:
- **Eval(ctx map[string]interface{}, env IEnvironment) (interface{}, error):**  
  Evaluates the expression using the provided context and execution environment.
- **Pos() (line, column int):**  
  Returns the source location for error reporting.
- **String() string:**  
  Returns a canonical string representation of the expression.

**Usage Workflow:**  
1. Create a lexer from your DSL string using  
   `github.com/SpecDrivenDesign/lql/pkg/lexer.NewLexer(expression)`.
2. Instantiate a parser with the lexer:  
   `p, err := parser.NewParser(lexer)`.
3. Parse your expression to obtain an AST:  
   `ast, err := p.ParseExpression()`.

The resulting AST is immutable and reusable for multiple evaluations.

---

## 2. Execution Environment and Standard Libraries

**Package:**  
`github.com/SpecDrivenDesign/lql/pkg/env`

**Purpose:**  
The execution environment holds the available namespaced libraries. To evaluate expressions, create an environment that is pre‑loaded with the standard libraries.

**Creating an Environment:**  
```go
environment := env.NewEnvironment()
```

**Standard Libraries Included:**  
- **time:** Provides functions such as `time.now()`, `time.parse()`, `time.add()`, etc.
- **math:** Includes arithmetic and aggregation functions like `math.abs()`, `math.sqrt()`, `math.sum()`, etc.
- **string:** Offers functions for string manipulation such as `string.concat()`, `string.toLower()`, etc.
- **regex:** Supports regular expression operations with functions like `regex.match()`, `regex.replace()`, etc.
- **array:** Contains functions for array handling, e.g., `array.contains()`, `array.sort()`, etc.
- **cond:** Provides conditional operations like `cond.ifExpr()` and `cond.coalesce()`.
- **type:** Offers type predicates and conversion functions such as `type.int()`, `type.string()`, etc.

Each library is invoked by its namespace when an expression is evaluated.

---

## 3. Evaluating DSL Expressions

**Steps to Evaluate an Expression:**

1. **Import Required Packages:**
   ```go
   import (
       "github.com/SpecDrivenDesign/lql/pkg/lexer"
       "github.com/SpecDrivenDesign/lql/pkg/parser"
       "github.com/SpecDrivenDesign/lql/pkg/env"
   )
   ```

2. **Parse the Expression:**
   Create a lexer and parser for your DSL expression:
   ```go
   lex := lexer.NewLexer(`$sensor.reading > 100 AND $user.age >= 18 AND $user.country == "US"`)
   p, err := parser.NewParser(lex)
   if err != nil {
       // Handle parser error.
   }
   ast, err := p.ParseExpression()
   if err != nil {
       // Handle parse error.
   }
   ```
   
3. **Evaluate the AST:**
   Prepare your context (a map containing your data) and create the environment:
   ```go
   context := map[string]interface{}{
       "sensor": map[string]interface{}{
           "reading": 123.45,
           "status":  "active",
       },
       "user": map[string]interface{}{
           "age":     25,
           "country": "US",
           "name":    "Alice",
       },
   }
   environment := env.NewEnvironment()
   result, err := ast.Eval(context, environment)
   if err != nil {
       // Handle evaluation error.
   }
   fmt.Printf("Evaluation Result: %v\n", result)
   ```
   In this example, the DSL expression checks that the sensor reading is greater than 100, the user is at least 18 years old, and the user’s country is `"US"`. The evaluation should yield `true`.

---

## 4. Compiling and Executing Bytecode

**Package:**  
`github.com/SpecDrivenDesign/lql/pkg/bytecode`

**Purpose:**  
The bytecode package allows you to compile DSL expressions into a compact binary format and execute them using the `ByteCodeReader`.

**Compiling DSL Expressions to Bytecode:**
- The lexer provides methods to export tokens as a binary stream:
  - `ExportTokens()` returns the token stream as a `[]byte`.
  - Optionally, `ExportTokensSigned()` can produce signed bytecode using an RSA private key.

The compiled bytecode includes:
- A header magic (e.g., `"STOK"`),
- A 4‑byte little‑endian length field,
- The token stream,
- Optionally, an RSA signature.

**Executing Compiled Bytecode:**
1. **Import the Bytecode Package:**
   ```go
   import "github.com/SpecDrivenDesign/lql/pkg/bytecode"
   ```

2. **Create a ByteCodeReader and Parse the Bytecode:**
   Suppose you have compiled bytecode stored in a variable `data`:
   ```go
   reader := bytecode.NewByteCodeReader(data)
   p, err := parser.NewParser(reader)
   if err != nil {
       // Handle parser error.
   }
   ast, err := p.ParseExpression()
   if err != nil {
       // Handle parse error.
   }
   ```
   
3. **Evaluate the Expression from Bytecode:**
   Use the same evaluation steps as above:
   ```go
   context := map[string]interface{}{
       "sensor": map[string]interface{}{
           "reading": 123.45,
       },
       "user": map[string]interface{}{
           "age":     25,
           "country": "US",
       },
   }
   environment := env.NewEnvironment()
   result, err := ast.Eval(context, environment)
   if err != nil {
       // Handle evaluation error.
   }
   fmt.Printf("Bytecode Evaluation Result: %v\n", result)
   ```
   
**Example Program – Compiling and Executing Bytecode with Real Context Variables:**

```go
package main

import (
    "fmt"
    "log"

    "github.com/SpecDrivenDesign/lql/pkg/lexer"
    "github.com/SpecDrivenDesign/lql/pkg/bytecode"
    "github.com/SpecDrivenDesign/lql/pkg/parser"
    "github.com/SpecDrivenDesign/lql/pkg/env"
)

func main() {
    // DSL expression that uses real context variables.
    expression := `$sensor.reading > 100 AND $user.age >= 18 AND $user.country == "US"`

    // Compile the expression into bytecode.
    lex := lexer.NewLexer(expression)
    tokenData, err := lex.ExportTokens()
    if err != nil {
        log.Fatalf("Error exporting tokens: %v", err)
    }

    // Create a ByteCodeReader from the compiled token data.
    reader := bytecode.NewByteCodeReader(tokenData)

    // Parse the bytecode to obtain an AST.
    p, err := parser.NewParser(reader)
    if err != nil {
        log.Fatalf("Parser error: %v", err)
    }
    ast, err := p.ParseExpression()
    if err != nil {
        log.Fatalf("Parse error: %v", err)
    }

    // Define a real evaluation context.
    context := map[string]interface{}{
        "sensor": map[string]interface{}{
            "reading": 123.45,
            "status":  "active",
        },
        "user": map[string]interface{}{
            "age":     25,
            "country": "US",
            "name":    "Alice",
        },
    }

    // Create the execution environment with standard libraries.
    environment := env.NewEnvironment()

    // Evaluate the expression parsed from bytecode.
    result, err := ast.Eval(context, environment)
    if err != nil {
        log.Fatalf("Evaluation error: %v", err)
    }

    fmt.Printf("Bytecode Evaluation Result: %v\n", result)
}
```
