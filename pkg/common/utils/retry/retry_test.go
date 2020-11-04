package retry

import (
	"github.com/juju/errgo"
	adTest "gopkg.in/check.v1"
	"testing"
	"time"
)

var customError = errgo.New("This is a custom error")
var customError3 = errgo.New("This is a custom error3")

type retryTestSuite struct {
}

var _ = adTest.Suite(&retryTestSuite{})

func Test(t *testing.T) {
	adTest.TestingT(t)
}

func (p *retryTestSuite) SetUpSuite(c *adTest.C) {

}

func (p *retryTestSuite) TearDownSuite(c *adTest.C) {
}

func retryFunction(retryOptions ...RetryOption) error {
	op := func() error {
		return nil
	}
	if len(retryOptions) == 0 {
		//not input retry options from out of this api
		return Do(op, Timeout(0), MaxTries(1), RetryChecker(errgo.Any))
	}
	return Do(op, retryOptions...)
}

func retryFunction2(retryOptions ...RetryOption) error {
	op := func() error {
		return customError
	}
	if len(retryOptions) == 0 {
		//not input retry options from out of this api
		return Do(op, Timeout(0), MaxTries(1), RetryChecker(errgo.Any))
	}
	return Do(op, retryOptions...)
}

func retryFunction3(retryOptions ...RetryOption) error {
	op := func() error {
		return customError3
	}
	if len(retryOptions) == 0 {
		//not input retry options from out of this api
		return Do(op, Timeout(0), MaxTries(1), RetryChecker(errgo.Any))
	}
	return Do(op, retryOptions...)
}

func (p *retryTestSuite) TestRetry(c *adTest.C) {
	err := retryFunction(Timeout(15*time.Second), MaxTries(10), Sleep(10))
	c.Assert(err, adTest.IsNil)
}

func (p *retryTestSuite) TestRetry2(c *adTest.C) {
	err := retryFunction2(Timeout(15*time.Second), MaxTries(10), Sleep(10))
	c.Assert(err, adTest.NotNil)
}

func (p *retryTestSuite) TestRetry3(c *adTest.C) {
	err := retryFunction3(Timeout(15*time.Second), MaxTries(10), Sleep(10), RetryChecker(errgo.Is(customError3)))
	c.Assert(err, adTest.NotNil)
}
