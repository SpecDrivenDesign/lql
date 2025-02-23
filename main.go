package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/SpecDrivenDesign/lql/pkg/ast/expressions"
	"github.com/SpecDrivenDesign/lql/pkg/bytecode"
	"github.com/SpecDrivenDesign/lql/pkg/env"
	"github.com/SpecDrivenDesign/lql/pkg/errors"
	"github.com/SpecDrivenDesign/lql/pkg/lexer"
	"github.com/SpecDrivenDesign/lql/pkg/parser"
	"github.com/SpecDrivenDesign/lql/pkg/signing"
	"github.com/SpecDrivenDesign/lql/pkg/testing"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"strings"
)

// Color constants
const (
	colorReset   = "\033[0m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorGreen   = "\033[32m"
	colorRed     = "\033[31m"
	colorYellow  = "\033[33m"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Subcommand required: test, compile, exec, repl, validate, or highlight")
		fmt.Println("Usage:")
		fmt.Println("  lql test [--test-file=testcases.yml] [--fail-fast] [--verbose] [--output text|yaml]")
		fmt.Println("  lql compile -expr \"<expression>\" -out <outfile> [-signed -private <private.pem>]")
		fmt.Println("  lql exec -in <infile> [-signed -public <public.pem>]")
		fmt.Println("  lql repl -expr \"<expression>\" [-format json|yaml]")
		fmt.Println("  lql validate -expr \"<expression>\" | -in <file>")
		fmt.Println("  lql highlight -expr \"<expression>\" [-theme mild|vivid|dracula|solarized]")
		os.Exit(1)
	}

	subcommand := os.Args[1]
	switch subcommand {
	case "test":
		runTestCmd()
	case "compile":
		runCompileCmd()
	case "exec":
		runExecCmd()
	case "repl":
		runReplCmd()
	case "validate":
		runValidateCmd()
	case "highlight":
		runHighlightCmd()
	case "export-contexts":
		runExportContextsCmd()
	default:
		fmt.Printf("Unknown subcommand: %s\n", subcommand)
		os.Exit(1)
	}
}

func runTestCmd() {
	testCmd := flag.NewFlagSet("test", flag.ExitOnError)
	helpPtr := testCmd.Bool("help", false, "Show help message")
	failFastPtr := testCmd.Bool("fail-fast", false, "Stop on first failure")
	verbosePtr := testCmd.Bool("verbose", false, "Verbose output")
	outputFormatPtr := testCmd.String("output", "text", "Output format: text or yaml")
	testFile := testCmd.String("test-file", "testcases.yml", "YAML file containing test cases")
	benchmarkPtr := testCmd.Bool("benchmark", false, "Run each expression 1000 times and print benchmark info (only for function calls)")
	if err := testCmd.Parse(os.Args[2:]); err != nil {
		fmt.Printf("Error reading command line args: %v\n", err)
		os.Exit(1)
	}
	if *helpPtr {
		testCmd.Usage()
		os.Exit(0)
	}

	data, err := os.ReadFile(*testFile)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	var testCases []testing.TestCase
	err = yaml.Unmarshal(data, &testCases)
	if err != nil {
		log.Fatalf("Error parsing YAML: %s", err)
	}

	env := env.NewEnvironment()
	suiteResult := testing.RunTests(testCases, env, *failFastPtr, *benchmarkPtr)

	// Output printing remains here.
	if strings.ToLower(*outputFormatPtr) == "yaml" {
		renderYAMLOutput(suiteResult)
	} else {
		renderTextOutput(suiteResult, *verbosePtr)
	}

	if suiteResult.Failed > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}

