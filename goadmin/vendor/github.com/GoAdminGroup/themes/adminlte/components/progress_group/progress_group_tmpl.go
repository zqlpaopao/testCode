package progress_group

var List = map[string]string{
	"progress-group": `{{define "progress-group"}}
    <div class="progress-group">
        <span class="progress-text">{{langHtml .Title}}</span>
        <span class="progress-number"><b>{{.Molecular}}</b>/{{.Denominator}}</span>

        <div class="progress sm">
            {{if .IsHexColor}}
                <div class="progress-bar" style="width: {{.Percent}}%;background-color: {{.Color}} !important;"></div>
            {{else}}
                <div class="progress-bar progress-bar-{{.Color}}" style="width: {{.Percent}}%"></div>
            {{end}}
        </div>
    </div>
{{end}}`,
}
