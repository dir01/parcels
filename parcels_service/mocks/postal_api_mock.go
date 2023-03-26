package mocks

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	mm_parcels_service "github.com/dir01/parcels/parcels_service"
	"github.com/gojuno/minimock/v3"
)

// PostalApiMock implements parcels_service.PostalApi
type PostalApiMock struct {
	t minimock.Tester

	funcFetch          func(ctx context.Context, trackingNumber string) (p1 mm_parcels_service.PostalApiResponse)
	inspectFuncFetch   func(ctx context.Context, trackingNumber string)
	afterFetchCounter  uint64
	beforeFetchCounter uint64
	FetchMock          mPostalApiMockFetch

	funcParse          func(rawResponse mm_parcels_service.PostalApiResponse) (tp1 *mm_parcels_service.TrackingInfo, err error)
	inspectFuncParse   func(rawResponse mm_parcels_service.PostalApiResponse)
	afterParseCounter  uint64
	beforeParseCounter uint64
	ParseMock          mPostalApiMockParse
}

// NewPostalApiMock returns a mock for parcels_service.PostalApi
func NewPostalApiMock(t minimock.Tester) *PostalApiMock {
	m := &PostalApiMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FetchMock = mPostalApiMockFetch{mock: m}
	m.FetchMock.callArgs = []*PostalApiMockFetchParams{}

	m.ParseMock = mPostalApiMockParse{mock: m}
	m.ParseMock.callArgs = []*PostalApiMockParseParams{}

	return m
}

type mPostalApiMockFetch struct {
	mock               *PostalApiMock
	defaultExpectation *PostalApiMockFetchExpectation
	expectations       []*PostalApiMockFetchExpectation

	callArgs []*PostalApiMockFetchParams
	mutex    sync.RWMutex
}

// PostalApiMockFetchExpectation specifies expectation struct of the PostalApi.Fetch
type PostalApiMockFetchExpectation struct {
	mock    *PostalApiMock
	params  *PostalApiMockFetchParams
	results *PostalApiMockFetchResults
	Counter uint64
}

// PostalApiMockFetchParams contains parameters of the PostalApi.Fetch
type PostalApiMockFetchParams struct {
	ctx            context.Context
	trackingNumber string
}

// PostalApiMockFetchResults contains results of the PostalApi.Fetch
type PostalApiMockFetchResults struct {
	p1 mm_parcels_service.PostalApiResponse
}

// Expect sets up expected params for PostalApi.Fetch
func (mmFetch *mPostalApiMockFetch) Expect(ctx context.Context, trackingNumber string) *mPostalApiMockFetch {
	if mmFetch.mock.funcFetch != nil {
		mmFetch.mock.t.Fatalf("PostalApiMock.Fetch mock is already set by Set")
	}

	if mmFetch.defaultExpectation == nil {
		mmFetch.defaultExpectation = &PostalApiMockFetchExpectation{}
	}

	mmFetch.defaultExpectation.params = &PostalApiMockFetchParams{ctx, trackingNumber}
	for _, e := range mmFetch.expectations {
		if minimock.Equal(e.params, mmFetch.defaultExpectation.params) {
			mmFetch.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmFetch.defaultExpectation.params)
		}
	}

	return mmFetch
}

// Inspect accepts an inspector function that has same arguments as the PostalApi.Fetch
func (mmFetch *mPostalApiMockFetch) Inspect(f func(ctx context.Context, trackingNumber string)) *mPostalApiMockFetch {
	if mmFetch.mock.inspectFuncFetch != nil {
		mmFetch.mock.t.Fatalf("Inspect function is already set for PostalApiMock.Fetch")
	}

	mmFetch.mock.inspectFuncFetch = f

	return mmFetch
}

