% Basic Prolog Examples for Logic MCP Server
% These examples demonstrate fundamental Prolog concepts

% =============================================================================
% FACTS - Basic building blocks
% =============================================================================

% Simple facts about animals
animal(cat).
animal(dog).
animal(bird).
animal(fish).
animal(elephant).

% Properties of animals
has_fur(cat).
has_fur(dog).
has_fur(elephant).
has_feathers(bird).
has_scales(fish).

% Sizes
small(cat).
small(bird).
small(fish).
large(dog).
large(elephant).

% =============================================================================
% RULES - Logical relationships
% =============================================================================

% A mammal is an animal that has fur
mammal(X) :- animal(X), has_fur(X).

% Something that can fly has feathers (simplified)
can_fly(X) :- animal(X), has_feathers(X).

% A pet is a small mammal or a bird
pet(X) :- mammal(X), small(X).
pet(X) :- animal(X), has_feathers(X).

% =============================================================================
% EXAMPLE QUERIES TO TRY
% =============================================================================

% Simple yes/no queries:
% ?- mammal(cat).           % true
% ?- mammal(bird).          % false
% ?- can_fly(bird).         % true
% ?- pet(dog).              % false (dog is large)

% Finding all solutions:
% ?- mammal(X).             % X = cat; X = dog; X = elephant
% ?- pet(X).                % X = cat; X = bird
% ?- small(X).              % X = cat; X = bird; X = fish

% Complex queries:
% ?- animal(X), has_fur(X), small(X).  % X = cat