# TODO
0. BFS GENERATIONAL SIEVE BASED SOLVER
1. When stepping through each tier of the graph we don't need to make a copy for the last value, we can instead just use that graph! I know it's a small thing but this is kind of the whole point of making it's safely solvable as a single entity in the first place.
2. I think We might want to isolate compacting and remapping of the nodes until after certain thresholds are met. We need a dirty state and a next edge count state for it to to give us the signals for when to skip the process since it hasn't changed since the last remap, and what the compaction threshold is for the next remap. ENSURE WE PROPAGATE THIS STATE TO THE NEW GRAPH.
    2a. Maybe we just chose to optimistically do this after a merge or normalization has been performed to a triggering levels that way we are not duplicating this work each time we copy the same dirty graph for each generation? this would DRASTICALLY simplify the copy function...


# BFS based Generational Sieve depth limited by the lowest confirmed solution count
1. Same current solver approach BUT since we know we are not likely to explodiate the stack if we recurse, because the most we gerneally see is 24-26 but even 108 is reasonable if we consider memory per tier carefully.
2. We can still do the BFS based approach but at the cut off of the KEEP COUNT we then turn into a generator that feeds into a new goroutine recursive call of ourself that acts as a consumer of states to solve the next tier of the graph, repeating this process until we have solved the graph.
3. Maybe treat the BFS as a generational Sieve of Eratosthenes where we window the solutions to each tier after sorting them based on a Pluggable Heuristic.
4. If a Generation finishes generating without a solution it can safely close the chanel to the next generation and block waiting for the next generation to report the results back to it.

## ISSUES
1. How do we limit the depth of the recursion?
    1a. We can use our prior "best" generated solutions as a direct and verifiable limit to the maximum depth of the recursion.
2. How do we ensure the BFS at a lower generation that we just have not generated yet has been falsely cut off?
    2a. just as we used the prior "best" solution as a starting cut off, the generation that got a solution communicates it back to it's generator and goes into a discard loop for graphs that are still in the pip until the chanel is close by the parent generator. ENSURING THAT THE GRAPHS ARE RETURNED TO THE POOLING GENERALIZATION.
    2b. We also need to ensure that we stop producing new generations on the generatoin that finds a solution and this needs to stop sending  graphs to the next generation and signal it's closer by ending it's channel and then waiting for it's result channel to also close, discarding any solution as this would be a ... step too far...
3. How do we ensure that we are not generating the same graph in any generations?
    3a. NOT A DAMN CLUE TBH

# Color focused solver
What if we build an additional graph representation for each color on the map?
We might create these as a different type of edge resolution.
OR we might try to optimize for each color in the graph as a specific solver that then feeds into the main graph solver.


Reason for graph:
1. Simplification of problem space removing explicit spatial representation.
We might want to look at additional dimensionality that could be codified into the graph solver.
1. Edges of neighboring edges.
2. Edges of color nodes.
    * This would likely provide additional nodal weightings to the graph to make use of during the solving process.
    * This may end up being a many to many and without correlation back to the original graph.
3. 







# Fractal Spatial Coordinate System with N-Resolution Voxelization of N-Dimensional Space
The following may only work as an N-Dimensional N-Resolution voxelization of space which is exciting but maybe not relevant to the problem at hand.
* What if we make use of several resolutions of Hilbert Space Filling curve to generalize the path walking?
* The H-tree (fractal antenna's use this) is also a good example of a consistent space filling tree instead of curve that also strongly preserves locality across iterations.
* We can generalize locality via Hilbert R-tree like semantics.
* We'd start by mapping it at the highest resolution and then generalizing it into lower tiers linking these as parents to the higher resolution children, likely skipping at least one resolution to keep the tree balanced.
This appears to be some way to consistently voxelize a given multidimensional space. In theory this should scale to any dimension count...

To generalize it fully you need:
1. A space filling curve that provides the same locality in all resolutions of it use.
2. A Hilbert R-tre like generalization but for each resolution of the space filling curve.
3. A Way to map the space filling curve as it relates between each level of dimensionality.
4. An efficient data structure to represent these tiers of the space filling curve.
5. A way to walk the space filling curve in a way that is efficient and can be generalized to any resolution of the curve.

To extend this:
* I'd like to bake in the concepts of general relativity into motion and distances described in this space.
* I'd like to explore how a space filling curve might represent space time curvature if this has been explored yet at all that is...

## ISSUES
1. THIS REQUIRES A CONSISTENT ORIGIN POINT AND ORIENTATION WITHIN THE SPACE FOR THIS TO WORK AND THIS WHOLE DAMN REALITY IS RELATIVE TO EACH OBSERVER!!!!
    1a. It can still be used as a coord space in a snapshot of the simulations as a given state to realte things to eacother spatilally...
    1b. Game engine don't usually care about this detail even if I do...
2. things move and thes curves and trees both have large locaity jumps without representative  overlap.. PERHAPS WE USE MULTIPLE REPRESENTATIONS SMAYBE EVEN CURVES AND TREES TO ALLOW FOR LOCALITY TO BE PRESEVED ACROSS MORE THEN ONE TYPE OF THRESHOLD. ALA BLOOM FILTER?



Batching framework with workers that have their own buffers and pools to work on a graph Solving

# 3 Stage Solver
What if we do a keep count to find the median and then keep all values to the right that are
equivalent to the median's heuristic value? This would make the node layout and quirks of the graph
more stable and predictable at the cost of more memory. Ultimately the heuristic cutoff is arbitrary
to the dataset itself and as such it fails to satisfactorily solve the graph... We might be able to
make use of it for weighting the graph's edges or nodes in some way for an A* style walker. A Fast
walker to find a baseline isn't a bad approach but it's not the end goal.

What if we solve until >50% (181+) solved via BFS keeping all current tier solutions available, AND
THEN, we swap to an A* heuristic search to find the best heuristic solution, THEN, we swap to a
bactracking DFS limiting the depth by the lowest solution count found sofar?
Maybe the above threshold can be moved around for testing and MEM per worker needs?

# Implemented

## Need a Generic Functional Optioner package

It justifies it's existence as the validation step for the Options.

It should support these types of Options:
* `Option[T any](*T) *T`
    * Implements a handler for `OptionError`
    * Implements a handler for `OptionRevertable`
    * Implements a handler for `OptionRevertableError`
* `OptionRevertable[T any](*T) (*T, func(*T))`
    * Implements a handler for `OptionRevertableError`
* `OptionError[T any](*T) (*T, error)`
    * Implements a handler for `OptionRevertableError`
* `OptionRevertableError[T any](*T) (*T, func(*T), error)`

It should support these types of Apply functions:
* `Apply[T any](*T, Option[T]...)`
* `ApplyRevertable[T any](*T, OptionRevertable[T]...) []func()`
* `ApplyError[T any](*T, OptionError[T]...) error`
* `ApplyRevertableError[T any](*T, OptionRevertableError[T]...) ([]func(), error)`
Options and have an error and non error validator with the later able to accept both for easy
ergonomics for the dev to chose between the two handling options.


# Implemented but maybe poorly

## Need per worker circular buffers for *graph.Graph