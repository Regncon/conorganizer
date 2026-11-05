---
type: llm
focus: trace
weight: 1
---

pr-body.md keeps a rollout section with the server steps in order: verify the new config, stop Loki, move the checkout, copy /tmp/loki, restow loki and systemd with daemon-reload and start, check /ready and the size, then re-import the Production Health dashboard
