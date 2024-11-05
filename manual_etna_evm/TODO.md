### Nodes spit error

```bash
./data/node0/logs/main.log:[11-05|07:40:12.933] ERROR chains/manager.go:404 error creating chain {"subnetID": "J1XKcNUAAfkSGvF9uRZaM5ngSUCQEWNVDv6XzNmZnSsSC2rzw", "chainID": "2LGZYfs4gUzUyhS5P929nCBCbauBiPWwwcqUYcjvGMFXJuqRa2", "chainAlias": "2LGZYfs4gUzUyhS5P929nCBCbauBiPWwwcqUYcjvGMFXJuqRa2", "vmID": "srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy", "error": "error while getting vmFactory: \"srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy\" was not found"}
```

This is from "Creating chain" logs:
```
2024/11/05 07:51:26 Using vmID: srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy
```

Why do we have this VMID?!

Avalanche CLI ends up with one file named `qDNkrRvkeM1GrdrEhDbGZGz6Hiq5nQ4utodRih4x111KBqiC8` in the plugins folder, which is different. Is the file name equal to vmID?

qDNkrRvkeM1GrdrEhDbGZGz6Hiq5nQ4utodRih4x111KBqiC8 is not googleable, srEXiWaHuhNyGwPUi444Tu47ZEDwxTWrbQiuD7FmgSAQ6X7Dy is in all the examples.

Is vmID a hash or any value?

For sure, need to populate plugins dir


### `cmd/11_check_subnet/check.go` is unfinished for node, needs more work
