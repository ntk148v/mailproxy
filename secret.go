package main

// secret special type for storing secrets.
type secret string

// MarshalYAML implements the yaml.Marshaler interface for Secrets.
func (s secret) MarshalYAML() (interface{}, error) {
	if s != "" {
		return "<secret>", nil
	}
	return nil, nil
}

//UnmarshalYAML implements the yaml.Unmarshaler interface for Secrets.
func (s *secret) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain secret
	return unmarshal((*plain)(s))
}
