Amphiptere 🐍

"Born to defend the sacred gold."

Amphiptere is a lightweight perimeter defense tool I built in Go for Amphiptere OS. 

The idea is simple: when you enable it, everything gets locked down. 
Want in? You have to hit port 9999 with the temporary code first. 
No code = no access. 

I got tired of leaving SSH and other ports wide open, so I made a sentinel.

Most firewall tools are a pain to configure. I wanted something I could run with one command that:
Locks the box immediately, Lets me auth in with a temp code, Auto-cleans everything when I’m done. . . 

No configs. No daemons to hunt down. Just `--enable` and `--disable`.

Features:
Your code lives only in RAM. Restart = it's gone. Hashed with SHA-256.
Drops all incoming traffic except port 9999 and established connections.
Enter the right code and your IP gets whitelisted for 5 minutes.
 `--disable` flushes everything. You won't brick your own server.

Built specifically for Amphiptere OS, but works on any Debian/Kali/Ubuntu box.

License
Copyright © 2026 Arjun Raj. All rights reserved. 
This project is proprietary software created for Amphiptere OS. See the [LICENSE](LICENSE) file for details.
