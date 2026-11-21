---
# Live test: does gh-aw's max-ai-credits stop a Cursor agent? Compiled by
# ghaw-compile into tests/live/compiled; never installed into a user's repo.
# The cap is 2 credits ($0.02): any real Cursor session passes it on its
# first request, so an enforced cap must cut the run short.
on:
  workflow_dispatch:
    inputs:
      model:
        description: A Cursor model id from `agent --list-models`
        required: true
        type: string
engine:
  id: cursor
  model: ${{ inputs.model }}
imports:
  - shared/cursor.md
permissions:
  contents: read
network:
  allowed:
    - defaults
timeout-minutes: 10
max-ai-credits: 2
tools:
  bash: [":*"]
safe-outputs:
  noop:
  threat-detection: false
---
# Cap test

Read every file in this repository one at a time with separate tool calls.
After each file, write a detailed 300-word explanation of it. Then read them
all again and compare each pair of files in detail. Keep going until you have
made at least 40 tool calls. Finally call noop with a one-line summary.
