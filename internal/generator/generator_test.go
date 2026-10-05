package generator

import (
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/erwinhermantodev/hexa-go/internal/config"
	"github.com/erwinhermantodev/hexa-go/internal/utils"
)

// generateSample builds a project with the default User model plus one of
// each component kind, mirroring what a user does after `hexa-go generate`.
func generateSample(t *testing.T, minimal bool) string {
	t.Helper()
	dir := t.TempDir()
	cfg := config.ProjectConfig{
		Name:        filepath.Join(dir, "demo"),
		ModuleName:  "example.com/demo",
		Description: "demo project",
		Author:      "tester",
	}
	if !minimal {
		cfg.Models = append(cfg.Models, config.DefaultUserModel())
		cfg.Services = append(cfg.Services, "Auth")
	}

	g := New()
	if err := g.CreateProject(cfg); err != nil {
		t.Fatalf("CreateProject: %v", err)
	}

	// Fields go through the same parser the CLI uses, covering the documented syntax.
	fields, err := utils.ParseFields([]string{
		"Name:string::required,min=2",
		"Slug:string:unique:required",
		"Price:float64::required,gt=0",
		"Published:bool:default=false",
		"PublishedAt:*time.Time",
	})
	if err != nil {
		t.Fatalf("ParseFields: %v", err)
	}
	product := config.ModelConfig{
		Name:    "Product",
		Fields:  utils.WithDefaultFields(fields),
		HasRepo: true, HasService: true, HasHandler: true,
	}
	if err := g.GenerateModelFiles(cfg, product); err != nil {
		t.Fatalf("GenerateModelFiles: %v", err)
	}
	for _, name := range []string{"order-item", "Category"} {
		m := config.ModelConfig{
			Name:    name,
			Fields:  utils.WithDefaultFields([]config.FieldConfig{{Name: "Name", Type: "string", Tag: "`json:\"name\"`", Validate: "required"}}),
			HasRepo: true, HasService: true, HasHandler: true,
		}
		if err := g.GenerateModelFiles(cfg, m); err != nil {
			t.Fatalf("GenerateModelFiles(%s): %v", name, err)
		}
	}
	if err := g.GenerateServiceFile(cfg, "NotificationService"); err != nil {
		t.Fatalf("GenerateServiceFile: %v", err)
	}
	if err := g.GenerateServiceFile(cfg, "Payment"); err != nil {
		t.Fatalf("GenerateServiceFile: %v", err)
	}
	if err := g.GenerateHandlerFile(cfg, "ReportHandler"); err != nil {
		t.Fatalf("GenerateHandlerFile: %v", err)
	}
	if err := g.GenerateModuleFiles(cfg, config.ModelConfig{
		Name: "blog-post",
		Fields: []config.FieldConfig{
			{Name: "ID", Type: "uint", Tag: "`gorm:\"primaryKey\" json:\"id\"`"},
			{Name: "Title", Type: "string", Tag: "`json:\"title\"`", Validate: "required"},
		},
		HasRepo: true, HasService: true, HasHandler: true,
	}, "blog-post"); err != nil {
		t.Fatalf("GenerateModuleFiles: %v", err)
	}
	return cfg.Name
}

// TestGeneratedFilesAreValid is fast and offline: every expected file exists,
// every .go file parses, and no template placeholder leaked into the output.
func TestGeneratedFilesAreValid(t *testing.T) {
	for _, minimal := range []bool{false, true} {
		name := "default"
		if minimal {
			name = "minimal"
		}
		t.Run(name, func(t *testing.T) {
			root := generateSample(t, minimal)

			for dest := range baseFiles {
				if _, err := os.Stat(filepath.Join(root, dest)); err != nil {
					t.Errorf("missing generated file %s", dest)
				}
			}

			fset := token.NewFileSet()
			err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return err
				}
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				if strings.Contains(string(content), "<no value>") {
					t.Errorf("%s contains an unresolved template value", path)
				}
				if strings.HasSuffix(path, ".go") {
					if _, err := parser.ParseFile(fset, path, content, parser.ParseComments); err != nil {
						t.Errorf("%s does not parse: %v", path, err)
					}
				}
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

// TestGeneratedProjectBuilds is the end-to-end check: tidy, build and vet the
// generated project (vet also compiles the generated tests). It needs network
// access to resolve dependencies, so it is skipped with -short.
func TestGeneratedProjectBuilds(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping: needs network for go mod tidy")
	}
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go toolchain not found")
	}

	for _, minimal := range []bool{false, true} {
		name := "default"
		if minimal {
			name = "minimal"
		}
		t.Run(name, func(t *testing.T) {
			root := generateSample(t, minimal)
			for _, args := range [][]string{{"mod", "tidy"}, {"build", "./..."}, {"vet", "./..."}} {
				cmd := exec.Command("go", args...)
				cmd.Dir = root
				if out, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("go %s failed: %v\n%s", strings.Join(args, " "), err, out)
				}
			}
		})
	}
}

