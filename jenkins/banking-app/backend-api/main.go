package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Account represents a bank account
type Account struct {
	ID      string  `json:"id"`
	Owner   string  `json:"owner"`
	Balance float64 `json:"balance"`
}

// Transaction represents a banking transaction
type Transaction struct {
	ID          string  `json:"id"`
	AccountID   string  `json:"account_id"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"` // "deposit", "withdrawal", "transfer"
	Description string  `json:"description"`
	Status      string  `json:"status"` // "pending", "completed", "failed"
}

// InMemoryDB is a simple in-memory database
type InMemoryDB struct {
	accounts     map[string]Account
	transactions map[string]Transaction
	accountID    int
	transactionID int
	mu           sync.RWMutex
}

var db = InMemoryDB{
	accounts:     make(map[string]Account),
	transactions: make(map[string]Transaction),
	accountID:    1000,
	transactionID: 5000,
}

func main() {
	// Register routes
	http.HandleFunc("/accounts", handleAccounts)
	http.HandleFunc("/accounts/", handleAccountById)
	http.HandleFunc("/transactions", handleTransactions)
	http.HandleFunc("/transactions/", handleTransactionById)
	
	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})
	
	// Start server
	fmt.Println("Banking API Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	switch r.Method {
	case "GET":
		// List all accounts
		db.mu.RLock()
		accounts := make([]Account, 0, len(db.accounts))
		for _, account := range db.accounts {
			accounts = append(accounts, account)
		}
		db.mu.RUnlock()
		
		json.NewEncoder(w).Encode(accounts)
		
	case "POST":
		// Create a new account
		var newAccount Account
		err := json.NewDecoder(r.Body).Decode(&newAccount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		db.mu.Lock()
		db.accountID++
		newAccount.ID = strconv.Itoa(db.accountID)
		db.accounts[newAccount.ID] = newAccount
		db.mu.Unlock()
		
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newAccount)
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleAccountById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Extract account ID from URL
	id := r.URL.Path[len("/accounts/"):]
	
	db.mu.RLock()
	account, exists := db.accounts[id]
	db.mu.RUnlock()
	
	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}
	
	switch r.Method {
	case "GET":
		// Get account by ID
		json.NewEncoder(w).Encode(account)
		
	case "PUT":
		// Update account
		var updatedAccount Account
		err := json.NewDecoder(r.Body).Decode(&updatedAccount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		// Preserve ID
		updatedAccount.ID = id
		
		db.mu.Lock()
		db.accounts[id] = updatedAccount
		db.mu.Unlock()
		
		json.NewEncoder(w).Encode(updatedAccount)
		
	case "DELETE":
		// Delete account
		db.mu.Lock()
		delete(db.accounts, id)
		db.mu.Unlock()
		
		w.WriteHeader(http.StatusNoContent)
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	switch r.Method {
	case "GET":
		// List all transactions
		db.mu.RLock()
		transactions := make([]Transaction, 0, len(db.transactions))
		for _, txn := range db.transactions {
			transactions = append(transactions, txn)
		}
		db.mu.RUnlock()
		
		json.NewEncoder(w).Encode(transactions)
		
	case "POST":
		// Create a new transaction
		var newTxn Transaction
		err := json.NewDecoder(r.Body).Decode(&newTxn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		// Verify account exists
		db.mu.RLock()
		_, accountExists := db.accounts[newTxn.AccountID]
		db.mu.RUnlock()
		
		if !accountExists {
			http.Error(w, "Account not found", http.StatusBadRequest)
			return
		}
		
		// Process transaction
		db.mu.Lock()
		db.transactionID++
		newTxn.ID = strconv.Itoa(db.transactionID)
		newTxn.Status = "pending" // Set initial status
		db.transactions[newTxn.ID] = newTxn
		
		// Update account balance for completed transactions
		if newTxn.Type == "deposit" {
			account := db.accounts[newTxn.AccountID]
			account.Balance += newTxn.Amount
			newTxn.Status = "completed"
			db.accounts[newTxn.AccountID] = account
			db.transactions[newTxn.ID] = newTxn
		} else if newTxn.Type == "withdrawal" {
			account := db.accounts[newTxn.AccountID]
			if account.Balance >= newTxn.Amount {
				account.Balance -= newTxn.Amount
				newTxn.Status = "completed"
				db.accounts[newTxn.AccountID] = account
				db.transactions[newTxn.ID] = newTxn
			} else {
				newTxn.Status = "failed"
				db.transactions[newTxn.ID] = newTxn
			}
		}
		db.mu.Unlock()
		
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTxn)
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleTransactionById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Extract transaction ID from URL
	id := r.URL.Path[len("/transactions/"):]
	
	db.mu.RLock()
	txn, exists := db.transactions[id]
	db.mu.RUnlock()
	
	if !exists {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}
	
	switch r.Method {
	case "GET":
		// Get transaction by ID
		json.NewEncoder(w).Encode(txn)
		
	case "PUT":
		// Update transaction status (e.g., for workflow purposes)
		var updatedTxn Transaction
		err := json.NewDecoder(r.Body).Decode(&updatedTxn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		// Only allow status updates
		txn.Status = updatedTxn.Status
		
		db.mu.Lock()
		db.transactions[id] = txn
		db.mu.Unlock()
		
		json.NewEncoder(w).Encode(txn)
		
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}