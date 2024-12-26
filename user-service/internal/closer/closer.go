package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// globalCloser - глобальный экземпляр Closer с сигналами завершения по умолчанию.
var globalCloser = New(syscall.SIGINT, syscall.SIGTERM)

// Add добавляет функции завершения в глобальный Closer.
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// Wait блокирует выполнение до получения сигнала завершения.
func Wait() {
	globalCloser.Wait()
}

// CloseAll выполняет все добавленные функции завершения.
func CloseAll() {
	globalCloser.CloseAll()
}

// Closer - структура для управления функциями завершения.
// Обеспечивает потокобезопасный доступ к функциям, которые необходимо выполнить перед завершением работы приложения.
type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

// New создает новый экземпляр Closer и устанавливает обработчик сигналов завершения.
// При получении сигнала вызывается метод CloseAll.
func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}
	return c
}

// Add добавляет функции завершения в Closer.
// Функции будут выполнены при получении сигнала завершения.
func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

// Wait блокирует выполнение до завершения работы всех функций завершения.
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll выполняет все добавленные функции завершения и уведомляет о завершении работы.
func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))
		log.Println("closing funcs:", len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Println("error returned from Closer")
			}
		}

		log.Println("all closed")
	})
}
