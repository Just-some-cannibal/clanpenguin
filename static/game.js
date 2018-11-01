let x = 0;
let y = 0;

let keys = {}
document.addEventListener("keydown", function(e) {
    keys[e.key] = true;
});
document.addEventListener("keyup", function(e) {
    keys[e.key] = false;
});

class Game {
    constructor (ctx) {
        this.ctx = ctx;
        this.prevTime = performance.now(); 
    }
  
    Draw() {
        this.ctx.fillStyle = "#000000";
        this.ctx.fillRect(0, 0, this.ctx.canvas.width, this.ctx.canvas.height);
        this.ctx.fillStyle = "#FFFFFF";
        this.ctx.fillRect(x, y, 100,100);
    }
  
    Physics(dt) {
      if (keys.d) {
          x += 20/dt;
      }
      if (keys.a) {
          x -= 20/dt;
      }
      if (keys.s) {
          y += 20/dt;
      }
      if (keys.w) {
          y -= 20/dt;
      }
    }
  
    RenderLoop() {
        let dt = performance.now() - this.prevTime;
        this.prevTime = performance.now();
        this.Physics(dt);
        this.Draw();
        requestAnimationFrame(this.RenderLoop.bind(this));
    }
}

window.onload = function() {
    let canvas = document.getElementById("game");
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
    
    let ctx = canvas.getContext('2d');
  
    game = new Game(ctx);
    game.RenderLoop();
}