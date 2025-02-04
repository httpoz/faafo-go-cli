# FAAFO Go cli
In the spirit of FAAFO (fuck around and find out), I created this small CLI app as part of a hacking hour with [tamirarnesty](https://github.com/tamirarnesty). The POC we were trying to build here is validating an OpenAPI json document.

## Planned improvements
- [ ] Function calling to actually apply the changes.
- [ ] Split between fixing syntax and improving spec documentation.
- [ ] Implement other tools like redocly to chunk larger schema files and run the clean up as a pipeline.
- [ ] Learn benchmarking and improve the app->response performance.
- [x] Make it run on my M1 Pro in a reasonable amount of time. This would make it a viable tool for most developers. My M4 Pro is a cheat code.

## How to run the app
1. Install Ollama with `brew install ollama && ollama pull `
2. Set up the model in docker locally `make ollama-up`
3. Run the app `make app`

# Tests (if you want to call it that)
The output file can be viewed in `files/output-fix.json`. You can also install `git-delta` and run the command:
```
delta files/sample-api.json files/output-fix.json
```

## Bill of Materials
1. Official OpenAI go client
2. Ollama official docker image
3. llama3.2 3b model


# Learnings
- Docker does not use the GPU on MacBooks. It is better to run inference on the native Ollama client. The initial implementation used the ollama docker image and ran the request unsuccessfully failing at 40minutes, after switching to the native client, the request was completed in 40 seconds consistently.
