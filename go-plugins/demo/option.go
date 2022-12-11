package demo

type HandleOption func(*Handler)

func WithDemoPath(demoPath string) HandleOption {
	return func(h *Handler) {
		if len(demoPath) > 0 {
			h.demoPath = demoPath
		}
	}
}
