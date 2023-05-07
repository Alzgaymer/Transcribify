Add to routes.GetVideoTranscription request to OPENAI

```go
	c := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	prompt, err := formatPrompt(video)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		route.logger.Info("Error while formatting video to OPENAI prompt", zap.Error(err))
		return
	}
	req := openai.CompletionRequest{
		Model:     openai.GPT3Ada,
		MaxTokens: 2,
		Prompt:    prompt,
	}
	resp, err := c.CreateCompletion(ctx, req) // HTTP 400 model`s max tokens 2048 in prompt  ~11`000
	if err != nil {
		route.logger.Info("Error while sending request to OPENAI", zap.Error(err))
		return
	}
	route.logger.Info("OPENAI response", zap.Any("resp", resp))
```