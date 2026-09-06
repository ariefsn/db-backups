package model

import (
	"net/url"
	"strings"
)

// ConnectionParts holds the connection details derived from a connection URI.
type ConnectionParts struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

// ParseConnectionURI extracts the connection details from a connection URI.
// It returns false when the URI is empty or cannot be parsed.
func ParseConnectionURI(uri string) (ConnectionParts, bool) {
	if uri == "" {
		return ConnectionParts{}, false
	}

	u, err := url.Parse(uri)
	if err != nil {
		return ConnectionParts{}, false
	}

	parts := ConnectionParts{
		Host: u.Hostname(),
		Port: u.Port(),
	}

	if u.User != nil {
		parts.Username = u.User.Username()
		parts.Password, _ = u.User.Password()
	}

	// Path is "/dbname" (mongo, postgres, mysql) or "/0" (redis db index).
	// Anything after the first segment is not part of the database name.
	if path := strings.TrimPrefix(u.Path, "/"); path != "" {
		parts.Database, _, _ = strings.Cut(path, "/")
	}

	return parts, true
}

// applyConnectionParts fills the empty targets from parts. Values that are
// already set are never overwritten.
func applyConnectionParts(parts ConnectionParts, host, port, username, password, database *string) {
	fill := func(target *string, value string) {
		if *target == "" {
			*target = value
		}
	}

	fill(host, parts.Host)
	fill(port, parts.Port)
	fill(username, parts.Username)
	fill(password, parts.Password)
	fill(database, parts.Database)
}

// ApplyConnectionURI fills the empty connection fields from ConnectionURI.
// Fields that are already set are left untouched, and an unparsable URI is
// silently ignored so a backup never fails because of it.
func (r *BackupRequest) ApplyConnectionURI() {
	parts, ok := ParseConnectionURI(r.ConnectionURI)
	if !ok {
		return
	}

	applyConnectionParts(parts, &r.Host, &r.Port, &r.Username, &r.Password, &r.Database)
}

// ApplyConnectionURI fills the empty connection fields from ConnectionURI so a
// database saved with only a URI still records its host and database name.
func (d *Database) ApplyConnectionURI() {
	parts, ok := ParseConnectionURI(d.ConnectionURI)
	if !ok {
		return
	}

	applyConnectionParts(parts, &d.Host, &d.Port, &d.Username, &d.Password, &d.Database)
}
