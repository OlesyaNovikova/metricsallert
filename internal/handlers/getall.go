package handlers

import (
	"fmt"
	"net/http"
	"text/template"
)

const tplStr = `<table>
<thead>
	<tr>
		<th>Metric</th>
		<th>Value</th>
	</tr>
</thead>
<tbody>
	{{range $name, $value := . }}
		<tr>
			<td>{{ $name }}</td>
			<td>{{ $value }}</td>
		</tr>
	{{ end }}
</tbody>
</table>`

func GetAllMems() http.HandlerFunc {
	fn := func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodGet {
			fmt.Print("Only GET requests are allowed!\n")
			http.Error(res, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
			return
		}
		ctx := req.Context()
		res.Header().Set("Content-Type", "text/html")
		table,err := memBase.s.GetAll(ctx)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
		tpl, err := template.New("table").Parse(tplStr)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "text/html")
		//res.WriteHeader(http.StatusOK)
		err = tpl.Execute(res, table)
		if err != nil {
			fmt.Println(err)
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
	return http.HandlerFunc(fn)
}
