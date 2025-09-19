% Family Tree Example - Demonstrating Complex Relationships
% This example shows how to model family relationships using Prolog

% =============================================================================
% BASE FACTS - Family members and their relationships
% =============================================================================

% People in our family tree
person(john).
person(mary).
person(bob).
person(alice).
person(charlie).
person(diana).
person(edward).
person(susan).

% Direct parent relationships
parent(john, bob).       % John is parent of Bob
parent(mary, bob).       % Mary is parent of Bob  
parent(bob, alice).      % Bob is parent of Alice
parent(bob, charlie).    % Bob is parent of Charlie
parent(alice, diana).    % Alice is parent of Diana
parent(charlie, edward). % Charlie is parent of Edward
parent(susan, edward).   % Susan is parent of Edward

% Gender information
male(john).
male(bob).
male(charlie).
male(edward).
female(mary).
female(alice).
female(diana).
female(susan).

% =============================================================================
% DERIVED RELATIONSHIPS - Rules that define family relationships
% =============================================================================

% Father: male parent
father(X, Y) :- parent(X, Y), male(X).

% Mother: female parent  
mother(X, Y) :- parent(X, Y), female(X).

% Child: reverse of parent
child(X, Y) :- parent(Y, X).

% Son: male child
son(X, Y) :- child(X, Y), male(X).

% Daughter: female child
daughter(X, Y) :- child(X, Y), female(X).

% Grandparent: parent of parent
grandparent(X, Z) :- parent(X, Y), parent(Y, Z).

% Grandfather: male grandparent
grandfather(X, Z) :- grandparent(X, Z), male(X).

% Grandmother: female grandparent
grandmother(X, Z) :- grandparent(X, Z), female(X).

% Sibling: same parents, but not the same person
sibling(X, Y) :- 
    parent(Z, X), 
    parent(Z, Y), 
    X \= Y.

% Brother: male sibling
brother(X, Y) :- sibling(X, Y), male(X).

% Sister: female sibling
sister(X, Y) :- sibling(X, Y), female(X).

% Uncle: brother of parent
uncle(X, Y) :- parent(Z, Y), brother(X, Z).

% Aunt: sister of parent
aunt(X, Y) :- parent(Z, Y), sister(X, Z).

% Cousin: child of uncle or aunt
cousin(X, Y) :- parent(Z, X), parent(W, Y), sibling(Z, W).

% Ancestor: parent or ancestor of parent (recursive)
ancestor(X, Y) :- parent(X, Y).
ancestor(X, Y) :- parent(X, Z), ancestor(Z, Y).

% Descendant: reverse of ancestor
descendant(X, Y) :- ancestor(Y, X).

% =============================================================================
% EXAMPLE QUERIES TO EXPLORE THE FAMILY TREE
% =============================================================================

% Basic relationships:
% ?- father(john, bob).          % Is John the father of Bob? -> true
% ?- mother(mary, alice).        % Is Mary the mother of Alice? -> false
% ?- parent(bob, X).             % Who are Bob's children? -> X = alice; X = charlie

% Finding specific roles:
% ?- father(X, bob).             % Who is Bob's father? -> X = john
% ?- mother(X, bob).             % Who is Bob's mother? -> X = mary
% ?- child(X, bob).              % Who are Bob's children? -> X = alice; X = charlie

% Extended family:
% ?- grandparent(john, X).       % Who are John's grandchildren? -> X = alice; X = charlie
% ?- grandfather(X, diana).      % Who is Diana's grandfather? -> X = bob
% ?- sibling(alice, charlie).    % Are Alice and Charlie siblings? -> true

% Complex queries:
% ?- uncle(X, diana).            % Who is Diana's uncle? -> X = charlie
% ?- cousin(X, edward).          % Who are Edward's cousins? -> X = diana
% ?- ancestor(john, X).          % Who are John's descendants? -> X = bob; X = alice; X = charlie; X = diana

% Multi-generational:
% ?- ancestor(X, diana), ancestor(Y, X), X \= Y.  % Find Diana's ancestors and their ancestors