package testing

import (
	stdErrors "errors"
	"fmt"
	astClass "github.com/RyanCopley/expression-parser/pkg/ast/expressions"
	"github.com/RyanCopley/expression-parser/pkg/env"
	"github.com/RyanCopley/expression-parser/pkg/errors"
	"github.com/RyanCopley/expression-parser/pkg/lexer"
	"github.com/RyanCopley/expression-parser/pkg/parser"
	"github.com/RyanCopley/expression-parser/pkg/types"
	"math"
	"strings"
	"time"
)

// TestCase represents a DSL test case.
type TestCase struct {
	Description          string                 `yaml:"description"`
	Context              map[string]interface{} `yaml:"context"`
	Expression           string                 `yaml:"expression"`
	ExpectedError        string                 `yaml:"expectedError"`
	ExpectedErrorMessage string                 `yaml:"expectedErrorMessage"`
	ExpectedResult       interface{}            `yaml:"expectedResult"`
	Skip                 bool                   `yaml:"skip"`
	Focus                bool                   `yaml:"focus"`
}

// TestResult represents the result of executing a test case.
type TestResult struct {
	TestID               int                    `yaml:"testId"`
	Description          string                 `yaml:"description"`
	Expression           string                 `yaml:"expression"`
	Context              map[string]interface{} `yaml:"context"`
	ExpectedResult       interface{}            `yaml:"expectedResult,omitempty"`
	ExpectedError        string                 `yaml:"expectedError,omitempty"`
	ExpectedErrorMessage string                 `yaml:"expectedErrorMessage,omitempty"`
	ActualResult         interface{}            `yaml:"actualResult,omitempty"`
	ActualError          error                  `yaml:"actualError,omitempty"`
	Status               string                 `yaml:"status"`
	ErrLine              int                    `yaml:"errorLine,omitempty"`
	ErrColumn            int                    `yaml:"errorColumn,omitempty"`
	ErrorContext         string                 `yaml:"errorSnippet,omitempty"`
	BenchmarkTime        string                 `yaml:"benchmarkTime,omitempty"`
	BenchmarkOpsSec      float64                `yaml:"benchmarkOpsSec,omitempty"`
}

// TestSuiteResult aggregates the results of a test suite.
type TestSuiteResult struct {
	Passed      int          `yaml:"passed"`
	Failed      int          `yaml:"failed"`
	Skipped     int          `yaml:"skipped"`
	Total       int          `yaml:"total"`
	TestResults []TestResult `yaml:"test_results"`
}

// RunTests processes test cases and returns a suite result.

