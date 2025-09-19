% Logic Puzzles - Advanced Prolog Problem Solving
% These examples demonstrate how to solve complex logic puzzles

% =============================================================================
% PUZZLE 1: Classic Einstein's Riddle (Simplified Version)
% =============================================================================

% There are 3 houses in a row, each with different colors, pets, and drinks
% Solve the puzzle based on the given clues

% Domain definitions
house_number(1). house_number(2). house_number(3).
color(red). color(blue). color(green).
pet(dog). pet(cat). pet(bird).
drink(tea). drink(coffee). drink(milk).

% Solution predicate - finds the complete assignment
solve_houses(Solution) :-
    % Solution format: [house(Number, Color, Pet, Drink), ...]
    Solution = [house(1, Color1, Pet1, Drink1),
                house(2, Color2, Pet2, Drink2),
                house(3, Color3, Pet3, Drink3)],
    
    % All colors must be different
    permutation([red, blue, green], [Color1, Color2, Color3]),
    
    % All pets must be different
    permutation([dog, cat, bird], [Pet1, Pet2, Pet3]),
    
    % All drinks must be different
    permutation([tea, coffee, milk], [Drink1, Drink2, Drink3]),
    
    % Constraints (clues):
    % 1. The red house is first
    Color1 = red,
    
    % 2. The person with the dog drinks coffee
    member(house(_, _, dog, coffee), Solution),
    
    % 3. The green house owner drinks tea
    member(house(_, green, _, tea), Solution),
    
    % 4. The cat lives in house 3
    Pet3 = cat,
    
    % 5. The blue house is next to the house with the bird
    (   (Color1 = blue, Pet2 = bird) ;
        (Color2 = blue, Pet1 = bird) ;
        (Color2 = blue, Pet3 = bird) ;
        (Color3 = blue, Pet2 = bird)
    ).

% =============================================================================
% PUZZLE 2: Map Coloring Problem
% =============================================================================

% Color a map so that no adjacent regions have the same color
% Using only 3 colors: red, blue, green

% Map structure (which regions are adjacent)
adjacent(region1, region2).
adjacent(region1, region3).
adjacent(region2, region3).
adjacent(region2, region4).
adjacent(region3, region4).
adjacent(region3, region5).
adjacent(region4, region5).

% Adjacency is symmetric
adjacent(X, Y) :- adjacent(Y, X).

% Available colors
map_color(red). map_color(blue). map_color(green).

% Solve map coloring
color_map(Coloring) :-
    Coloring = [color(region1, C1), color(region2, C2), color(region3, C3),
                color(region4, C4), color(region5, C5)],
    
    % Assign colors
    map_color(C1), map_color(C2), map_color(C3), map_color(C4), map_color(C5),
    
    % No adjacent regions can have the same color
    \+ (adjacent(R1, R2), member(color(R1, Color), Coloring), 
        member(color(R2, Color), Coloring)).

% =============================================================================
% PUZZLE 3: Sudoku (4x4 Mini Version)
% =============================================================================

% Solve a 4x4 Sudoku puzzle
% Each row, column, and 2x2 box must contain numbers 1-4

% Valid numbers
sudoku_num(1). sudoku_num(2). sudoku_num(3). sudoku_num(4).

% Check if a list contains all different elements
all_different([]).
all_different([H|T]) :- \+ member(H, T), all_different(T).

% Solve 4x4 Sudoku - grid represented as list of 16 elements
solve_sudoku_4x4(Grid) :-
    Grid = [A1,A2,A3,A4, B1,B2,B3,B4, C1,C2,C3,C4, D1,D2,D3,D4],
    
    % All cells must have valid numbers
    maplist(sudoku_num, Grid),
    
    % Rows must be all different
    all_different([A1,A2,A3,A4]),
    all_different([B1,B2,B3,B4]),
    all_different([C1,C2,C3,C4]),
    all_different([D1,D2,D3,D4]),
    
    % Columns must be all different
    all_different([A1,B1,C1,D1]),
    all_different([A2,B2,C2,D2]),
    all_different([A3,B3,C3,D3]),
    all_different([A4,B4,C4,D4]),
    
    % 2x2 boxes must be all different
    all_different([A1,A2,B1,B2]),  % Top-left box
    all_different([A3,A4,B3,B4]),  % Top-right box
    all_different([C1,C2,D1,D2]),  % Bottom-left box
    all_different([C3,C4,D3,D4]).  % Bottom-right box

% Solve with some pre-filled cells (partial puzzle)
solve_partial_sudoku(Grid) :-
    Grid = [1,_,_,4, _,3,_,_, _,_,2,_, 4,_,_,1],
    solve_sudoku_4x4(Grid).

% =============================================================================
% PUZZLE 4: N-Queens Problem (4-Queens)
% =============================================================================

% Place 4 queens on a 4x4 chessboard so none attack each other

% Safe position check - queens don't attack each other
safe_queens([]).
safe_queens([Q|Qs]) :- safe_queens(Qs), \+ attacks(Q, Qs, 1).

% Check if queen attacks any other queen
attacks(_, [], _).
attacks(Q, [Q1|Qs], D) :-
    (   Q =:= Q1                    % Same column
    ;   Q =:= Q1 + D                % Diagonal attack
    ;   Q =:= Q1 - D                % Other diagonal
    ), !.
attacks(Q, [_|Qs], D) :-
    D1 is D + 1,
    attacks(Q, Qs, D1).

% Solve 4-Queens
solve_4_queens(Solution) :-
    Solution = [Q1, Q2, Q3, Q4],
    permutation([1,2,3,4], Solution),
    safe_queens(Solution).

% =============================================================================
% EXAMPLE QUERIES FOR THE PUZZLES
% =============================================================================

% Einstein's Riddle:
% ?- solve_houses(Solution).

% Map Coloring:
% ?- color_map(Coloring).

% Sudoku:
% ?- solve_partial_sudoku(Grid).

% N-Queens:
% ?- solve_4_queens(Queens).