// readFile is a test helper.
func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestNamingInGeneratedCode(t *testing.T) {
	root := generateSample(t, true)

	for _, f := range []string{
		"model/order_item.go", "repository/order_item.go", "service/order_item.go",
		"transport/http/handler/order_item_handler.go",
		"model/category.go", "service/notification.go", "service/payment.go",
		"transport/http/handler/report_handler.go", "internal/modules/blogpost/handler.go",
	} {
		if _, err := os.Stat(filepath.Join(root, f)); err != nil {
			t.Errorf("expected %s: %v", f, err)
		}
	}

	routes := readFile(t, filepath.Join(root, "transport/http/routes/routes.go"))
	for _, want := range []string{
		`api.Group("/order-items")`, `api.Group("/categories")`, `api.Group("/blog-posts")`,
		"r.OrderItemHandler.GetAllOrderItems", "r.CategoryHandler.GetAllCategories",
	} {
		if !strings.Contains(routes, want) {
			t.Errorf("routes.go missing %q", want)
		}
	}
	main := readFile(t, filepath.Join(root, "main.go"))
	for _, want := range []string{"orderItemRepo := repository.NewOrderItemRepository(db)", "notificationService := service.NewNotificationService()"} {
		if !strings.Contains(main, want) {
			t.Errorf("main.go missing %q", want)
		}
	}
}

func TestRefusesToOverwrite(t *testing.T) {
	root := generateSample(t, false)
	cfg := config.ProjectConfig{Name: root, ModuleName: "example.com/demo"}
	g := New()

	if err := g.GenerateModelFiles(cfg, config.DefaultUserModel()); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Errorf("model: got %v, want already-exists error", err)
	}
	if err := g.GenerateServiceFile(cfg, "Payment"); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Errorf("service: got %v, want already-exists error", err)
	}
	if err := g.GenerateHandlerFile(cfg, "Report"); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Errorf("handler: got %v, want already-exists error", err)
	}

	// A model and a module with the same name collide on the Router field.
	m := config.ModelConfig{Name: "BlogPost", Fields: config.DefaultModelFields(), HasRepo: true, HasService: true, HasHandler: true}
	if err := g.GenerateModelFiles(cfg, m); err == nil || !strings.Contains(err.Error(), "already registered") {
		t.Errorf("router field: got %v, want already-registered error", err)
	}

	// --force goes through and must not duplicate wiring.
	before := readFile(t, filepath.Join(root, "main.go"))
	g.Force = true
	if err := g.GenerateServiceFile(cfg, "Payment"); err != nil {
		t.Fatalf("forced service: %v", err)
	}
	if after := readFile(t, filepath.Join(root, "main.go")); after != before {
		t.Errorf("forced re-generation changed main.go:\n%s", after)
	}
}

func TestFailedGenerationRollsBack(t *testing.T) {
	root := generateSample(t, true)
	cfg := config.ProjectConfig{Name: root, ModuleName: "example.com/demo"}

	// Break routes.go so wiring fails after the files were already written.
	routesPath := filepath.Join(root, "transport/http/routes/routes.go")
	routes := readFile(t, routesPath)
	broken := strings.Replace(routes, "// [ROUTES-INIT]", "", 1)
	if err := os.WriteFile(routesPath, []byte(broken), 0644); err != nil {
		t.Fatal(err)
	}
	mainBefore := readFile(t, filepath.Join(root, "main.go"))

	m := config.ModelConfig{Name: "Invoice", Fields: config.DefaultModelFields(), HasRepo: true, HasService: true, HasHandler: true}
	if err := New().GenerateModelFiles(cfg, m); err == nil {
		t.Fatal("expected wiring to fail")
	}

	for _, f := range []string{"model/invoice.go", "repository/invoice.go", "service/invoice.go",
		"transport/http/handler/invoice_handler.go", "service/invoice_test.go"} {
		if _, err := os.Stat(filepath.Join(root, f)); !os.IsNotExist(err) {
			t.Errorf("%s was left behind after a failed generation", f)
		}
	}
	if got := readFile(t, filepath.Join(root, "main.go")); got != mainBefore {
		t.Error("main.go was not restored after a failed generation")
	}
	if got := readFile(t, routesPath); got != broken {
		t.Error("routes.go was not restored after a failed generation")
	}
}

func TestInvalidNamesAreRejected(t *testing.T) {
	root := generateSample(t, true)
	cfg := config.ProjectConfig{Name: root, ModuleName: "example.com/demo"}
	g := New()

	for _, name := range []string{"../evil", "a/b", "1abc", "type", ""} {
		if err := g.GenerateServiceFile(cfg, name); err == nil {
			t.Errorf("service %q: expected error", name)
		}
	}
	if err := g.GenerateModuleFiles(cfg, config.ModelConfig{Fields: config.DefaultModelFields()}, "utils"); err == nil {
		t.Error("module named utils: expected clash error")
	}
	m := config.ModelConfig{Name: "Orphan", Fields: config.DefaultModelFields(), HasService: true}
	if err := g.GenerateModelFiles(cfg, m); err == nil {
		t.Error("service without repository: expected error")
	}
}

func TestCreateProjectCleansUpOnFailure(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "demo")
	cfg := config.ProjectConfig{
		Name: dir, ModuleName: "example.com/demo",
		Models: []config.ModelConfig{{Name: "bad name!", HasRepo: true}},
	}
	if err := New().CreateProject(cfg); err == nil {
		t.Fatal("expected error")
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Error("project directory was left behind")
	}

	// A non-empty existing directory is refused and left untouched.
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "keep.txt"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg.Models = nil
	if err := New().CreateProject(cfg); err == nil || !strings.Contains(err.Error(), "not empty") {
		t.Errorf("got %v, want not-empty error", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "keep.txt")); err != nil {
		t.Error("existing directory content was removed")
	}
}
