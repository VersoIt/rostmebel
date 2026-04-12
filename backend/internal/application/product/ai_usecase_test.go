package product

import (
	"testing"

	domProduct "github.com/rostmebel/backend/internal/domain/product"
)

func TestFilterByRequestedCategoryKeepsKitchenQueriesFocused(t *testing.T) {
	kitchenCategoryID := int64(1)
	wardrobeCategoryID := int64(2)
	projects := []*domProduct.Project{
		{ID: 1, ProjectCategoryID: &wardrobeCategoryID, Name: "Шкаф"},
		{ID: 2, ProjectCategoryID: &kitchenCategoryID, Name: "Кухня"},
	}
	categories := map[int64]string{
		kitchenCategoryID:  "Кухни",
		wardrobeCategoryID: "Шкафы-купе",
	}

	filtered := filterByRequestedCategory("светлая кухня до 400000", projects, categories)
	if len(filtered) != 1 || filtered[0].ID != 2 {
		t.Fatalf("expected only kitchen project, got %#v", filtered)
	}
}

func TestFilterByRequestedCategoryFallsBackWhenCategoryIsUnknown(t *testing.T) {
	categoryID := int64(1)
	projects := []*domProduct.Project{{ID: 1, ProjectCategoryID: &categoryID, Name: "Проект"}}

	filtered := filterByRequestedCategory("мебель в прихожую", projects, map[int64]string{categoryID: "Кухни"})
	if len(filtered) != len(projects) {
		t.Fatalf("expected original projects for unknown category query")
	}
}
