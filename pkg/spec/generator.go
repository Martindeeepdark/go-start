package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Generator represents the code generator
type Generator struct {
	spec      *Spec
	outputDir string
}

// NewGenerator creates a new code generator
func NewGenerator(spec *Spec, outputDir string) *Generator {
	return &Generator{
		spec:      spec,
		outputDir: outputDir,
	}
}

// Generate generates all code from the spec
func (g *Generator) Generate() error {
	fmt.Printf("\n🚀 开始生成代码...\n\n")

	// 1. Generate models
	if err := g.generateModels(); err != nil {
		return fmt.Errorf("生成模型失败: %w", err)
	}

	// 2. Generate repositories
	if err := g.generateRepositories(); err != nil {
		return fmt.Errorf("生成仓储失败: %w", err)
	}

	// 3. Generate services
	if err := g.generateServices(); err != nil {
		return fmt.Errorf("生成服务失败: %w", err)
	}

	// 4. Generate controllers
	if err := g.generateControllers(); err != nil {
		return fmt.Errorf("生成控制器失败: %w", err)
	}

	// 5. Generate request validators (if any)
	if len(g.spec.Requests) > 0 {
		if err := g.generateValidators(); err != nil {
			fmt.Printf("⚠️  生成验证器跳过（模板未实现）\n")
		}
	}

	// 6. Generate routes
	if err := g.generateRoutes(); err != nil {
		fmt.Printf("⚠️  生成路由跳过（模板未实现）\n")
	}

	fmt.Printf("\n✅ 代码生成完成！\n")
	return nil
}

// generateModels generates model files
func (g *Generator) generateModels() error {
	fmt.Println("📦 生成数据模型...")

	for _, model := range g.spec.Models {
		outputPath := filepath.Join(g.outputDir, "internal/model", strings.ToLower(model.Name)+".go")

		if err := g.generateFile("model.go.tmpl", outputPath, map[string]interface{}{
			"Spec":  g.spec,
			"Model": model,
		}); err != nil {
			return err
		}

		fmt.Printf("  ✓ %s\n", model.Name)
	}

	return nil
}

// generateRepositories generates repository files
func (g *Generator) generateRepositories() error {
	fmt.Println("\n📦 生成数据访问层...")

	for _, model := range g.spec.Models {
		outputPath := filepath.Join(g.outputDir, "internal/repository", strings.ToLower(model.Name)+".go")

		if err := g.generateFile("repository.go.tmpl", outputPath, map[string]interface{}{
			"Spec":  g.spec,
			"Model": model,
		}); err != nil {
			return err
		}

		fmt.Printf("  ✓ %sRepository\n", model.Name)
	}

	return nil
}

// generateServices generates service files
func (g *Generator) generateServices() error {
	fmt.Println("\n📦 生成业务逻辑层...")

	for _, model := range g.spec.Models {
		outputPath := filepath.Join(g.outputDir, "internal/service", strings.ToLower(model.Name)+".go")

		// Get endpoints for this model
		endpoints := g.spec.GetEndpointsByModel(model.Name)

		var getCacheEnabled bool
		var listCacheEnabled bool
		var getCacheTTL int
		var listCacheTTL int
		for _, ep := range endpoints {
			if strings.EqualFold(ep.Method, "GET") {
				if ep.Cache != nil && ep.Cache.Enabled {
					if strings.Contains(ep.Path, ":") {
						getCacheEnabled = true
						getCacheTTL = ep.Cache.TTL
					} else {
						listCacheEnabled = true
						listCacheTTL = ep.Cache.TTL
					}
				}
			}
		}

		if err := g.generateFile("service.go.tmpl", outputPath, map[string]interface{}{
			"Spec":             g.spec,
			"Model":            model,
			"Endpoints":        endpoints,
			"GetCacheEnabled":  getCacheEnabled,
			"ListCacheEnabled": listCacheEnabled,
			"GetCacheTTL":      getCacheTTL,
			"ListCacheTTL":     listCacheTTL,
		}); err != nil {
			return err
		}

		fmt.Printf("  ✓ %sService\n", model.Name)
	}

	return nil
}