func runCompileCmd() {
	compileCmd := flag.NewFlagSet("compile", flag.ExitOnError)
	expr := compileCmd.String("expr", "", "DSL expression to compile")
	inFile := compileCmd.String("in", "", "File containing a DSL expression to compile")
	outFile := compileCmd.String("out", "", "Output filename for compiled byteCode")
	signed := compileCmd.Bool("signed", false, "Whether to sign the compiled byteCode")
	privateKeyFile := compileCmd.String("private", "private.pem", "Path to RSA private key for signing (required if -signed is true)")

	if err := compileCmd.Parse(os.Args[2:]); err != nil {
		fmt.Printf("Error reading command line args: %v\n", err)
		os.Exit(1)
	}
	var expression string
	if *inFile != "" {
		data, err := os.ReadFile(*inFile)
		if err != nil {
			log.Fatalf("Error reading expression file: %v", err)
		}
		expression = strings.TrimSpace(string(data))
	} else if *expr != "" {
		expression = *expr
	} else {
		fmt.Println("Either -expr or -in flag must be provided.")
		compileCmd.Usage()
		os.Exit(1)
	}

	if *outFile == "" {
		fmt.Println("The -out flag is required.")
		compileCmd.Usage()
		os.Exit(1)
	}

	lex := lexer.NewLexer(expression)
	var byteCode []byte
	var err error
	if *signed {
		if *privateKeyFile == "" {
			fmt.Println("Private key file must be provided when -signed is true.")
			compileCmd.Usage()
			os.Exit(1)
		}
		privateKey, err := signing.LoadPrivateKey(*privateKeyFile)
		if err != nil {
			log.Fatalf("Error loading private key: %v", err)
		}
		byteCode, err = lex.ExportTokensSigned(privateKey)
		if err != nil {
			log.Fatalf("Error exporting signed tokens: %v", err)
		}
	} else {
		byteCode, err = lex.ExportTokens()
		if err != nil {
			log.Fatalf("Error exporting tokens: %v", err)
		}
	}

	err = os.WriteFile(*outFile, byteCode, 0600)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}
	fmt.Printf("Compilation successful. Bytecode written to %s\n", *outFile)
}

func runExecCmd() {
	execCmd := flag.NewFlagSet("exec", flag.ExitOnError)
	inFile := execCmd.String("in", "", "Input filename of compiled bytecode")
	expr := execCmd.String("expr", "", "Raw DSL expression to execute")
	signed := execCmd.Bool("signed", false, "Indicate if the bytecode is signed (only used with -in)")
	publicKeyFile := execCmd.String("public", "", "Path to RSA public key for signature verification (required if -signed is true)")
	contextFormat := execCmd.String("format", "yaml", "Format of context input from stdin: json or yaml")
	if err := execCmd.Parse(os.Args[2:]); err != nil {
		fmt.Printf("Error reading command line args: %v\n", err)
		os.Exit(1)
	}
	contextData, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading context from stdin: %v", err)
	}
	var ctx map[string]interface{}
	if len(strings.TrimSpace(string(contextData))) > 0 {
		if strings.ToLower(*contextFormat) == "json" {
			err = json.Unmarshal(contextData, &ctx)
		} else {
			err = yaml.Unmarshal(contextData, &ctx)
		}
		if err != nil {
			log.Fatalf("Error parsing context: %v", err)
		}
	} else {
		ctx = make(map[string]interface{})
	}

	if *expr == "" && *inFile == "" {
		fmt.Println("Either -expr or -in flag must be provided.")
		execCmd.Usage()
		os.Exit(1)
	}

	if *expr != "" {
		lex := lexer.NewLexer(*expr)
		p, err := parser.NewParser(lex)
		if err != nil {
			log.Fatalf("Error creating p: %v", err)
		}
		ast, err := p.ParseExpression()
		if err != nil {
			log.Fatalf("Error parsing expression: %v", err)
		}
		env := env.NewEnvironment()
		result, err := ast.Eval(ctx, env)
		if err != nil {
			log.Fatalf("Error executing expression: %v", err)
		}
		fmt.Printf("Execution result: %v\n", result)
		return
	}

	data, err := os.ReadFile(*inFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	var tokenStream parser.TokenStream
	if *signed {
		if *publicKeyFile == "" {
			fmt.Println("Public key file must be provided when -signed is true.")
			execCmd.Usage()
			os.Exit(1)
		}
		pubKey, err := signing.LoadPublicKey(*publicKeyFile)
		if err != nil {
			log.Fatalf("Error loading public key: %v", err)
		}
		tokenStream, err = bytecode.NewByteCodeReaderFromSignedData(data, pubKey)
		if err != nil {
			log.Fatalf("Error verifying signed bytecode: %v", err)
		}
	} else {
		tokenStream = bytecode.NewByteCodeReader(data)
	}

	p, err := parser.NewParser(tokenStream)
	if err != nil {
		log.Fatalf("Error creating p: %v", err)
	}
	ast, err := p.ParseExpression()
	if err != nil {
		log.Fatalf("Error parsing expression from bytecode: %v", err)
	}
	env := env.NewEnvironment()
	result, err := ast.Eval(ctx, env)
	if err != nil {
		log.Fatalf("Error executing bytecode: %v", err)
	}
	fmt.Printf("Execution result: %v\n", result)
}

