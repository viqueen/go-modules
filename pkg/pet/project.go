// Package pet provides functionality to manage and retrieve project names.
package pet

// Project is an interface that defines a method to get the project name.
type Project interface {
	Get() string
}
