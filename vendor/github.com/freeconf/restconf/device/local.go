package device

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/freeconf/yang/fc"
	"github.com/freeconf/yang/meta"
	"github.com/freeconf/yang/node"
	"github.com/freeconf/yang/nodeutil"
	"github.com/freeconf/yang/parser"
	"github.com/freeconf/yang/source"
)

type Local struct {
	browsers     map[string]*node.Browser
	schemaSource source.Opener
	uiSource     source.Opener
}

func New(schemaSource source.Opener) *Local {
	return &Local{
		schemaSource: schemaSource,
		browsers:     make(map[string]*node.Browser),
	}
}

func NewWithUi(schemaSource source.Opener, uiSource source.Opener) *Local {
	return &Local{
		schemaSource: schemaSource,
		uiSource:     uiSource,
		browsers:     make(map[string]*node.Browser),
	}
}

func (self *Local) SchemaSource() source.Opener {
	return self.schemaSource
}

func (self *Local) UiSource() source.Opener {
	return self.uiSource
}

func (self *Local) Modules() map[string]*meta.Module {
	mods := make(map[string]*meta.Module)
	for _, b := range self.browsers {
		mods[b.Meta.Ident()] = b.Meta
	}
	return mods
}

func (self *Local) Browser(module string) (*node.Browser, error) {
	return self.browsers[module], nil
}

func (self *Local) Close() {
}

func (self *Local) Add(module string, n node.Node) error {
	m, err := parser.LoadModule(self.schemaSource, module)
	if err != nil {
		return err
	}
	self.browsers[module] = node.NewBrowser(m, n)
	return nil
}

func (self *Local) AddSource(module string, src func() node.Node) error {
	m, err := parser.LoadModule(self.schemaSource, module)
	if err != nil {
		return err
	}
	self.browsers[module] = node.NewBrowserSource(m, src)
	return nil
}

func (self *Local) AddBrowser(b *node.Browser) {
	self.browsers[b.Meta.Ident()] = b
}

func (self *Local) ApplyStartupConfig(config io.Reader) error {
	var cfg map[string]interface{}
	if err := json.NewDecoder(config).Decode(&cfg); err != nil {
		return err
	}
	return self.ApplyStartupConfigData(cfg)
}

func (self *Local) ApplyStartupConfigData(config map[string]interface{}) error {
	for module, data := range config {
		b, err := self.Browser(module)
		if err != nil {
			return err
		}
		if b == nil {
			return fmt.Errorf("%w. browser not found: %s", fc.NotFoundError, module)
		}
		moduleCfg := data.(map[string]interface{})
		if err := b.Root().UpsertFromSetDefaults(nodeutil.ReflectChild(moduleCfg)).LastErr; err != nil {
			return err
		}
	}
	return nil
}

func (self *Local) ApplyStartupConfigFile(fname string) error {
	cfgRdr, err := os.Open(fname)
	defer cfgRdr.Close()
	if err != nil {
		panic(err)
	}
	return self.ApplyStartupConfig(cfgRdr)
}
