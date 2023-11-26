package mocks

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	mm_service "github.com/dir01/parcels/service"
	"github.com/gojuno/minimock/v3"
)

// PostalAPIMock implements service.PostalAPI
type PostalAPIMock struct {
	t minimock.Tester

	funcFetch          func(ctx context.Context, trackingNumber string) (p1 mm_service.PostalApiResponse)
	inspectFuncFetch   func(ctx context.Context, trackingNumber string)
	afterFetchCounter  uint64
	beforeFetchCounter uint64
	FetchMock          mPostalAPIMockFetch

	funcParse          func(rawResponse mm_service.PostalApiResponse) (tp1 *mm_service.TrackingInfo, err error)
	inspectFuncParse   func(rawResponse mm_service.PostalApiResponse)
	afterParseCounter  uint64
	beforeParseCounter uint64
	ParseMock          mPostalAPIMockParse
}

// NewPostalAPIMock returns a mock for service.PostalAPI
func NewPostalAPIMock(t minimock.Tester) *PostalAPIMock {
	m := &PostalAPIMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FetchMock = mPostalAPIMockFetch{mock: m}
	m.FetchMock.callArgs = []*PostalAPIMockFetchParams{}

	m.ParseMock = mPostalAPIMockParse{mock: m}
	m.ParseMock.callArgs = []*PostalAPIMockParseParams{}

	return m
}

type mPostalAPIMockFetch struct {
	mock               *PostalAPIMock
	defaultExpectation *PostalAPIMockFetchExpectation
	expectations       []*PostalAPIMockFetchExpectation

	callArgs []*PostalAPIMockFetchParams
	mutex    sync.RWMutex
}

// PostalAPIMockFetchExpectation specifies expectation struct of the PostalAPI.Fetch
type PostalAPIMockFetchExpectation struct {
	mock    *PostalAPIMock
	params  *PostalAPIMockFetchParams
	results *PostalAPIMockFetchResults
	Counter uint64
}

// PostalAPIMockFetchParams contains parameters of the PostalAPI.Fetch
type PostalAPIMockFetchParams struct {
	ctx            context.Context
	trackingNumber string
}

// PostalAPIMockFetchResults contains results of the PostalAPI.Fetch
type PostalAPIMockFetchResults struct {
	p1 mm_service.PostalApiResponse
}

// Expect sets up expected params for PostalAPI.Fetch
func (mmFetch *mPostalAPIMockFetch) Expect(ctx context.Context, trackingNumber string) *mPostalAPIMockFetch {
	if mmFetch.mock.funcFetch != nil {
		mmFetch.mock.t.Fatalf("PostalAPIMock.Fetch mock is already set by Set")
	}

	if mmFetch.defaultExpectation == nil {
		mmFetch.defaultExpectation = &PostalAPIMockFetchExpectation{}
	}

	mmFetch.defaultExpectation.params = &PostalAPIMockFetchParams{ctx, trackingNumber}
	for _, e := range mmFetch.expectations {
		if minimock.Equal(e.params, mmFetch.defaultExpectation.params) {
			mmFetch.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmFetch.defaultExpectation.params)
		}
	}

	return mmFetch
}

// Inspect accepts an inspector function that has same arguments as the PostalAPI.Fetch
func (mmFetch *mPostalAPIMockFetch) Inspect(f func(ctx context.Context, trackingNumber string)) *mPostalAPIMockFetch {
	if mmFetch.mock.inspectFuncFetch != nil {
		mmFetch.mock.t.Fatalf("Inspect function is already set for PostalAPIMock.Fetch")
	}

	mmFetch.mock.inspectFuncFetch = f

	return mmFetch
}

// Return sets up results that will be returned by PostalAPI.Fetch
func (mmFetch *mPostalAPIMockFetch) Return(p1 mm_service.PostalApiResponse) *PostalAPIMock {
	if mmFetch.mock.funcFetch != nil {
		mmFetch.mock.t.Fatalf("PostalAPIMock.Fetch mock is already set by Set")
	}

	if mmFetch.defaultExpectation == nil {
		mmFetch.defaultExpectation = &PostalAPIMockFetchExpectation{mock: mmFetch.mock}
	}
	mmFetch.defaultExpectation.results = &PostalAPIMockFetchResults{p1}
	return mmFetch.mock
}

