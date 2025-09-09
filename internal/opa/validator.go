package opa

import (
	"context"
	"fmt"

	"github.com/open-policy-agent/opa/v1/ast"
	"github.com/open-policy-agent/opa/v1/rego"
	"go.uber.org/zap"
)

// Validator handles policy compilation and validation
type Validator struct {
	preparedQuery rego.PreparedEvalQuery
}

func NewValidatorFromDir(policiesDir string) (*Validator, error) {
	reader := NewPolicyReader()

	policies, err := reader.ReadPolicies(policiesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read policies: %w", err)
	}

	return NewValidator(policies)
}

func NewValidator(policies map[string]string) (*Validator, error) {
	if len(policies) == 0 {
		return nil, fmt.Errorf("no policies provided for validation")
	}

	validator := &Validator{}

	if err := validator.compilePolicies(policies); err != nil {
		return nil, fmt.Errorf("failed to compile policies: %w", err)
	}

	zap.S().Named("opa").Infof("OPA validator initialized with %d policies", len(policies))
	return validator, nil
}

// compilePolicies Compile the provided policy content and prepares the query
func (v *Validator) compilePolicies(policies map[string]string) error {
	compiler := ast.NewCompiler()
	modules := make(map[string]*ast.Module)

	for filename, content := range policies {
		// Parse with v0 for compatibility with existing Forklift policies
		module, err := ast.ParseModuleWithOpts(filename, content, ast.ParserOptions{
			RegoVersion: ast.RegoV0,
		})
		if err != nil {
			return fmt.Errorf("failed to parse policy %s: %w", filename, err)
		}
		modules[filename] = module
	}

	compiler.Compile(modules)
	if compiler.Failed() {
		return fmt.Errorf("policy compilation failed: %v", compiler.Errors)
	}

	// TODO:
	r := rego.New()

	preparedQuery, err := r.PrepareForEval(context.Background())
	if err != nil {
		return fmt.Errorf("failed to prepare rego query: %w", err)
	}

	v.preparedQuery = preparedQuery
	zap.S().Named("opa").Infof("Successfully compiled %d policy files", len(policies))
	return nil
}