// generateControllers generates controller files
func (g *Generator) generateControllers() error {
	fmt.Println("\n📦 生成控制器层...")

	for _, model := range g.spec.Models {
		outputPath := filepath.Join(g.outputDir, "internal/controller", strings.ToLower(model.Name)+".go")

		// Get endpoints for this model
		endpoints := g.spec.GetEndpointsByModel(model.Name)

		// Derive validator and auth/permission per operation from endpoints
		var createValidator, updateValidator string
		var createAuth, updateAuth, getAuth, deleteAuth, listAuth bool
		var createPerm, updatePerm, getPerm, deletePerm, listPerm string
		for _, ep := range endpoints {
			m := strings.ToUpper(ep.Method)
			switch m {
			case "POST":
				if ep.Validate != "" && createValidator == "" {
					createValidator = ep.Validate
				}
				if ep.Auth {
					createAuth = true
				}
				if ep.Permission != "" && createPerm == "" {
					createPerm = ep.Permission
				}
			case "PUT":
				if ep.Validate != "" && updateValidator == "" {
					updateValidator = ep.Validate
				}
				if ep.Auth {
					updateAuth = true
				}
				if ep.Permission != "" && updatePerm == "" {
					updatePerm = ep.Permission
				}
			case "GET":
				// distinguish list vs get by path heuristic
				if strings.Contains(ep.Path, ":") { // detail
					if ep.Auth {
						getAuth = true
					}
					if ep.Permission != "" && getPerm == "" {
						getPerm = ep.Permission
					}
				} else { // list
					if ep.Auth {
						listAuth = true
					}
					if ep.Permission != "" && listPerm == "" {
						listPerm = ep.Permission
					}
				}
			case "DELETE":
				if ep.Auth {
					deleteAuth = true
				}
				if ep.Permission != "" && deletePerm == "" {
					deletePerm = ep.Permission
				}
			}
		}

		if err := g.generateFile("controller.go.tmpl", outputPath, map[string]interface{}{
			"Spec":            g.spec,
			"Model":           model,
			"Endpoints":       endpoints,
			"CreateValidator": createValidator,
			"UpdateValidator": updateValidator,
			"CreateAuth":      createAuth,
			"UpdateAuth":      updateAuth,
			"GetAuth":         getAuth,
			"DeleteAuth":      deleteAuth,
			"ListAuth":        listAuth,
			"CreatePerm":      createPerm,
			"UpdatePerm":      updatePerm,
			"GetPerm":         getPerm,
			"DeletePerm":      deletePerm,
			"ListPerm":        listPerm,
		}); err != nil {
			return err
		}

		fmt.Printf("  ✓ %sController\n", model.Name)
	}

	return nil
}

// generateValidators generates validator files
func (g *Generator) generateValidators() error {
	fmt.Println("\n📦 生成请求验证器...")

	for _, req := range g.spec.Requests {
		outputPath := filepath.Join(g.outputDir, "internal/validator", strings.ToLower(req.Name)+".go")

		if err := g.generateFile("validator.go.tmpl", outputPath, map[string]interface{}{
			"Spec":    g.spec,
			"Request": req,
		}); err != nil {
			return err
		}

		fmt.Printf("  ✓ %s\n", req.Name)
	}

	return nil
}

// generateRoutes generates route registration
func (g *Generator) generateRoutes() error {
	fmt.Println("\n📦 生成路由注册...")

	outputPath := filepath.Join(g.outputDir, "internal/routes", "auto_routes.go")

	// Build per-model auth/permission info for routes
	type routeInfo struct {
		ModelName  string
		ModelVar   string
		CreateAuth bool
		UpdateAuth bool
		GetAuth    bool
		DeleteAuth bool
		ListAuth   bool
		CreatePerm string
		UpdatePerm string
		GetPerm    string
		DeletePerm string
		ListPerm   string
	}

	var modelsInfo []routeInfo
	for _, model := range g.spec.Models {
		endpoints := g.spec.GetEndpointsByModel(model.Name)
		ri := routeInfo{ModelName: model.Name, ModelVar: toLowerCamelCase(model.Name)}
		for _, ep := range endpoints {
			m := strings.ToUpper(ep.Method)
			switch m {
			case "POST":
				if ep.Auth {
					ri.CreateAuth = true
				}
				if ep.Permission != "" && ri.CreatePerm == "" {
					ri.CreatePerm = ep.Permission
				}
			case "PUT":
				if ep.Auth {
					ri.UpdateAuth = true
				}
				if ep.Permission != "" && ri.UpdatePerm == "" {
					ri.UpdatePerm = ep.Permission
				}
			case "GET":
				if strings.Contains(ep.Path, ":") {
					if ep.Auth {
						ri.GetAuth = true
					}
					if ep.Permission != "" && ri.GetPerm == "" {
						ri.GetPerm = ep.Permission
					}
				} else {
					if ep.Auth {
						ri.ListAuth = true
					}
					if ep.Permission != "" && ri.ListPerm == "" {
						ri.ListPerm = ep.Permission
					}
				}
			case "DELETE":
				if ep.Auth {
					ri.DeleteAuth = true
				}
				if ep.Permission != "" && ri.DeletePerm == "" {
					ri.DeletePerm = ep.Permission
				}
			}
		}
		modelsInfo = append(modelsInfo, ri)
	}

	if err := g.generateFile("routes.go.tmpl", outputPath, map[string]interface{}{
		"Spec":       g.spec,
		"ModelsInfo": modelsInfo,
	}); err != nil {
		return err
	}

	fmt.Printf("  ✓ 自动路由注册\n")

	return nil
}

