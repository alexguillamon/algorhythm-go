package orderedset

import (
	"reflect"
	"testing"
)

func TestAddAndContains(t *testing.T) {
	os := NewSet[string]()
	os.Add("a")
	os.Add("b")
	os.Add("c")

	// Check that elements exist
	if !os.Contains("a") {
		t.Errorf("Expected set to contain 'a'")
	}
	if !os.Contains("b") {
		t.Errorf("Expected set to contain 'b'")
	}
	if !os.Contains("c") {
		t.Errorf("Expected set to contain 'c'")
	}

	// Adding a duplicate should not alter the order or size.
	before := os.ToSlice()
	os.Add("a")
	after := os.ToSlice()
	if !reflect.DeepEqual(before, after) {
		t.Errorf("Duplicate addition should not change the set. Got before: %v, after: %v", before, after)
	}
}

// New test for variadic Add functionality
func TestAddMultiple(t *testing.T) {
	os := NewSet[string]()

	// Add multiple items at once.
	os.Add("a", "b", "c")
	expected := []string{"a", "b", "c"}
	if got := os.ToSlice(); !reflect.DeepEqual(got, expected) {
		t.Errorf("After Add('a', 'b', 'c'), expected %v, got %v", expected, got)
	}

	// Add more items, including a duplicate, and check order is maintained.
	os.Add("b", "d", "e")
	expected = []string{"a", "b", "c", "d", "e"}
	if got := os.ToSlice(); !reflect.DeepEqual(got, expected) {
		t.Errorf("After Add('b', 'd', 'e'), expected %v, got %v", expected, got)
	}
}

func TestRemove(t *testing.T) {
	os := NewSet("a", "b", "c")
	os.Remove("b")

	if os.Contains("b") {
		t.Errorf("Expected set to not contain 'b' after removal")
	}

	expected := []string{"a", "c"}
	result := os.ToSlice()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected slice %v after removal, got %v", expected, result)
	}

	// Removing an element that doesn't exist should have no effect.
	os.Remove("non-existent")
	if !reflect.DeepEqual(os.ToSlice(), expected) {
		t.Errorf("Removal of a non-existent element should not change the set")
	}
}

func TestToSlice(t *testing.T) {
	items := []string{"apple", "banana", "cherry"}
	os := NewSet(items...)
	result := os.ToSlice()

	if !reflect.DeepEqual(result, items) {
		t.Errorf("Expected slice %v, got %v", items, result)
	}
}

func TestUnion(t *testing.T) {
	// Create two sets with some overlapping elements.
	left := NewSet("a", "b", "c")
	right := NewSet("b", "c", "d")

	// The union should preserve the order: all left elements then only new right elements.
	unionSet := left.Union(right)
	expected := []string{"a", "b", "c", "d"}
	result := unionSet.ToSlice()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected union %v, got %v", expected, result)
	}
}

func TestDifference(t *testing.T) {
	left := NewSet("a", "b", "c", "d")
	right := NewSet("b", "d")

	// The difference should only include elements in left that are not in right.
	diffSet := left.Difference(right)
	expected := []string{"a", "c"}
	result := diffSet.ToSlice()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected difference %v, got %v", expected, result)
	}
}