// Set uses given function f to mock the PostalAPI.Fetch method
func (mmFetch *mPostalAPIMockFetch) Set(f func(ctx context.Context, trackingNumber string) (p1 mm_service.PostalApiResponse)) *PostalAPIMock {
	if mmFetch.defaultExpectation != nil {
		mmFetch.mock.t.Fatalf("Default expectation is already set for the PostalAPI.Fetch method")
	}

	if len(mmFetch.expectations) > 0 {
		mmFetch.mock.t.Fatalf("Some expectations are already set for the PostalAPI.Fetch method")
	}

	mmFetch.mock.funcFetch = f
	return mmFetch.mock
}

// When sets expectation for the PostalAPI.Fetch which will trigger the result defined by the following
// Then helper
func (mmFetch *mPostalAPIMockFetch) When(ctx context.Context, trackingNumber string) *PostalAPIMockFetchExpectation {
	if mmFetch.mock.funcFetch != nil {
		mmFetch.mock.t.Fatalf("PostalAPIMock.Fetch mock is already set by Set")
	}

	expectation := &PostalAPIMockFetchExpectation{
		mock:   mmFetch.mock,
		params: &PostalAPIMockFetchParams{ctx, trackingNumber},
	}
	mmFetch.expectations = append(mmFetch.expectations, expectation)
	return expectation
}

// Then sets up PostalAPI.Fetch return parameters for the expectation previously defined by the When method
func (e *PostalAPIMockFetchExpectation) Then(p1 mm_service.PostalApiResponse) *PostalAPIMock {
	e.results = &PostalAPIMockFetchResults{p1}
	return e.mock
}

// Fetch implements service.PostalAPI
func (mmFetch *PostalAPIMock) Fetch(ctx context.Context, trackingNumber string) (p1 mm_service.PostalApiResponse) {
	mm_atomic.AddUint64(&mmFetch.beforeFetchCounter, 1)
	defer mm_atomic.AddUint64(&mmFetch.afterFetchCounter, 1)

	if mmFetch.inspectFuncFetch != nil {
		mmFetch.inspectFuncFetch(ctx, trackingNumber)
	}

	mm_params := &PostalAPIMockFetchParams{ctx, trackingNumber}

	// Record call args
	mmFetch.FetchMock.mutex.Lock()
	mmFetch.FetchMock.callArgs = append(mmFetch.FetchMock.callArgs, mm_params)
	mmFetch.FetchMock.mutex.Unlock()

	for _, e := range mmFetch.FetchMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.p1
		}
	}

	if mmFetch.FetchMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmFetch.FetchMock.defaultExpectation.Counter, 1)
		mm_want := mmFetch.FetchMock.defaultExpectation.params
		mm_got := PostalAPIMockFetchParams{ctx, trackingNumber}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmFetch.t.Errorf("PostalAPIMock.Fetch got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmFetch.FetchMock.defaultExpectation.results
		if mm_results == nil {
			mmFetch.t.Fatal("No results are set for the PostalAPIMock.Fetch")
		}
		return (*mm_results).p1
	}
	if mmFetch.funcFetch != nil {
		return mmFetch.funcFetch(ctx, trackingNumber)
	}
	mmFetch.t.Fatalf("Unexpected call to PostalAPIMock.Fetch. %v %v", ctx, trackingNumber)
	return
}

// FetchAfterCounter returns a count of finished PostalAPIMock.Fetch invocations
func (mmFetch *PostalAPIMock) FetchAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmFetch.afterFetchCounter)
}

// FetchBeforeCounter returns a count of PostalAPIMock.Fetch invocations
func (mmFetch *PostalAPIMock) FetchBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmFetch.beforeFetchCounter)
}

// Calls returns a list of arguments used in each call to PostalAPIMock.Fetch.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmFetch *mPostalAPIMockFetch) Calls() []*PostalAPIMockFetchParams {
	mmFetch.mutex.RLock()

	argCopy := make([]*PostalAPIMockFetchParams, len(mmFetch.callArgs))
	copy(argCopy, mmFetch.callArgs)

	mmFetch.mutex.RUnlock()

	return argCopy
}

// MinimockFetchDone returns true if the count of the Fetch invocations corresponds
// the number of defined expectations
func (m *PostalAPIMock) MinimockFetchDone() bool {
	for _, e := range m.FetchMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.FetchMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterFetchCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcFetch != nil && mm_atomic.LoadUint64(&m.afterFetchCounter) < 1 {
		return false
	}
	return true
}

