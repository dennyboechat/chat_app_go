# A simple chat application using Go language

## Features  
-	Users send messages
-	Users receive messages
-	Persist chat history
-	Message filtering
-	Search by user or keyword

## Structure  
-	Database to store chat history
-  -	timestamp and user ID
-	Built-in concurrency features with goroutines
-	Channels for message passing
-	Focus on performance

## Manageable components
### 1. User Simulator 
  - Responsible for generating user IDs and input.
  - Can be CLI-driven or simulated in code.
### 2. Message Handler
  - Takes user input and formats it into a Message struct with timestamp.
### 3. Message Storage
  - In-memory map[string][]Message to store messages per user.
  - Search and filter logic built on top of this.
### 4. Concurrency Controller
  - Goroutines for user simulations.
  - Channels for message passing between users.
### 5. Search & filter Engine
  - Search by user ID or keyword.
  - Basic keyword match in message content. 

## Challenges
### 1. Concurrency Race conditions
  - Managing access to shared message storage with multiple goroutines.
### 2. Channel Management
  - Ensuring proper channel closing and avoiding deadlocks or leaks.
### 3. Search Scalability
  - Filtering messages efficiently as message history grows (even in memory).
### 4. Testing Simultaneous Users
   - Simulating and debugging concurrent user activity can get tricky

## GitHub repository  
-	https://github.com/dennyboechat/chat_app_go

## Tasks assignment  (Tenative, can be changed)
-	Aanish - Set up core messages structs, in-memory storage, and filtering. Implement goroutines, channels and user simulation logic (These things can be split up or if you have a different idea let me know.)
-	Denny - 

## Timeline  
-	24-May (Sat) – Send deliverable 1
-	07-Jun (Sat) – Send deliverable 2
-	? – Prepare video presentation
-	21-Jun (Sat) - Send deliverable 3

## Above timeslines can be filled based on these
### Week                         Task
Week 3                   Set up repo, base structs, user simulator
Week 4                   Implement messaging flow & concurrency
Week 5                   Finalize features, prepare comparison report
Week 6                   Polish filtering & prepare demo
Week 7                   Final testing and video demonstration
