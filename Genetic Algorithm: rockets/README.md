## Implementation of Genetic Algorithm(Paradigm) in Golang

- Thanks to Daniel Shiffman (Nature of Code) for the inspiration.
- This is a simulation, where given a target, a rocket will try to reach the target using genetic algorithm.
- The rocket will have a DNA, which will have genes that will determine the movement of the rocket.(i.e. force in x and y direction during each frame)
- The rocket will have a fitness score, which will be calculated based on the distance from the target.
- error will be calculated based on the distance from the target, and if the result is satifactory, the simulation will stop.

You can expect a trajectory i.e. Force that the thrusters need to apply to reach the target in each frame.

### Demo


https://github.com/CrypticMessenger/cool-algorithms/assets/75074904/b82f2570-dc70-4a17-b334-7a06f3a9678a


### Todo

- [ ] Add obstacles and change fitness function accordingly.
