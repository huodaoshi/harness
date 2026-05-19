package session_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/huodaoshi/harness/backend/internal/session"
	"github.com/huodaoshi/harness/backend/internal/store"
)

func TestProfileInject_TwentyCases(t *testing.T) {
	ctx := context.Background()
	cases := buildInjectCases()
	if len(cases) != 20 {
		t.Fatalf("want 20 cases, got %d", len(cases))
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mem := store.NewMemoryStore()
			seedStore(t, ctx, mem, tc)

			exec, err := session.NewExecutorWithStore(ctx, mem)
			if err != nil {
				t.Fatal(err)
			}

			out, err := exec.RunTurn(ctx, session.Input{
				UserID:  tc.userID,
				Message: tc.message,
				Mode:    tc.mode,
			})
			if err != nil {
				t.Fatal(err)
			}
			if out.Crisis != nil || out.Medical != nil || out.Block != nil {
				t.Fatal("unexpected gate branch")
			}

			blob := out.Chat + "\n" + out.InjectBlock
			for _, want := range tc.mustContain {
				if !strings.Contains(blob, want) {
					t.Fatalf("missing %q in output:\n%s", want, blob)
				}
			}
			for _, forbid := range tc.mustNotContain {
				if strings.Contains(blob, forbid) {
					t.Fatalf("forbidden %q in output:\n%s", forbid, blob)
				}
			}
		})
	}
}

