package utils

import "testing"

func TestMap(t *testing.T) {
	items := []int{1, 2, 3, 6}

	contexts := Map(&items, func(ctx *MapContext[int]) *MapContext[int] {
		return ctx
	})

	for index, ctx := range contexts {
		if ctx.Index != index {
			t.Errorf("Index should be %d but is %d instead", index, ctx.Index)
		}

		if *ctx.Item != items[index] {
			t.Errorf("Item should be %d but is %d instead", items[index], *ctx.Item)
		}

		if len(*ctx.Items) != len(items) {
			t.Errorf("*ctx.Items and items should have the same length but *ctx.Items contains %d items while items contains %d items", len(*ctx.Items), len(items))
		}

		if index == 0 {
			if !ctx.FirstItem {
				t.Error("For first item ctx.FirstItem should be true but it's false instead")
			}
		} else {
			if ctx.FirstItem {
				t.Error("For every item that's not the first one ctx.FirstItem should be false but it's true instead")
			}
		}

		if index == len(items)-1 {
			if !ctx.LastItem {
				t.Error("For last item ctx.LastItem should be true but it's false instead")
			}
		} else {
			if ctx.LastItem {
				t.Error("For every item that's not the last one ctx.LastItem should be false but it's true instead")
			}
		}
	}
}
