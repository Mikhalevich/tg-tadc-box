# tg-tadc-box

A Telegram bot application for managing collectible card boxes, rewards, and user interactions.

## Overview

This project implements a Telegram bot with three main components:

- **Bot** (`cmd/bot`) — the Telegram bot that handles user commands, button callbacks, card viewing, box opening, shop purchases, likes, and referral links.
- **Poller** (`cmd/poller`) — a background worker that monitors for boxes ready to be opened and sends notifications to users.
- **Outbox Poller** (`cmd/outboxpoller`) — a background worker that processes pending outbox messages and delivers them via the Telegram bot.

## Prerequisites

- Go 1.26+
- Docker & Docker Compose (for local development)
- Make

## Quick Start

### Local Development with Docker Compose

```bash
make compose-up
```

This starts:
- The bot, poller, and outboxpoller services
- PostgreSQL 16.3-alpine (port 5432)
- Jaeger tracing UI (ports 16686, 14268, 4318)
- SQL migration job (runs on startup)

### Building Binaries

```bash
make build
```

Binaries are placed in `./bin/`:
- `bin/bot`
- `bin/poller`
- `bin/outboxpoller`

## Configuration

Each component loads its configuration from YAML files in `./config/`:

| File | Component |
|------|-----------|
| `config-bot.yaml` | Bot |
| `config-poller.yaml` | Poller |
| `config-outboxpoller.yaml` | Outbox Poller |
| `dbconfig.yml` | Database migration |

Example configuration files are provided with `-example` suffix. Copy them to the non-example names and fill in real values before running.

### Bot Configuration

Key fields in `config-bot.yaml`:

```yaml
log_level: "info"
tracing:
  endpoint: "jaeger:4318"
  service_name: "box-bot"
bot:
  token: "TELEGRAM_BOT_TOKEN"
  name: "bot_name"
  webhook_token: "webhook_token or empty if running in polling mode"
postgres:
  connection: "host=postgres port=5432 user=bot password=bot dbname=bot sslmode=disable"
box_reward_percent:
  common:
    legendary: 1
    epic: 5
    rare: 20
  # ...
box_costs:
  common: { amount: 0 }
  rare: { amount: 25 }
  epic: { amount: 125 }
  legendary: { amount: 500 }
box_wait_period:
  common: 10s
  rare: 10s
  epic: 10s
  legendary: 10s
bonus_box_attempts:
  common: 3
  rare: 3
  epic: 3
abstraction_costs:
  common: 1
  rare: 5
  epic: 25
  legendary: 100
```

### Poller Configuration

```yaml
log_level: "info"
tracing:
  endpoint: "jaeger:4318"
  service_name: "poller"
postgres:
  connection: "host=postgres port=5432 user=bot password=bot dbname=bot sslmode=disable"
box_ready_notification_worker:
  count: 1
  interval: 5s
  batch_size: 10
```

### Outbox Poller Configuration

```yaml
log_level: "info"
tracing:
  endpoint: "jaeger:4318"
  service_name: "outbox-poller"
bot:
  token: "TELEGRAM_BOT_TOKEN"
postgres:
  connection: "host=postgres port=5432 user=bot password=bot dbname=bot sslmode=disable"
worker:
  count: 2
  interval: 1s
  batch_size: 10
  max_retry_count: 10
```

## Make Targets

| Target | Description |
|--------|-------------|
| `make build` | Build all binaries into `./bin/` |
| `make test` | Run all Go tests |
| `make bench` | Run benchmarks |
| `make lint` | Run golangci-lint |
| `make fmt` | Format code with golangci-lint |
| `make vendor` | Run `go mod tidy` and `go mod vendor` |
| `make generate` | Run `go generate ./...` |
| `make compose-up` | Start all services with Docker Compose |
| `make compose-down` | Stop Docker Compose services |
| `make load-assets` | Load initial DB assets |
| `make minikube-*` | Minikube Helm deployment targets |
| `make do-helm-*` | DigitalOcean Helm deployment targets |

## Testing

```bash
make test
```

## Project Structure

```
cmd/
├── bot/              # Telegram bot main entrypoint
│   └── internal/
│       ├── app/      # Bot application setup & routing
│       ├── config/   # Bot configuration types
│       └── setup/    # Bot initialization
├── poller/           # Box-ready notification worker
│   └── internal/
│       ├── app/
│       ├── config/
│       └── setup/
└── outboxpoller/     # Outbox message delivery worker
    └── internal/
        ├── app/
        ├── config/
        └── setup/

internal/
├── adapter/              # External adapters
│   ├── imageprovider/    # Card image assets
│   ├── markdownescaper/  # Markdown escaping
│   ├── messagesender/    # Telegram message sending
│   ├── repository/       # PostgreSQL repositories
│   └── timeprovider/     # Time utilities
├── domain/
│   ├── model/            # Domain models (box, card, player, reward, etc.)
│   └── service/          # Domain services & use cases
│       ├── boxsvc/       # Box operations
│       ├── likesvc/      # Like mechanics
│       ├── messagesvc/   # Message formatting
│       ├── notificationsvc/ # User notifications
│       ├── outbox/       # Outbox message processing
│       ├── playersvc/    # Player operations
│       ├── referralsvc/  # Referral links
│       └── rewardsvc/    # Reward logic
└── infra/
    ├── application/      # Application lifecycle (config, logger, tracing)
    ├── logger/           # Logging (Logrus + OpenTelemetry)
    └── tracing/          # OpenTelemetry tracer

config/                   # YAML configuration files
script/
├── docker/               # Docker Compose & Dockerfiles
├── db/                   # Database scripts & asset SQL
└── k8s/                  # Kubernetes Helm charts (minikube & DO)
```

## CI/CD

GitHub Actions workflows are defined in `.github/workflows/`:

- **build_and_test.yml** — builds and tests on every push and PR
- **golangci-lint.yml** — runs linting on master and PRs
- **deploy_by_tag.yml** — builds images and deploys to Kubernetes on tag push (`v*`)

## Technologies

- Go 1.26
- Telegram Bot API (`github.com/go-telegram/bot`)
- PostgreSQL (`pgx`, `sqlx`, `sql-migrate`)
- OpenTelemetry (tracing)
- Logrus (logging)
- Docker Compose / Kubernetes Helm (deployment)
- golangci-lint (linting)

## License

No license file is present in the repository.
