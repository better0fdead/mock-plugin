package parser

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/better0fdead/mock-plugin/parser/annotation"
)

// MethodInfo holds information about an interface method
// MethodPackageInfo holds information about an interface method
type MethodPkgInfo struct {
	Name        string          `json:"name,omitempty"`
	Parameters  []FieldPkgInfo  `json:"parameters,omitempty"`
	Returns     []FieldPkgInfo  `json:"returns,omitempty"`
	Annotations annotation.Tags `json:"annotations,omitempty"`
}

// FieldInfo holds information about a struct field
type FieldPkgInfo struct {
	Name        string          `json:"name,omitempty"`
	Kind        string          `json:"kind,omitempty"`
	IsVariadic  bool            `json:"isVariadic,omitempty"`
	IsPointer   bool            `json:"isPointer,omitempty"`
	IsScalar    bool            `json:"isScalar,omitempty"`
	IsInline    bool            `json:"isInline,omitempty"`
	Annotations annotation.Tags `json:"annotations,omitempty"`
}

// TypePkgInfo holds information about a type
type TypePkgInfo struct {
	Kind      string `json:"kind,omitempty"`
	IsScalar  bool   `json:"isScalar,omitempty"`
	IsPointer bool   `json:"isPointer,omitempty"`
}

// InterfacePackageInfo hold information about an interface
type InterfacePkgInfo struct {
	Name        string          `json:"name,omitempty"`
	Methods     []MethodPkgInfo `json:"methods,omitempty"`
	Annotations annotation.Tags `json:"annotations,omitempty"`
	Pkg         string          `json:"pkg,omitempty"`
}

// InterfaceInfo holds information about an interface
type InterfaceInfo struct {
	Name        string          `json:"name,omitempty"`
	Methods     []MethodInfo    `json:"methods,omitempty"`
	Annotations annotation.Tags `json:"annotations,omitempty"`
}

// MethodInfo holds information about an interface method
type MethodInfo struct {
	Name        string          `json:"name,omitempty"`
	Parameters  []FieldInfo     `json:"parametrs,omitempty"`
	Returns     []FieldInfo     `json:"returns,omitempty"`
	Annotations annotation.Tags `json:"annotations,omitempty"`
}

// FieldInfo holds information about a struct field
type FieldInfo struct {
	Name        string          `json:"name,omitempty"`
	Type        TypeInfo        `json:"type,omitempty"`
	IsVariadic  bool            `json:"isVariadic,omitempty"`
	IsPointer   bool            `json:"isPointer,omitempty"`
	Annotations annotation.Tags `json:"annotations,omitempty"`
	IsInline    bool            `json:"isInline,omitempty"`
}

// TypeInfo holds information about a type
type TypeInfo struct {
	Name        string          `json:"name,omitempty"`
	Detail      TypeDetail      `json:"detail,omitempty"`
	Inline      bool            `json:"inline,omitempty"`
	IsPointer   bool            `json:"isPointer,omitempty"`
	IsAlias     bool            `json:"isAlias,omitempty"`
	IsEmbedded  bool            `json:"isEmbedded,omitempty"`
	IsScalar    bool            `json:"isScalar,omitempty"`
	Annotations annotation.Tags `json:"annotations,omitempty"`
	Pkg         string          `json:"pkg,omitempty"`
}

func (t TypeInfo) MarshalJSON() ([]byte, error) {
	type Alias TypeInfo
	aux := struct {
		Alias
		Kind        string           `json:"kind,omitempty"`
		KeyType     *TypeDescriptor  `json:"keyType,omitempty"`
		ValueType   *TypeDescriptor  `json:"valueType,omitempty"`
		ElementType *TypeDescriptor  `json:"elementType,omitempty"`
		Length      int              `json:"length,omitempty"`
		Fields      []FieldTypeInfo  `json:"fields,omitempty"`
		Methods     []MethodPkgInfo  `json:"methods,omitempty"`
		Parameters  []TypeDescriptor `json:"parameters,omitempty"`
		Returns     []TypeDescriptor `json:"returns,omitempty"`
	}{
		Alias: (Alias)(t),
	}

	// Populate the fields from TypeDetail based on its actual type
	switch detail := t.Detail.(type) {
	case MapTypeDetail:
		aux.Kind = detail.Kind
		if detail.KeyType != (TypeDescriptor{}) {
			aux.KeyType = &detail.KeyType
		}
		if detail.ValueType != (TypeDescriptor{}) {
			aux.ValueType = &detail.ValueType
		}
	case SliceTypeDetail:
		aux.Kind = detail.Kind
		if detail.ElementType != (TypeDescriptor{}) {
			aux.ElementType = &detail.ElementType
		}
	case ChannelTypeDetail:
		aux.Kind = detail.Kind
		if detail.ValueType != (TypeDescriptor{}) {
			aux.ValueType = &detail.ValueType
		}
	case StructTypeDetail:
		aux.Kind = detail.Kind
		aux.Fields = detail.Fields
	case InterfaceTypeDetail:
		aux.Kind = detail.Kind
		aux.Methods = detail.Methods
	case ArrayTypeDetail:
		aux.Kind = detail.Kind
		if detail.ElementType != (TypeDescriptor{}) {
			aux.ElementType = &detail.ElementType
		}
		aux.Length = detail.Length
	case FunctionTypeDetail:
		aux.Kind = detail.Kind
		aux.Parameters = detail.Parameters
		aux.Returns = detail.Returns
	case BasicTypeDetail:
		aux.Kind = detail.Kind
	default:
		return nil, nil
	}
	aux.Detail = nil

	return json.Marshal(aux)
}

