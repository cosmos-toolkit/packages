# cosmos-toolkit/pkgs

Pacotes Go reutilizáveis para CLI, workers, crons e APIs. Pensado para monólito modular, microserviço ou worker-only.

**Princípio:** não criar tudo de uma vez. Mínimos bem feitos, evoluir conforme os projetos pedem. Pacote pequeno, interface clara, poucas dependências, fácil de remover.

A adição de pacotes a projetos existentes é feita **somente pelo [cosmos-cli](https://github.com/cosmos-toolkit/cosmos-cli)**; este repositório contém apenas o código dos pacotes e o `manifest.yaml`.

---

## Instalação via Cosmos CLI

Na raiz do seu projeto (onde está o `go.mod`):

```bash
cosmos pkg logger    # copia pkg/logger + copy_deps (ex.: contextx) e reescreve imports
cosmos pkg config    # copia pkg/config e instala dependências (ex.: godotenv)
cosmos pkg validator # copia pkg/validator e instala validator/v10
```

O CLI baixa o repositório [cosmos-toolkit/packages](https://github.com/cosmos-toolkit/packages), copia o pacote para `pkg/<name>`, reescreve `github.com/cosmos-toolkit/pkgs` pelo módulo do seu projeto e roda `go get` / `go mod tidy`.

Listar pacotes disponíveis:

```bash
cosmos list pkgs
```

---

## Uso como módulo (monorepo / replace)

```go
import (
    "github.com/cosmos-toolkit/pkgs/pkg/logger"
    "github.com/cosmos-toolkit/pkgs/pkg/config"
    "github.com/cosmos-toolkit/pkgs/pkg/errors"
    "github.com/cosmos-toolkit/pkgs/pkg/contextx"
    "github.com/cosmos-toolkit/pkgs/pkg/validator"
)
```

No mesmo monorepo (ex.: projeto em `templates/monorepo-starter`):

```bash
go mod edit -replace github.com/cosmos-toolkit/pkgs=../packages
go mod tidy
```

---

## Núcleo transversal (disponível)

| Pacote        | Descrição                                                                                       |
| ------------- | ----------------------------------------------------------------------------------------------- |
| **logger**    | Wrapper sobre slog. Modo CLI / API / Worker. Context-aware (trace_id, request_id via contextx). |
| **config**    | Carregar env (dotenv), validação de config, defaults.                                           |
| **errors**    | Erros tipados, Is/As, mapeamento para HTTP / exit code / retry, stack opcional.                 |
| **contextx**  | Helpers para context: timeout padrão, metadata (tenant, trace, user), cancelamento.             |
| **validator** | Wrapper validator/v10, mensagens padronizadas, reuso API/CLI.                                   |

---

## Manifest (`manifest.yaml`)

Define, por pacote, dependências usadas pelo Cosmos CLI:

- **copy_deps:** pacotes deste repositório copiados junto (ex.: `logger` → `[contextx]`).
- **go_get:** dependências externas instaladas com `go get` após a cópia.

---

## Roadmap (evoluir conforme necessidade)

- **APIs:** httpx, httperrors, router, auth, pagination
- **Workers:** worker, queue, cron, retry, idempotency
- **Domínio:** result, mapper, clock, uuid
- **Infra:** db, tx, cache, outbox
- **Observabilidade:** metrics, tracing, health
- **CLI:** cli, prompt, output
- **Testes:** testkit, fixture

---

## Estrutura

```
pkg/
├── logger/    # slog + context
├── config/    # env + validação
├── errors/    # tipados + HTTP/exit/retry
├── contextx/  # timeout + metadata
└── validator/ # go-playground validator
```
