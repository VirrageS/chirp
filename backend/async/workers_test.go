package async

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Workers", func() {
	var (
		pool *WorkerPool

		taskTests []int64
	)

	BeforeEach(func() {
		pool = NewWorkerPool(func(task Task) *Result {
			id := task.(int64)
			return &Result{
				Value: id + 1,
				Error: nil,
			}
		})

		taskTests = []int64{12, 14, 200, -1, 0}
	})

	AfterEach(func() {
		pool.Close()
	})

	It("should post tasks", func() {
		for _, task := range taskTests {
			pool.PostTask(task)
		}
	})

	It("should get result", func() {
		for _, task := range taskTests {
			pool.PostTask(task)
		}

		for range taskTests {
			result := pool.GetResult()
			Expect(result.Error).NotTo(HaveOccurred())
			Expect(taskTests).To(ContainElement(result.Value.(int64) - 1)) // it is easier that way...
		}
	})

	It("should block when result is not ready", func() {
		started := make(chan bool, 1)
		done := make(chan bool, 1)

		go func() {
			defer GinkgoRecover()

			started <- true
			result := pool.GetResult()
			Expect(result.Error).NotTo(HaveOccurred())
			done <- true
		}()
		<-started

		select {
		case <-done:
			Fail("GetResult is not blocked")
		case <-time.After(time.Millisecond * 500):
			// ok
		}

		pool.PostTask(taskTests[0])

		select {
		case <-done:
			// ok
		case <-time.After(time.Second):
			Fail("GetResult is not unblocked")
		}
	})
})
