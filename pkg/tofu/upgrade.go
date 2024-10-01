package tofu

import "context"

// Upgrade runs a tofu apply, just like Install()
func (m *Mixin) Upgrade(ctx context.Context) error {
	return m.Install(ctx)
}
