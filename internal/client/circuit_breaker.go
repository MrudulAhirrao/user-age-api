package client

import (
	"errors"
	"log"
	"sync"
	"time"
)

type State string

const(
	StateClosed State = "CLOSED"
	StateOpen   State = "OPEN"
	StateHalfOpen State = "HALF_OPEN"
)

type CircuitBreaker struct{
	mu sync.Mutex
	state State
	failureCount int
	failureThreshold int
	timeout time.Duration
	lastFailureTime time.Time
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker{
	return &CircuitBreaker{
		state: StateClosed,
		failureThreshold: threshold,
		timeout: timeout,
	}
}	

func (cb *CircuitBreaker) Execute(action func() error) error{
	cb.mu.Lock()

	if cb.state == StateOpen{
		if time.Since(cb.lastFailureTime) > cb.timeout{
			log.Println("Circuit HalfOpen: Testing Connection...")
			cb.state = StateHalfOpen
		} else{
			cb.mu.Unlock()
			return errors.New("circuit breaker is open: email service unavailable")
		}
	}else if cb.state == StateHalfOpen{
		
	}
	cb.mu.Unlock()
	err:= action()
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if err != nil{
		cb.failureCount++
		cb.lastFailureTime = time.Now()
		if cb.state == StateHalfOpen{
			log.Println("Circuit ReOpened: Test Failed ")
			cb.state = StateOpen
			cb.failureCount = 0
		}else if cb.failureCount >= cb.failureThreshold{
			if cb.state == StateClosed{
				log.Printf("Circuit Tripped: %d Failures\n", cb.failureCount)
				cb.state = StateOpen
			}
		}
		return err
	}
	if cb.state == StateHalfOpen{
		log.Println("Circuit Closed: Service Restarted")
		cb.state = StateClosed
		cb.failureCount = 0
	}else if cb.failureCount > 0{
		cb.failureCount = 0
	}
	return nil
}