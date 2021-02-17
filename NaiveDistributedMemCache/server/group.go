package server

type Group struct {
	name       string
	client     Client
	localCache Cache
	Getter     func(key string) []byte
}

func NewGroup(name string, virtualReplicas int, getter func(key string) []byte) *Group {
	return &Group{
		name:       name,
		client:     *NewClient(virtualReplicas, nil),
		localCache: *NewCache(),
		Getter:     getter,
	}
}

func (group *Group) Put(key string, value []byte) bool {
	return group.localCache.Put(key, value)
}

func (group *Group) Get(key string) []byte {
	// First, try to get from local cache
	if value := group.localCache.Get(key); value != nil {
		return value
	}
	// Second, try to get from remote cache
	if value, err := group.client.getFromRemotePeer(group.name, key); err == nil {
		return value
	}
	// Third, try to get by user-defined Getter
	if group.Getter != nil {
		if value := group.Getter(key); value != nil {
			if ok := group.localCache.Put(key, value); ok {
				return value
			}
		}
	}
	return nil
}

func (group *Group) RegisterPeers(self string, peers map[string]string) {
	for name, address := range peers {
		group.client.consistentHash.AddNodes(name)
		if name != self {
			group.client.peer2AddressMap[name] = address
		}
	}
}
