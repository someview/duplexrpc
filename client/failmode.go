package client

// FailMode decides how clients action when clients fail to invoke services
type FailMode int

const (
	// Failover selects another server automaticaly
	Failover FailMode = iota
	// Failfast returns error immediately
	Failfast
	// Failtry use current client again
	Failtry
	// Failbackup select another server if the first server doesn't respond in specified time and use the fast response.
	Failbackup
)
