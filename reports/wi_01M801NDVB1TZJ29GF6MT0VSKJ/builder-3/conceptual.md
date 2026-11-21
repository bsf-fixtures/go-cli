A timesheet line with `-` as its project is now treated as malformed, like a bad day or duration. It fails with an error that names the line instead of creating a bogus `-` row. Nothing else changes.
