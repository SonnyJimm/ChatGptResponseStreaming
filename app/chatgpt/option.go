package chatgpt

type Option func(c *ChatGptClient)

func WithModel(model string) Option {
	return func(c *ChatGptClient) {
		c.Model = model
	}
}
func WithPrompts(prompts []ChatGptPrompt) Option {
	return func(c *ChatGptClient) {
		c.Messages = prompts
	}
}
func WithPrompt(prompt ChatGptPrompt) Option {
	return func(c *ChatGptClient) {
		c.Messages = append(c.Messages, prompt)
	}
}