// Return sets up results that will be returned by PostalApi.Fetch
func (mmFetch *mPostalApiMockFetch) Return(p1 mm_parcels_service.PostalApiResponse) *PostalApiMock {
	if mmFetch.mock.funcFetch != nil {
		mmFetch.mock.t.Fatalf("PostalApiMock.Fetch mock is already set by Set")
	}

	if mmFetch.defaultExpectation == nil {
		mmFetch.defaultExpectation = &PostalApiMockFetchExpectation{mock: mmFetch.mock}
	}
	mmFetch.defaultExpectation.results = &PostalApiMockFetchResults{p1}
	return mmFetch.mock
}

// Set uses given function f to mock the PostalApi.Fetch method
func (mmFetch *mPostalApiMockFetch) Set(f func(ctx context.Context, trackingNumber string) (p1 mm_parcels_service.PostalApiResponse)) *PostalApiMock {
	if mmFetch.defaultExpectation != nil {
		mmFetch.mock.t.Fatalf("Default expectation is already set for the PostalApi.Fetch method")
	}

	if len(mmFetch.expectations) > 0 {
		mmFetch.mock.t.Fatalf("Some expectations are already set for the PostalApi.Fetch method")
	}

	mmFetch.mock.funcFetch = f
	return mmFetch.mock
}

// When sets expectation for the PostalApi.Fetch which will trigger the result defined by the following
// Then helper
func (mmFetch *mPostalApiMockFetch) When(ctx context.Context, trackingNumber string) *PostalApiMockFetchExpectation {
	if mmFetch.mock.funcFetch != nil {
		mmFetch.mock.t.Fatalf("PostalApiMock.Fetch mock is already set by Set")
	}

	expectation := &PostalApiMockFetchExpectation{
		mock:   mmFetch.mock,
		params: &PostalApiMockFetchParams{ctx, trackingNumber},
	}
	mmFetch.expectations = append(mmFetch.expectations, expectation)
	return expectation
}

// Then sets up PostalApi.Fetch return parameters for the expectation previously defined by the When method
func (e *PostalApiMockFetchExpectation) Then(p1 mm_parcels_service.PostalApiResponse) *PostalApiMock {
	e.results = &PostalApiMockFetchResults{p1}
	return e.mock
}

// Fetch implements parcels_service.PostalApi
func (mmFetch *PostalApiMock) Fetch(ctx context.Context, trackingNumber string) (p1 mm_parcels_service.PostalApiResponse) {
	mm_atomic.AddUint64(&mmFetch.beforeFetchCounter, 1)
	defer mm_atomic.AddUint64(&mmFetch.afterFetchCounter, 1)

	if mmFetch.inspectFuncFetch != nil {
		mmFetch.inspectFuncFetch(ctx, trackingNumber)
	}

	mm_params := &PostalApiMockFetchParams{ctx, trackingNumber}

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
		mm_got := PostalApiMockFetchParams{ctx, trackingNumber}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmFetch.t.Errorf("PostalApiMock.Fetch got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmFetch.FetchMock.defaultExpectation.results
		if mm_results == nil {
			mmFetch.t.Fatal("No results are set for the PostalApiMock.Fetch")
		}
		return (*mm_results).p1
	}
	if mmFetch.funcFetch != nil {
		return mmFetch.funcFetch(ctx, trackingNumber)
	}
	mmFetch.t.Fatalf("Unexpected call to PostalApiMock.Fetch. %v %v", ctx, trackingNumber)
	return
}

// FetchAfterCounter returns a count of finished PostalApiMock.Fetch invocations
func (mmFetch *PostalApiMock) FetchAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmFetch.afterFetchCounter)
}

// FetchBeforeCounter returns a count of PostalApiMock.Fetch invocations
func (mmFetch *PostalApiMock) FetchBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmFetch.beforeFetchCounter)
}

// Calls returns a list of arguments used in each call to PostalApiMock.Fetch.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmFetch *mPostalApiMockFetch) Calls() []*PostalApiMockFetchParams {
	mmFetch.mutex.RLock()

	argCopy := make([]*PostalApiMockFetchParams, len(mmFetch.callArgs))
	copy(argCopy, mmFetch.callArgs)

	mmFetch.mutex.RUnlock()

	return argCopy
}

