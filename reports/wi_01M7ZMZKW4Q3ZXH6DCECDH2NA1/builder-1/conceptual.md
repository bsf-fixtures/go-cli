# Conceptual

Without `--project`, tally already lists each project's time. Billers and reporters still need one number for the whole timesheet. The grand-total line is that number: same column layout, labelled `total`, summing every entry.

With `--project`, the tool answers a narrower question ("how much on site?"), so a grand total would be noise and is omitted.

Summing stays in `internal/tally` (`GrandTotal`), not in `main`, matching how per-project totals already work.
