# Implementation brief — AKR-GITHUB-REPO-TOOLS

## GitHub adapter

Métodos nuevos en `internal/adapter/external/github.Client` (mismo auth):
- SearchCode(ctx, account, owner, name, query)
- ReadFile(ctx, account, owner, name, path, ref) — rechaza `..` / paths absolutos
- ListRecentCommits(ctx, account, owner, name, branch, limit)
- ReadCommit(ctx, account, owner, name, sha)
- ReadDiff(ctx, account, owner, name, sha) — patch vía commit files o compare

Todo vía API remota; credenciales solo en punto de uso (wipe after).

## Repo resolution

Misma cadena que EvidenceAssembler: Investigation.IncidentID → Incident.ProjectID
→ Project.GitHubRepository (owner/name/default_branch + GitHubAccountID) →
GitHubAccountReader.Get → tools.

## QVAC tools

Wrappers que exponen las 5 tools al registry del runner; resultados como texto
JSON/saneado enviado como DATA al modelo.
