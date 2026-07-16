package config

import "testing"

func TestLoadUsesRoleBasedModelsAndLegacyAliases(t *testing.T) {
	t.Setenv("QWEN_API_KEY", "test-key")
	t.Setenv("QWEN_EXTRACTION_MODEL", "role-extraction")
	t.Setenv("QWEN_REASONING_MODEL", "role-reasoning")
	t.Setenv("QWEN_TURBO_MODEL", "legacy-turbo")
	t.Setenv("QWEN_MAX_MODEL", "legacy-max")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.QwenExtractionModel != "role-extraction" || cfg.QwenReasoningModel != "role-reasoning" {
		t.Fatalf("role models = %q, %q", cfg.QwenExtractionModel, cfg.QwenReasoningModel)
	}
	if cfg.QwenTurboModel != "role-extraction" || cfg.QwenMaxModel != "role-reasoning" {
		t.Fatalf("compatibility aliases = %q, %q", cfg.QwenTurboModel, cfg.QwenMaxModel)
	}
}

func TestLoadFallsBackToLegacyModelVariables(t *testing.T) {
	t.Setenv("QWEN_API_KEY", "test-key")
	t.Setenv("QWEN_EXTRACTION_MODEL", "")
	t.Setenv("QWEN_REASONING_MODEL", "")
	t.Setenv("QWEN_TURBO_MODEL", "legacy-turbo")
	t.Setenv("QWEN_MAX_MODEL", "legacy-max")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.QwenExtractionModel != "legacy-turbo" || cfg.QwenReasoningModel != "legacy-max" {
		t.Fatalf("legacy migration models = %q, %q", cfg.QwenExtractionModel, cfg.QwenReasoningModel)
	}
}

func TestLoadDefaultsToTextEmbeddingV4(t *testing.T) {
	t.Setenv("QWEN_API_KEY", "test-key")
	t.Setenv("QWEN_EMBEDDING_MODEL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.QwenEmbeddingModel != "text-embedding-v4" {
		t.Fatalf("embedding model = %q, want text-embedding-v4", cfg.QwenEmbeddingModel)
	}
}

func TestLoadDefaultsToOpenAIProtocolAndNativeDashScopeSettings(t *testing.T) {
	t.Setenv("QWEN_API_KEY", "test-key")
	t.Setenv("LLM_PROTOCOL", "")
	t.Setenv("QWEN_NATIVE_BASE_URL", "")
	t.Setenv("QWEN_RERANK_MODEL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.LLMProtocol != "openai" {
		t.Fatalf("LLMProtocol = %q, want openai", cfg.LLMProtocol)
	}
	if cfg.QwenNativeBaseURL != "https://dashscope-intl.aliyuncs.com" {
		t.Fatalf("QwenNativeBaseURL = %q, want native DashScope host", cfg.QwenNativeBaseURL)
	}
	if cfg.QwenRerankModel != "qwen3-rerank" {
		t.Fatalf("QwenRerankModel = %q, want qwen3-rerank", cfg.QwenRerankModel)
	}
}

func TestLoadAcceptsDashScopeProtocolAndRejectsUnknownProtocol(t *testing.T) {
	t.Setenv("QWEN_API_KEY", "test-key")
	t.Setenv("LLM_PROTOCOL", "DaShScOpE")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load dashscope: %v", err)
	}
	if cfg.LLMProtocol != "dashscope" {
		t.Fatalf("LLMProtocol = %q, want dashscope", cfg.LLMProtocol)
	}

	t.Setenv("LLM_PROTOCOL", "unsupported")
	if _, err := Load(); err == nil {
		t.Fatal("Load unsupported protocol: expected an error")
	}
}
