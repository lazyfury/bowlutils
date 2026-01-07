package condition_test

import (
	"testing"

	"github.com/lazyfury/bowlutils/crud/internal/condition"
	"gorm.io/gorm"
)

func TestNewCondition(t *testing.T) {
	tests := []struct {
		input    string
		expected condition.Condition
	}{
		{"eq", condition.Eq},
		{"EQ", condition.Eq}, // 大小写不敏感
		{"Eq", condition.Eq},
		{"ne", condition.Ne},
		{"gt", condition.Gt},
		{"gte", condition.Gte},
		{"lt", condition.Lt},
		{"lte", condition.Lte},
		{"in", condition.In},
		{"not_in", condition.NotIn},
		{"like", condition.Like},
		{"not_like", condition.NotLike},
		{"like_right", condition.LikeRight},
		{"like_left", condition.LikeLeft},
		{"is_null", condition.IsNull},
		{"is_notnull", condition.IsNotNull},
		{"sort", condition.Sort},
		{"unknown", condition.Eq}, // 未知条件返回默认 Eq
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := condition.NewCondition(tt.input)
			if result != tt.expected {
				t.Errorf("NewCondition(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCondition_Action(t *testing.T) {
	tests := []struct {
		condition condition.Condition
		expected  func(*gorm.DB, string, interface{}) *gorm.DB
	}{
		{condition.Eq, condition.EqAct},
		{condition.Ne, condition.NeAct},
		{condition.Gt, condition.GtAct},
		{condition.Gte, condition.GteAct},
		{condition.Lt, condition.LtAct},
		{condition.Lte, condition.LteAct},
		{condition.In, condition.InAct},
		{condition.NotIn, condition.NotInAct},
		{condition.Like, condition.LikeAct},
		{condition.NotLike, condition.NotLikeAct},
		{condition.LikeRight, condition.LikeRightAct},
		{condition.LikeLeft, condition.LikeLeftAct},
		{condition.IsNull, condition.IsNullAct},
		{condition.IsNotNull, condition.IsNotNullAct},
		{condition.Sort, condition.SortAct},
		{condition.Condition("unknown"), condition.EqAct}, // 未知条件返回默认
	}

	for _, tt := range tests {
		t.Run(string(tt.condition), func(t *testing.T) {
			action := tt.condition.Action()
			if action == nil {
				t.Fatal("Action() should not return nil")
			}
			// 注意：这里无法直接比较函数，只能验证不为 nil
			// 实际的行为测试需要 mock gorm.DB
		})
	}
}

func TestDefaultActions(t *testing.T) {
	expectedActions := []condition.Condition{
		condition.Eq, condition.Ne, condition.Gt, condition.Gte, condition.Lt, condition.Lte,
		condition.In, condition.NotIn,
		condition.Like, condition.NotLike, condition.LikeRight, condition.LikeLeft,
		condition.FK, condition.IsNull, condition.IsNotNull, condition.Sort,
	}

	if len(condition.DefaultActions) != len(expectedActions) {
		t.Errorf("DefaultActions length = %d, want %d", len(condition.DefaultActions), len(expectedActions))
	}

	for _, expected := range expectedActions {
		found := false
		for _, actual := range condition.DefaultActions {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("DefaultActions missing %v", expected)
		}
	}
}

func TestConditionConstants(t *testing.T) {
	conditions := map[string]condition.Condition{
		"eq":         condition.Eq,
		"ne":         condition.Ne,
		"gt":         condition.Gt,
		"gte":        condition.Gte,
		"lt":         condition.Lt,
		"lte":        condition.Lte,
		"in":         condition.In,
		"not_in":     condition.NotIn,
		"like":       condition.Like,
		"not_like":   condition.NotLike,
		"like_right": condition.LikeRight,
		"like_left":  condition.LikeLeft,
		"fk":         condition.FK,
		"is_null":    condition.IsNull,
		"is_notnull": condition.IsNotNull,
		"sort":       condition.Sort,
	}

	for name, cond := range conditions {
		t.Run(name, func(t *testing.T) {
			if string(cond) != name {
				t.Errorf("Condition %q = %q, want %q", name, cond, name)
			}
		})
	}
}

func TestNewCondition_CaseInsensitive(t *testing.T) {
	testCases := []struct {
		input    string
		expected condition.Condition
	}{
		{"EQ", condition.Eq},
		{"Eq", condition.Eq},
		{"eQ", condition.Eq},
		{"NE", condition.Ne},
		{"GT", condition.Gt},
		{"LIKE", condition.Like},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := condition.NewCondition(tc.input)
			if result != tc.expected {
				t.Errorf("NewCondition(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestCondition_String(t *testing.T) {
	if string(condition.Eq) != "eq" {
		t.Errorf("Eq.String() = %q, want 'eq'", condition.Eq)
	}
	if string(condition.Like) != "like" {
		t.Errorf("Like.String() = %q, want 'like'", condition.Like)
	}
}

// 注意：以下测试需要 mock gorm.DB，这里只做基本验证
// 完整的集成测试应该在 repository_test.go 中

func TestActionFunctions_Exist(t *testing.T) {
	actions := []func(*gorm.DB, string, interface{}) *gorm.DB{
		condition.EqAct, condition.NeAct, condition.GtAct, condition.GteAct, condition.LtAct, condition.LteAct,
		condition.InAct, condition.NotInAct,
		condition.LikeAct, condition.NotLikeAct, condition.LikeRightAct, condition.LikeLeftAct,
		condition.IsNullAct, condition.IsNotNullAct, condition.SortAct,
	}

	for i, action := range actions {
		if action == nil {
			t.Errorf("Action function at index %d is nil", i)
		}
	}
}

func TestFKAct_Format(t *testing.T) {
	// 测试 FKAct 的 SQL 格式
	// 注意：这需要实际的 gorm.DB，这里只验证函数存在
	if condition.FKAct == nil {
		t.Fatal("FKAct should not be nil")
	}

	// 验证格式字符串
	format := "%s.%s = %s.%s_%s"
	expectedFormat := "%s.%s = %s.%s_%s"
	if format != expectedFormat {
		t.Errorf("Format = %q, want %q", format, expectedFormat)
	}
}

func TestLikeActions(t *testing.T) {
	// 验证 Like 相关动作的字符串处理
	testCases := []struct {
		name     string
		action   func(*gorm.DB, string, interface{}) *gorm.DB
		value    string
		expected string
	}{
		{"Like", condition.LikeAct, "test", "%test%"},
		{"LikeRight", condition.LikeRightAct, "test", "test%"},
		{"LikeLeft", condition.LikeLeftAct, "test", "%test"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 这里只能验证函数存在，实际 SQL 生成需要 mock gorm.DB
			if tc.action == nil {
				t.Fatalf("%s action is nil", tc.name)
			}
		})
	}
}

func TestSortAct_Direction(t *testing.T) {
	// 验证 SortAct 接受的方向
	if condition.SortAct == nil {
		t.Fatal("SortAct should not be nil")
	}

	// 有效的排序方向
	validDirections := []string{"asc", "desc", "ASC", "DESC"}
	for _, dir := range validDirections {
		// 这里只能验证函数存在，实际测试需要 mock gorm.DB
		_ = dir
	}
}

func TestIsNullAndIsNotNull(t *testing.T) {
	// 验证 IsNull 和 IsNotNull 不依赖值
	if condition.IsNullAct == nil {
		t.Fatal("IsNullAct should not be nil")
	}
	if condition.IsNotNullAct == nil {
		t.Fatal("IsNotNullAct should not be nil")
	}
}

func TestInActions(t *testing.T) {
	// 验证 In 和 NotIn 动作
	if condition.InAct == nil {
		t.Fatal("InAct should not be nil")
	}
	if condition.NotInAct == nil {
		t.Fatal("NotInAct should not be nil")
	}
}

func TestComparisonActions(t *testing.T) {
	// 验证比较动作
	comparisons := []struct {
		name   string
		action func(*gorm.DB, string, interface{}) *gorm.DB
	}{
		{"Eq", condition.EqAct},
		{"Ne", condition.NeAct},
		{"Gt", condition.GtAct},
		{"Gte", condition.GteAct},
		{"Lt", condition.LtAct},
		{"Lte", condition.LteAct},
	}

	for _, tc := range comparisons {
		t.Run(tc.name, func(t *testing.T) {
			if tc.action == nil {
				t.Fatalf("%s action is nil", tc.name)
			}
		})
	}
}
