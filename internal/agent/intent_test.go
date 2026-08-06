package agent

import "testing"

func TestClassifyIntent_GreetingsGoToChat(t *testing.T) {
	cases := []string{"hi", "Hi!", "hello", "hey there", "thanks", "thank you", "ok", "bye"}
	for _, c := range cases {
		mode, ok := ClassifyIntent(c)
		if !ok || mode != ModeChat {
			t.Fatalf("%q → (%q, %v), want (chat, true)", c, mode, ok)
		}
	}
}

func TestClassifyIntent_BuildVerbs(t *testing.T) {
	cases := []string{
		"fix the login bug",
		"implement user auth",
		"create a new README",
		"refactor the agent loop",
		"add a /mode command",
	}
	for _, c := range cases {
		mode, ok := ClassifyIntent(c)
		if !ok || mode != ModeBuild {
			t.Fatalf("%q → (%q, %v), want (build, true)", c, mode, ok)
		}
	}
}

func TestClassifyIntent_PlanVsChat(t *testing.T) {
	mode, ok := ClassifyIntent("explain this codebase architecture")
	if !ok || mode != ModePlan {
		t.Fatalf("codebase question → (%q, %v), want plan", mode, ok)
	}

	mode, ok = ClassifyIntent("what is a mutex?")
	if !ok || mode != ModeChat {
		t.Fatalf("simple conceptual Q → (%q, %v), want chat", mode, ok)
	}
}

func TestIsTrivialPrompt(t *testing.T) {
	if !IsTrivialPrompt("hi") {
		t.Fatal("hi should be trivial")
	}
	if IsTrivialPrompt("fix the nil pointer in agent loop") {
		t.Fatal("engineering task should not be trivial")
	}
}
