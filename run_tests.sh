#!/bin/bash

echo "=== Database Service Test Suite ==="
echo

echo "1. Running All Services Test..."
go run cmd/test_all_services/main.go
echo

echo "2. Running Client Test..."
go run cmd/test_client/main.go
echo

echo "3. Running Production DB Test..."
go run cmd/test_prod_db/main.go
echo

echo "=== Test Suite Complete ==="