package swagvalidator

type (
	// locale is an interface for defining custom error strings
	locale interface {
		False() string
		Required() string
		InvalidType() string
		NumberAnyOf() string
		NumberOneOf() string
		NumberAllOf() string
		NumberNot() string
		MissingDependency() string
		Internal() string
		Const() string
		Enum() string
		ArrayNotEnoughItems() string
		ArrayNoAdditionalItems() string
		ArrayMinItems() string
		ArrayMaxItems() string
		Unique() string
		ArrayContains() string
		ArrayMinProperties() string
		ArrayMaxProperties() string
		AdditionalPropertyNotAllowed() string
		InvalidPropertyPattern() string
		InvalidPropertyName() string
		StringGTE() string
		StringLTE() string
		DoesNotMatchPattern() string
		DoesNotMatchFormat() string
		MultipleOf() string
		NumberGTE() string
		NumberGT() string
		NumberLTE() string
		NumberLT() string

		// Schema validations
		RegexPattern() string
		GreaterThanZero() string
		MustBeOfA() string
		MustBeOfAn() string
		CannotBeUsedWithout() string
		CannotBeGT() string
		MustBeOfType() string
		MustBeValidRegex() string
		MustBeValidFormat() string
		MustBeGTEZero() string
		KeyCannotBeGreaterThan() string
		KeyItemsMustBeOfType() string
		KeyItemsMustBeUnique() string
		ReferenceMustBeCanonical() string
		NotAValidType() string
		Duplicated() string
		HttpBadStatus() string
		ParseError() string

		ConditionThen() string
		ConditionElse() string

		// ErrorFormat
		ErrorFormat() string
	}

	// CustomLocale is a locale for schema validator
	CustomLocale struct{}
)

// False returns a format-string for "false" schema validation errors
func (l CustomLocale) False() string {
	return "False always fails validation"
}

// Required returns a format-string for "required" schema validation errors
func (l CustomLocale) Required() string {
	return `{{.property}} is required`
}

// InvalidType ...
func (l CustomLocale) InvalidType() string {
	return `Invalid type. Expected: {{.expected}}, given: {{.given}}`
}

// NumberAnyOf ...
func (l CustomLocale) NumberAnyOf() string {
	return `Must validate at least one schema (anyOf)`
}

// NumberOneOf ...
func (l CustomLocale) NumberOneOf() string {
	return `Must validate one and only one schema (oneOf)`
}

// NumberAllOf ...
func (l CustomLocale) NumberAllOf() string {
	return `Must validate all the schemas (allOf)`
}

// NumberNot ...
func (l CustomLocale) NumberNot() string {
	return `Must not validate the schema (not)`
}

// MissingDependency ...
func (l CustomLocale) MissingDependency() string {
	return `Has a dependency on {{.dependency}}`
}

// Internal ...
func (l CustomLocale) Internal() string {
	return `Internal Error {{.error}}`
}

// Const ...
func (l CustomLocale) Const() string {
	return `Does not match: {{.allowed}}`
}

// Enum ...
func (l CustomLocale) Enum() string {
	return `Must be one of the following: {{.allowed}}`
}

// ArrayNoAdditionalItems ...
func (l CustomLocale) ArrayNoAdditionalItems() string {
	return `No additional items allowed on array`
}

// ArrayNotEnoughItems ...
func (l CustomLocale) ArrayNotEnoughItems() string {
	return `Not enough items on array to match positional list of schema`
}

// ArrayMinItems ...
func (l CustomLocale) ArrayMinItems() string {
	return `Array must have at least {{.min}} items`
}

// ArrayMaxItems ...
func (l CustomLocale) ArrayMaxItems() string {
	return `Array must have at most {{.max}} items`
}

// Unique ...
func (l CustomLocale) Unique() string {
	return `{{.type}} items[{{.i}},{{.j}}] must be unique`
}

// ArrayContains ...
func (l CustomLocale) ArrayContains() string {
	return `At least one of the items must match`
}

// ArrayMinProperties ...
func (l CustomLocale) ArrayMinProperties() string {
	return `Must have at least {{.min}} properties`
}

// ArrayMaxProperties ...
func (l CustomLocale) ArrayMaxProperties() string {
	return `Must have at most {{.max}} properties`
}

// AdditionalPropertyNotAllowed ...
func (l CustomLocale) AdditionalPropertyNotAllowed() string {
	return `Is not allowed as an additional property`
}

// InvalidPropertyPattern ...
func (l CustomLocale) InvalidPropertyPattern() string {
	return `Property does not match pattern {{.pattern}}`
}

