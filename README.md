# Cockroach Trigrams
Let's use CockroachDB's trigram indexes to search USDA's FDCD!

A naive approach to using trigram indexes and the similarity operator (`%`) on
the FDCD dataset has proved to be slower and less useful than would be ideal.

This repository explores different techniques of indexing and querying data to
find the ideal way to make the FDCD dataset discoverable.

Each technique will be evaluated by its speed and the relevance of the search
results it returns.


## Running locally

### Start up CockroachDB
```bash
❯❯❯ cockroach version
Build Tag:        v22.2.0
Build Time:       2022/12/05 16:56:56
Distribution:     CCL
Platform:         darwin arm64 (aarch64-apple-darwin21.2)
Go Version:       go1.19.1
C Compiler:       Clang 10.0.0
Build Commit ID:  77667a1b0101cd323090011f50cf910aaa933654
Build Type:       release

❯❯❯ cockroach start-single-node --insecure
```

### Index some data
```bash
# This will automatically download the FDCD dataset to ./data and load it into
# CRDB.
make crdb-trgrm && ./crdb-trgrm load
```

### Explore the querying techniques
The `query` command will query the dataset with each configured querier.
```
make crdb-trgrm && ./crdb-trgrm query "melk"
```

## Sample Results

### Simple Query

```
Searched "eggs" via "Raw query % operator on raw text" (626.196625ms):
        1: "NICE! NICE!, LARGE EGGS"
        2: "NICE LARGE EGGS"
        3: "FREE RANGE EGGS"
        4: "EXTRA LARGE EGGS"
        5: "EXTRA LARGE EGGS"
        6: "NICE! LARGE EGGS"
        7: "NICE! LARGE EGGS"
        8: "NICE! LARGE EGGS"
        9: "LARGE DUCK EGGS"
        10: "NICE! LARGE EGGS"


Searched "eggs" via "Analyzed query % operator on raw text" (462.954416ms):
        1: "SURPRISE EGG"
        2: "SNICKERS EGG"
        3: "THE EGG BAR"
        4: "GALERIE EGG"
        5: "EGG BEATERS EGG"
        6: "EGGNOG"
        7: "EGG PASTA"
        8: "EGG SALAD"
        9: "EGG SALAD"
        10: "JUST EGG"


Searched "eggs" via "Raw query % operator on analyzed text" (440.267625ms):
        1: "EGGNOG"
        2: "EGGS"
        3: "EGGS"
        4: "EGGS"
        5: "EGGS"
        6: "EGGS"
        7: "EGGS"


Searched "eggs" via "Analyzed query % operator on analyzed text" (448.280709ms):
        1: "RED BEET EGGS, RED BEET"
        2: "Egg, yolk, dried"
        3: "IGA LARGE EGGS"
        4: "IGA LARGE EGGS"
        5: "KIRKLAND EGGS"
        6: "SNICKERS EGG"
        7: "SUNSHINE SUNSHINE, EGGS"
        8: "SURPRISE EGG"
        9: "THE EGG BAR"
        10: "GRADE A EGGS"


Searched "eggs" via "% operator query on (food_id, token) table with join" (385.287625ms):
        1: "KELLOGG EGGO SAVORY BACON, EGG & CHEESE MEDLEYS 33.5OZ"
        2: "Kellogg's Eggo Savory Handheld Sausage Egg & Cheese 27.6oz"
        3: "Kellogg's Eggo Savory Handheld Bacon, Egg & Cheese 21.7oz"
        4: "CADBURY DAIRY MILK CHOCOLATE EGG EGGHEADS"
        5: "SHOPRITE THE GREAT EGGSCAPE 99% EGG PRODUCT"
        6: "EGGLAND'S BEST EGGLAND'S BEST, HARD-COOKED PEELED EGGS"
        7: "EGGLAND'S BEST CAGE FREE LARGE WHITE EGGS"
        8: "EGG-LAND'S BEST EXTRA LARGE BROWN EGGS"
        9: "EGG-LAND'S BEST FARM FRESH EXTRA LARGE BROWN EGGS"
        10: "EGG-LAND'S BEST LARGE BROWN EGGS"


Searched "eggs" via "ILIKE on analyzed ordered by similarity" (82.679084ms):
        1: "EGGS"
        2: "EGGS"
        3: "EGGS"
        4: "EGGS"
        5: "EGGS"
        6: "EGGS"
        7: "EGGNOG"
        8: "LARGE EGGS"
        9: "LARGE EGGS"
        10: "LARGE EGGS"


Searched "eggs" via "DIY trigram search using ILIKE ordered by similarity" (78.384834ms):
        1: "EGGS"
        2: "EGGS"
        3: "EGGS"
        4: "EGGS"
        5: "EGGS"
        6: "EGGS"
        7: "EGGNOG"
        8: "LARGE EGGS"
        9: "LARGE EGGS"
        10: "LARGE EGGS"
```

