package lib

import "testing"

func TestMapReduce(t *testing.T) {
	RunSubO(t, `basic`, `
		var data = [
			{ group: "foo", x: 0.1, y: 70 }
			{ group: "foo", x: 0.2, y: 60 }
			{ group: "foo", x: 0.3, y: 50 }
			{ group: "bar", x: 0.4, y: 40 }
			{ group: "bar", x: 0.5, y: 30 }
			{ group: "bar", x: 0.6, y: 20 }
			{ group: "bar", x: 0.7, y: 10 }
			{ x: 3.14 }
			{ x: 2.718, y: 13 }
		]

		var res = map_reduce(data, ".group", [
			{ reduce: count }
			{ reduce: min, pick: ".x" }
			{ reduce: max, pick: &v -> v.y? }
		])

		res |
			map(&group,res -> {
				GROUP: group
				COUNT: res[0]
				MINX: res[1]
				MAXY: res[2]
			}) |
			print_table
`, `
GROUP COUNT MINX  MAXY
foo   3     0.1   70
bar   4     0.4   40
<nil> 2     2.718 13
`)
}