// generateFile generates a single file from template
func (g *Generator) generateFile(templateName, outputPath string, data interface{}) error {
	// Create output directory
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// Get template content
	templateContent := getBuiltinTemplate(templateName)
	if templateContent == "" {
		return fmt.Errorf("模板 %s 不存在", templateName)
	}

	// Parse template with custom functions
	funcMap := template.FuncMap{
		"ToCamelCase":      toCamelCase,
		"ToLowerCamelCase": toLowerCamelCase,
		"pluralize":        pluralize,
		"getGoType":        getGoType,
		"getGoType2":       getGoType2,
		"getGormTag":       getGormTag,
		"getJSONTag":       getJSONTag,
		"getIndexTags":     getIndexTags,
		"getReqFieldType":  getReqFieldType,
	}

	tmpl, err := template.New(templateName).Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("解析模板失败: %w", err)
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败: %w", err)
	}
	defer outputFile.Close()

	// Execute template
	if err := tmpl.Execute(outputFile, data); err != nil {
		return fmt.Errorf("执行模板失败: %w", err)
	}

	return nil
}

// Helper functions for templates

func toCamelCase(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		if i > 0 || len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, "")
}

func toLowerCamelCase(s string) string {
	camel := toCamelCase(s)
	return strings.ToLower(camel[:1]) + camel[1:]
}

func pluralize(s string) string {
	if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	return s + "s"
}

func getGoType(fieldType string) string {
	typeMap := map[string]string{
		"uint":      "uint",
		"int":       "int",
		"string":    "string",
		"text":      "string",
		"bool":      "bool",
		"float":     "float64",
		"double":    "float64",
		"decimal":   "float64",
		"timestamp": "time.Time",
		"date":      "time.Time",
		"datetime":  "time.Time",
		"json":      "string",
	}

	if goType, ok := typeMap[fieldType]; ok {
		return goType
	}

	return "string" // default
}

func getGoType2(field FieldDef) string {
	base := getGoType(field.Type)
	if field.NotNull {
		return base
	}
	switch base {
	case "uint":
		return "*uint"
	case "int":
		return "*int"
	case "string":
		return "*string"
	case "bool":
		return "*bool"
	case "float64":
		return "*float64"
	case "time.Time":
		return "*time.Time"
	default:
		return "*" + base
	}
}

func getGormTag(field FieldDef) string {
	var tags []string

	if field.PrimaryKey {
		tags = append(tags, "primarykey")
	}

	if field.AutoIncrement {
		tags = append(tags, "autoIncrement")
	}

	if field.Size > 0 {
		tags = append(tags, fmt.Sprintf("size:%d", field.Size))
	}

	if field.NotNull {
		tags = append(tags, "not null")
	}

	if field.Unique {
		tags = append(tags, "uniqueIndex")
	}

	if field.Index {
		tags = append(tags, "index")
	}

	if field.Default != "" {
		tags = append(tags, fmt.Sprintf("default:%s", field.Default))
	}

	if field.ForeignKey != "" {
		tags = append(tags, fmt.Sprintf("foreignKey:%s", field.ForeignKey))
		// configurable FK strategies
		onUpdate := "CASCADE"
		onDelete := "SET NULL"
		if field.NotNull {
			onDelete = "CASCADE"
		}
		if field.OnUpdate != "" {
			onUpdate = field.OnUpdate
		}
		if field.OnDelete != "" {
			onDelete = field.OnDelete
		}
		tags = append(tags, fmt.Sprintf("constraint:OnUpdate:%s,OnDelete:%s", onUpdate, onDelete))
	}

	if field.AutoCreateTime {
		tags = append(tags, "autoCreateTime")
	}

	if field.AutoUpdateTime {
		tags = append(tags, "autoUpdateTime")
	}

	if field.Comment != "" {
		tags = append(tags, fmt.Sprintf("comment:%s", field.Comment))
	}

	return strings.Join(tags, ";")
}

// getReqFieldType 根据规则字符串推断请求校验结构的字段类型
// 简化映射：包含 array -> []string；包含 numeric -> int；包含 bool -> bool；否则 string
func getReqFieldType(rules string) string {
	rs := strings.ToLower(rules)
	if strings.Contains(rs, "array") {
		return "[]string"
	}
	if strings.Contains(rs, "numeric") {
		return "int"
	}
	if strings.Contains(rs, "bool") || strings.Contains(rs, "boolean") {
		return "bool"
	}
	return "string"
}

func getIndexTags(fieldName string, indexes []IndexDef) string {
	var parts []string
	for _, idx := range indexes {
		for _, f := range idx.Fields {
			if f == fieldName {
				if idx.Unique {
					parts = append(parts, "uniqueIndex:"+idx.Name)
				} else {
					parts = append(parts, "index:"+idx.Name)
				}
				break
			}
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return ";" + strings.Join(parts, ";")
}

func getJSONTag(fieldName string, customJSON string) string {
	if customJSON != "" {
		return customJSON
	}
	return toLowerCamelCase(fieldName)
}
