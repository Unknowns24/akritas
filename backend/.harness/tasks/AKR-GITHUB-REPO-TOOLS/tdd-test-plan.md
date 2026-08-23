# TDD test plan — AKR-GITHUB-REPO-TOOLS

*(Aprobado vía instrucción de ejecución del plan H3 restante.)*

- search_code / read_file / commits / diff contra httptest GitHub API.
- path traversal (`../`, absoluto) → error, sin request peligroso.
- scope: requests solo a owner/name del project.
- tool wrappers invocables desde el registry del runner.
- credencial no aparece en logs ni en respuestas de tool hacia QVAC (solo contenido de código).
