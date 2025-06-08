🧾 PRD — Product Requirements Document

Módulo: Wallet (Carteira de Usuário)
Objetivo: Gerenciar o saldo dos usuários do sistema de casino, permitindo consultar, creditar, debitar e registrar histórico de operações.

⸻

🎯 Visão Geral

O módulo wallet será responsável por armazenar e manipular o saldo dos usuários. Ele se integrará com outros módulos como transaction, game e user, permitindo atualizações em tempo real no saldo.

⸻

📚 Funcionalidades

1. Obter saldo atual
	•	O sistema deve expor uma função para consultar o saldo atual de um usuário.
	•	Entrada: userID
	•	Saída: saldo atual (int64 em centavos)

⸻

2. Creditar valor
	•	Deve ser possível adicionar crédito ao saldo de um usuário (ex: após vitória ou depósito).
	•	Entrada: userID, valor positivo
	•	Saída: saldo atualizado

⸻

3. Debitar valor
	•	Deve ser possível subtrair valor do saldo (ex: ao realizar uma aposta).
	•	Entrada: userID, valor positivo
	•	Validação: saldo não pode ficar negativo
	•	Saída: saldo atualizado ou erro de saldo insuficiente

⸻

4. Registrar operações no histórico (opcional)
	•	Armazenar registros de crédito/débito com:
	•	userID
	•	valor
	•	tipo (crédito/débito)
	•	origem (ex: “depósito”, “ganho”, “aposta”)
	•	timestamp
	•	Útil para relatórios e auditoria

⸻

🧩 Regras de Negócio
	•	Todas as operações devem ser atômicas
	•	Saldos devem ser armazenados em centavos (int64)
	•	Não permitir saldo negativo
	•	Usuários inativos devem ter carteira bloqueada (futuro)

⸻

📡 Integrações
	•	Pode ser chamado via REST ou consumer Kafka (ex: após processar transação)

⸻

📈 Métricas de sucesso
	•	100% de cobertura de testes nos serviços
	•	Suporte a concorrência (ex: múltiplos créditos concorrentes)
	•	Histórico visível e rastreável

⸻
