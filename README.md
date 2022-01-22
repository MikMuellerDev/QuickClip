# QuickClip
 A self hosted, simple alternative to PasteBin.
### Issues
On Debian `./QuickClip-0.4.0-Beta: /lib/x86_64-linux-gnu/libc.so.6: version GLIBC_2.32' not found (required by ./QuickClip-0.4.0-Beta)` can happen when running the executable.
#### Fix
```
sudo apt install golang
```
```
cd QuickClip
```
```
make
```
Your custom-built executable can now be found in `.//build/QuicklipX.X.X`.