// MinimockFetchInspect logs each unmet expectation
func (m *PostalAPIMock) MinimockFetchInspect() {
	for _, e := range m.FetchMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to PostalAPIMock.Fetch with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.FetchMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterFetchCounter) < 1 {
		if m.FetchMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to PostalAPIMock.Fetch")
		} else {
			m.t.Errorf("Expected call to PostalAPIMock.Fetch with params: %#v", *m.FetchMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcFetch != nil && mm_atomic.LoadUint64(&m.afterFetchCounter) < 1 {
		m.t.Error("Expected call to PostalAPIMock.Fetch")
	}
}

type mPostalAPIMockParse struct {
	mock               *PostalAPIMock
	defaultExpectation *PostalAPIMockParseExpectation
	expectations       []*PostalAPIMockParseExpectation

	callArgs []*PostalAPIMockParseParams
	mutex    sync.RWMutex
}

// PostalAPIMockParseExpectation specifies expectation struct of the PostalAPI.Parse
type PostalAPIMockParseExpectation struct {
	mock    *PostalAPIMock
	params  *PostalAPIMockParseParams
	results *PostalAPIMockParseResults
	Counter uint64
}

// PostalAPIMockParseParams contains parameters of the PostalAPI.Parse
type PostalAPIMockParseParams struct {
	rawResponse mm_service.PostalApiResponse
}

// PostalAPIMockParseResults contains results of the PostalAPI.Parse
type PostalAPIMockParseResults struct {
	tp1 *mm_service.TrackingInfo
	err error
}

// Expect sets up expected params for PostalAPI.Parse
func (mmParse *mPostalAPIMockParse) Expect(rawResponse mm_service.PostalApiResponse) *mPostalAPIMockParse {
	if mmParse.mock.funcParse != nil {
		mmParse.mock.t.Fatalf("PostalAPIMock.Parse mock is already set by Set")
	}

	if mmParse.defaultExpectation == nil {
		mmParse.defaultExpectation = &PostalAPIMockParseExpectation{}
	}

	mmParse.defaultExpectation.params = &PostalAPIMockParseParams{rawResponse}
	for _, e := range mmParse.expectations {
		if minimock.Equal(e.params, mmParse.defaultExpectation.params) {
			mmParse.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmParse.defaultExpectation.params)
		}
	}

	return mmParse
}

// Inspect accepts an inspector function that has same arguments as the PostalAPI.Parse
func (mmParse *mPostalAPIMockParse) Inspect(f func(rawResponse mm_service.PostalApiResponse)) *mPostalAPIMockParse {
	if mmParse.mock.inspectFuncParse != nil {
		mmParse.mock.t.Fatalf("Inspect function is already set for PostalAPIMock.Parse")
	}

	mmParse.mock.inspectFuncParse = f

	return mmParse
}

// Return sets up results that will be returned by PostalAPI.Parse
func (mmParse *mPostalAPIMockParse) Return(tp1 *mm_service.TrackingInfo, err error) *PostalAPIMock {
	if mmParse.mock.funcParse != nil {
		mmParse.mock.t.Fatalf("PostalAPIMock.Parse mock is already set by Set")
	}

	if mmParse.defaultExpectation == nil {
		mmParse.defaultExpectation = &PostalAPIMockParseExpectation{mock: mmParse.mock}
	}
	mmParse.defaultExpectation.results = &PostalAPIMockParseResults{tp1, err}
	return mmParse.mock
}

// Set uses given function f to mock the PostalAPI.Parse method
func (mmParse *mPostalAPIMockParse) Set(f func(rawResponse mm_service.PostalApiResponse) (tp1 *mm_service.TrackingInfo, err error)) *PostalAPIMock {
	if mmParse.defaultExpectation != nil {
		mmParse.mock.t.Fatalf("Default expectation is already set for the PostalAPI.Parse method")
	}

	if len(mmParse.expectations) > 0 {
		mmParse.mock.t.Fatalf("Some expectations are already set for the PostalAPI.Parse method")
	}

	mmParse.mock.funcParse = f
	return mmParse.mock
}