### Multi-Word Query

```
Searched "kodiak cakes" via "Raw query % operator on raw text" (1.912001041s):
        1: "KODIAK CAKES CHOCOLATE CHIP BLONDIE BROWNIE MIX, CHOCOLATE CHIP BLONDIE"
        2: "KODIAK CAKES CHOCOLATE CHIP CRUNCHY GRANOLA BAR, CHOCOLATE CHIP"
        3: "KODIAK CAKES RED RASPBERRY SUPER FRUIT SYRUP, RED RASPBERRY"
        4: "KODIAK CAKES CINNAMON ROLL UNLEASHED MUFFIN, CINNAMON ROLL"
        5: "KODIAK CAKES CHOCOLATE CHIP POWER CUP OATMEAL, CHOCOLATE CHIP"
        6: "KODIAK CAKES HOMESTEAD STYLE POWER WAFFLES, HOMESTEAD STYLE"
        7: "5 CRAB CAKES"
        8: "KODIAK CAKES DOUBLE CHOCOLATE BROWNIE MIX, DOUBLE CHOCOLATE"
        9: "KODIAK CAKES CHOCOLATE CHIP POWER CUP MUFFIN, CHOCOLATE CHIP"
        10: "KODIAK CAKES JALAPENO UNLEASHED CORNBREAD, JALAPENO"


Searched "kodiak cakes" via "Analyzed query % operator on raw text" (1.842518666s):
        1: "KODIAK CAKES HONEY GRAHAM BEAR BITES, HONEY"
        2: "CUPCAKE CAKE"
        3: "CUPCAKE CAKE"
        4: "CUPCAKE CAKE"
        5: "CUPCAKE CAKE"
        6: "CUPCAKE CAKE"
        7: "KODIAK CAKES CINNAMON POWER WAFFLES, CINNAMON"
        8: "KODIAK CAKES CHOCOLATE CHIP OATMEAL, CHOCOLATE CHIP"
        9: "KODIAK CAKES PEACHES & CREAM OATMEAL, PEACHES & CREAM"
        10: "KODIAK CAKES CINNAMON OATMEAL, CINNAMON"


Searched "kodiak cakes" via "Raw query % operator on analyzed text" (1.683989125s):
        1: "KODIAK CAKES BLUEBERRY POWER WAFFLES, BLUEBERRY"
        2: "KODIAK CAKES HONEY GRAHAM BEAR BITES, HONEY"
        3: "KODIAK CAKES CARAMEL UNLEASHED OATMEAL, CARAMEL"
        4: "KODIAK CAKES CHOCOLATE CHIP POWER WAFFLES, CHOCOLATE CHIP"
        5: "KODIAK CAKES ORGANIC OAT AND BAKING MIX"
        6: "KODIAK CAKES CINNAMON POWER WAFFLES, CINNAMON"
        7: "KODIAK CAKES PEACHES & CREAM OATMEAL, PEACHES & CREAM"
        8: "KODIAK CAKES CHOCOLATE CHIP OATMEAL, CHOCOLATE CHIP"
        9: "KODIAK CAKES CINNAMON OATMEAL, CINNAMON"


Searched "kodiak cakes" via "Analyzed query % operator on analyzed text" (1.6860825s):
        1: "KODIAK CAKES COOKIE MIX, OATMEAL DARK CHOCOLATE"
        2: "KODIAK CAKES CHOCOLATE CHIP PROTEIN-PACKED FLAPJACKS, CHOCOLATE CHIP"
        3: "KODIAK CAKES CINNAMON ROLL UNLEASHED MUFFIN, CINNAMON ROLL"
        4: "KODIAK CAKES BIRTHDAY CAKE WITH SPRINKLES BAKING MIX, BIRTHDAY CAKE WITH SPRINKLES"
        5: "KODIAK CAKES FRONTIER FLAPJACK AND BAKING MIX"
        6: "KODIAK CAKES HOMESTEAD STYLE POWER WAFFLES, HOMESTEAD STYLE"
        7: "KODIAK CAKES APPLE CINNAMON UNLEASHED MUFFIN, APPLE CINNAMON"
        8: "KODIAK CAKES CHOCOLATE CHIP CRUNCHY GRANOLA BAR, CHOCOLATE CHIP"
        9: "KODIAK CAKES CHOCOLATE CHIP CRUNCHY GRANOLA BARS, CHOCOLATE CHIP"
        10: "KODIAK CAKES CORNBREAD MIX, HOMESTEAD STYLE"


Searched "kodiak cakes" via "% operator query on (food_id, token) table with join" (3.20705775s):
        1: "ORIGINAL CAKERIE AUTUMN SPICE CAKE FOR 2"
        2: "CONFETTI CAKEIRE BAR CAKE"
        3: "CAKEBALLZ CAKEBALLZ, CAKE BALLS, BIRTHDAY CAKE, BIRTHDAY CAKE"
        4: "CAKEBALLZ CAKEBALLZ, CAKE BALLS, CHOCOLATE, CHOCOLATE"
        5: "CAKEBALLZ CAKEBALLZ, CAKE BALLS, RED VELVET, RED VELVET"
        6: "KODIAK CAKES KODIAK CAKES, GRANOLA UNLEASHED, FRENCH VANILLA ALMOND, FRENCH VANILLA ALMOND"
        7: "KODIAK CAKES KODIAK CAKES, GRANOLA UNLEASHED, VERMONT MAPLE PECAN, VERMONT MAPLE PECAN"
        8: "KODIAK CAKES BLUEBERRY UNLEASHED MUFFIN, BLUEBERRY"
        9: "KODIAK CAKES APPLE CINNAMON UNLEASHED MUFFIN, APPLE CINNAMON"
        10: "KODIAK CAKES PUMPKIN DARK CHOCOLATE UNLEASHED MUFFIN, PUMPKIN DARK CHOCOLATE"


Searched "kodiak cakes" via "ILIKE on analyzed ordered by similarity" (122.429792ms):
        1: "KODIAK CAKES CINNAMON OATMEAL, CINNAMON"
        2: "KODIAK CAKES CHOCOLATE CHIP OATMEAL, CHOCOLATE CHIP"
        3: "KODIAK CAKES PEACHES & CREAM OATMEAL, PEACHES & CREAM"
        4: "KODIAK CAKES CINNAMON POWER WAFFLES, CINNAMON"
        5: "KODIAK CAKES ORGANIC OAT AND BAKING MIX"
        6: "KODIAK CAKES CHOCOLATE CHIP POWER WAFFLES, CHOCOLATE CHIP"
        7: "KODIAK CAKES HONEY GRAHAM BEAR BITES, HONEY"
        8: "KODIAK CAKES CARAMEL UNLEASHED OATMEAL, CARAMEL"
        9: "KODIAK CAKES BLUEBERRY POWER WAFFLES, BLUEBERRY"
        10: "KODIAK CAKES DARK CHOCOLATE POWER WAFFLES, DARK CHOCOLATE"


Searched "kodiak cakes" via "DIY trigram search using ILIKE ordered by similarity" (41.377292ms):
        1: "KODIAK CAKES CINNAMON OATMEAL, CINNAMON"
        2: "KODIAK CAKES CHOCOLATE CHIP OATMEAL, CHOCOLATE CHIP"
        3: "KODIAK CAKES PEACHES & CREAM OATMEAL, PEACHES & CREAM"
        4: "KODIAK CAKES CINNAMON POWER WAFFLES, CINNAMON"
        5: "KODIAK CAKES ORGANIC OAT AND BAKING MIX"
        6: "KODIAK CAKES CHOCOLATE CHIP POWER WAFFLES, CHOCOLATE CHIP"
        7: "KODIAK CAKES BLUEBERRY POWER WAFFLES, BLUEBERRY"
        8: "KODIAK CAKES HONEY GRAHAM BEAR BITES, HONEY"
        9: "KODIAK CAKES CARAMEL UNLEASHED OATMEAL, CARAMEL"
        10: "KODIAK CAKES DARK CHOCOLATE POWER WAFFLES, DARK CHOCOLATE"
```

