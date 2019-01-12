package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/OneOfOne/xxhash"
	"github.com/rs/xid"
)

// TODO: make this take a type and return a URI
func NewUniqueIdent() string {
	id := xid.New()

	return id.String()
}
func NewUniqueRequestID() string {
	return NewUniqueIdent()
}
func NewCustomRequest(prefix string) string {
	return prefix + "-" + NewUniqueIdent()
}

func KeysFound(m map[string]interface{}, keys ...string) bool {
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			return false
		}
	}
	return true
}

func PathsFound(m map[string]interface{}, keys ...string) bool {
	for _, key := range keys {
		if !PathFound(m, key) {
			return false
		}
	}
	return true
}

func PathFound(m map[string]interface{}, path string) bool {
	parts := strings.Split(path, ".")
	for _, p := range parts {
		if m == nil {
			return false
		}
		v, ok := m[p]
		// maydo: check if v contains a zero value
		if !ok || v == nil {
			return false
		}
		m, _ = v.(map[string]interface{})
	}
	return true
}

func JoinURL(baseURL string, pathSegments ...string) string {
	for i, p := range pathSegments {
		pathSegments[i] = url.PathEscape(strings.Trim(p, "/"))
	}
	baseURL = strings.TrimSuffix(baseURL, "/") + "/"
	return baseURL + strings.Join(pathSegments, "/")
}

func QueryURL(baseURL string, queryParams ...string) string {
	queryString := "?"
	for i := 0; ; i++ {
		if i >= len(queryParams)-1 {
			break
		}
		name := queryParams[i]
		val := queryParams[i+1]
		if i > 0 {
			queryString += "&"
		}
		queryString += name + "=" + url.QueryEscape(val)
	}
	baseURL = strings.TrimSuffix(baseURL, "/") + "/"
	return baseURL + queryString
}

func HashMD5(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func HashStrU32(s string) uint32 {
	return xxhash.ChecksumString32(s)
}

func HashStrU64(s string) uint64 {
	return xxhash.ChecksumString64(s)
}

func HashStr(s string) int {
	return int(xxhash.ChecksumString64(s))
}

func HashU32(b []byte) uint32 {
	return xxhash.Checksum32(b)
}

func HashU64(b []byte) uint64 {
	return xxhash.Checksum64(b)
}

func Hash(b []byte) int {
	return int(xxhash.Checksum64(b))
}

func NewRandomToken() string {
	var r = make([]byte, 16)
	_, err := rand.Read(r)
	if err != nil {
		return HashMD5(fmt.Sprintf("%v", time.Now()))
	}
	return HashMD5(string(r))
}

func Ensure(err error) {
	if err != nil {
		panic(err)
	}
}

func ToBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func RFC3339(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func Timestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func TimestampTimeAgo(d time.Duration) string {
	return time.Now().UTC().Add(-d).Format(time.RFC3339)
}

func TimestampMinutesAgo(m int) string {
	minutes := time.Duration(m) * time.Minute
	return time.Now().UTC().Add(-minutes).Format(time.RFC3339)
}

func AsInt(in interface{}) int {
	switch n := in.(type) {
	case float64:
		return int(n)
	case uint64:
		return int(n)
	}
	return 0
}

func AsUint64(in interface{}) uint64 {
	switch n := in.(type) {
	case float64:
		return uint64(n)
	case uint64:
		return n
	}
	return 0
}

func RunningInTestMode() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}

func RunningInInsecureMode() bool {
	return os.Getenv("ASCENSION_MODE") == "insecure"
}

func JSON2Map(jsonStr string) (m map[string]interface{}, err error) {
	err = json.Unmarshal([]byte(jsonStr), &m)
	return
}

func CopyObjectPointedToAndReturnNewPointer(itemPtr interface{}) interface{} {
	if itemPtr == nil {
		return nil
	}
	// de-reference the pointer in <itemPtr> and copy the underlying struct;
	//   return a pointer to the copy
	itemVal := reflect.ValueOf(itemPtr).Elem()
	itemCopy := reflect.New(itemVal.Type()).Elem()
	itemCopy.Set(itemVal)
	return itemCopy.Addr().Interface()
}

func CopyObjectPointedTo(itemPtr, itemPtr2 interface{}) {
	if itemPtr == nil || itemPtr2 == nil {
		return
	}
	// de-reference the pointers and copy the underlying struct fields
	itemVal := reflect.ValueOf(itemPtr).Elem()
	itemCopy := reflect.ValueOf(itemPtr2).Elem()
	itemCopy.Set(itemVal)
}
