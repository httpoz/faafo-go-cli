# FAAFO Go cli
In the spirit of FAAFO (fuck around and find out), I created this small CLI app as part of a hacking hour with [tamirarnesty](https://github.com/tamirarnesty). The POC we were trying to build here is validating an OpenAPI json document.

## Planned improvements
- [ ] Function calling to actually apply the changes.
- [ ] Split between fixing syntax and improving spec documentation.
- [ ] Implement other tools like redocly to chunk larger schema files and run the clean up as a pipeline.
- [ ] Learn benchmarking and improve the app->response performance.
- [ ] Make it run on my M1 Pro in a reasonable amount of time. This would make it a viable tool for most developers. My M4 Pro is a cheat code.

## How to run the app
1. Set up the model in docker locally `make ollama-up`
3. Run the app `make app`

# Tests (if you want to call it that)
The output file can be viewed in `files/output-fix.json`. You can also install `git-delta` and run the command:
```
delta files/sample-api.json files/output-fix.json
```

### M1 MacBook Pro with 32GB RAM
Well the first run went as you would expect. I may have failed to convince Tamir to join the LocalLLaMA movement.
![first-run](.yuck/bpytop-m1.png)

And after peaking at 2000% CPU with all cores fired up, it just died after 20 minutes.
![ollama out](.yuck/docker-m1.png)

### M4 MacBook Pro with 24GB RAM (runtime: ~1.5 minutes)
Under the same conditions as test 1, it ran in under a minute with all cylinders firing with the initial prompt. After some "prompt engineering", the runtime was 1.5 minutes.
![first-run](.yuck/bpytop-m4.png)
![easy](.yuck/docker-m4.png)

## Bill of Materials
1. Official OpenAI go client
2. Ollama official docker image
3. llama3.2 3b model
