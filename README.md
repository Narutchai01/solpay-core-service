# solpay-core-service

```
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=solpay_user
DB_PASSWORD=solpay_password
DB_NAME=solpay_db
```

### TopUp Flow (Event-Driven)

```mermaid
sequenceDiagram
    title TopUp Flow
    actor User
    participant API as Handler & Service
    participant DB as Postgres DB
    participant Orch as Transaction Orchestrator
    participant Solana as Solana Worker
    participant Balance as Balance Worker

    User->>+API: POST /api/v1/topups
    API->>API: Extract AccountID & Validate
    API->>+DB: Create Transaction Record (Status: PENDING)
    DB-->>-API: Transaction Entity
    API->>+Orch: Publish TransactionMessage (SOLANA_SUBMITTED)
    API-->>-User: 200 OK (TransactionDTO)

    Note over API, Balance: Async Event-Driven Flow starts

    Orch->>+DB: Get & Update Tx Status to SOLANA_SUBMITTED
    DB-->>-Orch: DB Updated
    Orch->>+Solana: Publish to SOLANA_WORK_QUEUE
    Solana->>Solana: Handle On-Chain Verification
    Solana->>Orch: Publish Status Update (SOLANA_SUCCESS)

    Orch->>+DB: Update Tx Status to SOLANA_SUCCESS
    DB-->>-Orch: DB Updated
    Orch->>+Balance: Publish to BALANCE_QUEUE (Action: DEPOSIT)
    Balance->>+DB: Add to User Balance
    DB-->>-Balance: Balance Updated
    Balance->>Orch: Publish Status Update (BALANCE_UPDATED)

    Orch->>+DB: Update Tx Status to COMPLETED
    DB-->>-Orch: Success
    Note right of Orch: TopUp Transaction is fully finalized
```

### Offchain Flow (Event-Driven)

```mermaid
sequenceDiagram
    title Offchain Flow (Withdrawal)
    actor User
    participant API as Handler & Service
    participant DB as Postgres DB
    participant Orch as Transaction Orchestrator
    participant Balance as Balance Worker
    participant Payment as Payment Worker
    participant Storage as Supabase

    User->>+API: POST /api/v1/offchains
    API->>API: Extract AccountID & Validate Quote
    API->>+DB: Create Transaction Record (Status: PENDING)
    DB-->>-API: Transaction Entity
    API->>+Orch: Publish TransactionMessage (BALANCE_WITHDRAWING)
    API-->>-User: 200 OK (TransactionDTO)

    Note over API, Storage: Async Event-Driven Flow starts

    Orch->>+DB: Get & Update Tx Status to BALANCE_WITHDRAWING
    DB-->>-Orch: DB Updated
    Orch->>+Balance: Publish to BALANCE_QUEUE (Action: WITHDRAW)
    Balance->>+DB: Deduct from User Balance
    DB-->>-Balance: Balance Updated
    Balance->>Orch: Publish Status Update (BALANCE_UPDATED)

    Orch->>+DB: Update Tx Status to BALANCE_UPDATED
    DB-->>-Orch: DB Updated
    Orch->>+Payment: Publish to PAYMENT_QUEUE
    Payment->>Payment: Process Payment (PromptPay/Omise)
    Payment->>Orch: Publish Status Update (PAYMENT_SUCCESS)

    Orch->>+DB: Get Transaction
    DB-->>-Orch: Tx Entity
    Orch->>Orch: Generate Slip Data
    Orch->>+Storage: Upload Slip (Supabase)
    Storage-->>-Orch: Slip URL
    Orch->>+DB: Update Offchain Record with Slip URL
    DB-->>-Orch: DB Updated
    Orch->>+DB: Update Tx Status to COMPLETED
    DB-->>-Orch: Success
    Note right of Orch: Offchain Transaction is fully finalized
```
