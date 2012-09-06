package main

import (
	"syscall"
	"unsafe"
)

const (
    NO_ERROR                         = 0
    ERROR_SUCCESS                    = 0
    ERROR_FILE_NOT_FOUND             = 2
    ERROR_PATH_NOT_FOUND             = 3
    ERROR_ACCESS_DENIED              = 5
    ERROR_INVALID_HANDLE             = 6
    ERROR_BAD_FORMAT                 = 11
    ERROR_INVALID_NAME               = 123
    ERROR_MORE_DATA                  = 234
    ERROR_NO_MORE_ITEMS              = 259
    ERROR_INVALID_SERVICE_CONTROL    = 1052
    ERROR_SERVICE_REQUEST_TIMEOUT    = 1053
    ERROR_SERVICE_NO_THREAD          = 1054
    ERROR_SERVICE_DATABASE_LOCKED    = 1055
    ERROR_SERVICE_ALREADY_RUNNING    = 1056
    ERROR_SERVICE_DISABLED           = 1058
    ERROR_SERVICE_DOES_NOT_EXIST     = 1060
    ERROR_SERVICE_CANNOT_ACCEPT_CTRL = 1061
    ERROR_SERVICE_NOT_ACTIVE         = 1062
    ERROR_DATABASE_DOES_NOT_EXIST    = 1065
    ERROR_SERVICE_DEPENDENCY_FAIL    = 1068
    ERROR_SERVICE_LOGON_FAILED       = 1069
    ERROR_SERVICE_MARKED_FOR_DELETE  = 1072
    ERROR_SERVICE_DEPENDENCY_DELETED = 1075
)

// Registry Key Security and Access Rights
const (
    KEY_ALL_ACCESS         = 0xF003F
    KEY_CREATE_SUB_KEY     = 0x0004
    KEY_ENUMERATE_SUB_KEYS = 0x0008
    KEY_NOTIFY             = 0x0010
    KEY_QUERY_VALUE        = 0x0001
    KEY_SET_VALUE          = 0x0002
    KEY_READ               = 0x20019
    KEY_WRITE              = 0x20006
)

type (
    HANDLE       uintptr
    HKEY         HANDLE
)

// Registry value types
const (
    RRF_RT_REG_NONE         = 0x00000001
    RRF_RT_REG_SZ           = 0x00000002
    RRF_RT_REG_EXPAND_SZ    = 0x00000004
    RRF_RT_REG_BINARY       = 0x00000008
    RRF_RT_REG_DWORD        = 0x00000010
    RRF_RT_REG_MULTI_SZ     = 0x00000020
    RRF_RT_REG_QWORD        = 0x00000040
    RRF_RT_DWORD            = (RRF_RT_REG_BINARY | RRF_RT_REG_DWORD)
    RRF_RT_QWORD            = (RRF_RT_REG_BINARY | RRF_RT_REG_QWORD)
    RRF_RT_ANY              = 0x0000ffff
    RRF_NOEXPAND            = 0x10000000
    RRF_ZEROONFAILURE       = 0x20000000
    REG_PROCESS_APPKEY      = 0x00000001
    REG_MUI_STRING_TRUNCATE = 0x00000001
)

// Registry Value Types
const (
    REG_NONE                       = 0
    REG_SZ                         = 1
    REG_EXPAND_SZ                  = 2
    REG_BINARY                     = 3
    REG_DWORD                      = 4
    REG_DWORD_LITTLE_ENDIAN        = 4
    REG_DWORD_BIG_ENDIAN           = 5
    REG_LINK                       = 6
    REG_MULTI_SZ                   = 7
    REG_RESOURCE_LIST              = 8
    REG_FULL_RESOURCE_DESCRIPTOR   = 9
    REG_RESOURCE_REQUIREMENTS_LIST = 10
    REG_QWORD                      = 11
    REG_QWORD_LITTLE_ENDIAN        = 11
)

// Registry predefined keys
const (
    HKEY_CLASSES_ROOT     HKEY = 0x80000000
    HKEY_CURRENT_USER     HKEY = 0x80000001
    HKEY_LOCAL_MACHINE    HKEY = 0x80000002
    HKEY_USERS            HKEY = 0x80000003
    HKEY_PERFORMANCE_DATA HKEY = 0x80000004
    HKEY_CURRENT_CONFIG   HKEY = 0x80000005
    HKEY_DYN_DATA         HKEY = 0x80000006
)

var (
    modadvapi32 = syscall.NewLazyDLL("advapi32.dll")

	RegOpenKeyEx = modadvapi32.NewProc("RegOpenKeyExW")
	RegCloseKey = modadvapi32.NewProc("RegCloseKey")
	RegSetValueEx = modadvapi32.NewProc("RegSetValueExW")
	RegQueryValueEx = modadvapi32.NewProc("RegQueryValueExW")
	RegDeleteValue = modadvapi32.NewProc("RegDeleteValueW")
)

func with(hkey HKEY, path string, handler func(HKEY)) {
	var subkey HKEY
	err, _, _ := RegOpenKeyEx.Call(
		uintptr(hkey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		0,
		uintptr(KEY_ALL_ACCESS),
		uintptr(unsafe.Pointer(&subkey)),
	)
	if err == ERROR_SUCCESS {
		defer RegCloseKey.Call(uintptr(subkey))
		handler(subkey)
	}
}

func get(hkey HKEY, path, key string) string {
	str := ""
	
	with(hkey, path, func(subkey HKEY) {
		var bufLen uint32
		
		RegQueryValueEx.Call(
			uintptr(subkey),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(key))),
			0,
			0,
			0,
			uintptr(unsafe.Pointer(&bufLen)),
		)
		
		if bufLen == 0 {
			return
		}
		
		buf := make([]uint16, bufLen)
		
		err, _, _ := RegQueryValueEx.Call(
			uintptr(subkey),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(key))),
			0,
			0,
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&bufLen)),
		)
		
		if err != ERROR_SUCCESS {
			return
		}
		
		str = syscall.UTF16ToString(buf)
	})
	
	return str
}

func set(hkey HKEY, path, key, value string) {
	v := syscall.StringToUTF16(value)
	with(hkey, path, func(subkey HKEY) {	
		RegSetValueEx.Call(
			uintptr(subkey),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(key))),
			0,
			uintptr(REG_EXPAND_SZ),
			uintptr(unsafe.Pointer(&v[0])),
			uintptr(len(v) * 2),			
		)
	})
}

func delete(hkey HKEY, path, key string) {
	with(hkey, path, func(subkey HKEY) {	
		err, _, _ := RegDeleteValue.Call(
			uintptr(subkey),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(key))),
		)
		
		if err != ERROR_SUCCESS {
			panic(err)
		}
	})
}