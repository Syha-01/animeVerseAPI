package validator

// We will create a new type named Validator
type Validator struct {
	Errors map[string]string
}

// Construct a new Validator and return a pointer to it
// All validation errors go into this one Validator instance
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// IsEmpty checks to see if the Validator's map contains any entries
func (v *Validator) IsEmpty() bool {
	return len(v.Errors) == 0
}

// AddError adds a new error entry to the Validator's error map
// Check first if an entry with the same key does not already exist
func (v *Validator) AddError(key string, message string) {
	_, exists := v.Errors[key]
	if !exists {
		v.Errors[key] = message
	}
}

// Check adds an error to the map only if a validation check is not 'ok'.
func (v *Validator) Check(ok bool, key string, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// PermittedValue checks if a string value is in a list of permitted values.
func PermittedValue(value string, permittedValues ...string) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
