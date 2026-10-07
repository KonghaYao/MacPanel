package mirrors

import (
	"net/url"
	"strings"
)

func validateField(ecosystem string, field Field, value string) error {
	if value == "" {
		return nil
	}
	switch field.Kind {
	case KindURL:
		return validateHTTPURL(ecosystem, field.Key, value, false)
	case KindHost:
		return validateHost(ecosystem, field.Key, value)
	case KindURLs:
		for _, item := range splitURLList(value) {
			if err := validateHTTPURL(ecosystem, field.Key, item, false); err != nil {
				return err
			}
		}
		return nil
	case KindGoproxy:
		return validateGoproxy(ecosystem, value)
	case KindGosumdb:
		return validateGosumdb(ecosystem, value)
	case KindCargo:
		return validateCargoRegistry(ecosystem, value)
	default:
		return invalidf("unknown field kind %q", field.Kind)
	}
}

func validateHTTPURL(ecosystem, field, raw string, allowSparse bool) error {
	if strings.ContainsAny(raw, " \t\r\n'\"\\") {
		return invalidf("%s %s contains unsupported characters", ecosystem, field)
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return invalidf("%s %s is not a valid URL", ecosystem, field)
	}
	switch parsed.Scheme {
	case "http", "https":
	case "sparse+http", "sparse+https":
		if !allowSparse {
			return invalidf("%s %s must use http or https", ecosystem, field)
		}
	default:
		if allowSparse {
			return invalidf("%s %s must use http, https, or a sparse+ registry URL", ecosystem, field)
		}
		return invalidf("%s %s must use http or https", ecosystem, field)
	}
	if parsed.Host == "" {
		return invalidf("%s %s is missing a host", ecosystem, field)
	}
	return nil
}

func validateHost(ecosystem, field, value string) error {
	if strings.ContainsAny(value, " \t\r\n/\\'\"@:") {
		return invalidf("%s %s must be a host name", ecosystem, field)
	}
	if value == "" || strings.Contains(value, "..") || !strings.Contains(value, ".") {
		return invalidf("%s %s must be a host name", ecosystem, field)
	}
	return nil
}

func validateCargoRegistry(ecosystem, raw string) error {
	if raw == "" {
		return nil
	}
	return validateHTTPURL(ecosystem, "registry", raw, true)
}

func validateGoproxy(ecosystem, value string) error {
	if strings.ContainsAny(value, " \t\r\n'\"\\") {
		return invalidf("%s goproxy contains unsupported characters", ecosystem)
	}
	for _, part := range strings.Split(value, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			return invalidf("%s goproxy contains an empty entry", ecosystem)
		}
		if part == "off" || part == "direct" {
			continue
		}
		if err := validateHTTPURL(ecosystem, "goproxy", part, false); err != nil {
			return err
		}
	}
	return nil
}

func validateGosumdb(ecosystem, value string) error {
	if value == "off" {
		return nil
	}
	if strings.ContainsAny(value, " \t\r\n/\\'\"") {
		return invalidf("%s gosumdb must be a host or off", ecosystem)
	}
	host := value
	if idx := strings.Index(value, "+"); idx >= 0 {
		host = value[:idx]
		if value[idx+1:] == "" {
			return invalidf("%s gosumdb is missing its public key", ecosystem)
		}
	}
	if host == "" || strings.Contains(host, "..") || !strings.Contains(host, ".") {
		return invalidf("%s gosumdb must be a host or off", ecosystem)
	}
	return nil
}