func RunTests(testCases []TestCase, env *env.Environment, failFast bool, benchmark bool) TestSuiteResult {
	suiteResult := TestSuiteResult{
		TestResults: []TestResult{},
	}
	// Determine if any test is marked as focused.
	focusMode := false
	for _, tc := range testCases {
		if tc.Focus {
			focusMode = true
			break
		}
	}

	// Process each test case.
	for i, tc := range testCases {
		testID := i + 1
		result := TestResult{
			TestID:               testID,
			Description:          tc.Description,
			Expression:           tc.Expression,
			Context:              tc.Context,
			ExpectedResult:       tc.ExpectedResult,
			ExpectedError:        tc.ExpectedError,
			ExpectedErrorMessage: tc.ExpectedErrorMessage,
		}

		// Skip tests that are not focused when focus mode is active.
		if focusMode && !tc.Focus {
			result.Status = "SKIPPED"
			suiteResult.Skipped++
			suiteResult.TestResults = append(suiteResult.TestResults, result)
			continue
		}

		// Skip tests explicitly marked as skipped.
		if tc.Skip {
			result.Status = "SKIPPED"
			suiteResult.Skipped++
			suiteResult.TestResults = append(suiteResult.TestResults, result)
			continue
		}

		// Only count tests that actually run.
		suiteResult.Total++

		// Parse the expression.
		lexer := lexer.NewLexer(tc.Expression)
		parser, err := parser.NewParser(lexer)
		if err != nil {
			var errorWithDetail errors.PositionalError
			hasErrorWithDetail := stdErrors.As(err, &errorWithDetail)

			errMsg := err.Error()
			result.ActualError = err
			errLine, errColumn := errors.GetErrorPosition(err)
			result.ErrLine = errLine
			result.ErrColumn = errColumn
			result.ErrorContext = errors.GetErrorContext(tc.Expression, errLine, errColumn, false)
			if (hasErrorWithDetail && tc.ExpectedError == errorWithDetail.Kind()) && strings.Contains(errMsg, tc.ExpectedErrorMessage) {
				result.Status = "PASSED"
				suiteResult.Passed++
			} else {
				result.Status = "FAILED"
				suiteResult.Failed++
				if failFast {
					suiteResult.TestResults = append(suiteResult.TestResults, result)
					break
				}
			}
			suiteResult.TestResults = append(suiteResult.TestResults, result)
			continue
		}

		ast, parseErr := parser.ParseExpression()
		if parseErr != nil {
			var errorWithDetail errors.PositionalError
			hasErrorWithDetail := stdErrors.As(parseErr, &errorWithDetail)
			errMsg := parseErr.Error()
			result.ActualError = parseErr
			errLine, errColumn := errors.GetErrorPosition(parseErr)
			result.ErrLine = errLine
			result.ErrColumn = errColumn
			result.ErrorContext = errors.GetErrorContext(tc.Expression, errLine, errColumn, false)
			if (hasErrorWithDetail && tc.ExpectedError == errorWithDetail.Kind()) && strings.Contains(errMsg, tc.ExpectedErrorMessage) {
				result.Status = "PASSED"
				suiteResult.Passed++
			} else {
				result.Status = "FAILED"
				suiteResult.Failed++
				if failFast {
					suiteResult.TestResults = append(suiteResult.TestResults, result)
					break
				}
			}
			suiteResult.TestResults = append(suiteResult.TestResults, result)
			continue
		}
		result.Expression = ast.String()

		// Evaluate the AST.
		evalResult, evalErr := ast.Eval(tc.Context, env)
		if evalErr != nil {
			var errorWithDetail errors.PositionalError
			hasErrorWithDetail := stdErrors.As(evalErr, &errorWithDetail)
			errMsg := evalErr.Error()
			result.ActualError = evalErr
			errLine, errColumn := errors.GetErrorPosition(evalErr)
			result.ErrLine = errLine
			result.ErrColumn = errColumn
			result.ErrorContext = errors.GetErrorContext(tc.Expression, errLine, errColumn, false)
			if (hasErrorWithDetail && tc.ExpectedError == errorWithDetail.Kind()) && strings.Contains(errMsg, tc.ExpectedErrorMessage) {
				result.Status = "PASSED"
				suiteResult.Passed++
			} else {
				result.Status = "FAILED"
				suiteResult.Failed++
				if failFast {
					suiteResult.TestResults = append(suiteResult.TestResults, result)
					break
				}
			}
			suiteResult.TestResults = append(suiteResult.TestResults, result)
			continue
		}

		// If an error was expected but evaluation produced a result.
		if tc.ExpectedError != "" {
			result.ActualResult = evalResult
			result.Status = "FAILED"
			suiteResult.Failed++
			if failFast {
				suiteResult.TestResults = append(suiteResult.TestResults, result)
				break
			}
			suiteResult.TestResults = append(suiteResult.TestResults, result)
			continue
		}

		// Compare the actual result with the expected result.
		result.ActualResult = evalResult
		var passTest bool
		if rVal, ok := types.ToFloat(evalResult); ok {
			if eVal, ok2 := types.ToFloat(tc.ExpectedResult); ok2 {
				passTest = math.Abs(rVal-eVal) < 1e-9
			} else {
				passTest = fmt.Sprintf("%v", evalResult) == fmt.Sprintf("%v", tc.ExpectedResult)
			}
		} else {
			var resultStr, expectedStr string
			if resStr, ok := evalResult.(string); ok {
				resultStr = strings.ReplaceAll(resStr, "\n", "\\n")
			} else {
				resultStr = fmt.Sprintf("%v", evalResult)
			}
			if expStr, ok := tc.ExpectedResult.(string); ok {
				expectedStr = strings.ReplaceAll(expStr, "\n", "\\n")
			} else {
				expectedStr = fmt.Sprintf("%v", tc.ExpectedResult)
			}
			passTest = resultStr == expectedStr
		}

		if passTest {
			result.Status = "PASSED"
			suiteResult.Passed++
		} else {
			result.Status = "FAILED"
			suiteResult.Failed++
			if failFast {
				suiteResult.TestResults = append(suiteResult.TestResults, result)
				break
			}
		}

		// --- BENCHMARKING ---
		// Only run benchmark if the flag is enabled,
		// the test passed and no error was expected.
		// And only benchmark if the top-level AST is a FunctionCallExpr.
		if benchmark && result.Status == "PASSED" && tc.ExpectedError == "" {
			if _, isFuncCall := ast.(*astClass.FunctionCallExpr); isFuncCall {
				iterations := 1000
				start := time.Now()
				for j := 0; j < iterations; j++ {
					// We ignore errors here since the single-run was already successful.
					_, _ = ast.Eval(tc.Context, env)
				}
				elapsed := time.Since(start)
				result.BenchmarkTime = elapsed.String()
				result.BenchmarkOpsSec = float64(iterations) / elapsed.Seconds()
			}
		}
		// --- end benchmark ---

		suiteResult.TestResults = append(suiteResult.TestResults, result)
	}
	return suiteResult
}
