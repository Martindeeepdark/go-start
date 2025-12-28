package spec

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "gopkg.in/yaml.v3"
)

// Spec represents the API specification
type Spec struct {
	Spec     string            `yaml:"spec"`
	Kind     string            `yaml:"kind"`
	Name     string            `yaml:"name"`
	Version  string            `yaml:"version"`
	Project  ProjectConfig     `yaml:"project"`
	Models   []ModelDefinition `yaml:"models"`
	APIs     []APIEndpoint     `yaml:"endpoints"`
	Requests []RequestDef      `yaml:"requests"`
	Rules    []BusinessRule    `yaml:"rules"`
}

// ProjectConfig represents project configuration
type ProjectConfig struct {
	Module      string `yaml:"module"`
	Author      string `yaml:"author"`
	Description string `yaml:"description"`
}

// ModelDefinition represents a data model definition
type ModelDefinition struct {
	Name    string     `yaml:"name"`
	Table   string     `yaml:"table"`
	Comment string     `yaml:"comment"`
	Fields  []FieldDef `yaml:"fields"`
	Indexes []IndexDef `yaml:"indexes"`
}

// FieldDef represents a field definition
type FieldDef struct {
    Name           string `yaml:"name"`
    Type           string `yaml:"type"`
    Size           int    `yaml:"size,omitempty"`
    PrimaryKey     bool   `yaml:"primary,omitempty"`
    AutoIncrement  bool   `yaml:"autoIncrement,omitempty"`
    NotNull        bool   `yaml:"notNull,omitempty"`
    Unique         bool   `yaml:"unique,omitempty"`
    Index          bool   `yaml:"index,omitempty"`
    Default        string `yaml:"default,omitempty"`
    JSON           string `yaml:"json,omitempty"` // 自定义 JSON tag
    ForeignKey     string `yaml:"foreignKey,omitempty"`
    OnDelete       string `yaml:"onDelete,omitempty"`
    OnUpdate       string `yaml:"onUpdate,omitempty"`
    Comment        string `yaml:"comment,omitempty"`
    AutoCreateTime bool   `yaml:"autoCreateTime,omitempty"`
    AutoUpdateTime bool   `yaml:"autoUpdateTime,omitempty"`
}

// IndexDef represents an index definition
type IndexDef struct {
	Name   string   `yaml:"name"`
	Fields []string `yaml:"fields"`
	Unique bool     `yaml:"unique"`
}

// APIEndpoint represents an API endpoint definition
type APIEndpoint struct {
	Method     string       `yaml:"method"`
	Path       string       `yaml:"path"`
	Handler    string       `yaml:"handler"`
	Auth       bool         `yaml:"auth"`
	Permission string       `yaml:"permission,omitempty"`
	Validate   string       `yaml:"validate,omitempty"`
	Comment    string       `yaml:"comment,omitempty"`
	Cache      *CacheConfig `yaml:"cache,omitempty"`
	Pagination interface{}  `yaml:"pagination,omitempty"` // 支持 bool 和 PaginationConfig
}

// CacheConfig represents cache configuration
type CacheConfig struct {
	Enabled bool `yaml:"enabled"`
	TTL     int  `yaml:"ttl,omitempty"`
}

// PaginationConfig represents pagination configuration
type PaginationConfig struct {
	Page        int `yaml:"page,omitempty"`
	PageSize    int `yaml:"pageSize,omitempty"`
	MaxPageSize int `yaml:"maxPageSize,omitempty"`
}

// RequestDef represents a request validation definition
type RequestDef struct {
	Name    string         `yaml:"name"`
	Comment string         `yaml:"comment"`
	Fields  []RequestField `yaml:"fields"`
}

// RequestField represents a request validation field
type RequestField struct {
	Name    string `yaml:"name"`
	Rules   string `yaml:"rules"`
	Comment string `yaml:"comment"`
}

// BusinessRule represents a business rule definition
type BusinessRule struct {
	Name    string `yaml:"name"`
	Comment string `yaml:"comment"`
	Trigger string `yaml:"trigger"`
	Action  string `yaml:"action"`
}

// Parser represents the spec parser
type Parser struct {
	specDir string
}

// New creates a new spec parser
func New(specDir string) *Parser {
	return &Parser{
		specDir: specDir,
	}
}