func TestProfileInject_NoProfileNoPanic(t *testing.T) {
	ctx := context.Background()
	mem := store.NewMemoryStore()
	exec, err := session.NewExecutorWithStore(ctx, mem)
	if err != nil {
		t.Fatal(err)
	}
	out, err := exec.RunTurn(ctx, session.Input{
		UserID:  "user-empty",
		Message: "只是随便聊聊",
		Mode:    "normal",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.Chat == "" {
		t.Fatal("expected chat reply")
	}
}

func TestProfileInject_DistressMode(t *testing.T) {
	ctx := context.Background()
	mem := store.NewMemoryStore()
	_ = mem.UpsertProfile(ctx, store.RelationshipProfile{
		UserID: "u09", CurrentIssue: "催婚压力",
	})
	exec, err := session.NewExecutorWithStore(ctx, mem)
	if err != nil {
		t.Fatal(err)
	}
	out, err := exec.RunTurn(ctx, session.Input{
		UserID: "u09", Message: "撑不住", Mode: "distress",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.Chat, "催婚压力") || !strings.Contains(out.Chat, "很难受是正常的") {
		t.Fatalf("got %q", out.Chat)
	}
}

type injectCase struct {
	name           string
	userID         string
	message        string
	mode           string
	profile        *store.RelationshipProfile
	summary        *store.SessionSummary
	mustContain    []string
	mustNotContain []string
}

func buildInjectCases() []injectCase {
	return []injectCase{
		{name: "01_profile_issue_only", userID: "u01", message: "我又想他了", mode: "normal", profile: &store.RelationshipProfile{
			UserID: "u01", CurrentIssue: "和伴侣反复冷战", Self: store.ProfileSelf{Note: "容易焦虑"},
		}, mustContain: []string{"和伴侣反复冷战", "容易焦虑"}},
		{name: "02_summary_only", userID: "u02", message: "继续昨晚的话题", mode: "normal", summary: &store.SessionSummary{
			UserID: "u02", Summary3: []string{"昨晚谈到母亲边界", "情绪一度失控", "最后稍微平静"},
		}, mustContain: []string{"昨晚谈到母亲边界", "情绪一度失控"}},
		{name: "03_profile_and_summary", userID: "u03", message: "还是很难受", mode: "normal", profile: &store.RelationshipProfile{
			UserID: "u03", CurrentIssue: "父亲控制欲强",
		}, summary: &store.SessionSummary{
			UserID: "u03", Summary3: []string{"提到想搬出去住", "害怕冲突", "需要支持"},
		}, mustContain: []string{"父亲控制欲强", "想搬出去住"}},
		{name: "04_no_data", userID: "u04", message: "你好", mode: "normal", mustContain: []string{"我听到你了"}, mustNotContain: []string{"[关系档案]", "[上次会话摘要]"}},
		{name: "05_person_label", userID: "u05", message: "今天他又这样", mode: "normal", profile: &store.RelationshipProfile{
			UserID: "u05", People: []store.ProfilePerson{{Label: "阿明", Relation: "伴侣", Note: "回避型"}},
			CurrentIssue: "沟通僵局",
		}, mustContain: []string{"阿明", "回避型", "沟通僵局"}},
		{name: "06_two_people", userID: "u06", message: "夹在中间", mode: "normal", profile: &store.RelationshipProfile{
			UserID: "u06", People: []store.ProfilePerson{
				{Label: "妈妈", Relation: "母亲"},
				{Label: "老公", Relation: "伴侣"},
			},
		}, mustContain: []string{"妈妈", "老公"}},
		{name: "07_self_note", userID: "u07", message: "很累", mode: "normal", profile: &store.RelationshipProfile{
			UserID: "u07", Self: store.ProfileSelf{Note: "长期失眠"},
		}, mustContain: []string{"长期失眠"}},
		{name: "08_latest_summary_wins", userID: "u08", message: "后来呢", mode: "normal", mustContain: []string{"最新一次复盘"}, mustNotContain: []string{"旧摘要不应出现"}},
		{name: "09_profile_distress_issue", userID: "u09b", message: "撑不住", mode: "distress", profile: &store.RelationshipProfile{
			UserID: "u09b", CurrentIssue: "催婚压力",
		}, mustContain: []string{"催婚压力"}},
		{name: "10_long_issue_snippet", userID: "u10", message: "怎么办", mode: "normal", profile: &store.RelationshipProfile{
			UserID: "u10", CurrentIssue: "婆婆同住矛盾升级",
		}, mustContain: []string{"婆婆同住"}},
		{name: "11_summary_line2", userID: "u11", message: "嗯", mode: "normal", summary: &store.SessionSummary{
			UserID: "u11", Summary3: []string{"第一句", "第二句关键词", "第三句"},
		}, mustContain: []string{"第二句关键词"}},
		{name: "12_profile_without_summary", userID: "u12", message: "聊几句", mode: "normal", profile: &store.RelationshipProfile{
			UserID: "u12", CurrentIssue: "异地信任问题",
		}, mustNotContain: []string{"[上次会话摘要]"}, mustContain: []string{"异地信任"}},
		{name: "13_summary_without_profile", userID: "u13", message: "继续", mode: "normal", summary: &store.SessionSummary{
			UserID: "u13", Summary3: []string{"仅摘要存在"},
		}, mustNotContain: []string{"[关系档案]"}, mustContain: []string{"仅摘要存在"}},
		{name: "14_empty_user_id", userID: "", message: "匿名倾诉", mode: "normal", mustContain: []string{"我听到你了"}, mustNotContain: []string{"【已读上下文】"}},
		{name: "15_partner_conflict", userID: "u15", message: "他又冷暴力", mode: "normal", profile: &store.RelationshipProfile{
			UserID: "u15", People: []store.ProfilePerson{{Label: "小陈", Relation: "伴侣"}}, CurrentIssue: "冷暴力",
		}, mustContain: []string{"小陈", "冷暴力"}},
		{name: "16_parent_boundary", userID: "u16", message: "回家就吵", mode: "normal", profile: &store.RelationshipProfile{
			UserID: "u16", CurrentIssue: "父母过度干涉",
		}, summary: &store.SessionSummary{UserID: "u16", Summary3: []string{"谈到经济独立"}},
			mustContain: []string{"父母过度干涉", "经济独立"}},
		{name: "17_three_line_summary", userID: "u17", message: "在吗", mode: "normal", summary: &store.SessionSummary{
			UserID: "u17", Summary3: []string{"A", "B", "C"},
		}, mustContain: []string{"A", "B", "C"}},
		{name: "18_self_and_issue", userID: "u18", message: "我又哭了", mode: "normal", profile: &store.RelationshipProfile{
			UserID: "u18", Self: store.ProfileSelf{Note: "高敏感"}, CurrentIssue: "分手未愈",
		}, mustContain: []string{"高敏感", "分手未愈"}},
		{name: "19_work_family", userID: "u19", message: "两边压力", mode: "normal", profile: &store.RelationshipProfile{
			UserID: "u19", CurrentIssue: "工作与家庭拉扯", Self: store.ProfileSelf{Note: "职场新人"},
		}, mustContain: []string{"工作与家庭拉扯", "职场新人"}},
		{name: "20_inject_marker", userID: "u20", message: "test", mode: "normal", profile: &store.RelationshipProfile{
			UserID: "u20", CurrentIssue: "标记测试议题",
		}, mustContain: []string{"【已读上下文】", "标记测试议题"}},
	}
}

func seedStore(t *testing.T, ctx context.Context, mem *store.MemoryStore, tc injectCase) {
	t.Helper()
	if tc.name == "08_latest_summary_wins" {
		_ = mem.SaveSummary(ctx, store.SessionSummary{
			UserID: "u08", Summary3: []string{"旧摘要不应出现"}, CreatedAt: time.Now().UTC().Add(-48 * time.Hour),
		})
		_ = mem.SaveSummary(ctx, store.SessionSummary{
			UserID: "u08", Summary3: []string{"最新一次复盘"}, CreatedAt: time.Now().UTC(),
		})
		return
	}
	if tc.profile != nil {
		if err := mem.UpsertProfile(ctx, *tc.profile); err != nil {
			t.Fatal(err)
		}
	}
	if tc.summary != nil {
		if err := mem.SaveSummary(ctx, *tc.summary); err != nil {
			t.Fatal(err)
		}
	}
}
