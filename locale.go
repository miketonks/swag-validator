package swagvalidator

type (
	// locale is an interface for defining custom error strings
	locale interface {
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

func (l CustomLocale) Required() string {
	return `Is required`
}

func (l CustomLocale) InvalidType() string {
	return `Invalid type. Expected: {{.expected}}, given: {{.given}}`
}

func (l CustomLocale) NumberAnyOf() string {
	return `Must validate at least one schema (anyOf)`
}

func (l CustomLocale) NumberOneOf() string {
	return `Must validate one and only one schema (oneOf)`
}

func (l CustomLocale) NumberAllOf() string {
	return `Must validate all the schemas (allOf)`
}

func (l CustomLocale) NumberNot() string {
	return `Must not validate the schema (not)`
}

func (l CustomLocale) MissingDependency() string {
	return `Has a dependency on {{.dependency}}`
}

func (l CustomLocale) Internal() string {
	return `Internal Error {{.error}}`
}

func (l CustomLocale) Const() string {
	return `Does not match: {{.allowed}}`
}

func (l CustomLocale) Enum() string {
	return `Must be one of the following: {{.allowed}}`
}

func (l CustomLocale) ArrayNoAdditionalItems() string {
	return `No additional items allowed on array`
}

func (l CustomLocale) ArrayNotEnoughItems() string {
	return `Not enough items on array to match positional list of schema`
}

func (l CustomLocale) ArrayMinItems() string {
	return `Array must have at least {{.min}} items`
}

func (l CustomLocale) ArrayMaxItems() string {
	return `Array must have at most {{.max}} items`
}

func (l CustomLocale) Unique() string {
	return `{{.type}} items[{{.i}},{{.j}}] must be unique`
}

func (l CustomLocale) ArrayContains() string {
	return `At least one of the items must match`
}

func (l CustomLocale) ArrayMinProperties() string {
	return `Must have at least {{.min}} properties`
}

func (l CustomLocale) ArrayMaxProperties() string {
	return `Must have at most {{.max}} properties`
}

func (l CustomLocale) AdditionalPropertyNotAllowed() string {
	return `Is not allowed as an additional property`
}

func (l CustomLocale) InvalidPropertyPattern() string {
	return `Property does not match pattern {{.pattern}}`
}

func (l CustomLocale) InvalidPropertyName() string {
	return `Property name of "{{.property}}" does not match`
}

func (l CustomLocale) StringGTE() string {
	return `String length must be greater than or equal to {{.min}}`
}

func (l CustomLocale) StringLTE() string {
	return `String length must be less than or equal to {{.max}}`
}

func (l CustomLocale) DoesNotMatchPattern() string {
	return `Does not match pattern '{{.pattern}}'`
}

func (l CustomLocale) DoesNotMatchFormat() string {
	return `Field does not match format '{{.format}}'`
}

func (l CustomLocale) MultipleOf() string {
	return `Must be a multiple of {{.multiple}}`
}

func (l CustomLocale) NumberGTE() string {
	return `Must be greater than or equal to {{.min}}`
}

func (l CustomLocale) NumberGT() string {
	return `Must be greater than {{.min}}`
}

func (l CustomLocale) NumberLTE() string {
	return `Must be less than or equal to {{.max}}`
}

func (l CustomLocale) NumberLT() string {
	return `Must be less than {{.max}}`
}

// Schema validators
func (l CustomLocale) RegexPattern() string {
	return `Invalid regex pattern '{{.pattern}}'`
}

func (l CustomLocale) GreaterThanZero() string {
	return `{{.number}} must be strictly greater than 0`
}

func (l CustomLocale) MustBeOfA() string {
	return `{{.x}} must be of a {{.y}}`
}

func (l CustomLocale) MustBeOfAn() string {
	return `{{.x}} must be of an {{.y}}`
}

func (l CustomLocale) CannotBeUsedWithout() string {
	return `{{.x}} cannot be used without {{.y}}`
}

func (l CustomLocale) CannotBeGT() string {
	return `{{.x}} cannot be greater than {{.y}}`
}

func (l CustomLocale) MustBeOfType() string {
	return `{{.key}} must be of type {{.type}}`
}

func (l CustomLocale) MustBeValidRegex() string {
	return `{{.key}} must be a valid regex`
}

func (l CustomLocale) MustBeValidFormat() string {
	return `{{.key}} must be a valid format {{.given}}`
}

func (l CustomLocale) MustBeGTEZero() string {
	return `{{.key}} must be greater than or equal to 0`
}

func (l CustomLocale) KeyCannotBeGreaterThan() string {
	return `{{.key}} cannot be greater than {{.y}}`
}

func (l CustomLocale) KeyItemsMustBeOfType() string {
	return `{{.key}} items must be {{.type}}`
}

func (l CustomLocale) KeyItemsMustBeUnique() string {
	return `{{.key}} items must be unique`
}

func (l CustomLocale) ReferenceMustBeCanonical() string {
	return `Reference {{.reference}} must be canonical`
}

func (l CustomLocale) NotAValidType() string {
	return `has a primitive type that is NOT VALID -- given: {{.given}} Expected valid values are:{{.expected}}`
}

func (l CustomLocale) Duplicated() string {
	return `{{.type}} type is duplicated`
}

func (l CustomLocale) HttpBadStatus() string {
	return `Could not read schema from HTTP, response status is {{.status}}`
}

// Replacement options: field, description, context, value
func (l CustomLocale) ErrorFormat() string {
	return `{{.field}}: {{.description}}`
}

//Parse error
func (l CustomLocale) ParseError() string {
	return `Expected: {{.expected}}, given: Invalid JSON`
}

//If/Else
func (l CustomLocale) ConditionThen() string {
	return `Must validate "then" as "if" was valid`
}

func (l CustomLocale) ConditionElse() string {
	return `Must validate "else" as "i"`
}