// InvalidPropertyName ...
func (l CustomLocale) InvalidPropertyName() string {
	return `Property name of "{{.property}}" does not match`
}

// StringGTE ...
func (l CustomLocale) StringGTE() string {
	return `String length must be greater than or equal to {{.min}}`
}

// StringLTE ...
func (l CustomLocale) StringLTE() string {
	return `String length must be less than or equal to {{.max}}`
}

// DoesNotMatchPattern ...
func (l CustomLocale) DoesNotMatchPattern() string {
	return `Does not match pattern '{{.pattern}}'`
}

// DoesNotMatchFormat ...
func (l CustomLocale) DoesNotMatchFormat() string {
	return `Field does not match format '{{.format}}'`
}

// MultipleOf ...
func (l CustomLocale) MultipleOf() string {
	return `Must be a multiple of {{.multiple}}`
}

// NumberGTE ...
func (l CustomLocale) NumberGTE() string {
	return `Must be greater than or equal to {{.min}}`
}

// NumberGT ...
func (l CustomLocale) NumberGT() string {
	return `Must be greater than {{.min}}`
}

// NumberLTE ...
func (l CustomLocale) NumberLTE() string {
	return `Must be less than or equal to {{.max}}`
}

// NumberLT ...
func (l CustomLocale) NumberLT() string {
	return `Must be less than {{.max}}`
}

// RegexPattern ...
func (l CustomLocale) RegexPattern() string {
	return `Invalid regex pattern '{{.pattern}}'`
}

// GreaterThanZero ...
func (l CustomLocale) GreaterThanZero() string {
	return `{{.number}} must be strictly greater than 0`
}

// MustBeOfA ...
func (l CustomLocale) MustBeOfA() string {
	return `{{.x}} must be of a {{.y}}`
}

// MustBeOfAn ...
func (l CustomLocale) MustBeOfAn() string {
	return `{{.x}} must be of an {{.y}}`
}

// CannotBeUsedWithout ...
func (l CustomLocale) CannotBeUsedWithout() string {
	return `{{.x}} cannot be used without {{.y}}`
}

// CannotBeGT ...
func (l CustomLocale) CannotBeGT() string {
	return `{{.x}} cannot be greater than {{.y}}`
}

// MustBeOfType ...
func (l CustomLocale) MustBeOfType() string {
	return `{{.key}} must be of type {{.type}}`
}

// MustBeValidRegex ...
func (l CustomLocale) MustBeValidRegex() string {
	return `{{.key}} must be a valid regex`
}

// MustBeValidFormat ...
func (l CustomLocale) MustBeValidFormat() string {
	return `{{.key}} must be a valid format {{.given}}`
}

// MustBeGTEZero ...
func (l CustomLocale) MustBeGTEZero() string {
	return `{{.key}} must be greater than or equal to 0`
}

// KeyCannotBeGreaterThan ...
func (l CustomLocale) KeyCannotBeGreaterThan() string {
	return `{{.key}} cannot be greater than {{.y}}`
}

// KeyItemsMustBeOfType ...
func (l CustomLocale) KeyItemsMustBeOfType() string {
	return `{{.key}} items must be {{.type}}`
}

// KeyItemsMustBeUnique ...
func (l CustomLocale) KeyItemsMustBeUnique() string {
	return `{{.key}} items must be unique`
}

// ReferenceMustBeCanonical ...
func (l CustomLocale) ReferenceMustBeCanonical() string {
	return `Reference {{.reference}} must be canonical`
}

// NotAValidType ...
func (l CustomLocale) NotAValidType() string {
	return `has a primitive type that is NOT VALID -- given: {{.given}} Expected valid values are:{{.expected}}`
}

// Duplicated ...
func (l CustomLocale) Duplicated() string {
	return `{{.type}} type is duplicated`
}

// HttpBadStatus ...
func (l CustomLocale) HttpBadStatus() string { // nolint
	return `Could not read schema from HTTP, response status is {{.status}}`
}

// ErrorFormat ...
func (l CustomLocale) ErrorFormat() string {
	return `{{.field}}: {{.description}}`
}

// ParseError ...
func (l CustomLocale) ParseError() string {
	return `Expected: {{.expected}}, given: Invalid JSON`
}

// ConditionThen ...
func (l CustomLocale) ConditionThen() string {
	return `Must validate "then" as "if" was valid`
}

// ConditionElse ...
func (l CustomLocale) ConditionElse() string {
	return `Must validate "else" as "i"`
}
