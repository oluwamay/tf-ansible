#!/usr/bin/env python3
"""
Transaction Processing Service for Core Banking Application - FastAPI Version
"""

import os
import time
import logging
from typing import Dict, List, Optional, Union

import requests
from fastapi import FastAPI, HTTPException, Query
from pydantic import BaseModel, Field

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="Transaction Processing Service",
    description=(
        "Handles complex transaction processing and analysis for core banking"
    ),
    version="1.0.0"
)

# Configuration
API_HOST = os.environ.get('API_HOST', 'api')
API_PORT = os.environ.get('API_PORT', '8080')
API_BASE_URL = f"http://{API_HOST}:{API_PORT}"


# Pydantic models for request/response validation
class Transaction(BaseModel):
    id: str
    account_id: Optional[str] = None
    amount: Optional[float] = None
    type: Optional[str] = None
    status: Optional[str] = None


class TransactionResponse(BaseModel):
    transaction_id: str
    status: str
    message: str


class BatchResponse(BaseModel):
    results: List[TransactionResponse]


class AnalyticsResponse(BaseModel):
    total_transactions: int
    completed: Optional[int] = None
    pending: Optional[int] = None
    failed: Optional[int] = None
    total_deposits: Optional[float] = None
    total_withdrawals: Optional[float] = None
    net_flow: Optional[float] = None
    message: Optional[str] = None
    error: Optional[str] = None


class HealthResponse(BaseModel):
    status: str


class TransactionProcessor:
    """Handle complex transaction processing and analysis"""

    def __init__(self):
        self.transaction_queue = []

    def add_transaction(self, transaction: Dict) -> bool:
        """Add a transaction to the processing queue"""
        self.transaction_queue.append(transaction)
        logger.info("Added transaction %s to queue", transaction['id'])
        return True

    def process_batch(self) -> List[Dict]:
        """Process all transactions in the queue"""
        if not self.transaction_queue:
            logger.info("No transactions to process")
            return []

        results = []
        for transaction in self.transaction_queue:
            try:
                # Process the transaction
                result = self._process_transaction(transaction)
                results.append(result)
            except Exception as e:
                logger.error(
                    "Error processing transaction %s: %s",
                    transaction['id'], str(e)
                )
                results.append({
                    "transaction_id": transaction['id'],
                    "status": "error",
                    "message": str(e)
                })

        # Clear the queue
        self.transaction_queue = []
        return results

    def _process_transaction(self, transaction: Dict) -> Dict:
        """Process an individual transaction"""
        # Simulate processing time
        time.sleep(0.1)

        # Update transaction status via API
        try:
            url = f"{API_BASE_URL}/transactions/{transaction['id']}"
            response = requests.put(
                url,
                json={"status": "completed"},
                timeout=10
            )

            if response.status_code == 200:
                logger.info(
                    "Successfully processed transaction %s", transaction['id']
                )
                return {
                    "transaction_id": transaction['id'],
                    "status": "success",
                    "message": "Transaction processed"
                }
            else:
                logger.warning(
                    "Failed to update transaction %s: %s",
                    transaction['id'], response.text
                )
                return {
                    "transaction_id": transaction['id'],
                    "status": "error",
                    "message": f"API error: {response.text}"
                }
        except Exception as e:
            logger.error("Exception during API call: %s", str(e))
            return {
                "transaction_id": transaction['id'],
                "status": "error",
                "message": f"API error: {str(e)}"
            }

    def analyze_transactions(self, account_id: Optional[str] = None) -> Dict:
        """Analyze transactions for an account or all accounts"""
        try:
            # Get transactions from API
            url = f"{API_BASE_URL}/transactions"
            response = requests.get(url, timeout=10)

            if response.status_code != 200:
                return {
                    "error": f"Failed to fetch transactions: {response.text}"
                }

            transactions = response.json()

            # Filter by account if specified
            if account_id:
                transactions = [
                    t for t in transactions if t.get('account_id') == account_id
                ]

            # Basic analytics
            total_count = len(transactions)
            if total_count == 0:
                return {
                    "total_transactions": 0,
                    "message": "No transactions found"
                }

            completed = len([t for t in transactions if t.get('status') == 'completed'])
            pending = len([t for t in transactions if t.get('status') == 'pending'])
            failed = len([t for t in transactions if t.get('status') == 'failed'])

            deposits = sum([t.get('amount', 0) for t in transactions if t.get('type') == 'deposit' and t.get('status') == 'completed'])
            withdrawals = sum([t.get('amount', 0) for t in transactions if t.get('type') == 'withdrawal' and t.get('status') == 'completed'])

            return {
                "total_transactions": total_count,
                "completed": completed,
                "pending": pending,
                "failed": failed,
                "total_deposits": deposits,
                "total_withdrawals": withdrawals,
                "net_flow": deposits - withdrawals
            }

        except Exception as e:
            logger.error("Error analyzing transactions: %s", str(e))
            return {"error": str(e)}


# Create processor instance
processor = TransactionProcessor()


@app.post("/process", status_code=202, response_model=Dict)
async def process_transaction(transaction: Transaction):
    """Add a transaction to the processing queue"""
    processor.add_transaction(transaction.dict())
    return {"status": "queued", "transaction_id": transaction.id}


@app.post("/process/batch", response_model=BatchResponse)
async def process_batch():
    """Process all transactions in the queue"""
    results = processor.process_batch()
    return {"results": results}


@app.get("/analytics", response_model=AnalyticsResponse)
async def get_analytics(account_id: Optional[str] = Query(None, description="Filter by account ID")):
    """Get transaction analytics"""
    analytics = processor.analyze_transactions(account_id)
    return analytics


@app.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy"}

if __name__ == '__main__':
    import uvicorn
    logger.info("Starting Transaction Processing Service...")
    uvicorn.run(app, host="0.0.0.0", port=5000)