# Frontend Files and Uploads Policy

For modules with file/document uploads:

- Validate file type before upload when the contract defines allowed types.
- Validate max size when defined.
- Show progress if upload may take noticeable time.
- Show per-file errors when uploading multiple files.
- Do not block the whole screen when only one file fails.
- Keep UI clear for download, preview and delete actions according to permissions.
- Use existing upload helpers/services when present.
- Do not duplicate file hashing or upload logic if `core/libs/documents`, `core/libs/files` or equivalent exists.