func runReplCmd() {
	replCmd := flag.NewFlagSet("repl", flag.ExitOnError)
	expr := replCmd.String("expr", "", "DSL expression to evaluate in REPL mode")
	if err := replCmd.Parse(os.Args[2:]); err != nil {
		fmt.Printf("Error reading command line args: %v\n", err)
		os.Exit(1)
	}
	if *expr == "" {
		fmt.Println("The -expr flag is required in repl mode.")
		replCmd.Usage()
		os.Exit(1)
	}

	lex := lexer.NewLexer(*expr)
	p, err := parser.NewParser(lex)
	if err != nil {
		log.Fatalf("Error creating p: %v", err)
	}
	ast, err := p.ParseExpression()
	if err != nil {
		log.Fatalf("Error parsing expression: %v", err)
	}
	env := env.NewEnvironment()

	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Fatalf("Error stating stdin: %v", err)
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) == "" {
				continue
			}
			var ctx map[string]interface{}
			if json.Unmarshal([]byte(line), &ctx) != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error parsing context: %v\n", err)
				continue
			}
			result, err := ast.Eval(ctx, env)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error executing expression: %v\n", err)
				continue
			}
			fmt.Println(result)
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading from stdin: %v", err)
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter context (empty line to exit): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("\nExiting REPL.")
				break
			}
			input = strings.TrimSpace(input)
			if input == "" {
				fmt.Println("Exiting REPL.")
				break
			}
			var ctx map[string]interface{}
			if json.Unmarshal([]byte(input), &ctx) != nil {
				fmt.Printf("Error parsing context: %v\n", err)
				continue
			}
			result, err := ast.Eval(ctx, env)
			if err != nil {
				fmt.Printf("Error executing expression: %v\n", err)
				continue
			}
			fmt.Printf("%v\n", result)
		}
	}
}

