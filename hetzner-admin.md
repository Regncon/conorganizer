# Hetzner admin documentation

## Add a user
To add a use to he Hetzner could vm add the user and assign it the correct groups

Add user:
```bash 
adduser <username>

```

Assign groups:
```bash
usermod -aG adm sudo www-data
```

Change the password:
```bash
passwd <username>
```

Set up SSH keys
**ToDo: add more details here**


