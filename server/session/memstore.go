package session

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type SessionProvider struct {
	lock     sync.Mutex
	sessions map[string]*list.Element //to save in memory
	list     *list.List               //to perform grabage collection. Reason why we separate this?
}

var provider = &SessionProvider{list: list.New()}

//MemStore = session store
type MemStore struct {
	sid          string
	items        Event
	timeAccessed time.Time
}

// Set value type must follow the struct field types, except for ResizeFrom and ResizeTo, which is a map[string]string.
func (ms *MemStore) Set(key, value interface{}) error {
	switch key {
	case "WebsiteUrl":
		ms.items.WebsiteUrl = value.(string)

	case "SessionId":
		ms.items.SessionId = value.(string)

	case "ResizeFrom":
		ms.items.ResizeFrom.Width = value.(map[string]string)["width"]
		ms.items.ResizeFrom.Height = value.(map[string]string)["height"]

	case "ResizeTo":
		ms.items.ResizeTo.Width = value.(map[string]string)["width"]
		ms.items.ResizeTo.Height = value.(map[string]string)["height"]

	case "CopyAndPaste":
		if ms.items.CopyAndPaste == nil {
			ms.items.CopyAndPaste = make(map[string]bool)
		}

		for fieldId, pasted := range value.(map[string]bool) {
			fieldId := fieldId
			pasted := pasted
			ms.items.CopyAndPaste[fieldId] = pasted
		}

	case "FormCompletionTime":
		ms.items.FormCompletionTime = int(value.(float64))

	default:
		return errors.New("key " + key.(string) + " is undefined")

	}

	provider.SessionUpdate(ms.sid)
	return nil
}

func (ms *MemStore) Get(key interface{}) interface{} {
	provider.SessionUpdate(ms.sid)
	switch key {
	case "items":
		return ms.items
	default:
		return nil
	}
}

func (ms *MemStore) Delete(key interface{}) error {
	// delete(ms.items, key), but commented out currently as it is not used in this challenge
	provider.SessionUpdate(ms.sid)

	return nil
}

func (ms *MemStore) SessionID() string {
	return ms.sid
}

func (provider *SessionProvider) SessionInit(sid string) (Session, error) {
	provider.lock.Lock()
	defer provider.lock.Unlock()
	v := Event{}
	newSession := &MemStore{sid: sid, timeAccessed: time.Now(), items: v}
	element := provider.list.PushBack(newSession) //for GC. "element" is the newly added element in the list
	provider.sessions[sid] = element              //to save in memory

	return newSession, nil
}

func (provider *SessionProvider) SessionRead(sid string) (Session, error) {
	if element, ok := provider.sessions[sid]; ok {
		return element.Value.(*MemStore), nil
	} else {
		session, err := provider.SessionInit(sid)
		return session, err
	}
}

func (provider *SessionProvider) SessionDestroy(sid string) error {
	if element, ok := provider.sessions[sid]; ok {
		delete(provider.sessions, sid)
		provider.list.Remove(element)
		return nil

	}
	return nil
}

func (provider *SessionProvider) SessionGC(maxLifetime int64) {
	provider.lock.Lock()
	defer provider.lock.Unlock()

	for {
		element := provider.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*MemStore).timeAccessed.Unix() + maxLifetime) < time.Now().Unix() {

			provider.list.Remove(element)
			delete(provider.sessions, element.Value.(*MemStore).sid)
		} else {
			break
		}
	}
}

// update timeAccessed, as well as repositioning session ot the front of the list.
func (provider *SessionProvider) SessionUpdate(sid string) error {
	provider.lock.Lock()
	defer provider.lock.Unlock()

	if element, ok := provider.sessions[sid]; ok {

		element.Value.(*MemStore).timeAccessed = time.Now()
		provider.list.MoveToFront(element)
		return nil
	}
	return nil

}

func init() {
	provider.sessions = make(map[string]*list.Element)
	Register("memory", provider)
}
