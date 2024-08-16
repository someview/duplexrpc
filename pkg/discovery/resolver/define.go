package resolver

import "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/pkg/discovery"

type Resolver interface {
	/// <summary>
	/// Starts listening to resolver for results with the specified callback. Can only be called once.
	/// <para>
	/// The <see cref="ResolverResult"/> passed to the callback has addresses when successful,
	/// otherwise a <see cref="Status"/> details the resolution error.
	/// </para>
	/// </summary>
	/// <param name="listener">The callback used to receive updates on the target.</param>
	Start(discovery.ServiceListener)
	/// <summary>
	/// Refresh resolution. Can only be called after <see cref="Start(Action{ResolverResult})"/>.
	/// The default implementation is no-op.
	/// <para>
	/// This is only a hint. Implementation takes it as a signal but may not start resolution.
	/// </para>
	/// </summary>
	Refresh()
	Stop()
}