### Multi-Word Query With Typo

```
Searched "kodik cakes" via "Raw query % operator on raw text" (1.902022083s):
        1: "KELKIN RICE CAKES"
        2: "CAKE"
        3: "CUPCAKES"
        4: "CUPCAKES"
        5: "CUPCAKES"
        6: "CUPCAKES"
        7: "CUPCAKES"
        8: "CUPCAKES"
        9: "CUPCAKES"
        10: "CUPCAKES"


Searched "kodik cakes" via "Analyzed query % operator on raw text" (1.816818s):
        1: "CAKES"
        2: "CREME CAKE"
        3: "CRUMB CAKE"
        4: "RICE CAKE"
        5: "CARROT CAKE, CARROT"
        6: "CUPCAKE CAKE"
        7: "CUPCAKE CAKE"
        8: "CUPCAKE CAKE"
        9: "CUPCAKE CAKE"
        10: "CUPCAKE CAKE"


Searched "kodik cakes" via "Raw query % operator on analyzed text" (1.707279041s):
        1: "CAKE"
        2: "CAKES"


Searched "kodik cakes" via "Analyzed query % operator on analyzed text" (1.664149875s):
        1: "KROGER CAKE CUPS"
        2: "KODIAK CAKES CINNAMON OATMEAL, CINNAMON"
        3: "KROGER CAKE CUPS"
        4: "CUPCAKE CAKE"
        5: "CUPCAKE CAKE"
        6: "CRUMB CAKE"
        7: "CUPCAKE CAKE"
        8: "CAKE BITES"
        9: "CREME CAKE"
        10: "CARROT CAKE, CARROT"


Searched "kodik cakes" via "% operator query on (food_id, token) table with join" (2.896549334s):
        1: "ORIGINAL CAKERIE AUTUMN SPICE CAKE FOR 2"
        2: "CONFETTI CAKEIRE BAR CAKE"
        3: "CAKEBALLZ CAKEBALLZ, CAKE BALLS, BIRTHDAY CAKE, BIRTHDAY CAKE"
        4: "CAKEBALLZ CAKEBALLZ, CAKE BALLS, CHOCOLATE, CHOCOLATE"
        5: "CAKEBALLZ CAKEBALLZ, CAKE BALLS, RED VELVET, RED VELVET"
        6: "KODIAK CAKES KODIAK CAKES, GRANOLA UNLEASHED, FRENCH VANILLA ALMOND, FRENCH VANILLA ALMOND"
        7: "KODIAK CAKES KODIAK CAKES, GRANOLA UNLEASHED, VERMONT MAPLE PECAN, VERMONT MAPLE PECAN"
        8: "KODIAK CAKES BLUEBERRY UNLEASHED MUFFIN, BLUEBERRY"
        9: "KODIAK CAKES APPLE CINNAMON UNLEASHED MUFFIN, APPLE CINNAMON"
        10: "KODIAK CAKES PUMPKIN DARK CHOCOLATE UNLEASHED MUFFIN, PUMPKIN DARK CHOCOLATE"


Searched "kodik cakes" via "ILIKE on analyzed ordered by similarity" (5.070084ms):
        No Results Found


Searched "kodik cakes" via "DIY trigram search using ILIKE ordered by similarity" (25.670958ms):
        1: "KODIAK CAKES CINNAMON OATMEAL, CINNAMON"
        2: "KODIAK CAKES CHOCOLATE CHIP OATMEAL, CHOCOLATE CHIP"
        3: "KODIAK CAKES PEACHES & CREAM OATMEAL, PEACHES & CREAM"
        4: "KODIAK CAKES CINNAMON POWER WAFFLES, CINNAMON"
        5: "KODIAK CAKES ORGANIC OAT AND BAKING MIX"
        6: "KODIAK CAKES CHOCOLATE CHIP POWER WAFFLES, CHOCOLATE CHIP"
        7: "KODIAK CAKES CARAMEL UNLEASHED OATMEAL, CARAMEL"
        8: "KODIAK CAKES BLUEBERRY POWER WAFFLES, BLUEBERRY"
        9: "KODIAK CAKES HONEY GRAHAM BEAR BITES, HONEY"
        10: "KODIAK CAKES DARK CHOCOLATE POWER WAFFLES, DARK CHOCOLATE"
```

### Summary

Unexpectedly, we have a very obvious winner here. Synthesizing our own trigram
in the application layer to build a query using the ILIKE operator results in
very performant, typo resistant, and relevant search results.

It's peculiar that DIY trigrams outperforms standard ILIKE queries in the
"Multi-Word" scenarior. Perhaps there's a subtle optimization that we've
stumbled across when the ILIKE argument is exactly 3 characters with wild cards?
