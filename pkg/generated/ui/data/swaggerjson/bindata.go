// Code generated by go-bindata.
// sources:
// assets/generated/swagger/api.swagger.json
// DO NOT EDIT!

package swaggerjson

import (
	"github.com/elazarl/go-bindata-assetfs"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _apiSwaggerJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5a\x4f\x6f\xdb\xb8\x12\xbf\xe7\x53\x0c\xf4\x1e\xf0\x2e\xad\x9d\xf6\x54\xe4\xf4\x8a\x04\x68\x83\x36\xef\x15\x9b\xa2\x39\x2c\x0a\x83\xa6\xc6\xd2\xd4\x14\xa9\x90\x94\xb3\xde\x45\xbe\xfb\x82\xa4\x64\xfd\xb1\x64\x3b\x72\x36\xeb\x5d\x6c\x2f\x4d\x4c\xce\xff\xe1\xcc\x6f\xc6\xf9\xed\x0c\x20\x32\x0f\x2c\x49\x50\x47\x17\x10\xbd\x9d\x9c\x47\xaf\xdc\x67\x24\x17\x2a\xba\x00\x77\x0e\x10\x59\xb2\x02\xdd\xf9\xa5\x28\x8c\x45\x0d\x37\x4c\xb2\x04\x35\x7c\xbb\xb9\x63\x1a\xe1\x23\x8a\x1c\x35\xbc\xff\x72\xed\xa9\x01\xa2\x15\x6a\x43\x4a\x3a\x9a\xd5\xf9\xe4\x4d\xc9\x16\x20\xe2\x4a\x5a\xc6\xed\x86\x37\x40\x24\x59\xe6\x99\xdf\x10\x4f\x19\x0a\xf8\x86\x12\x7f\x25\x56\x52\x00\x44\x85\x16\xee\x3c\xb5\x36\x37\x17\xd3\x69\x42\x56\xb0\xf9\x84\xab\x6c\x9a\xad\xc2\xdd\x29\xcf\xd8\xeb\x55\xf6\xc0\x34\xd6\x64\x98\x31\xf2\x84\xd5\xad\xff\x26\xee\x13\x47\x18\xf9\x3b\x8f\x67\x00\x8f\xde\x5c\xc3\x53\xcc\xd0\x44\x17\xf0\x73\x50\xd3\xcb\xaa\x74\x76\xbf\x38\x8a\xef\xfe\x2e\x57\xd2\x14\xad\xcb\x2c\xcf\x05\x71\x66\x49\xc9\xe9\x0f\xa3\x64\x7d\x37\xd7\x2a\x2e\xf8\x81\x77\x99\x4d\x4d\xed\xf3\x29\xcb\x69\xba\x7a\x33\xe5\xc1\xe5\x4d\x87\x25\xd8\xf4\x9f\x53\xbf\xc8\x32\xa6\xd7\xce\xd6\x3b\x12\x02\x34\x5a\x4d\xb8\x42\xb0\x29\x82\xb1\xcc\x16\x06\xd4\x02\x18\x94\xcc\x80\xc9\x18\xc8\x1a\x58\x16\x73\xe4\x4a\x2e\x28\x81\x85\xd2\xc0\x95\x94\xc8\x2d\xad\xc8\xae\x37\x7e\x04\x88\x54\x8e\xda\xab\x7c\x1d\x3b\x19\x1f\xd0\x96\x89\xd0\xbc\xa4\xd1\xe4\x4a\x1a\x34\x2d\xdd\x00\xa2\xb7\xe7\xe7\x9d\x8f\x00\xa2\x18\x0d\xd7\x94\xdb\x32\x49\x1a\x8c\x82\x45\x2e\x20\x6c\x8b\x0c\x20\xfa\xb7\xc6\x85\xa3\xf8\xd7\x34\xc6\x05\x49\x72\x1c\x8c\x8b\x7f\x08\x7f\xad\xdb\x4f\x98\x8b\x75\xd4\x22\x7f\x3c\xeb\xfb\xf9\xb1\x61\x44\xce\x34\xcb\xd0\xa2\xae\x43\x16\xfe\x75\xd4\xaf\x72\xd6\xff\xff\x6a\xa7\x69\xff\x63\x19\x3a\xef\xbb\x58\x54\xfe\xb7\x0a\xe6\x08\x42\xa9\x25\xc6\x50\xe4\x93\x2e\x0b\xf2\x94\xf7\x05\xea\x75\xf7\x48\xe3\x7d\x41\x1a\x5d\x20\x16\x4c\x18\xec\x1c\xdb\x75\xee\x15\x33\x56\x93\x4c\xa2\x5e\x83\xbf\x37\x0c\xb6\x2c\xe9\x9a\x5a\xbd\xf2\x9a\xf8\xfb\x59\xc7\x53\x51\x8c\x02\x2d\xee\xce\xc2\x70\xa7\xce\xba\x1d\x19\x75\xe5\xaf\x9e\x6c\x52\xb5\xd4\x3b\x95\xbc\xba\x4b\x99\x05\x32\xcd\xbc\xfa\x8f\x01\x47\xe8\xd2\x2b\x46\x63\xb5\x5a\xff\xf5\x32\x2b\x2f\xf6\x14\xb7\x5c\xab\x15\xb9\xe6\x72\x50\x66\x5d\x6a\x64\x27\x9c\x59\x2d\xf5\x5e\x24\xb3\xe6\x2a\xde\x8a\x7c\x48\x8a\xbe\x93\x46\x4e\x58\x5d\x74\x53\xe2\x79\xcc\xbe\x31\xc9\x21\x46\x8f\xcf\xad\xb3\x86\xcf\xba\x7d\x75\x2a\xc8\xd8\x71\xcd\x95\x81\xa3\x75\xa5\xbd\xe4\x65\x0e\xea\x99\x9f\x9d\xc0\x93\x4b\xc4\xb6\x7e\xa3\x32\xf1\x19\x82\x52\xc3\xc6\x27\xc5\xa3\xd0\x12\x4a\x52\x70\xc8\x55\x67\xde\xf3\xc0\xe6\xaa\xb0\xc0\x72\x02\x83\x7a\xb5\xb3\x4e\x7c\x40\xfb\x2d\x70\xb8\xae\x19\x9c\x64\x98\x4a\x35\x5f\x2c\x44\x1b\x90\xdc\x50\xa8\x86\xa9\x1d\x85\x7a\x7c\xd8\x88\x65\xd5\x42\xd4\xfc\x07\xf2\xfa\x11\x38\x98\x9c\xa3\xb6\xd4\x71\x6e\x94\x90\x9d\x6d\xa7\x04\xf4\x74\xa3\x57\xad\xb3\x6a\x5c\xf9\x9a\x22\x58\x96\x80\x92\xbe\x4b\x26\x64\x41\x63\xae\x0c\x59\xa5\x1b\xbe\x6b\x7a\xc8\x89\xe4\x2a\xcb\xc8\x8e\x96\x98\x32\x93\x56\x80\xcf\x89\x2c\xd9\x0d\x8a\xb3\x1a\x71\xe6\x50\x3a\x8e\x12\x79\x97\xa2\x4d\x51\x83\xd2\x20\x95\xf5\x52\x1d\x47\x78\x60\x06\xb8\x40\x26\xe1\x21\x45\x09\xf3\x82\xc4\x80\x12\xee\x28\x9e\xc5\x63\x15\xb8\x62\xd6\x03\x5c\xcf\x66\xc0\x4c\x75\x54\x1c\xcb\xac\x72\x42\x12\x05\x85\xc1\xd8\x41\x1c\xae\xb2\x9c\x04\xf6\x4b\x2c\x0f\xf5\x28\x79\x97\x25\xb1\x17\xd5\xcf\x3f\x17\xcc\xba\x1c\x1f\xc5\xff\x4b\x49\x0c\x64\x43\x98\x82\xbc\xd8\x0f\x61\x53\xd0\x85\x94\x24\x5d\xda\x36\x64\xf7\xd6\xcd\xba\xa1\x86\x47\x7c\x85\x96\x91\xb8\xb6\x98\x1d\xf3\xec\x28\x1e\x65\xd5\xf5\x55\x67\xcc\xe9\xf7\x5c\x09\x44\x9e\xce\xbf\x67\x90\xea\x97\x10\x26\xde\x91\x8f\xa9\x06\xd5\xf5\xe0\xbc\x57\x62\x3d\x47\x1f\x2d\xb5\x31\x92\xfb\x1c\xf7\x13\xb9\xfb\xb1\x5f\x89\x83\xd2\xe2\x9f\x84\x78\x91\x84\xd8\x13\x8b\x2e\xe6\x3d\x22\x20\xcf\xe9\xb2\x72\x19\xb1\x19\xab\x06\x2b\x9e\xbb\x10\x6f\x57\xd4\xa7\x40\xfd\x2f\x25\x8f\xdb\x1c\xf9\x8e\x0e\x5a\x89\x02\x93\x23\xa7\x45\xb9\x2f\x1b\xe7\xe9\x96\xc8\x3f\xc3\xe5\xcd\x0c\xda\x18\xf6\x1a\x04\x2d\x11\xca\x95\x65\x7f\x4d\x79\x67\xc6\x34\xcd\x0e\xf2\x74\xee\x5c\xd5\xed\xf3\x53\x31\x47\x2d\xd1\xa2\xf1\xbd\xe6\x41\xe9\x25\x3a\xdc\x10\xa3\x99\xc0\xa5\x92\x56\x2b\x01\xb9\x60\x72\x43\x65\x80\x69\x84\xd8\x4d\x98\x19\x49\x8c\x61\xbe\xf6\xa6\xdc\x30\x9e\x92\x44\xe7\xd6\x49\xbf\x01\xa5\x71\x47\xa4\x4b\x58\x2f\xef\x49\x96\x72\x07\x3d\x94\x2a\x0d\x85\x52\x4a\xd2\x19\x5b\x31\x12\x6c\x4e\x82\xec\x7a\xc8\xaf\x73\xa5\x1c\x70\x6a\x0b\x0d\x88\x76\xf0\x78\x17\x1e\xab\x5e\x1a\x19\xf8\xf8\x7e\xa0\x08\xa2\x75\xc1\x98\x2d\xd8\x5c\x13\x1f\x8d\x3d\x03\x79\xf9\xa4\x3b\xe8\xe5\x09\x8f\x26\x0c\x17\x47\xbc\x16\xb5\x7c\x69\xdf\x3a\x1c\xd5\xa8\x61\x21\x4d\xc9\x80\xc6\xfb\x02\xcd\x00\xf2\xdd\x5e\xae\x1f\x94\xa3\x8d\xae\x3a\x1c\x8a\xd8\x63\xb1\x6e\xbb\xa8\xf4\x81\x6a\x9e\x1c\x17\xa0\xc6\xcb\x38\xb2\xa6\x99\x9c\xf1\x71\x85\xed\xab\x73\xef\x86\x05\x30\xa1\x64\x02\x0f\x64\xd3\xa6\xbd\x7e\x1f\xe9\x3f\x1c\xac\xca\x50\x48\xba\x2f\x50\xac\x81\x62\x94\x96\x16\x6b\x60\x90\xf9\x6f\x98\xe2\xdd\x7d\x3d\x0b\x35\x68\xb0\xb3\x33\xad\x59\x7b\xa5\x15\x91\xc5\xac\x7b\x7f\x7f\xc4\x1b\xc5\xae\xb5\xf8\xec\x77\x4c\x79\xdb\xc0\x43\x4a\x3c\xf5\xe8\x5e\x93\xc1\x51\xa8\xa1\x67\xf5\x7c\x7a\xcf\xf2\x52\x15\x22\x6e\x65\xf9\x1c\xab\x0d\xf4\x10\x9c\x38\x02\x93\xdd\xb6\x70\xd8\xf6\x03\xdf\xed\xd1\xbe\x4d\xd7\xe9\xb9\xf4\xda\xb4\x3b\x47\xd8\x60\x98\xb5\x71\x45\x67\x57\x2d\xfb\xe3\x1f\x43\xb3\xfc\xed\x7f\x0c\x9f\xbb\x3b\xd2\xa7\xc7\xe9\x6f\x17\xa3\x71\xfd\xa6\x31\xdc\x3f\xc1\x87\xad\x75\xe1\x09\xfa\x70\x01\x9b\x4d\xac\xef\xe1\xff\xff\x34\x80\x26\x83\x1d\x33\xea\x5d\x2d\xee\x70\xe3\xfe\x0d\xe5\xee\xa5\x53\xf3\xe6\xb6\xdf\x37\x7e\xdb\xd0\x79\x29\x1e\x5d\x57\xb8\xbb\x55\xa1\xb6\x62\xd4\xec\x2d\x47\xc4\xa7\x30\xa1\xd9\x8e\x86\x8e\x15\x03\xaf\xfa\xed\xed\x47\x60\x9c\xa3\x31\x03\x50\x5a\x99\x23\x16\xa4\xca\xd8\x83\xa4\xe4\x4a\x8f\x97\xe2\x88\x0f\x93\xa2\x69\xc5\x2c\xce\x96\x38\x38\x10\xec\x17\x16\x78\xc0\x12\xd7\x95\xcc\x09\x78\x80\x94\x15\xc6\x86\xf9\xba\x5c\xeb\x14\x1a\xe3\x6a\x23\x5d\xce\x2e\x24\x8d\x65\x92\xe3\x80\x82\x3c\xcc\x65\x33\x3f\x97\x3d\xd3\x4c\xb8\x7c\x67\x36\xf9\xe9\x14\xf6\x85\xac\x39\xff\x95\xea\xbb\xd1\xd0\xcd\x2d\x4a\x8a\x00\xcb\x42\xa9\x5b\x04\x70\xbd\x20\x14\xb1\x3b\xf6\x0f\x0e\xe3\xc9\x41\x2f\xc4\xc9\xe7\x1a\x3d\xd8\x63\xc2\xb8\x69\xa5\x30\x21\xef\xb8\x03\x89\x24\x93\x6d\xac\x54\x7f\x0f\x81\xbf\x78\x5c\x29\xae\x14\x6f\x7c\x11\xd1\xb1\xf1\x46\x69\x2c\xbf\xff\x39\xf8\x8f\x96\x9e\xf6\x77\x46\x4e\x9f\xb3\xc7\xb3\xdf\x03\x00\x00\xff\xff\x18\x17\x45\xa7\x43\x25\x00\x00")

func apiSwaggerJsonBytes() ([]byte, error) {
	return bindataRead(
		_apiSwaggerJson,
		"api.swagger.json",
	)
}

func apiSwaggerJson() (*asset, error) {
	bytes, err := apiSwaggerJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "api.swagger.json", size: 9539, mode: os.FileMode(420), modTime: time.Unix(1536876582, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"api.swagger.json": apiSwaggerJson,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"api.swagger.json": &bintree{apiSwaggerJson, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}


func assetFS() *assetfs.AssetFS {
	assetInfo := func(path string) (os.FileInfo, error) {
		return os.Stat(path)
	}
	for k := range _bintree.Children {
		return &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: assetInfo, Prefix: k}
	}
	panic("unreachable")
}