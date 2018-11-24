# README

This is the read me (must read)


# Apple .app packages

To create an application correctly for apple devices, adhere to the followign folder structure

```
Filemaster.app/
└── Contents
    ├── Info.plist
    ├── MacOS
    │   └── filemaster
    └── Resources
        └── icon.icns
```

This [article](https://medium.com/@mattholt/packaging-a-go-application-for-macos-f7084b00f6b5) was very helpful.

# Info.plist

The info.plist must include the `LSItemContentTypes` node and an associated array in order to allow drag and drop functionality.

```
<key>LSItemContentTypes</key>
<array>
    <string>public.jpg</string>
</array>
```

Here is the [documentation](https://developer.apple.com/library/archive/documentation/General/Reference/InfoPlistKeyReference/Articles/CoreFoundationKeys.html#//apple_ref/doc/uid/TP40009249-101685-TPXREF107) for info.plist key and vlaue pairs.