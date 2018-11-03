let keys = {}
document.addEventListener("keydown", function(e) {
    keys[e.code] = true;
});
document.addEventListener("keyup", function(e) {
    keys[e.code] = false;
});

let temp = 0;

class Rect {
    constructor(x, y, w, h, color="#FFFFFF") {
        this.x = x;
        this.y = y;
        this.h = h;
        this.w = w;
      
        this.vx = 0;
        this.vy = 0;
        this.frict=0;
        this.color=color;
    }
  
    Physics (dt) {
       this.x += this.vx*dt;
       this.y += this.vy*dt;
    }
  
    Draw(ctx) {
        ctx.fillStyle = this.color;
        ctx.fillRect(this.x - this.w/2, this.y - this.h/2, this.w, this.h);
    }
  
    Collide(rect, dt) {
        let x1 = this.x + this.w/2 > rect.x - rect.w/2
        let x2 = this.x - this.w/2 < rect.x + rect.w/2
        let y1 = this.y + this.h/2 > rect.y - rect.h/2 
        let y2 = this.y - this.h/2 < rect.y + rect.h/2
        
      
        if (x1 && x2 && y1 && y2) {
            let x=(this.x - rect.x)/rect.w;
            let y=(this.y - rect.y)/rect.h;
            if (Math.abs(x) > Math.abs(y)) {
                let vy = Math.abs(this.vx) * rect.frict * Math.sign(this.vy)*dt;
              
                if (Math.abs(vy) < Math.abs(this.vx) ) {
                    this.vx -= vy;
                } else this.vy = 0;
                this.vx = 0;
                if(x > 0) {
                   this.x = rect.x+rect.w/2+this.w/2;
                } else {
                   this.x = rect.x-rect.w/2-this.w/2;
                  
                }
            } else if (y == x) {
                this.vx = 0;
                this.vy = 0;
            } else {
                let vx = Math.abs(this.vy) * rect.frict * Math.sign(this.vx)*dt;
              
                if (Math.abs(vx) < Math.abs(this.vx) )
                  this.vx -= vx;
                else this.vx = 0;
              
                this.vy = 0;
                if(y > 0) {
                   this.y = rect.y+rect.h/2+this.h/2;
                } else {
                   this.y = rect.y-rect.h/2-this.h/2;
                                  this.ground = true;
                }
            }
        } 
    }
}

rect = new Rect(500,300, 50, 50, "#FF0000")
rect2 = new Rect(700, 600, 1400, 100)
rect2.frict = 1;
rect3 = new Rect(800, 300, 300, 100)
rect3.frict = 1;
class Game {
    constructor (ctx) {
        this.ctx = ctx;
        this.prevTime = performance.now(); 
    }
  
    Draw() {
        this.ctx.fillStyle = "#000000";
        this.ctx.fillRect(0, 0, this.ctx.canvas.width, this.ctx.canvas.height);
        rect.Draw(this.ctx);
        rect2.Draw(this.ctx)
        rect3.Draw(this.ctx)
    }
  
    Physics(dt) {
      rect.vy += 5*dt;
      rect.Physics(dt)
      rect.Collide(rect2, dt)
      rect.Collide(rect3, dt)
      if (keys.KeyD) {
          rect.vx += 5*dt;
      }
      if (keys.KeyA) {
          rect.vx += -5*dt;
      }
      if (keys.KeyW) {
        if (rect.ground)
        rect.vy = -60;
      }
      if (keys.KeyS) {
        rect.vy += 5*dt;
      }
      rect.ground = false;
    }
  
    RenderLoop() {
        let dt = performance.now() - this.prevTime;
        dt /= 1000/17
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