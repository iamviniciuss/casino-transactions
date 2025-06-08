
✅ TDD — Technical Design Document

Módulo: Wallet
Responsável técnico: [Seu Nome]
Data: [Hoje]

⸻

🔧 Arquitetura

Estrutura de pacotes:

internal/
  wallet/
    entity.go
    repository.go
    usecase.go
    service.go
    history.go
    tests/
      service_test.go


⸻

📌 Interfaces

Repository

type Repository interface {
	GetByUserID(ctx context.Context, userID string) (*Wallet, error)
	UpdateBalance(ctx context.Context, userID string, amount int64) error
	RecordHistory(ctx context.Context, op WalletOperation) error
}

Usecase

type Usecase interface {
	GetBalance(ctx context.Context, userID string) (*Wallet, error)
	Credit(ctx context.Context, userID string, amount int64, origin string) error
	Debit(ctx context.Context, userID string, amount int64, origin string) error
}


⸻

⚙️ Lógica principal

Credit
	•	Soma o valor ao saldo
	•	Atualiza saldo
	•	Salva histórico (se ativado)

Debit
	•	Valida saldo ≥ valor
	•	Subtrai do saldo
	•	Atualiza
	•	Salva histórico (se ativado)

⸻

🧪 Testes

service_test.go — Casos:
	•	✅ GetBalance retorna valor correto
	•	✅ Credit adiciona corretamente
	•	✅ Debit subtrai corretamente
	•	❌ Debit com saldo insuficiente retorna erro
	•	✅ Histórico é salvo quando ativado

⸻

🧱 Banco

Tabela: wallets

user_id     TEXT (PK)
balance     BIGINT
updated_at  TIMESTAMP

Tabela: wallet_operations

id          UUID
user_id     TEXT
amount      BIGINT
type        TEXT -- credit | debit
origin      TEXT
created_at  TIMESTAMP

