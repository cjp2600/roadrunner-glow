package glow

// SequenceType
type SequenceType string

const (
	Parallel SequenceType = "parallel"
	Sync     SequenceType = "sync"
)

// Option
type Option struct {
	Debug bool `json:"debug"`
}

// Job
type Job struct {
	Id     string                 `json:"id,omitempty"`
	Url    string                 `json:"url,omitempty"`
	Method string                 `json:"method,omitempty"`
	Body   map[string]interface{} `json:"body,omitempty"`
	Header map[string]interface{} `json:"header,omitempty"`
	Var    []*Var                 `json:"var,omitempty"`
}

// Var
type Var struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	JPath string `json:"jPath,omitempty"`
}

// Sequence
type Sequence struct {
	Type SequenceType `json:"type,omitempty"`
	Jobs []*Job       `json:"jobs,omitempty"`
}

// ExecuteRequest
type ExecuteRequest struct {
	Options  *Option
	Sequence []*Sequence `json:"sequence,omitempty"`
}

// ExecuteResponse
type ExecuteResponse struct {
	Jobs map[string]string `json:"jobs,omitempty"`
}
