package naming

import "testing"

func TestCases(t *testing.T) {
	tests := []struct{ in, pascal, camel, snake, kebab, pkg string }{
		{"user", "User", "user", "user", "user", "user"},
		{"User", "User", "user", "user", "user", "user"},
		{"order-item", "OrderItem", "orderItem", "order_item", "order-item", "orderitem"},
		{"order_item", "OrderItem", "orderItem", "order_item", "order-item", "orderitem"},
		{"OrderItem", "OrderItem", "orderItem", "order_item", "order-item", "orderitem"},
		{"orderItem", "OrderItem", "orderItem", "order_item", "order-item", "orderitem"},
		{"Order Item", "OrderItem", "orderItem", "order_item", "order-item", "orderitem"},
		{"UserID", "UserID", "userID", "user_id", "user-id", "userid"},
		{"IDCard", "IDCard", "idCard", "id_card", "id-card", "idcard"},
		{"HTTPServer", "HTTPServer", "httpServer", "http_server", "http-server", "httpserver"},
		{"blog", "Blog", "blog", "blog", "blog", "blog"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			for name, pair := range map[string][2]string{
				"Pascal": {Pascal(tt.in), tt.pascal}, "Camel": {Camel(tt.in), tt.camel},
				"Snake": {Snake(tt.in), tt.snake}, "Kebab": {Kebab(tt.in), tt.kebab},
				"Package": {Package(tt.in), tt.pkg},
			} {
				if pair[0] != pair[1] {
					t.Errorf("%s(%q) = %q, want %q", name, tt.in, pair[0], pair[1])
				}
			}
		})
	}
}

func TestPlural(t *testing.T) {
	for in, want := range map[string]string{
		"Product": "Products", "Category": "Categories", "Address": "Addresses",
		"Box": "Boxes", "Match": "Matches", "Day": "Days", "Key": "Keys",
		"OrderItem": "OrderItems", "UserCategory": "UserCategories", "person": "People",
		"Status": "Statuses", "Equipment": "Equipment", "ID": "IDs", "Bus": "Buses",
		"order-item": "OrderItems",
	} {
		if got := Plural(in); got != want {
			t.Errorf("Plural(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestValidate(t *testing.T) {
	for _, ok := range []string{"User", "order-item", "order_item", "Order Item", "A1", "blog2"} {
		if err := Validate(ok); err != nil {
			t.Errorf("Validate(%q) = %v, want nil", ok, err)
		}
	}
	for _, bad := range []string{"", "1abc", "a--b", "-a", "a-", "a.b", "a/b", "../x", "ünï", "type", "func", "a b!"} {
		if err := Validate(bad); err == nil {
			t.Errorf("Validate(%q) = nil, want error", bad)
		}
	}
	for _, bad := range []string{"utils", "Routes", "handler", "Model"} {
		if err := ValidatePackage(bad); err == nil {
			t.Errorf("ValidatePackage(%q) = nil, want error", bad)
		}
	}
	if err := ValidatePackage("blog"); err != nil {
		t.Errorf("ValidatePackage(blog) = %v", err)
	}
}

func TestTrimSuffix(t *testing.T) {
	for _, tt := range [][3]string{
		{"NotificationService", "Service", "Notification"},
		{"ReportHandler", "Handler", "Report"},
		{"reportHANDLER", "Handler", "report"},
		{"Service", "Service", "Service"},
		{"Payment", "Service", "Payment"},
	} {
		if got := TrimSuffix(tt[0], tt[1]); got != tt[2] {
			t.Errorf("TrimSuffix(%q,%q) = %q, want %q", tt[0], tt[1], got, tt[2])
		}
	}
}
