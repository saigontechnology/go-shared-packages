# Transaction Context

Transaction context is a helper package that handle transaction in database
and inject into the Context, allowing the repository to check if 
the context contains any existing transaction to prioritize the given 
transaction over the connection pool given by the initialization


## Usage

### Transaction Control
- Usage across service

```
// initialize the package first
ctxTransaction := ctxtransaction.New()

// inject to context
ctx, err := ctxTransaction.BeginWithDB(ctx, gorm.DB)

if err != nil {
    // handle err
}

// 
if err := ctxTransaction.CommitFromContext(ctx); err != nil {
    // handle commit error
}

if err := ctxTransaction.RollbackFromContext(ctx); err != nil {
    // handle rollback error
}

```

### Usage Injection
- There are two ways of injection: either inject as initialized object or 
inject as closure function
- By using object initialization and closure, we can still support the usage of
custom transaction key
- But if using direct function to get session, only default transaction key is
supported, in this way, it's easier and repository doesn't need to keep any
helper object or helper closure 

1. Inject into repository with object initialization

```
// repository.go
type Repository struct {
    ctxTransaction ctxTransaction.TransactionContext
}

func NewRepository(db *gorm.DB) Repository {
    return Repository{
        ctxTransaction: ctxtransaction.NewWithDB(db),
    }
}

func(r repository) Create(ctx context.Context, model *model.Model) error {
    db := r.ctxTransaction.Session(ctx)
}
```

2. Inject as closure function

```
// repository.go
type Repository struct {
    dbSession ctxtransaction.GetSessionFunc
}

func NewRepository(transactionKey string, db *gorm.DB) Repository {
    return Repository{
        dbSession: ctxtransaction.SessionWithFallback(transactionKey, db),
    }
}

func(r repository) Create(ctx context.Context, model *model.Model) error {
    result := r.dbSession(ctx).Where(...).Find(...)
    return result.Error
}
```


3. Direct usage of helper func

```
// repository.go
type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return Repository{
        db: db,
    }
}

func(r repository) Create(ctx context.Context, model *model.Model) error {
    result := ctxtransaction.SessionFromContext(ctx, r.db).
                Where(...).
                Find(...)
    return result.Error
}
```
