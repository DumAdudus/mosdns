/*
 * Copyright (C) 2020-2022, IrineSistiana
 *
 * This file is part of mosdns.
 *
 * mosdns is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * mosdns is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package has_wanted_ans

import (
	"context"

	"github.com/IrineSistiana/mosdns/v5/pkg/query_context"
	"github.com/IrineSistiana/mosdns/v5/plugin/executable/sequence"
	"go.uber.org/zap"
)

const PluginType = "has_wanted_ans"

func init() {
	sequence.MustRegMatchQuickSetup(PluginType, QuickSetup)
}

type hasQuestionAns struct {
	L *zap.Logger
}

func (h hasQuestionAns) Match(_ context.Context, qCtx *query_context.Context) (bool, error) {
	q := qCtx.Q()
	if len(q.Question) == 0 {
		return false, nil
	}
	r := qCtx.R()
	if r == nil || len(r.Answer) == 0 {
		return false, nil
	}
	h.L.Debug("has_wanted_ans processing", zap.String("question", q.Question[0].String()))

	question := q.Question[0]
	for _, rr := range r.Answer {
		header := rr.Header()
		if header.Rrtype == question.Qtype && header.Class == question.Qclass {
			h.L.Debug("has_wanted_ans hit", zap.String("question", q.Question[0].String()))
			return true, nil
		}
	}
	return false, nil
}

func QuickSetup(bq sequence.BQ, _ string) (sequence.Matcher, error) {
	return hasQuestionAns{L: bq.L()}, nil
}
