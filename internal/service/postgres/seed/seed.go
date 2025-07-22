package main

import (
	"adel/internal/service/postgres"
	"database/sql"
	"fmt"
	"log"
)

func main() {
	db, err := sql.Open("pgx", "host=localhost user=admin password=adminpassword dbname=adel port=5440 sslmode=disable")
	if err != nil {
		log.Fatalf("db open %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("db: ping %v", err)
	}
	fmt.Println("connected to database")

	problemStore := postgres.NewPostgresProblemStore(db)

	err = seedProblems(problemStore)
	if err != nil {
		log.Fatalf("failed to seed problems: %v", err)
	}

	log.Println("seeding completed successfully")
}

func seedProblems(store *postgres.PostgresProblemStore) error {
	problems := []postgres.Problem{
		{
			UserID:      1,
			Title:       "Simple Addition",
			Description: "Write a program that reads two integers and outputs their sum.\n\nInput: Two space-separated integers a and b (1 ≤ a, b ≤ 1000)\nOutput: The sum of a and b",
			Difficulty:  "easy",
			TimeLimit:   3000,
			MemoryLimit: 256,
			IsActive:    true,
			TestCases: []postgres.TestCase{
				{InputData: "1 2", OutputData: "3", IsActive: true, UserID: 1},
				{InputData: "2 3", OutputData: "5", IsActive: true, UserID: 1},
				{InputData: "10 15", OutputData: "25", IsActive: true, UserID: 1},
			},
		},
		{
			UserID:      1,
			Title:       "Even or Odd",
			Description: "Write a program that determines if a given number is even or odd.\n\nInput: A single integer n (1 ≤ n ≤ 100)\nOutput: 'Even' if the number is even, 'Odd' if the number is odd",
			Difficulty:  "easy",
			TimeLimit:   2000,
			MemoryLimit: 128,
			IsActive:    true,
			TestCases: []postgres.TestCase{
				{InputData: "4", OutputData: "Even", IsActive: true, UserID: 1},
				{InputData: "7", OutputData: "Odd", IsActive: true, UserID: 1},
				{InputData: "100", OutputData: "Even", IsActive: true, UserID: 1},
			},
		},
		{
			UserID:      1,
			Title:       "Count Digits",
			Description: "Write a program that counts the number of digits in a given positive integer.\n\nInput: A single positive integer n (1 ≤ n ≤ 1000000)\nOutput: The number of digits in n",
			Difficulty:  "easy",
			TimeLimit:   2000,
			MemoryLimit: 128,
			IsActive:    true,
			TestCases: []postgres.TestCase{
				{InputData: "123", OutputData: "3", IsActive: true, UserID: 1},
				{InputData: "7", OutputData: "1", IsActive: true, UserID: 1},
				{InputData: "1000", OutputData: "4", IsActive: true, UserID: 1},
			},
		},
		{
			UserID:      1,
			Title:       "Reverse Number",
			Description: "Write a program that reverses the digits of a given positive integer.\n\nInput: A single positive integer n (1 ≤ n ≤ 100000)\nOutput: The number with its digits reversed",
			Difficulty:  "easy",
			TimeLimit:   2000,
			MemoryLimit: 128,
			IsActive:    true,
			TestCases: []postgres.TestCase{
				{InputData: "123", OutputData: "321", IsActive: true, UserID: 1},
				{InputData: "7", OutputData: "7", IsActive: true, UserID: 1},
				{InputData: "1230", OutputData: "321", IsActive: true, UserID: 1},
			},
		},
		{
			UserID:      1,
			Title:       "Prime Checker",
			Description: "Write a program that determines if a given number is prime.\n\nInput: A single integer n (2 ≤ n ≤ 10000)\nOutput: 'Prime' if the number is prime, 'Not Prime' if it's not prime",
			Difficulty:  "medium",
			TimeLimit:   3000,
			MemoryLimit: 256,
			IsActive:    true,
			TestCases: []postgres.TestCase{
				{InputData: "7", OutputData: "Prime", IsActive: true, UserID: 1},
				{InputData: "4", OutputData: "Not Prime", IsActive: true, UserID: 1},
				{InputData: "97", OutputData: "Prime", IsActive: true, UserID: 1},
				{InputData: "100", OutputData: "Not Prime", IsActive: true, UserID: 1},
			},
		},
		{
			UserID:      1,
			Title:       "Longest Common Subsequence",
			Description: "Given two strings, find the length of their longest common subsequence.\n\nInput: Two lines, each containing a string (1 ≤ length ≤ 100, containing only lowercase letters)\nOutput: The length of the longest common subsequence",
			Difficulty:  "hard",
			TimeLimit:   5000,
			MemoryLimit: 512,
			IsActive:    true,
			TestCases: []postgres.TestCase{
				{InputData: "abcde\nace", OutputData: "3", IsActive: true, UserID: 1},
				{InputData: "abc\ndef", OutputData: "0", IsActive: true, UserID: 1},
				{InputData: "programming\ngrading", OutputData: "4", IsActive: true, UserID: 1},
			},
		},
	}

	log.Println("seeding problems...")
	for _, problem := range problems {
		_, err := store.CreateProblem(&problem)
		if err != nil {
			return fmt.Errorf("could not create problem %q: %w", problem.Title, err)
		}
	}
	log.Println("problems seeded")

	return nil
}
