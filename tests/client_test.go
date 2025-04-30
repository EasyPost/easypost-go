package easypost_test

import (
	"context"
	"github.com/EasyPost/easypost-go/v5"
	"github.com/google/uuid"
)

func (c *ClientTests) TestCreateClient() {
	require := c.Require()

	client := easypost.New("")

	_, err := client.ListAddresses(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.Error(err) // Throw an error when no API key present
}

func (c *ClientTests) TestHooks() {
	assert, require := c.Assert(), c.Require()

	client := c.TestClient()

	requestHookCalled := false
	responseHookCalled := false

	var pairId uuid.UUID

	client.Hooks.AddRequestEventSubscriber(
		easypost.RequestHookEventSubscriber{
			Callback: func(ctx context.Context, event easypost.RequestHookEvent) error {
				assert.Equal("GET", event.Method)
				assert.Regexp("^https://api.easypost.com/v2/addresses", event.Url.String())
				pairId = event.Id
				requestHookCalled = true
				return nil
			},
			HookEventSubscriber: easypost.HookEventSubscriber{
				ID: "hook_1",
			},
		})
	client.Hooks.AddResponseEventSubscriber(
		easypost.ResponseHookEventSubscriber{
			Callback: func(ctx context.Context, event easypost.ResponseHookEvent) error {
				assert.Equal(200, event.HttpStatus)
				assert.Equal("GET", event.Method)
				assert.Equal(pairId, event.Id)
				responseHookCalled = true
				return nil
			},
			HookEventSubscriber: easypost.HookEventSubscriber{
				ID: "hook_2",
			},
		})

	// verify that the hooks were added
	assert.Equal(1, len(client.Hooks.RequestHookEventSubscriptions))
	assert.Equal(1, len(client.Hooks.ResponseHookEventSubscriptions))

	// fire off a request
	_, err := client.ListAddresses(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	// verify that the hooks were called
	assert.True(requestHookCalled)
	assert.True(responseHookCalled)
}

func (c *ClientTests) TestMultipleHooks() {
	assert, require := c.Assert(), c.Require()

	client := c.TestClient()

	requestHook1Called := false
	requestHook2Called := false

	client.Hooks.AddRequestEventSubscriber(
		easypost.RequestHookEventSubscriber{
			Callback: func(ctx context.Context, event easypost.RequestHookEvent) error {
				requestHook1Called = true
				return nil
			},
			HookEventSubscriber: easypost.HookEventSubscriber{
				ID: "hook_1",
			},
		})
	client.Hooks.AddRequestEventSubscriber(
		easypost.RequestHookEventSubscriber{
			Callback: func(ctx context.Context, event easypost.RequestHookEvent) error {
				requestHook2Called = true
				return nil
			},
			HookEventSubscriber: easypost.HookEventSubscriber{
				ID: "hook_2",
			},
		})

	// verify that the hooks were added
	assert.Equal(2, len(client.Hooks.RequestHookEventSubscriptions))

	// fire off a request
	_, err := client.ListAddresses(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	// verify that the hooks were called
	assert.True(requestHook1Called)
	assert.True(requestHook2Called)
}

func (c *ClientTests) TestHooksUnsubscribing() {
	assert, require := c.Assert(), c.Require()

	client := c.TestClient()

	requestCallbackCallCount := 0
	hookName := "hook_1"

	subscriber := easypost.RequestHookEventSubscriber{
		Callback: func(ctx context.Context, event easypost.RequestHookEvent) error {
			requestCallbackCallCount++
			return nil
		},
		HookEventSubscriber: easypost.HookEventSubscriber{
			ID: hookName,
		},
	}

	client.Hooks.AddRequestEventSubscriber(subscriber)

	// verify that the hook was added
	assert.Equal(1, len(client.Hooks.RequestHookEventSubscriptions))

	// fire off a request
	_, err := client.ListAddresses(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	// verify that the hook was called
	assert.Equal(1, requestCallbackCallCount)

	// unsubscribe the hook
	client.Hooks.RemoveRequestEventSubscriber(subscriber)

	// fire off a request
	// fire off a request
	_, err = client.ListAddresses(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	// verify that the hook was not called again
	assert.Equal(1, requestCallbackCallCount)
}