// MinimockFetchDone returns true if the count of the Fetch invocations corresponds
// the number of defined expectations
func (m *PostalApiMock) MinimockFetchDone() bool {
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
func (m *PostalApiMock) MinimockFetchInspect() {
	for _, e := range m.FetchMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to PostalApiMock.Fetch with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.FetchMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterFetchCounter) < 1 {
		if m.FetchMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to PostalApiMock.Fetch")
		} else {
			m.t.Errorf("Expected call to PostalApiMock.Fetch with params: %#v", *m.FetchMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcFetch != nil && mm_atomic.LoadUint64(&m.afterFetchCounter) < 1 {
		m.t.Error("Expected call to PostalApiMock.Fetch")
	}
}

type mPostalApiMockParse struct {
	mock               *PostalApiMock
	defaultExpectation *PostalApiMockParseExpectation
	expectations       []*PostalApiMockParseExpectation

	callArgs []*PostalApiMockParseParams
	mutex    sync.RWMutex
}

// PostalApiMockParseExpectation specifies expectation struct of the PostalApi.Parse
type PostalApiMockParseExpectation struct {
	mock    *PostalApiMock
	params  *PostalApiMockParseParams
	results *PostalApiMockParseResults
	Counter uint64
}

// PostalApiMockParseParams contains parameters of the PostalApi.Parse
type PostalApiMockParseParams struct {
	rawResponse mm_parcels_service.PostalApiResponse
}

// PostalApiMockParseResults contains results of the PostalApi.Parse
type PostalApiMockParseResults struct {
	tp1 *mm_parcels_service.TrackingInfo
	err error
}

// Expect sets up expected params for PostalApi.Parse
func (mmParse *mPostalApiMockParse) Expect(rawResponse mm_parcels_service.PostalApiResponse) *mPostalApiMockParse {
	if mmParse.mock.funcParse != nil {
		mmParse.mock.t.Fatalf("PostalApiMock.Parse mock is already set by Set")
	}

	if mmParse.defaultExpectation == nil {
		mmParse.defaultExpectation = &PostalApiMockParseExpectation{}
	}

	mmParse.defaultExpectation.params = &PostalApiMockParseParams{rawResponse}
	for _, e := range mmParse.expectations {
		if minimock.Equal(e.params, mmParse.defaultExpectation.params) {
			mmParse.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmParse.defaultExpectation.params)
		}
	}

	return mmParse
}

// Inspect accepts an inspector function that has same arguments as the PostalApi.Parse
func (mmParse *mPostalApiMockParse) Inspect(f func(rawResponse mm_parcels_service.PostalApiResponse)) *mPostalApiMockParse {
	if mmParse.mock.inspectFuncParse != nil {
		mmParse.mock.t.Fatalf("Inspect function is already set for PostalApiMock.Parse")
	}

	mmParse.mock.inspectFuncParse = f

	return mmParse
}

// Return sets up results that will be returned by PostalApi.Parse
func (mmParse *mPostalApiMockParse) Return(tp1 *mm_parcels_service.TrackingInfo, err error) *PostalApiMock {
	if mmParse.mock.funcParse != nil {
		mmParse.mock.t.Fatalf("PostalApiMock.Parse mock is already set by Set")
	}

	if mmParse.defaultExpectation == nil {
		mmParse.defaultExpectation = &PostalApiMockParseExpectation{mock: mmParse.mock}
	}
	mmParse.defaultExpectation.results = &PostalApiMockParseResults{tp1, err}
	return mmParse.mock
}

// Set uses given function f to mock the PostalApi.Parse method
func (mmParse *mPostalApiMockParse) Set(f func(rawResponse mm_parcels_service.PostalApiResponse) (tp1 *mm_parcels_service.TrackingInfo, err error)) *PostalApiMock {
	if mmParse.defaultExpectation != nil {
		mmParse.mock.t.Fatalf("Default expectation is already set for the PostalApi.Parse method")
	}

	if len(mmParse.expectations) > 0 {
		mmParse.mock.t.Fatalf("Some expectations are already set for the PostalApi.Parse method")
	}

	mmParse.mock.funcParse = f
	return mmParse.mock
}

// When sets expectation for the PostalApi.Parse which will trigger the result defined by the following
// Then helper
func (mmParse *mPostalApiMockParse) When(rawResponse mm_parcels_service.PostalApiResponse) *PostalApiMockParseExpectation {
	if mmParse.mock.funcParse != nil {
		mmParse.mock.t.Fatalf("PostalApiMock.Parse mock is already set by Set")
	}

	expectation := &PostalApiMockParseExpectation{
		mock:   mmParse.mock,
		params: &PostalApiMockParseParams{rawResponse},
	}
	mmParse.expectations = append(mmParse.expectations, expectation)
	return expectation
}

// Then sets up PostalApi.Parse return parameters for the expectation previously defined by the When method
func (e *PostalApiMockParseExpectation) Then(tp1 *mm_parcels_service.TrackingInfo, err error) *PostalApiMock {
	e.results = &PostalApiMockParseResults{tp1, err}
	return e.mock
}

// Parse implements parcels_service.PostalApi
func (mmParse *PostalApiMock) Parse(rawResponse mm_parcels_service.PostalApiResponse) (tp1 *mm_parcels_service.TrackingInfo, err error) {
	mm_atomic.AddUint64(&mmParse.beforeParseCounter, 1)
	defer mm_atomic.AddUint64(&mmParse.afterParseCounter, 1)

	if mmParse.inspectFuncParse != nil {
		mmParse.inspectFuncParse(rawResponse)
	}

	mm_params := &PostalApiMockParseParams{rawResponse}

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
		mm_got := PostalApiMockParseParams{rawResponse}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmParse.t.Errorf("PostalApiMock.Parse got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmParse.ParseMock.defaultExpectation.results
		if mm_results == nil {
			mmParse.t.Fatal("No results are set for the PostalApiMock.Parse")
		}
		return (*mm_results).tp1, (*mm_results).err
	}
	if mmParse.funcParse != nil {
		return mmParse.funcParse(rawResponse)
	}
	mmParse.t.Fatalf("Unexpected call to PostalApiMock.Parse. %v", rawResponse)
	return
}

// ParseAfterCounter returns a count of finished PostalApiMock.Parse invocations
func (mmParse *PostalApiMock) ParseAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmParse.afterParseCounter)
}

