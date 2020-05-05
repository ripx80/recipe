package recipe

import (
	"math"

	"github.com/ripx80/recipe/pkgs/plato"
)

/*MaltSum calculate sum of all malts
Schüttung(malt) muss prozentual im anschluss berechnet werden
calc when parsing and set in a unexproted mash var with % content of sum
*/
func (r *Recipe) MaltSum() float64 {
	var sum float64
	for _, v := range r.Mash.Malts {
		sum += v.Amount
	}
	return sum
}

/*WaterSum returns the total water sum*/
func (r *Recipe) WaterSum() float64 {
	return r.Water.MainCast + r.Water.Grouting
}

/*Round rounds the float numbers given*/
func Round(x float64) float64 {
	return math.Round(x/1) * 1
}

/*Filling calculates the complete filling in kg*/
func Filling(water, originalWort, sg, sudyield float64) float64 {
	return (((water * sg) * originalWort) / sudyield)
}

/*SudYieldF calculates the sudyield for your brewery with factor
Included the temp factor by water temp of 100C = 0.96
M3 calculates the Temperaturfactor by 100C (Wort measured by 100 after cooking)
(SG WATER 100°C 958.4 kg/m³ / SG Water 20°C 998.2 = 0.96)
*/
func SudYieldF(water, originalWort, sg, filling float64) float64 {
	return (((water * sg) * originalWort) * 0.96) / filling
}

/*SudYield calculates the sudyield for your brewery*/
func SudYield(water, originalWort, sg, filling float64) float64 {
	return ((water * sg) * originalWort) / filling
}

/*ScaleAmount scale the amount by faktor of sums*/
func ScaleAmount(amount, sum, newsum float64) float64 {
	return (newsum * (amount / sum))
}

/*Scale to new water amount and yield
recipe to json
*/
func (r *Recipe) Scale(water, yield float64) (*Recipe, error) {
	Plato := plato.New()
	rScale := &Recipe{}
	if _, err := rScale.Load(r.String()); err != nil {
		return nil, err
	}

	rScale.Global.DecisiveSeasoning = water
	entry, _ := Plato.SW(r.Global.OriginalWort)

	filling := Filling(rScale.Global.DecisiveSeasoning, (r.Global.OriginalWort)*0.96, entry.SG, r.Global.SudYield) // switch from rec.Global.SudYield check it!
	rScale.Global.SudYield = Round(SudYield(rScale.Global.DecisiveSeasoning, r.Global.OriginalWort, entry.SG, filling))

	// scale malt
	maltsum := r.MaltSum()
	for k, v := range r.Mash.Malts {
		rScale.Mash.Malts[k].Amount = Round((ScaleAmount(v.Amount, maltsum, filling) * 1000))
	}
	// scale water
	waterSum := r.WaterSum()
	waterSumNew := ScaleAmount(rScale.Global.DecisiveSeasoning, r.Global.DecisiveSeasoning, waterSum)
	factor := waterSumNew / waterSum
	rScale.Water.MainCast = ScaleAmount(r.Water.MainCast, waterSum, waterSumNew)
	rScale.Water.Grouting = ScaleAmount(r.Water.Grouting, waterSum, waterSumNew)

	// scale hops, creat scaleHop func for all
	for k, v := range r.Cook.FontHops {
		rScale.Cook.FontHops[k].Amount = (factor * v.Amount)
	}
	for k, v := range r.Cook.Hops {
		rScale.Cook.Hops[k].Amount = Round((factor * v.Amount))
	}
	for k, v := range r.Cook.Whirlpool {
		rScale.Cook.Whirlpool[k].Amount = Round((factor * v.Amount))
	}
	for k, v := range r.Fermentation.Hops {
		rScale.Fermentation.Hops[k].Amount = Round((factor * v.Amount))
	}

	// scale incredients
	for k, v := range r.Cook.Ingredients {
		rScale.Cook.Ingredients[k].Amount = Round((factor * v.Amount))
	}

	// scale fermentation
	for k, v := range r.Fermentation.Ingredients {
		rScale.Fermentation.Ingredients[k].Amount = Round((factor * v.Amount))
	}
	return rScale, nil
}
