package secret

/*
#cgo CFLAGS: -I/opt/meituan/kms/include/
#cgo LDFLAGS: -L/usr/lib/gcc/x86_64-redhat-linux/4.8.2/ -L/opt/meituan/kms/lib -lkms -lstdc++ -lkms_comm -lcryptopp -lthrift -llog4cplus -lm -lc -lstdc++
char* getKeyFromKms(char* test_appkey, char* test_name);
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

func GetValueByeKey(ak, key string) string {
	appkey := C.CString(ak)
	defer C.free(unsafe.Pointer(appkey))
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	v := C.getKeyFromKms(appkey, k)
	return C.GoString(v)
}