// ParseBeforeCounter returns a count of PostalApiMock.Parse invocations
func (mmParse *PostalApiMock) ParseBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmParse.beforeParseCounter)
}

// Calls returns a list of arguments used in each call to PostalApiMock.Parse.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmParse *mPostalApiMockParse) Calls() []*PostalApiMockParseParams {
	mmParse.mutex.RLock()

	argCopy := make([]*PostalApiMockParseParams, len(mmParse.callArgs))
	copy(argCopy, mmParse.callArgs)

	mmParse.mutex.RUnlock()

	return argCopy
}

// MinimockParseDone returns true if the count of the Parse invocations corresponds
// the number of defined expectations
func (m *PostalApiMock) MinimockParseDone() bool {
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
func (m *PostalApiMock) MinimockParseInspect() {
	for _, e := range m.ParseMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to PostalApiMock.Parse with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ParseMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterParseCounter) < 1 {
		if m.ParseMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to PostalApiMock.Parse")
		} else {
			m.t.Errorf("Expected call to PostalApiMock.Parse with params: %#v", *m.ParseMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcParse != nil && mm_atomic.LoadUint64(&m.afterParseCounter) < 1 {
		m.t.Error("Expected call to PostalApiMock.Parse")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *PostalApiMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockFetchInspect()

		m.MinimockParseInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *PostalApiMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *PostalApiMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockFetchDone() &&
		m.MinimockParseDone()
}
