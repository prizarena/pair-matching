package pairmodels

import (
	"testing"
	"strings"
)

func TestPairsBoardEntity_DrawBoard_ascii(t *testing.T) {
	board := PairsBoardEntity{
		Cells: "123456789abc",
		SizeX: 3,
		SizeY: 4,
	}
	expects := strings.Join([]string{"", "123", "456", "789", "abc", ""}, "\n")
	if result := board.DrawBoard(); result != expects {
		t.Error("Unexpected result:\n" + result)
	}
}

func TestPairsBoardEntity_DrawBoard_emoji(t *testing.T) {
	board := PairsBoardEntity{
		Cells: "🍇🍈🍉🍊🍋🍌🍍🍎🍏🍐🍑🍒",
		SizeX: 3,
		SizeY: 4,
	}
	expects := strings.Join([]string{"", "🍇🍈🍉", "🍊🍋🍌", "🍍🍎🍏", "🍐🍑🍒", ""}, "\n")
	if result := board.DrawBoard(); result != expects {
		t.Error("Unexpected result:\n" + result)
	}


	//🍓 Strawberry
	//🥝 Kiwi Fruit
	//🍅 Tomato
	//🥥 Coconut
	//🥑 Avocado
	//🍆 Eggplant
	//🥔 Potato
	//🥕 Carrot
	//🌽 Ear of Corn
	//🌶 Hot Pepper
	//🥒 Cucumber
	//🥦 Broccoli
	//🍄 Mushroom
	//🥜 Peanuts
	//🌰 Chestnut
	//🍞 Bread
	//🥐 Croissant
	//🥖 Baguette Bread
	//🥨 Pretzel
	//🥞 Pancakes
	//🧀 Cheese Wedge
	//🍖 Meat on Bone
	//🍗 Poultry Leg
	//🥩 Cut of Meat
	//🥓 Bacon
	//🍔 Hamburger
	//🍟 French Fries
	//🍕 Pizza
	//🌭 Hot Dog
	//🥪 Sandwich
	//🌮 Taco
	//🌯 Burrito
	//🍳 Cooking
	//🍲 Pot of Food
	//🥣 Bowl With Spoon
	//🥗 Green Salad
	//🍿 Popcorn
	//🥫 Canned Food
	//🍱 Bento Box
	//🍘 Rice Cracker
	//🍙 Rice Ball
	//🍚 Cooked Rice
	//🍛 Curry Rice
	//🍜 Steaming Bowl
	//🍝 Spaghetti
	//🍠 Roasted Sweet Potato
	//🍢 Oden
	//🍣 Sushi
	//🍤 Fried Shrimp
	//🍥 Fish Cake With Swirl
	//🍡 Dango
	//🥟 Dumpling
	//🥠 Fortune Cookie
	//🥡 Takeout Box
	//🍦 Soft Ice Cream
	//🍧 Shaved Ice
	//🍨 Ice Cream
	//🍩 Doughnut
	//🍪 Cookie
	//🎂 Birthday Cake
	//🍰 Shortcake
	//🥧 Pie
	//🍫 Chocolate Bar
	//🍬 Candy
	//🍭 Lollipop
	//🍮 Custard
	//🍯 Honey Pot
	//🍼 Baby Bottle
	//🥛 Glass of Milk
	//☕ Hot Beverage
	//🍵 Teacup Without Handle
	//🍶 Sake
	//🍾 Bottle With Popping Cork
	//🍷 Wine Glass
	//🍸 Cocktail Glass
	//🍹 Tropical Drink
	//🍺 Beer Mug
	//🍻 Clinking Beer Mugs
	//🥂 Clinking Glasses
	//🥃 Tumbler Glass
	//🥤 Cup With Straw
	//🥢 Chopsticks
	//🍽 Fork and Knife With Plate
	//🍴 Fork and Knife
	//🥄 Spoon
	//
	//Categories
	//😃 Smileys & People

}
