package app

type PostProcessor interface {
	BeforeInitialization(factory interface{})
	AfterInitialization(factory interface{})
}

type postProcessor struct{
}

var (
	postProcessors []PostProcessor
)

func init() {

}

func RegisterPostProcessor(p ...PostProcessor)  {
	postProcessors = append(postProcessors, p...)
}

func (p *postProcessor) BeforeInitialization(factory interface{})  {
	for _, processor := range postProcessors {
		processor.BeforeInitialization(factory)
	}
}

func (p *postProcessor) AfterInitialization(factory interface{})  {
	for _, processor := range postProcessors {
		processor.AfterInitialization(factory)
	}
}