func runValidateCmd() {
	validateCmd := flag.NewFlagSet("validate", flag.ExitOnError)
	expr := validateCmd.String("expr", "", "DSL expression to validate")
	inFile := validateCmd.String("in", "", "File containing a DSL expression to validate")
	if err := validateCmd.Parse(os.Args[2:]); err != nil {
		fmt.Printf("Error reading command line args: %v\n", err)
		os.Exit(1)
	}

	var expression string
	if *inFile != "" {
		data, err := os.ReadFile(*inFile)
		if err != nil {
			fmt.Printf("Error reading expression file: %v\n", err)
			os.Exit(1)
		}
		expression = strings.TrimSpace(string(data))
	} else if *expr != "" {
		expression = *expr
	} else {
		fmt.Println("Either -expr or -in flag must be provided.")
		validateCmd.Usage()
		os.Exit(1)
	}

	lex := lexer.NewLexer(expression)
	p, err := parser.NewParser(lex)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	_, err = p.ParseExpression()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func renderTextOutput(suite testing.TestSuiteResult, verbose bool) {
	for _, res := range suite.TestResults {
		if !verbose && res.Status == "PASSED" && res.BenchmarkTime == "" {
			continue
		}
		if res.Status == "SKIPPED" {
			continue
		}
		fmt.Printf("%s[Test #%d] %s%s\n", colorBlue, res.TestID, res.Description, colorReset)
		fmt.Printf("    Expression : %s\n", res.Expression)
		fmt.Printf("    Context    : %v\n", res.Context)
		if res.ExpectedError != "" || res.ActualError != nil {
			if res.ExpectedResult != nil {
				fmt.Printf("    Expected   : %v\n", res.ExpectedResult)
				fmt.Printf("    Actual     : %v\n", res.ActualResult)
			}
			if res.ExpectedError != "" {
				fmt.Printf("    Expected Error Message: %s: %s\n", res.ExpectedError, res.ExpectedErrorMessage)
			}
			fmt.Printf("    Actual Error Message  : %v\n", res.ActualError)
		} else {
			fmt.Printf("    Expected   : %v\n", res.ExpectedResult)
			fmt.Printf("    Actual     : %v\n", res.ActualResult)
		}
		if res.BenchmarkTime != "" {
			fmt.Printf("    Benchmark  : %s (%0.2f ops/sec)\n", res.BenchmarkTime, res.BenchmarkOpsSec)
		}
		if res.ActualError != nil && res.Status != "PASSED" {
			if res.ErrLine > 0 && res.ErrColumn > 0 {
				fmt.Println(errors.GetErrorContext(res.Expression, res.ErrLine, res.ErrColumn, true))
			}
		}
		statusColor := ""
		switch res.Status {
		case "PASSED":
			statusColor = colorGreen
		case "FAILED":
			statusColor = colorRed
		case "SKIPPED":
			statusColor = colorYellow
		}
		fmt.Printf("    Status     : %s%s%s\n\n", statusColor, res.Status, colorReset)
	}
	fmt.Println("==============================================")
	fmt.Println("Test Suite Completed")
	fmt.Printf("  %sPASSED  %s: %d\n  %sSKIPPED %s: %d\n  %sFAILED  %s: %d\n  TOTAL   : %d\n",
		colorGreen, colorReset, suite.Passed,
		colorYellow, colorReset, suite.Skipped,
		colorRed, colorReset, suite.Failed,
		suite.Total)
	fmt.Println("==============================================")
}

func renderYAMLOutput(suite testing.TestSuiteResult) {
	out, err := yaml.Marshal(suite)
	if err != nil {
		log.Fatalf("Error marshaling YAML: %s", err)
	}
	fmt.Println(string(out))
}

func runHighlightCmd() {
	highlightCmd := flag.NewFlagSet("highlight", flag.ExitOnError)
	exprPtr := highlightCmd.String("expr", "", "Expression to highlight")
	themePtr := highlightCmd.String("theme", "mild", "Color theme: mild|vivid|dracula|solarized")

	if err := highlightCmd.Parse(os.Args[2:]); err != nil {
		fmt.Printf("Error reading command line args: %v\n", err)
		os.Exit(1)
	}

	if *exprPtr == "" {
		fmt.Println("The -expr flag is required.")
		highlightCmd.Usage()
		os.Exit(1)
	}

	expressions.ColorEnabled = true

	// 1) Parse the user expression into an AST.
	lex := lexer.NewLexer(*exprPtr)
	p, err := parser.NewParser(lex)
	if err != nil {
		log.Fatalf("Error creating parser: %v", err)
	}
	ast, err := p.ParseExpression()
	if err != nil {
		log.Fatalf("Error parsing expression: %v", err)
	}

	// 3) Apply the chosen color theme.
	switch strings.ToLower(*themePtr) {
	case "mild":
		expressions.ApplyMildPalette()
	case "vivid":
		expressions.ApplyVividPalette()
	case "dracula":
		expressions.ApplyDraculaPalette()
	case "solarized":
		expressions.ApplySolarizedPalette()
	default:
		fmt.Printf("Unknown theme '%s'. Using mild.\n", *themePtr)
		expressions.ApplyMildPalette()
	}

	// 2) Get the canonical string from the AST.
	highlighted := ast.String()
	// 5) Print out the final colorized output
	fmt.Println(highlighted)
}

func runExportContextsCmd() {
	exportCmd := flag.NewFlagSet("export-contexts", flag.ExitOnError)
	expr := exportCmd.String("expr", "", "DSL expression to extract context identifiers from")
	inFile := exportCmd.String("in", "", "File containing a DSL expression")
	if err := exportCmd.Parse(os.Args[2:]); err != nil {
		fmt.Printf("Error reading command line args: %v\n", err)
		os.Exit(1)
	}
	var expression string
	if *inFile != "" {
		data, err := os.ReadFile(*inFile)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
		expression = string(data)
	} else if *expr != "" {
		expression = *expr
	} else {
		fmt.Println("Either -expr or -in flag must be provided.")
		exportCmd.Usage()
		os.Exit(1)
	}
	lex := lexer.NewLexer(expression)
	identifiers, err := lex.ExtractContextIdentifiers()
	if err != nil {
		fmt.Printf("Error extracting context identifiers: %v\n", err)
		os.Exit(1)
	}
	for _, id := range identifiers {
		fmt.Println(id)
	}
}