// ParseFile parses a spec file
func (p *Parser) ParseFile(specPath string) (*Spec, error) {
	// Read the spec file
	data, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("读取规范文件失败: %w", err)
	}

	// Parse YAML
	var spec Spec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("解析 YAML 失败: %w", err)
	}

	// Validate spec
	if err := p.validateSpec(&spec); err != nil {
		return nil, fmt.Errorf("规范验证失败: %w", err)
	}

	return &spec, nil
}

// ParseDir parses all spec files in a directory
func (p *Parser) ParseDir(dir string) ([]*Spec, error) {
	var specs []*Spec

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Only process .spec.yaml or .spec.yml files
		if !isSpecFile(file.Name()) {
			continue
		}

		specPath := filepath.Join(dir, file.Name())
		spec, err := p.ParseFile(specPath)
		if err != nil {
			return nil, fmt.Errorf("解析 %s 失败: %w", file.Name(), err)
		}

		specs = append(specs, spec)
	}

	return specs, nil
}

// validateSpec validates the spec
func (p *Parser) validateSpec(spec *Spec) error {
	// Check required fields
	if spec.Spec == "" {
		return fmt.Errorf("缺少 spec 版本号")
	}
	if spec.Kind == "" {
		return fmt.Errorf("缺少 kind 类型")
	}
	if spec.Name == "" {
		return fmt.Errorf("缺少 API 名称")
	}
	if spec.Project.Module == "" {
		return fmt.Errorf("缺少项目模块名")
	}

	// Validate models
	for _, model := range spec.Models {
		if err := p.validateModel(&model); err != nil {
			return fmt.Errorf("模型 %s 验证失败: %w", model.Name, err)
		}
	}

	// Validate endpoints
	for _, endpoint := range spec.APIs {
		if err := p.validateEndpoint(&endpoint); err != nil {
			return fmt.Errorf("端点 %s %s 验证失败: %w", endpoint.Method, endpoint.Path, err)
		}
	}

	return nil
}

// validateModel validates a model definition
func (p *Parser) validateModel(model *ModelDefinition) error {
	if model.Name == "" {
		return fmt.Errorf("模型名称不能为空")
	}

	if model.Table == "" {
		return fmt.Errorf("模型 %s 缺少表名", model.Name)
	}

	// Check for primary key
	hasPrimaryKey := false
	for _, field := range model.Fields {
		if field.PrimaryKey {
			hasPrimaryKey = true
			break
		}
	}

	if !hasPrimaryKey {
		return fmt.Errorf("模型 %s 缺少主键", model.Name)
	}

	return nil
}

// validateEndpoint validates an endpoint definition
func (p *Parser) validateEndpoint(endpoint *APIEndpoint) error {
	if endpoint.Method == "" {
		return fmt.Errorf("HTTP 方法不能为空")
	}

	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true,
		"DELETE": true, "PATCH": true,
	}

	if !validMethods[endpoint.Method] {
		return fmt.Errorf("无效的 HTTP 方法: %s", endpoint.Method)
	}

	if endpoint.Path == "" {
		return fmt.Errorf("路径不能为空")
	}

	if endpoint.Handler == "" {
		return fmt.Errorf("Handler 不能为空")
	}

	return nil
}

// isSpecFile checks if a file is a spec file
func isSpecFile(filename string) bool {
	ext := filepath.Ext(filename)
	base := filename[:len(filename)-len(ext)]
	return (ext == ".yaml" || ext == ".yml") &&
		(len(base) > 10 && base[len(base)-10:] == ".spec")
}

// GetModelByName gets a model by name
func (s *Spec) GetModelByName(name string) (*ModelDefinition, bool) {
	for _, model := range s.Models {
		if model.Name == name {
			return &model, true
		}
	}
	return nil, false
}

// GetEndpointsByModel gets all endpoints for a specific model
func (s *Spec) GetEndpointsByModel(modelName string) []APIEndpoint {
	var endpoints []APIEndpoint
	for _, endpoint := range s.APIs {
		// Check if handler contains model name
		if containsModelName(endpoint.Handler, modelName) {
			endpoints = append(endpoints, endpoint)
		}
	}
	return endpoints
}

// containsModelName checks if handler name contains model name
func containsModelName(handler, modelName string) bool {
    h := strings.ToLower(handler)
    m := strings.ToLower(modelName)
    // basic contains
    if strings.Contains(h, m) {
        return true
    }
    // plural forms
    plural := m + "s"
    if strings.HasSuffix(m, "y") {
        plural = m[:len(m)-1] + "ies"
    }
    return strings.Contains(h, plural)
}
