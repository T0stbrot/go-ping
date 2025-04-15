# go-ping
Go Module for ICMP Pings with TTL Support on Windows without Admin Privileges

# It doesn't work for me, what do i do?
Sometimes you need to allow ICMP(v6) trough Windows Firewall rules:

`netsh advfirewall firewall add rule name=AllowICMP protocol=ICMPv4 dir=in action=allow`

`netsh advfirewall firewall add rule name=AllowICMPv6 protocol=ICMPv6 dir=in action=allow`

Im looking to "fixing" this so this maybe works without those rules.