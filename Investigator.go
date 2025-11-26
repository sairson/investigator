package Investigator

import (
	"Investigator/helper/llm/vars"
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
)

func Investigate(ctx context.Context, state ModelGraphState) (ModelStateOutput, error) {
	workflow := compose.NewWorkflow[ModelGraphState, ModelStateOutput]()
	workflow.AddLambdaNode("topic_research", compose.InvokableLambda(
		func(ctx context.Context, input ModelGraphState) (output vars.LLMResponse, err error) {
			return MediaResearch(ctx, input)
		},
	)).AddInput(compose.START)
	workflow.AddLambdaNode("formatting", compose.InvokableLambda(
		func(ctx context.Context, input vars.LLMResponse) (output ModelStateOutput, err error) {
			return InputFormatting(ctx, input)
		},
	)).AddInput("topic_research")
	workflow.End().AddInput("formatting")
	run, err := workflow.Compile(ctx)
	if err != nil {
		return ModelStateOutput{}, fmt.Errorf("investigator workflow compile failed: %w", err)
	}
	result, err := run.Invoke(ctx, state)
	if err != nil {
		return ModelStateOutput{}, fmt.Errorf("investigator invoke failed: %w", err)
	}
	return ModelStateOutput{SEO: result.SEO, Analysis: result.Analysis, Creation: result.Creation, Response: result.Response}, nil
}
