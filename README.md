# L3MultiplexVPN
Multiplex up to 255 separate layer3 networks through one tunnel

The program requires the CAP_NET_RAW capability or root.

Usage:
```
./vpn <PeerIP> <Local network id list>
```
Example:
```
./vpn 192.168.5.5 0,2
```
