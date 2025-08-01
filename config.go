package plugin

import (
	"time"

	"github.com/notioncodes/types"
)

// CommonSettings is the configuration that is common to all plugins.
type CommonSettings struct {
	Reporter        bool          `json:"reporter" yaml:"reporter"`
	Workers         int           `json:"workers" yaml:"workers"`
	RuntimeTimeout  time.Duration `json:"runtime_timeout" yaml:"runtime_timeout"`
	RequestDelay    time.Duration `json:"request_delay" yaml:"request_delay"`
	ContinueOnError bool          `json:"continue_on_error" yaml:"continue_on_error"`
}

// ContentSettings is the configuration that determines what content is exported.
type ContentSettings struct {
	Flush     bool               `json:"flush" yaml:"flush"`
	Types     []types.ObjectType `json:"object_types" yaml:"object_types"`
	Databases DatabaseSettings   `json:"databases" yaml:"databases"`
	Pages     PageSettings       `json:"pages" yaml:"pages"`
	Blocks    BlockSettings      `json:"blocks" yaml:"blocks"`
}

// DatabaseSettings is the configuration that determines what databases are exported.
type DatabaseSettings struct {
	IDs    []types.DatabaseID `json:"ids" yaml:"ids"`
	Pages  bool               `json:"pages" yaml:"pages"`
	Blocks bool               `json:"blocks" yaml:"blocks"`
}

// PageSettings is the configuration that determines what pages are exported.
type PageSettings struct {
	IDs         []types.PageID `json:"ids" yaml:"ids"`
	Blocks      bool           `json:"blocks" yaml:"blocks"`
	Comments    bool           `json:"comments" yaml:"comments"`
	Attachments bool           `json:"attachments" yaml:"attachments"`
}

// BlockSettings is the configuration that determines what blocks are exported.
type BlockSettings struct {
	IDs      []types.BlockID `json:"ids" yaml:"ids"`
	Children bool            `json:"children" yaml:"children"`
}
