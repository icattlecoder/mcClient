package mcClient

import (
	"sync"
	"bytes"
	"encoding/gob"
	"github.com/bradfitz/gomemcache/memcache"
)


//======================================================================
//gob encoding
//gob的特点是解码快，适合Get>Set的情况
type gobMCClient struct {
	Client *memcache.Client
	keyPrefix string
	sync.Mutex
}

func (g *gobMCClient)realKey(key string) string{ 
	return g.keyPrefix + key 
}

// 如果共享mc，需要保证keyPrefix不同
func NewGobMCClient(keyPrefix string, host ...string) *gobMCClient {
	client := memcache.New(host...)
	return &gobMCClient{Client: client, keyPrefix: keyPrefix}
}

func (g *gobMCClient) Get(key string, val interface{}) (err error) {
	g.Lock()
	item, err := g.Client.Get(g.realKey(key))
	g.Unlock()
	if err == nil {
		reader := bytes.NewReader(item.Value)
		decoder := gob.NewDecoder(reader)
		err = decoder.Decode(val)
	}
	return
}

func (g *gobMCClient) GetString(key string)(val string,err error) {

	g.Lock()
	item,err := g.Client.Get(g.realKey(key))
	g.Unlock()
	if err ==nil{
		val = string(item.Value)
	}
	return
}

func (g *gobMCClient) SetString(key string,val string)(err error) {

	g.Lock()
	item := memcache.Item{Key:g.realKey(key),Value:[]byte(val)}
	g.Unlock()
	err = g.Client.Set(&item)
	return
}


func (g *gobMCClient) Set(key string, val interface{}) (err error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(val)
	if err != nil {
		return
	}
	item := memcache.Item{Key: g.realKey(key), Value: buffer.Bytes()}
	g.Lock()
	err = g.Client.Set(&item)
	g.Unlock()
	return
}

func (g *gobMCClient) Delete(key string) (err error) {
	g.Lock()
	defer g.Unlock()
	return g.Client.Delete(g.realKey(key))
}