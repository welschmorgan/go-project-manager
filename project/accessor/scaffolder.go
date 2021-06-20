package accessor

type Scaffolder interface {
	Name() string

	Scaffold(ctx *FinalizationContext) error
}
