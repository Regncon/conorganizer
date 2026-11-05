---
tags: [conorganizer-commits]
max_turns: 20
timeout_seconds: 600
allowed_tools: [Read, Glob, Grep, Skill, Write, Edit, Bash]
---

commit these changes in fixture-repo. the retry is there because checkin answers 502 for a few seconds after each of their deploys, and the login sync then dropped every ticket for whoever logged in right then. this fixes #712
