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

## Pacotes disponíveis

### Núcleo transversal

| Pacote        | Descrição                                                                                       |
| ------------- | ----------------------------------------------------------------------------------------------- |
| **logger**    | Wrapper sobre slog. Modo CLI / API / Worker. Context-aware (trace_id, request_id via contextx). |
| **config**    | Carregar env (dotenv), validação de config, defaults.                                           |
| **errors**    | Erros tipados, Is/As, mapeamento para HTTP / exit code / retry, stack opcional.                 |
| **contextx**  | Helpers para context: timeout padrão, metadata (tenant, trace, user), cancelamento.             |
| **validator** | Wrapper validator/v10, mensagens padronizadas, reuso API/CLI.                                   |
| **clock**     | Abstração de tempo (Clock interface + Real/Fake). Facilita testes.                              |
| **retry**     | Retry com backoff exponencial, jitter e max attempts.                                           |
| **testkit**   | NopLogger, ContextWithIDs. Helpers para testes.                                                 |
| **cli**       | Exit codes padronizados (ExitOK, ExitErr, …). Uso com pkg/errors.ExitCode.                      |

### APIs / HTTP

| Pacote    | Descrição                                                |
| --------- | -------------------------------------------------------- |
| **httpx** | Server bootstrap, graceful shutdown, healthcheck padrão. |

### Workers / Jobs / Crons

| Pacote     | Descrição                                                                               |
| ---------- | --------------------------------------------------------------------------------------- |
| **worker** | Worker pool, concurrency configurável, retry (base para SQS, cron, fila).               |
| **queue**  | Interface Publish/Consume + implementação in-memory (SQS/Rabbit podem ser adicionados). |
| **cron**   | Wrapper robfig/cron para agendamento de jobs.                                           |
| **outbox** | Padrão outbox: persist + publish (event-driven).                                        |

### Persistência / Infra

| Pacote    | Descrição                                                                 |
| --------- | ------------------------------------------------------------------------- |
| **db**    | Bootstrap DB, pool, healthcheck (driver importado pelo usuário).          |
| **cache** | Interface Get/Set/Delete + in-memory com TTL (Redis pode ser adicionado). |

### Observabilidade

| Pacote      | Descrição                                                      |
| ----------- | -------------------------------------------------------------- |
| **metrics** | Helpers Prometheus: Counter, Histogram, Handler para /metrics. |
| **tracing** | Wrapper OpenTelemetry: Tracer, StartSpan, context propagation. |

---

## Manifest (`manifest.yaml`)

Define, por pacote, dependências usadas pelo Cosmos CLI:

- **copy_deps:** pacotes deste repositório copiados junto (ex.: `logger` → `[contextx]`).
- **go_get:** dependências externas instaladas com `go get` após a cópia.

---

## Roadmap (evoluir conforme necessidade)

- **APIs:** httperrors, router, auth, pagination
- **Workers:** idempotency
- **Domínio:** result, mapper, uuid
- **Infra:** tx, health
- **CLI:** prompt, output (wrapper Cobra)
- **Testes:** fixture

---

## Estrutura

```
pkg/
├── logger/    # slog + context
├── config/    # env + validação
├── errors/    # tipados + HTTP/exit/retry
├── contextx/  # timeout + metadata
├── validator/ # go-playground validator
├── clock/     # abstração de tempo (Real/Fake)
├── retry/     # backoff + jitter + max attempts
├── testkit/   # NopLogger, ContextWithIDs
├── cli/       # exit codes padronizados
├── httpx/     # server + graceful shutdown + health
├── worker/    # worker pool + retry
├── queue/     # interface + in-memory
├── cron/      # scheduler (robfig/cron)
├── db/        # pool + healthcheck
├── cache/     # interface + in-memory
├── metrics/   # Prometheus helpers
├── tracing/   # OpenTelemetry wrapper
└── outbox/    # persist + publish
```
