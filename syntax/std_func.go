package syntax

// func callSelf(f rel.Value) rel.Value {
// 	result, err := f.(*rel.NativeFunction).Call(f, rel.EmptyScope)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return result
// }

// func λ(f func(g rel.Value) rel.Value) rel.Value {
// 	return rel.NewNativeLambda(f)
// }

// // (\f f(f))(\f \g \n g(f(f)(g))(n))
// var fixedPointCombinator = callSelf(λ(func(f rel.Value) rel.Value {
// 	F := f.(rel.Set)
// 	return λ(func(g rel.Value) rel.Value {
// 		G := g.(rel.Set)
// 		return λ(func(n rel.Value) rel.Value {
// 			return G.Call(F.Call(F).(rel.Set).Call(G)).(rel.Set).Call(n)
// 		})
// 	})
// }))

// // (\f f(f))(\f \t t :> \g \n g(f(f)(t))(n))
// var fixedPointTupleCombinator = callSelf(λ(func(f rel.Value) rel.Value {
// 	F := f.(rel.Set)
// 	return λ(func(t rel.Value) rel.Value {
// 		return t.(rel.Tuple).Map(func(g rel.Value) rel.Value {
// 			G := g.(rel.Set)
// 			return λ(func(n rel.Value) rel.Value {
// 				return G.Call(F.Call(F).(rel.Set).Call(t)).(rel.Set).Call(n)
// 			})
// 		})
// 	})
// }))
