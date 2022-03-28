package modproxy

import (
	"archive/zip"
	"sync"
)

// A Modver holds a module path and version.
type modver struct {
	path    string
	version string
}

func (mv modver) String() string {
	return mv.path + "@" + mv.version
}

// cache caches proxy info, mod and zip calls.
type cache struct {
	mu sync.Mutex

	infoCache map[modver]*versionInfo
	modCache  map[modver][]byte

	// One-element zip cache, to avoid a double download.
	// See TestFetchAndUpdateStateCacheZip in internal/worker/fetch_test.go.
	zipKey    modver
	zipReader *zip.Reader
}

func (c *cache) getInfo(modulePath, version string) *versionInfo {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.infoCache[modver{path: modulePath, version: version}]
}

func (c *cache) putInfo(modulePath, version string, v *versionInfo) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.infoCache == nil {
		c.infoCache = map[modver]*versionInfo{}
	}
	c.infoCache[modver{path: modulePath, version: version}] = v
}

func (c *cache) getMod(modulePath, version string) []byte {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.modCache[modver{path: modulePath, version: version}]
}

func (c *cache) putMod(modulePath, version string, b []byte) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.modCache == nil {
		c.modCache = map[modver][]byte{}
	}
	c.modCache[modver{path: modulePath, version: version}] = b
}

func (c *cache) getZip(modulePath, version string) *zip.Reader {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.zipKey == (modver{path: modulePath, version: version}) {
		return c.zipReader
	}
	return nil
}

func (c *cache) putZip(modulePath, version string, r *zip.Reader) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.zipKey = modver{path: modulePath, version: version}
	c.zipReader = r
}
