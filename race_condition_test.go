package golanggoroutine

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestRaceCondition(t *testing.T) {
	var writeMutex sync.Mutex
	x := 0
	for i := 1; i <= 1000; i++ {
		go func() {
			for i := 1; i <= 100; i++ {
				writeMutex.Lock()
				x += 1
				writeMutex.Unlock()
			}
		}()
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Counter : ", x)
}

type BankAccount struct {
	RWMutex sync.RWMutex
	Balance int
}

func (account *BankAccount) AddBalance(amount int) {
	account.RWMutex.Lock() // write mutex
	account.Balance = account.Balance + amount
	account.RWMutex.Unlock()
}

func (account *BankAccount) GetBalance() int {
	account.RWMutex.RLock()
	balance := account.Balance
	account.RWMutex.RUnlock()
	return balance
}

func TestRWMutex(t *testing.T) {
	account := BankAccount{}

	for i := 0; i < 100; i++ {
		go func ()  {
			for j := 0; j < 100; j++{
				account.AddBalance(1)
				fmt.Println(account.GetBalance())				
			}	
		}()
	}

	time.Sleep(5 * time.Second)
	fmt.Println("Total Balance : ", account.GetBalance())
}

type UserBalance struct{
	Mutex sync.Mutex
	Name string
	Balance int
}

func (user *UserBalance) Lock(){
	user.Mutex.Lock()
}

func (user *UserBalance) Unlock(){
	user.Mutex.Unlock()
}

func (user *UserBalance) Change(amount int) {
	user.Balance = user.Balance + amount
}

func Transfer(user1 *UserBalance, user2 *UserBalance, amount int) {
	user1.Lock()
	fmt.Println("Lock -> ", user1.Name)
	user1.Change(-amount)

	time.Sleep(1 * time.Second)

	user2.Lock()
	fmt.Println("Lock -> ", user2.Name)
	user2.Change(amount)

	time.Sleep(1 * time.Second)

	user1.Unlock()
	user2.Unlock()
}

func TestDeadLock(t *testing.T){
	user1 := UserBalance{
		Name: "Fauzan",
		Balance: 10000,
	}
	user2 := UserBalance{
		Name: "Susi",
		Balance: 10000,
	}

	go Transfer(&user1, &user2, 1000)
	go Transfer(&user2, &user1, 1000)

	time.Sleep(5 * time.Second)
	fmt.Println(&user1)
	fmt.Println(&user2)
}