package adapters

import "github.com/3-lines-studio/datafrost/internal/models"

// TreeLister is an optional interface that adapters can implement to
// return a hierarchical database → schema → table structure.
//
// Adapters that don't support it can omit the method; callers can use a
// type assertion to detect support and fall back to flat table listing.
type TreeLister interface {
	ListTree() ([]models.TreeNode, error)
}
