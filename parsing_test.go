package Investigator

import (
    "Investigator/helper/llm/vars"
    "context"
    "strings"
    "testing"
)

func TestInputFormatting_ComplexJSON(t *testing.T) {
    jsonContent := "{\n" +
        "\t\"Analysis\": {\n" +
        "\t\t\"事件梳理与时间线\": {\n" +
        "\t\t\t\"发布节点\": \"2025年11月19日，Google正式发布Gemini 3系列模型，包含Gemini 3 Pro预览版和Gemini 3 Deep Think增强推理模式\",\n" +
        "\t\t\t\"核心定位\": \"Google CEO Sundar Pichai称其为'迄今最聪明的模型'，DeepMind CEO Demis Hassabis强调其为'世界上最好的多模态理解模型'\",\n" +
        "\t\t\t\"生态布局\": \"同步推出智能体开发平台Google Antigravity，并集成至Google搜索、AI Studio、Vertex AI及Cursor、GitHub等第三方平台\"\n" +
        "\t\t},\n" +
        "\t\t\"核心能力评估\": {\n" +
        "\t\t\t\"推理能力突破\": {\n" +
        "\t\t\t\t\"基准测试表现\": \"LMArena以1501分历史最高分登顶，Humanity's Last Exam达37.5%（博士级水平），GPQA Diamond 91.9%，MathArena Apex 23.4%创下数学能力新标杆\",\n" +
        "\t\t\t\t\"Deep Think模式\": \"在Humanity's Last Exam达41.0%，GPQA Diamond 93.8%，ARC-AGI-2达45.1%，全面超越GPT-5.1和Claude Sonnet 4.5\"\n" +
        "\t\t\t},\n" +
        "\t\t\t\"多模态理解\": {\n" +
        "\t\t\t\t\"原生多模态架构\": \"支持文本、图像、音频、视频、PDF混合输入，单次可处理900张图像或PDF页面，MMMU-Pro 81%，Video-MMMU 87.6%\",\n" +
        "\t\t\t\t\"空间理解\": \"改进对屏幕和空间理解，视觉理解准确率达72.7%，为AI Agent操作计算机奠定基础\"\n" +
        "\t\t\t},\n" +
        "\t\t\t\"编程与Agent能力\": {\n" +
        "\t\t\t\t\"Vibe Coding\": \"WebDev Arena以1487 Elo排名第一，支持自然语言描述生成完整前端应用，支持草图/视频转代码\",\n" +
        "\t\t\t\t\"智能体执行\": \"Terminal-Bench 2.0达54.2%，SWE-bench Verified 76.2%，Google Antigravity平台实现端到端自主开发（写代码→测试→修复→部署）\"\n" +
        "\t\t\t},\n" +
        "\t\t\t\"技术规格\": \"100万Token输入上下文窗口，64k输出限制，稀疏MoE架构，基于Google TPU训练，知识截止2025年1月\"\n" +
        "\t\t},\n" +
        "\t\t\"局限与挑战\": {\n" +
        "\t\t\t\"稳定性问题\": \"生成后端或嵌入式代码时易出现逻辑漏洞，Python/Go项目一次通过率较低；长上下文超过100K Token时响应质量下降\",\n" +
        "\t\t\t\"基准测试与实际场景差距\": \"加州大学戴维斯分校团队反馈，在安防摄像头、门铃等真实场景视觉理解中性能反而不如上一代，暴露公榜benchmark与实际落地场景的Gap\",\n" +
        "\t\t\t\"竞品对比劣势\": \"在SWE-bench测试中落后于Claude 4.5，处理底层Bug时易陷入死循环；复杂多跳搜索任务中稳定性不如GPT-5 Pro\",\n" +
        "\t\t\t\"幻觉与可靠性\": \"有报道称幻觉率高达88%，在'知之为知之'的自我认知层面尚未突破；用户反馈其像'高学历实习生'，知识素养足够但实际干活错漏百出\",\n" +
        "\t\t\t\"市场定位矛盾\": \"ToC属性爆棚但ToB属性不足，被建议'玩一玩很不错，但尽量不要用于生产环境'，有用户报告使用时被删除800G重要文件\"\n" +
        "\t\t},\n" +
        "\t\t\"用户舆情分析\": {\n" +
        "\t\t\t\"正面评价\": \"开发者社区对其Vibe Coding和前端生成能力高度认可，认为'巨好玩'，跨模态能力有趣，在数学和推理基准上实现碾压\",\n" +
        "\t\t\t\"负面反馈\": \"企业级用户质疑其可靠性，认为谷歌来不及将模型调教得足够可靠；部分用户反映2.5版本在特定任务上更顺手\",\n" +
        "\t\t\t\"中立观点\": \"行业观察者指出这是从'追赶'到'反超'的关键一步，但需警惕为刷榜优化导致的实际性能偏差\"\n" +
        "\t\t}\n" +
        "\t},\n" +
        "\t\"Creation\": [\n" +
        "\t\t{\n" +
        "\t\t\t\"inspiration\": \"技术突破深度解析：从基准测试王者到实际应用鸿沟\",\n" +
        "\t\t\t\"content\": \"Gemini 3以1501分登顶LMArena，在Humanity's Last Exam、GPQA Diamond等测试展现博士级推理，但其真实场景表现却出现反直觉下降。本文深度剖析稀疏MoE架构、原生多模态设计与强化学习优化如何造就基准测试王者，同时揭示安防场景性能不如上一代、幻觉率高达88%等现实挑战，探讨AI行业'刷榜优化'与'落地实用'的深层矛盾。\",\n" +
        "\t\t\t\"origin\": [\n" +
        "\t\t\t\t\" `https://finance.sina.com.cn/tech/roll/2025-11-19/doc-infxwzqk2016225.shtml` \",\n" +
        "\t\t\t\t\" `https://m.huxiu.com/article/4808859.html` \",\n" +
        "\t\t\t\t\" `https://wallstreetcn.com/articles/3760053` \"\n" +
        "\t\t\t]\n" +
        "\t\t},\n" +
        "\t\t{\n" +
        "\t\t\t\"inspiration\": \"开发者实测报告：Vibe Coding革命还是'高学历实习生'陷阱？\",\n" +
        "\t\t\t\"content\": \"基于Google Antigravity平台实测，Gemini 3的Vibe Coding能力可一键生成赛博朋克风格待办应用，草图转代码大幅降低门槛。但深入测试发现，复杂后端项目漏洞频发，长上下文重构质量断崖式下跌，甚至有用户被删800G文件。对比Claude 4.5在SWE-bench的稳定性优势，本文给出生产环境使用风险评级与规避策略。\",\n" +
        "\t\t\t\"origin\": [\n" +
        "\t\t\t\t\" `https://cj.sina.cn/articles/view/7857201853/1d45362bd06801oucy?froms=ggmp&vt=4` \",\n" +
        "\t\t\t\t\" `https://user.guancha.cn/main/content?id=1552385` \",\n" +
        "\t\t\t\t\" `https://zhuanlan.zhihu.com/p/1974376875217416222` \"\n" +
        "\t\t\t]\n" +
        "\t\t},\n" +
        "\t\t{\n" +
        "\t\t\t\"inspiration\": \"多模态能力双刃剑：从视频理解到屏幕操作的跃迁与局限\",\n" +
        "\t\t\t\"content\": \"Gemini 3在Video-MMMU达87.6%，可分析教学视频生成交互式单词卡，屏幕截图理解准确率72.7%达SOTA两倍。但Vision组实测发现，门铃摄像头等真实场景理解反而退步，暴露训练数据与分布外场景的错配。结合其支持900张图片/单次、百万Token上下文的技术规格，解析原生多模态架构的潜力与当前边界。\",\n" +
        "\t\t\t\"origin\": [\n" +
        "\t\t\t\t\" `https://www.stcn.com/article/detail/3501141.html` \",\n" +
        "\t\t\t\t\" `https://m.huxiu.com/article/4808859.html` \",\n" +
        "\t\t\t\t\" `https://www.cnblogs.com/wintersun/p/19244057` \"\n" +
        "\t\t\t]\n" +
        "\t\t},\n" +
        "\t\t{\n" +
        "\t\t\t\"inspiration\": \"AI Agent新纪元：Google Antigravity能否挑战Cursor统治地位？\",\n" +
        "\t\t\t\"content\": \"Google Antigravity平台赋予Gemini 3直接操作编辑器、终端、浏览器的权限，实现'自动驾驶式开发'。Terminal-Bench 2.0得分54.2%，但相比Claude在底层Bug修复的成熟度仍有差距。本文对比测试端到端任务执行能力，分析其智能体生态整合策略，探讨在Cursor、GitHub Copilot主导的市场中，Google的'主动伙伴'定位能否撕开缺口。\",\n" +
        "\t\t\t\"origin\": [\n" +
        "\t\t\t\t\" `https://lpcomment.com/2025/11/19/google%E6%AD%A3%E5%BC%8F%E6%8E%A8%E5%87%BAgemini-3-ai%E6%A8%A1%E5%9E%8B%EF%BC%8C%E8%83%BD%E5%8A%9B%E5%BE%97%E5%88%B0%E5%A4%A7%E5%B9%85%E6%8F%90%E5%8D%87/` \",\n" +
        "\t\t\t\t\" `https://www.shengwang.cn/blog/blogdetail/gemini-3/` \",\n" +
        "\t\t\t\t\" `https://inews.hket.com/article/4040039/Google%EF%BD%9C%E6%96%B0%E6%8E%A8%20Gemini%203%E6%8E%A8%E7%90%86%E8%83%BD%E5%8A%9B%E6%9B%B4%E5%BC%BA%E3%80%80%E9%A6%99%E6%B8%AF%E5%8F%AF%E7%94%A8%E5%93%AA%E4%BA%9B%E5%8A%9F%E8%83%BD%EF%BC%9F?mtc=20040` \"\n" +
        "\t\t\t]\n" +
        "\t\t},\n" +
        "\t\t{\n" +
        "\t\t\t\"inspiration\": \"定价策略与市场博弈：百万Token 2美元能否撼动GPT-5.1地位？\",\n" +
        "\t\t\t\"content\": \"Gemini 3 Pro定价为每百万Token输入2美元（<20万Token）或4美元（>20万Token），输出12-18美元，显著低于GPT-5.1 Pro。配合100万上下文窗口和64k输出能力，Google显然瞄准企业级市场。但用户反馈其ToB属性不足，稳定性不如GPT-5 Pro。结合6.5亿月活和20亿AI Overviews用户基础，分析Google的'基础模型+全栈产品'策略能否在OpenAI和Anthropic夹击下实现反超。\",\n" +
        "\t\t\t\"origin\": [\n" +
        "\t\t\t\t\" `https://www.51cto.com/article/830114.html` \",\n" +
        "\t\t\t\t\" `https://finance.sina.com.cn/tech/roll/2025-11-19/doc-infxwzqk2016225.shtml` \",\n" +
        "\t\t\t\t\" `https://www.stcn.com/article/detail/3501141.html` \"\n" +
        "\t\t\t]\n" +
        "\t\t}\n" +
        "\t],\n" +
        "\t\"SEO\": [\n" +
        "\t\t\"Gemini 3能力评测\",\n" +
        "\t\t\"Gemini 3 vs GPT-5.1\",\n" +
        "\t\t\"Gemini 3编程实测\",\n" +
        "\t\t\"Google Antigravity平台\",\n" +
        "\t\t\"Gemini 3幻觉问题\",\n" +
        "\t\t\"多模态AI模型对比\",\n" +
        "\t\t\"Vibe Coding革命\",\n" +
        "\t\t\"AI Agent开发\",\n" +
        "\t\t\"Gemini 3定价策略\",\n" +
        "\t\t\"基准测试与实际应用差距\"\n" +
        "\t]\n" +
        "}"

    input := vars.LLMResponse{Content: jsonContent}
    out, err := InputFormatting(context.Background(), input)
    if err != nil {
        t.Fatalf("error: %v", err)
    }
    if out.Analysis == "" {
        t.Fatalf("analysis empty")
    }
    if !strings.Contains(out.Analysis, "\"事件梳理与时间线\"") {
        t.Fatalf("analysis missing expected key")
    }
    if len(out.Creation) != 5 {
        t.Fatalf("creation len=%d", len(out.Creation))
    }
    if len(out.Creation[0].Origin) != 3 {
        t.Fatalf("origin len=%d", len(out.Creation[0].Origin))
    }
    if out.Creation[0].Origin[0] != "https://finance.sina.com.cn/tech/roll/2025-11-19/doc-infxwzqk2016225.shtml" {
        t.Fatalf("origin cleaned wrong: %s", out.Creation[0].Origin[0])
    }
    if len(out.SEO) != 10 {
        t.Fatalf("seo len=%d", len(out.SEO))
    }
    if out.SEO[6] != "Vibe Coding革命" {
        t.Fatalf("seo item mismatch: %s", out.SEO[6])
    }
}