// When sets expectation for the PostalAPI.Parse which will trigger the result defined by the following
// Then helper
func (mmParse *mPostalAPIMockParse) When(rawResponse mm_service.PostalApiResponse) *PostalAPIMockParseExpectation {
	if mmParse.mock.funcParse != nil {
		mmParse.mock.t.Fatalf("PostalAPIMock.Parse mock is already set by Set")
	}

	expectation := &PostalAPIMockParseExpectation{
		mock:   mmParse.mock,
		params: &PostalAPIMockParseParams{rawResponse},
	}
	mmParse.expectations = append(mmParse.expectations, expectation)
	return expectation
}

// Then sets up PostalAPI.Parse return parameters for the expectation previously defined by the When method
func (e *PostalAPIMockParseExpectation) Then(tp1 *mm_service.TrackingInfo, err error) *PostalAPIMock {
	e.results = &PostalAPIMockParseResults{tp1, err}
	return e.mock
}

// Parse implements service.PostalAPI
func (mmParse *PostalAPIMock) Parse(rawResponse mm_service.PostalApiResponse) (tp1 *mm_service.TrackingInfo, err error) {
	mm_atomic.AddUint64(&mmParse.beforeParseCounter, 1)
	defer mm_atomic.AddUint64(&mmParse.afterParseCounter, 1)

	if mmParse.inspectFuncParse != nil {
		mmParse.inspectFuncParse(rawResponse)
	}

	mm_params := &PostalAPIMockParseParams{rawResponse}

	// Record call args
	mmParse.ParseMock.mutex.Lock()
	mmParse.ParseMock.callArgs = append(mmParse.ParseMock.callArgs, mm_params)
	mmParse.ParseMock.mutex.Unlock()

	for _, e := range mmParse.ParseMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.tp1, e.results.err
		}
	}

	if mmParse.ParseMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmParse.ParseMock.defaultExpectation.Counter, 1)
		mm_want := mmParse.ParseMock.defaultExpectation.params
		mm_got := PostalAPIMockParseParams{rawResponse}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmParse.t.Errorf("PostalAPIMock.Parse got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmParse.ParseMock.defaultExpectation.results
		if mm_results == nil {
			mmParse.t.Fatal("No results are set for the PostalAPIMock.Parse")
		}
		return (*mm_results).tp1, (*mm_results).err
	}
	if mmParse.funcParse != nil {
		return mmParse.funcParse(rawResponse)
	}
	mmParse.t.Fatalf("Unexpected call to PostalAPIMock.Parse. %v", rawResponse)
	return
}

// ParseAfterCounter returns a count of finished PostalAPIMock.Parse invocations
func (mmParse *PostalAPIMock) ParseAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmParse.afterParseCounter)
}

// ParseBeforeCounter returns a count of PostalAPIMock.Parse invocations
func (mmParse *PostalAPIMock) ParseBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmParse.beforeParseCounter)
}

// Calls returns a list of arguments used in each call to PostalAPIMock.Parse.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmParse *mPostalAPIMockParse) Calls() []*PostalAPIMockParseParams {
	mmParse.mutex.RLock()

	argCopy := make([]*PostalAPIMockParseParams, len(mmParse.callArgs))
	copy(argCopy, mmParse.callArgs)

	mmParse.mutex.RUnlock()

	return argCopy
}

// MinimockParseDone returns true if the count of the Parse invocations corresponds
// the number of defined expectations
func (m *PostalAPIMock) MinimockParseDone() bool {
	for _, e := range m.ParseMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ParseMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterParseCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcParse != nil && mm_atomic.LoadUint64(&m.afterParseCounter) < 1 {
		return false
	}
	return true
}

// MinimockParseInspect logs each unmet expectation
func (m *PostalAPIMock) MinimockParseInspect() {
	for _, e := range m.ParseMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to PostalAPIMock.Parse with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ParseMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterParseCounter) < 1 {
		if m.ParseMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to PostalAPIMock.Parse")
		} else {
			m.t.Errorf("Expected call to PostalAPIMock.Parse with params: %#v", *m.ParseMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcParse != nil && mm_atomic.LoadUint64(&m.afterParseCounter) < 1 {
		m.t.Error("Expected call to PostalAPIMock.Parse")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *PostalAPIMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockFetchInspect()

		m.MinimockParseInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *PostalAPIMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *PostalAPIMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockFetchDone() &&
		m.MinimockParseDone()
}
