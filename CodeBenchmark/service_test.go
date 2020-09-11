import (
	"context"
	"github.com/golang/mock/gomock"
	"service/mock_service"
	"testing"
)


// this test results in 
// === RUN   TestThatCausesDeadlock
// fatal error: all goroutines are asleep - deadlock!
func TestThatCausesDeadlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mock_service.NewMockFoo(ctrl)

        //We have neglected to set up the expected calls in this test
	new(Service).Baz(context.Background(), mock)

	ctrl.Finish()
}

// This test passes as expected 
func TestThatWorksFine(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mock_playgroundf.NewMockFoo(ctrl)

	mock.EXPECT().Bar(gomock.Any(), gomock.Any())

	new(Service).Baz(context.Background(), mock)

	ctrl.Finish()
}

