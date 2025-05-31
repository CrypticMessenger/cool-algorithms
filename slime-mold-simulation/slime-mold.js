        const numAgents = 1000;
        let agents = [];
        let pheromoneGrid;
        let prevPheromoneGrid;
        let foodSources = [];
        const gridResolution = 10;
        let cols, rows;

        // --- Canvas and Layout Parameters ---
        let simulationWidth;
        let graphPanelWidth;
        let totalCanvasWidth;
        let canvasHeight;
        const SIMULATION_TO_TOTAL_WIDTH_RATIO = 0.7;

        // --- Tunable Parameters (Simulation) --- // UPDATED FOR SLIME MOLD BEHAVIOR
        const SENSOR_ANGLE_SPREAD = Math.PI / 5;    // Slightly narrower focus: ~36 degrees
        const SENSOR_OFFSET_DISTANCE = 10;          // Look a bit less far ahead
        const TURN_SPEED = 0.33;                    // Slightly more responsive turning
        const RANDOM_WOBBLE = 0.6;                  // Reduced base random wobble
        const TRAIL_WEIGHT = 2.5;                   // Significantly increased pheromone influence
        const MOVE_SPEED = 0.8;                     // Slightly slower movement
        const SATIATION_TIME = 50;                  // General satiation reference (SATIATION_TIME_ON_FINISH is key)
        const DEPOSIT_STRENGTH = 12;                // Slightly increased base deposit
        const FOOD_FINISH_DEPOSIT_BOOST = 100;      // Stronger boost when food is consumed
        const EVAPORATION_RATE = 0.994;             // SLOWER evaporation (0.6% decay per update)
        const DIFFUSION_RATE = 0.12;                // Reduced diffusion (sharper trails)
        const FOOD_ATTRACTION_MULTIPLIER = 2.0;     // Slightly reduced direct food pull
        const FOOD_PERCEPTION_RADIUS = 60;          // Slightly reduced perception radius
        const FOOD_CONSUME_RADIUS = 6;              // Agent must be very close to consume
        const INITIAL_FOOD_AMOUNT = 100;            // Starting health/quantity of food
        const AGENT_CONSUMPTION_RATE = 0.009;         // Amount agent eats per frame when in range
        const SATIATION_TIME_ON_FINISH = 250;       // Longer "satiated exploration" after finishing food
        const NEW_FOOD_SPAWN_INTERVAL = 350;        // Slightly longer for network to stabilize
        const MAX_FOOD_SOURCES = 10;                // Fewer food sources for clearer network structure
        const PHEROMONE_DRAW_THRESHOLD = 1.0;       // Slightly lower to see fainter forming trails

        // --- Agent Energy Parameters --- // UPDATED FOR SLIME MOLD BEHAVIOR
        const AGENT_INITIAL_ENERGY = 1200;          // Slightly more energy
        const AGENT_MAX_ENERGY = 1200;              // Max storable energy
        const AGENT_ENERGY_DECAY_PER_FRAME = 0.15;  // Slightly reduced decay
        const ENERGY_GAINED_PER_FOOD_UNIT = 200;    // Energy gained per AGENT_CONSUMPTION_RATE unit of food
        const LOW_ENERGY_THRESHOLD_RATIO = 0.3;     // Ratio of AGENT_INITIAL_ENERGY. Below this, agent is "hungry".
        let LOW_ENERGY_THRESHOLD;
        const AGENT_CRITICAL_ENERGY_THRESHOLD_RATIO = 0.1; // Below this, agent is "critically hungry"
        let CRITICAL_ENERGY_THRESHOLD;
        const PHEROMONE_INFLUENCE_HUNGRY_MULTIPLIER = 0.1; // Even less influence from old trails when very hungry
        const EXPLORATION_WOBBLE_MULTIPLIER_HUNGRY = 3.0;  // More aggressive random exploration when hungry


        let frameCounter = 0;

        // --- Analytics and Graphing ---
        let analytics = {
            liveAgents: 0, avgEnergy: 0, hungryAgents: 0, satiatedAgents: 0,
            numFoodSources: 0, totalFoodAmount: 0, activePheromoneCells: 0, fps: 0
        };
        const ANALYTICS_UPDATE_INTERVAL = 10;
        const MAX_HISTORY_LENGTH = 150;
        let analyticsHistory = {
            liveAgents: [], avgEnergy: [], totalFoodAmount: [], activePheromoneCells: []
        };
        const GRAPH_COLORS = {
            liveAgents: [100, 200, 255],
            avgEnergy: [255, 200, 100],
            totalFoodAmount: [255, 100, 100],
            activePheromoneCells: [100, 255, 100]
        };


        function setup() {
            totalCanvasWidth = windowWidth * 0.95;
            canvasHeight = windowHeight * 0.85;
            simulationWidth = totalCanvasWidth * SIMULATION_TO_TOTAL_WIDTH_RATIO;
            graphPanelWidth = totalCanvasWidth * (1 - SIMULATION_TO_TOTAL_WIDTH_RATIO);

            createCanvas(totalCanvasWidth, canvasHeight);

            cols = floor(simulationWidth / gridResolution);
            rows = floor(height / gridResolution);

            LOW_ENERGY_THRESHOLD = AGENT_INITIAL_ENERGY * LOW_ENERGY_THRESHOLD_RATIO;
            CRITICAL_ENERGY_THRESHOLD = AGENT_INITIAL_ENERGY * AGENT_CRITICAL_ENERGY_THRESHOLD_RATIO;

            initializeSimulationState();
            initializeAnalyticsHistory();
        }

        function initializeSimulationState() {
            pheromoneGrid = create2DArray(cols, rows);
            prevPheromoneGrid = create2DArray(cols, rows);

            agents = [];
            for (let i = 0; i < numAgents; i++) {
                agents.push(new Agent());
            }

            foodSources = [];
            for (let i = 0; i < 3; i++) {
                addFoodSource(random(simulationWidth * 0.1, simulationWidth * 0.9), random(height * 0.1, height * 0.9));
            }
            frameCounter = 0;
        }

        function initializeAnalyticsHistory() {
            for (let key in analyticsHistory) {
                analyticsHistory[key] = new Array(MAX_HISTORY_LENGTH).fill(0);
            }
        }


        function mousePressed() {
            if (mouseX > 0 && mouseX < simulationWidth && mouseY > 0 && mouseY < height) {
                addFoodSource(mouseX, mouseY);
            }
        }

        function addFoodSource(x, y) {
            if (x > 0 && x < simulationWidth && y > 0 && y < height && foodSources.length < MAX_FOOD_SOURCES + 3) {
                foodSources.push({
                    pos: createVector(x, y),
                    amount: INITIAL_FOOD_AMOUNT,
                    initialAmount: INITIAL_FOOD_AMOUNT
                });
            }
        }

        function draw() {
            background(10, 10, 15);

            push();
            updatePheromones();
            displayPheromonesOptimized();

            for (let i = agents.length - 1; i >= 0; i--) {
                const agent = agents[i];
                agent.update();

                if (agent.energy <= 0) {
                    agents.splice(i, 1);
                } else {
                    agent.display();
                }
            }

            displayFood();

            frameCounter++;
            if (frameCounter % NEW_FOOD_SPAWN_INTERVAL === 0 && foodSources.length < MAX_FOOD_SOURCES) {
                addFoodSource(random(simulationWidth * 0.1, simulationWidth * 0.9), random(height * 0.1, height * 0.9));
            }

            if (frameCount % ANALYTICS_UPDATE_INTERVAL === 0 || frameCount < ANALYTICS_UPDATE_INTERVAL) {
                 updateAnalytics();
            }
            displayAnalytics();
            pop();

            displayGraphs();
        }

        function updateAnalytics() {
            analytics.liveAgents = agents.length;
            let totalEnergy = 0;
            let currentHungryAgents = 0;
            let currentSatiatedAgents = 0;
            if (agents.length > 0) {
                for (let agent of agents) {
                    totalEnergy += agent.energy;
                    if (agent.satiationTimer <= 0 && agent.energy < LOW_ENERGY_THRESHOLD) {
                        currentHungryAgents++;
                    }
                    if (agent.satiationTimer > 0) {
                        currentSatiatedAgents++;
                    }
                }
                analytics.avgEnergy = agents.length > 0 ? totalEnergy / agents.length : 0;
            } else {
                analytics.avgEnergy = 0;
            }
            analytics.hungryAgents = currentHungryAgents;
            analytics.satiatedAgents = currentSatiatedAgents;
            analytics.numFoodSources = foodSources.length;
            let currentTotalFood = 0;
            for (let food of foodSources) {
                currentTotalFood += food.amount;
            }
            analytics.totalFoodAmount = currentTotalFood;
            let currentActivePheromoneCells = 0;
            for (let x = 0; x < cols; x++) {
                for (let y = 0; y < rows; y++) {
                    if (pheromoneGrid[x][y] > PHEROMONE_DRAW_THRESHOLD) {
                        currentActivePheromoneCells++;
                    }
                }
            }
            analytics.activePheromoneCells = currentActivePheromoneCells;
            analytics.fps = frameRate();

            for (let key in analyticsHistory) {
                if (analytics.hasOwnProperty(key)) {
                    analyticsHistory[key].push(analytics[key]);
                    if (analyticsHistory[key].length > MAX_HISTORY_LENGTH) {
                        analyticsHistory[key].shift();
                    }
                }
            }
        }

        function displayAnalytics() {
            push();
            fill(220, 220, 220, 200);
            textSize(12);
            textAlign(LEFT, TOP);
            textFont('monospace');
            let x = 10;
            let y = 10;
            const lineHeight = 15;
            text(`FPS: ${analytics.fps.toFixed(1)}`, x, y); y += lineHeight;
            text(`Live Agents: ${analytics.liveAgents}`, x, y); y += lineHeight;
            text(`Avg Energy: ${analytics.avgEnergy.toFixed(1)}`, x, y); y += lineHeight;
            text(`Hungry Agents: ${analytics.hungryAgents}`, x, y); y += lineHeight;
            text(`Satiated Agents: ${analytics.satiatedAgents}`, x, y); y += lineHeight;
            text(`Food Sources: ${analytics.numFoodSources}`, x, y); y += lineHeight;
            text(`Total Food: ${analytics.totalFoodAmount.toFixed(1)}`, x, y); y += lineHeight;
            text(`Active Pheromone Cells: ${analytics.activePheromoneCells}`, x, y);
            pop();
        }

        function displayGraphs() {
            push();
            translate(simulationWidth, 0);

            const numGraphs = Object.keys(analyticsHistory).length;
            if (numGraphs === 0) { pop(); return; }

            const graphHeight = (height - 20 * (numGraphs +1)) / numGraphs ;
            const graphWidth = graphPanelWidth - 40;
            let currentY = 20;

            drawMetricGraph(
                analyticsHistory.liveAgents, "Live Agents",
                20, currentY, graphWidth, graphHeight,
                0, numAgents, color(GRAPH_COLORS.liveAgents)
            );
            currentY += graphHeight + 20;

            drawMetricGraph(
                analyticsHistory.avgEnergy, "Avg Agent Energy",
                20, currentY, graphWidth, graphHeight,
                0, AGENT_MAX_ENERGY, color(GRAPH_COLORS.avgEnergy)
            );
            currentY += graphHeight + 20;

            let maxFoodEver = 0;
             for(let val of analyticsHistory.totalFoodAmount) maxFoodEver = max(maxFoodEver, val);
             maxFoodEver = max(maxFoodEver, INITIAL_FOOD_AMOUNT * MAX_FOOD_SOURCES * 0.5);
             if(maxFoodEver === 0) maxFoodEver = 100;

            drawMetricGraph(
                analyticsHistory.totalFoodAmount, "Total Food Amount",
                20, currentY, graphWidth, graphHeight,
                0, maxFoodEver, color(GRAPH_COLORS.totalFoodAmount)
            );
            currentY += graphHeight + 20;

            let maxPheroCells = cols * rows * 0.5;
            drawMetricGraph(
                analyticsHistory.activePheromoneCells, "Active Pheromone Cells",
                20, currentY, graphWidth, graphHeight,
                0, maxPheroCells, color(GRAPH_COLORS.activePheromoneCells)
            );

            pop();
        }

        function drawMetricGraph(data, label, x, y, w, h, minVal, maxVal, lineColor) {
            push();
            translate(x, y);
            stroke(100);
            fill(20, 20, 30, 200);
            rect(0, 0, w, h);
            fill(220);
            noStroke();
            textSize(10);
            textAlign(CENTER);
            text(label, w / 2, -5);
            textAlign(RIGHT, CENTER);
            text(maxVal.toFixed(0), -2, 0);
            textAlign(RIGHT, CENTER);
            text(minVal.toFixed(0), -2, h);
            noFill();
            stroke(lineColor);
            strokeWeight(1.5);
            beginShape();
            for (let i = 0; i < data.length; i++) {
                let xPos = map(i, 0, MAX_HISTORY_LENGTH - 1, 0, w);
                let yPos = map(data[i], minVal, maxVal, h, 0);
                yPos = constrain(yPos, 0, h);
                vertex(xPos, yPos);
            }
            endShape();
            pop();
        }


        function create2DArray(c, r) {
            let arr = new Array(c);
            for (let i = 0; i < arr.length; i++) {
                arr[i] = new Array(r).fill(0);
            }
            return arr;
        }

        function updatePheromones() {
            for (let x = 0; x < cols; x++) {
                for (let y = 0; y < rows; y++) {
                    prevPheromoneGrid[x][y] = pheromoneGrid[x][y];
                }
            }
            let nextPheromoneGrid = create2DArray(cols, rows);
            for (let x = 1; x < cols - 1; x++) {
                for (let y = 1; y < rows - 1; y++) {
                    let currentVal = pheromoneGrid[x][y];
                    let diffusedSum = 0;
                    diffusedSum += pheromoneGrid[x-1][y-1]; diffusedSum += pheromoneGrid[x][y-1]; diffusedSum += pheromoneGrid[x+1][y-1];
                    diffusedSum += pheromoneGrid[x-1][y  ];                                       diffusedSum += pheromoneGrid[x+1][y  ];
                    diffusedSum += pheromoneGrid[x-1][y+1]; diffusedSum += pheromoneGrid[x][y+1]; diffusedSum += pheromoneGrid[x+1][y+1];
                    diffusedSum /= 8.0;
                    let diffusedValue = currentVal * (1 - DIFFUSION_RATE) + diffusedSum * DIFFUSION_RATE;
                    nextPheromoneGrid[x][y] = constrain(diffusedValue * EVAPORATION_RATE, 0, 255);
                }
            }
             for(let x=0; x<cols; x++){
                if(x===0 || x === cols-1){ for(let y=0; y<rows; y++) nextPheromoneGrid[x][y] = constrain(pheromoneGrid[x][y] * EVAPORATION_RATE, 0, 255); }
                else { nextPheromoneGrid[x][0] = constrain(pheromoneGrid[x][0] * EVAPORATION_RATE, 0, 255); nextPheromoneGrid[x][rows-1] = constrain(pheromoneGrid[x][rows-1] * EVAPORATION_RATE, 0, 255); }
            }
            pheromoneGrid = nextPheromoneGrid;
        }

        function displayPheromonesOptimized() {
            noStroke();
            for (let x = 0; x < cols; x++) {
                for (let y = 0; y < rows; y++) {
                    let currentPVal = pheromoneGrid[x][y];
                    let prevPVal = prevPheromoneGrid[x][y];
                    if (currentPVal > PHEROMONE_DRAW_THRESHOLD || (prevPVal > PHEROMONE_DRAW_THRESHOLD && currentPVal <= PHEROMONE_DRAW_THRESHOLD) ) {
                        if (currentPVal > PHEROMONE_DRAW_THRESHOLD) {
                            fill(40, 180, 80, currentPVal * 0.8);
                        } else {
                            fill(10, 10, 15);
                        }
                        rect(x * gridResolution, y * gridResolution, gridResolution, gridResolution);
                    }
                }
            }
        }

        function displayFood() {
            for (let food of foodSources) {
                let size = map(food.amount, 0, food.initialAmount, gridResolution * 0.5, gridResolution * 1.5);
                size = max(size, gridResolution * 0.3);
                let alpha = map(food.amount, 0, food.initialAmount, 100, 220);
                fill(255, 80, 80, alpha);
                stroke(255, 150, 150, alpha*0.8);
                strokeWeight(1);
                ellipse(food.pos.x, food.pos.y, size, size);
            }
        }

        class Agent {
            constructor() {
                this.pos = createVector(random(simulationWidth), random(height));
                this.angle = random(TWO_PI);
                this.vel = p5.Vector.fromAngle(this.angle).mult(MOVE_SPEED);
                this.defaultAgentColor = color(230, 230, 100, 180);
                this.satiationTimer = 0;
                this.energy = AGENT_INITIAL_ENERGY;
            }

            update() {
                if (this.satiationTimer > 0) this.satiationTimer--;
                this.energy -= AGENT_ENERGY_DECAY_PER_FRAME;
                this.energy = max(0, this.energy);
                if (this.energy <= 0) return;

                this.senseAndSteer();
                this.move();
                this.depositPheromone();
                this.consumeFood();
                this.edges();
            }

            senseAndSteer() {
                let sensorAngleRad = SENSOR_ANGLE_SPREAD;
                let fwd = p5.Vector.fromAngle(this.angle).mult(SENSOR_OFFSET_DISTANCE).add(this.pos);
                let lft = p5.Vector.fromAngle(this.angle - sensorAngleRad).mult(SENSOR_OFFSET_DISTANCE).add(this.pos);
                let rgt = p5.Vector.fromAngle(this.angle + sensorAngleRad).mult(SENSOR_OFFSET_DISTANCE).add(this.pos);

                let valFwd = this.getStimulusAt(fwd.x, fwd.y);
                let valLft = this.getStimulusAt(lft.x, lft.y);
                let valRgt = this.getStimulusAt(rgt.x, rgt.y);

                let desiredAngle = this.angle;
                let currentRandomWobble = RANDOM_WOBBLE;
                let currentTurnSpeed = TURN_SPEED;

                if (this.satiationTimer > SATIATION_TIME_ON_FINISH * 0.75) {
                    desiredAngle += random(-currentTurnSpeed * 1.8, currentTurnSpeed * 1.8);
                } else if (this.satiationTimer > SATIATION_TIME * 0.5) { // Using general SATIATION_TIME here for intermediate satiation
                    desiredAngle += random(-currentTurnSpeed * 1.2, currentTurnSpeed * 1.2);
                } else if (this.energy < LOW_ENERGY_THRESHOLD && this.satiationTimer <= 0) {
                    currentRandomWobble *= EXPLORATION_WOBBLE_MULTIPLIER_HUNGRY;
                    if (valFwd > valLft && valFwd > valRgt && valFwd > 0.01) {
                    } else if (valLft > valRgt && valLft > 0.01) {
                        desiredAngle -= currentTurnSpeed;
                    } else if (valRgt > valLft && valRgt > 0.01) {
                        desiredAngle += currentTurnSpeed;
                    } else {
                        desiredAngle += random(-currentTurnSpeed * 1.2, currentTurnSpeed * 1.2);
                    }
                } else {
                    if (valFwd > valLft && valFwd > valRgt && valFwd > 0.1) {
                    } else if (valLft > valRgt) {
                        desiredAngle -= currentTurnSpeed;
                    } else if (valRgt > valLft) {
                        desiredAngle += currentTurnSpeed;
                    } else {
                        desiredAngle += random(-currentTurnSpeed, currentTurnSpeed) * 0.5;
                    }
                }
                desiredAngle += random(-currentRandomWobble, currentRandomWobble);
                this.angle = lerpAngle(this.angle, desiredAngle, 0.7);
            }

            getStimulusAt(x, y) {
                let stimulus = 0;
                if (x >= 0 && x < simulationWidth && y >= 0 && y < height) {
                    let gridX = floor(x / gridResolution);
                    let gridY = floor(y / gridResolution);
                    gridX = constrain(gridX, 0, cols - 1);
                    gridY = constrain(gridY, 0, rows - 1);
                    let pheromoneValue = pheromoneGrid[gridX][gridY] * TRAIL_WEIGHT;
                    if (this.energy < LOW_ENERGY_THRESHOLD && this.satiationTimer <= 0) {
                        pheromoneValue *= PHEROMONE_INFLUENCE_HUNGRY_MULTIPLIER;
                    }
                    stimulus += pheromoneValue;
                }
                let foodAttractionFactor = map(this.satiationTimer, 0, SATIATION_TIME_ON_FINISH, 1, 0.1, true); // Use SATIATION_TIME_ON_FINISH for food attraction decay
                for (let food of foodSources) {
                    let d = dist(x, y, food.pos.x, food.pos.y);
                    if (d < FOOD_PERCEPTION_RADIUS) {
                        let foodStrength = (food.amount / food.initialAmount);
                        stimulus += (1 / (d * d / (FOOD_PERCEPTION_RADIUS * 0.5) + 1)) *
                                    FOOD_ATTRACTION_MULTIPLIER * 100 * foodStrength * foodAttractionFactor;
                    }
                }
                return stimulus;
            }

            move() {
                this.vel = p5.Vector.fromAngle(this.angle).mult(MOVE_SPEED);
                this.pos.add(this.vel);
            }

            depositPheromone(strength = DEPOSIT_STRENGTH) {
                if (this.energy <= 0) return;
                if (this.pos.x >=0 && this.pos.x < simulationWidth && this.pos.y >=0 && this.pos.y < height) {
                    let gridX = floor(this.pos.x / gridResolution);
                    let gridY = floor(this.pos.y / gridResolution);
                    if (gridX >= 0 && gridX < cols && gridY >= 0 && gridY < rows) {
                        pheromoneGrid[gridX][gridY] += strength;
                        pheromoneGrid[gridX][gridY] = min(pheromoneGrid[gridX][gridY], 255);
                    }
                }
            }

            consumeFood() {
                for (let i = foodSources.length - 1; i >= 0; i--) {
                    let food = foodSources[i];
                    let d = dist(this.pos.x, this.pos.y, food.pos.x, food.pos.y);
                    if (d < FOOD_CONSUME_RADIUS) {
                        let amountToConsume = min(food.amount, AGENT_CONSUMPTION_RATE);
                        food.amount -= amountToConsume;
                        this.energy += amountToConsume * ENERGY_GAINED_PER_FOOD_UNIT;
                        this.energy = min(this.energy, AGENT_MAX_ENERGY);
                        this.depositPheromone(DEPOSIT_STRENGTH * 0.5);
                        if (food.amount <= 0) {
                            foodSources.splice(i, 1);
                            this.depositPheromone(DEPOSIT_STRENGTH + FOOD_FINISH_DEPOSIT_BOOST);
                            this.satiationTimer = SATIATION_TIME_ON_FINISH;
                            this.angle += PI + random(-PI / 2.5, PI / 2.5);
                        } else {
                             if (this.satiationTimer < SATIATION_TIME * 0.2) this.satiationTimer = SATIATION_TIME * 0.2; // Use general SATIATION_TIME for partial eating
                        }
                        return true;
                    }
                }
                return false;
            }

            display() {
                push();
                translate(this.pos.x, this.pos.y);
                rotate(this.angle + PI / 2);
                let agentFill;
                if (this.satiationTimer > 0) {
                    if (this.satiationTimer > SATIATION_TIME_ON_FINISH * 0.75) { // Visual cue based on longer satiation
                        agentFill = color(180, 180, 255, 160);
                    } else {
                        agentFill = color(200, 200, 200, 170);
                    }
                } else {
                    if (this.energy <= CRITICAL_ENERGY_THRESHOLD) {
                        agentFill = color(150, 0, 0, 150);
                    } else if (this.energy < LOW_ENERGY_THRESHOLD) {
                        agentFill = color(255, 100, 0, 170);
                    } else {
                        agentFill = this.defaultAgentColor;
                    }
                }
                fill(agentFill);
                noStroke();
                beginShape(); vertex(0, -3.5); vertex(-2, 2.5); vertex(2, 2.5); endShape(CLOSE);
                pop();
            }

            edges() {
                if (this.pos.x > simulationWidth) { this.pos.x = 0; this.angle += random(-0.3,0.3); }
                else if (this.pos.x < 0) { this.pos.x = simulationWidth; this.angle += random(-0.3,0.3); }
                if (this.pos.y > height) { this.pos.y = 0; this.angle += random(-0.3,0.3); }
                else if (this.pos.y < 0) { this.pos.y = height; this.angle += random(-0.3,0.3); }
            }
        }

        function lerpAngle(startAngle, endAngle, amount) {
            let difference = endAngle - startAngle;
            while (difference < -PI) difference += TWO_PI;
            while (difference > PI) difference -= TWO_PI;
            return startAngle + difference * amount;
        }

        function windowResized() {
            totalCanvasWidth = windowWidth * 0.95;
            canvasHeight = windowHeight * 0.85;
            simulationWidth = totalCanvasWidth * SIMULATION_TO_TOTAL_WIDTH_RATIO;
            graphPanelWidth = totalCanvasWidth * (1 - SIMULATION_TO_TOTAL_WIDTH_RATIO);

            resizeCanvas(totalCanvasWidth, canvasHeight);

            cols = floor(simulationWidth / gridResolution);
            rows = floor(height / gridResolution);

            LOW_ENERGY_THRESHOLD = AGENT_INITIAL_ENERGY * LOW_ENERGY_THRESHOLD_RATIO;
            CRITICAL_ENERGY_THRESHOLD = AGENT_INITIAL_ENERGY * AGENT_CRITICAL_ENERGY_THRESHOLD_RATIO;

            initializeSimulationState();
            initializeAnalyticsHistory();
        }