// TypeDetail is an interface for detailed type information
type TypeDetail interface {
	GetKind() string
}

// MapTypeDetail holds information about a map type
type MapTypeDetail struct {
	Kind      string         `json:"kind,omitempty"`
	KeyType   TypeDescriptor `json:"keyType,omitempty"`
	ValueType TypeDescriptor `json:"valueType,omitempty"`
}

func (m MapTypeDetail) GetKind() string {
	return m.Kind
}

// SliceTypeDetail holds information about a slice type
type SliceTypeDetail struct {
	Kind        string         `json:"kind,omitempty"`
	ElementType TypeDescriptor `json:"elementType,omitempty"`
}

func (s SliceTypeDetail) GetKind() string {
	return s.Kind
}

// ChannelTypeDetail holds information about a channel type
type ChannelTypeDetail struct {
	Kind      string         `json:"kind,omitempty"`
	ValueType TypeDescriptor `json:"valueType,omitempty"`
}

func (c ChannelTypeDetail) GetKind() string {
	return c.Kind
}

// FieldInfo holds information about a struct field
type FieldTypeInfo struct {
	Name       string `json:"name,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Inline     bool   `json:"inline,omitempty"`
	IsPointer  bool   `json:"isPointer,omitempty"`
	IsAlias    bool   `json:"isAlias,omitempty"`
	IsEmbedded bool   `json:"isEmbedded,omitempty"`
	IsScalar   bool   `json:"isScalar,omitempty"`
}

// TypePkgInfo holds information about a type
type TypeDescriptor struct {
	Kind       string `json:"kind,omitempty"`
	Inline     bool   `json:"inline,omitempty"`
	IsPointer  bool   `json:"isPointer,omitempty"`
	IsAlias    bool   `json:"isAlias,omitempty"`
	IsEmbedded bool   `json:"isEmbedded,omitempty"`
	IsScalar   bool   `json:"isScalar,omitempty"`
}

// StructTypeDetail holds information about a struct type
type StructTypeDetail struct {
	Kind   string          `json:"kind,omitempty"`
	Fields []FieldTypeInfo `json:"fields,omitempty"`
}

func (s StructTypeDetail) GetKind() string {
	return s.Kind
}

// InterfaceTypeDetail holds information about an interface type
type InterfaceTypeDetail struct {
	Kind    string          `json:"kind,omitempty"`
	Methods []MethodPkgInfo `json:"methods,omitempty"`
}

func (i InterfaceTypeDetail) GetKind() string {
	return i.Kind
}

// ArrayTypeDetail holds information about an array type
type ArrayTypeDetail struct {
	Kind        string         `json:"kind,omitempty"`
	ElementType TypeDescriptor `json:"elementType,omitempty"`
	Length      int            `json:"length,omitempty"`
}

func (a ArrayTypeDetail) GetKind() string {
	return a.Kind
}

// FunctionTypeDetail holds information about a function type
type FunctionTypeDetail struct {
	Kind       string           `json:"kind,omitempty"`
	Parameters []TypeDescriptor `json:"parameters,omitempty"`
	Returns    []TypeDescriptor `json:"returns,omitempty"`
}

func (f FunctionTypeDetail) GetKind() string {
	return f.Kind
}

// BasicTypeDetail holds information about a basic type
type BasicTypeDetail struct {
	Kind string `json:"kind,omitempty"`
}

func (b BasicTypeDetail) GetKind() string {
	return b.Kind
}

// PackageInfo holds information about a package
type PackageInfo struct {
	Imports     map[string]string   `json:"imports,omitempty"`
	Services    []InterfacePkgInfo  `json:"services,omitempty"`
	Types       map[string]TypeInfo `json:"types,omitempty"`
	Annotations annotation.Tags     `json:"annotations,omitempty"`
}

// SanitizeKey sanitizes a given key by replacing non-alphanumeric characters with underscores
func SanitizeKey(key string) string {
	// Check if the key represents a slice type
	if strings.HasPrefix(key, "[]") {
		elementType := key[2:]
		return "slice_" + SanitizeKey(elementType)
	}

	// Check if the key represents an array type
	if strings.HasPrefix(key, "[") && strings.Contains(key, "]") {
		endIndex := strings.Index(key, "]")
		elementType := key[endIndex+1:]
		return "array_" + key[1:endIndex] + "_" + SanitizeKey(elementType)
	}

	// Check if the key represents a function type
	if strings.HasPrefix(key, "func(") {
		return "func_" + regexp.MustCompile(`[^a-zA-Z0-9_.]`).ReplaceAllString(key, "_")
	}

	// Check if the key represents a map type
	if strings.HasPrefix(key, "map[") {
		endIndex := strings.Index(key, "]")
		keyType := key[4:endIndex]
		valueType := key[endIndex+1:]
		return "map_" + SanitizeKey(keyType) + "_" + SanitizeKey(valueType)
	}

	// Check if the key represents a pointer type
	if strings.HasPrefix(key, "*") {
		elementType := key[1:]
		return "ptr_" + SanitizeKey(elementType)
	}

	// Check if the key represents a channel type
	if strings.HasPrefix(key, "chan ") {
		elementType := key[5:]
		return "chan_" + SanitizeKey(elementType)
	}

	// Use a regular expression to replace any non-alphanumeric characters with underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9_.]`)
	return re.ReplaceAllString(key, "_")
}
