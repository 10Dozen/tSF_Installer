package handlers

import (
	"fmt"
	"html/template"
	"slices"
	"strings"
	"tSF_Installer/pkg/application"
)

type LogHandlerData struct {
	Message template.HTML
	Tag     string
	Badge   string
}

type DoneHandlerData struct {
	Logs      []*LogHandlerData
	TargetDir template.HTML
}

var tagToBadge = map[string]string{
	application.LOG_TAG_INFO:  "secondary",
	application.LOG_TAG_ERROR: "danger",
}

func prepareLogs(logs []*application.LogMessage, order int) []*LogHandlerData {
	preparedLines := []*LogHandlerData{}
	for _, log := range logs {
		preparedLines = append(
			preparedLines,
			&LogHandlerData{template.HTML(log.Message), log.Level, tagToBadge[log.Level]},
		)
	}
	if order < 0 {
		slices.Reverse(preparedLines)
	}

	return preparedLines
}

func NextUnresolvedDiffInfo(app *application.Application) *DiffInfo {
	var nextDiffIdx int
	var nextDiff *application.Diff
	for idx, diff := range app.Diffs {
		if diff.Resolution == "" {
			nextDiffIdx = idx
			nextDiff = diff
			break
		}
	}

	if nextDiff == nil {
		return nil
	}

	info := DiffInfo{
		Idx:         nextDiffIdx,
		Position:    fmt.Sprintf("%d/%d", nextDiffIdx+1, len(app.Diffs)),
		File:        strings.Replace(nextDiff.OriginalFile, app.TargetDir, "", -1),
		DiffContent: template.HTML(nextDiff.Diff),
		HasVSC:      app.HasVSC,
		BlockUI:     app.State != application.APP_STATE_DIFF,
	}

	return &info
}

func NewDoneHandlerData(app *application.Application) *DoneHandlerData {
	return &DoneHandlerData{
		TargetDir: template.HTML(app.TargetDir),
		Logs:      prepareLogs(app.Logs, 1),
	}
}
