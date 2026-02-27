<script setup lang="ts">
import createRegl from 'regl'
import { onMounted, onUnmounted, ref } from 'vue'

const canvasRef = ref<HTMLCanvasElement | null>(null)
const mouse = ref([0, 0])

const lerp = (x: number, y: number, a: number) => x * (1 - a) + y * a

function lerpMatrix(mat1: number[][], mat2: number[][], v: number) {
  const result: number[][] = []

  for (let x = 0; x < mat1.length; x++) {
    result[x] = []
    for (let y = 0; y < mat1[x].length; y++) {
      result[x][y] = lerp(mat1[x][y], mat2[x][y], v)
    }
  }

  return result
}

const step = 0.15
const lineWidth = 0.025
const halfWidth = lineWidth * 0.37
const edge = 0.07

const left = step + 0.08 - halfWidth
const right = step - 0.08 + halfWidth
const nw = Math.PI * 1.75
const nn = Math.PI * 1.5
const ss = nn + Math.PI
const ww = Math.PI
const ee = ww + Math.PI
const sw = Math.PI * 2.25
const se = Math.PI * 2.75
const ne = Math.PI * 3.25

const dots = [
  [0, -0.07],
  [0, 0.07],
  [0.15, -0.07],
  [0.15, 0.07],
  [0.15, 0.0],
]

// 8
const EIGHT = [
  [-0.08, -0.08, 0.08, 0.08],
  [0.08, -0.08, -0.08, 0.08],

  [
    step + Math.cos(nw) * 0.065,
    Math.sin(nw) * 0.065,
    step + Math.cos(se) * 0.065,
    Math.sin(se) * 0.065,
  ],

  [
    step + Math.cos(ne) * 0.065,
    Math.sin(ne) * 0.065,
    step + Math.cos(sw) * 0.065,
    Math.sin(sw) * 0.065,
  ],

  [
    step + Math.cos(ee) * 0.065,
    Math.sin(ee) * 0.065,
    step + Math.cos(ww) * 0.065,
    Math.sin(ww) * 0.065,
  ],

  [
    step + Math.cos(ee) * 0.065,
    Math.sin(ee) * 0.065,
    step + Math.cos(ww) * 0.065,
    Math.sin(ww) * 0.065,
  ],
]

const i_x = step + 0.035
const t_x = step - 0.035

// IT
const IT = [
  [-0.04, -0.04, 0.04, 0.04],
  [0.04, -0.04, -0.04, 0.04],

  [
    t_x + Math.cos(ee) * 0.04,
    -0.03 + Math.sin(ee) * 0.04,
    t_x + Math.cos(ww) * 0.04,
    -0.03 + Math.sin(ww) * 0.04,
  ],
  [
    t_x + Math.cos(nn) * halfWidth,
    -0.03 + Math.sin(nn) * halfWidth,
    t_x + Math.cos(ss) * 0.08,
    -0.03 + Math.sin(ss) * 0.08,
  ],
  [
    i_x + Math.cos(nn) * 0.04,
    Math.sin(nn) * 0.04,
    i_x + Math.cos(ss) * 0.05,
    Math.sin(ss) * 0.05,
  ],
  [
    i_x + Math.cos(nn) * 0.04,
    Math.sin(nn) * 0.04,
    i_x + Math.cos(ss) * 0.05,
    Math.sin(ss) * 0.05,
  ],
]

// <*>
const FIRST = [
  [-0.08, -0.08, 0.08, 0.08],
  [0.08, -0.08, -0.08, 0.08],

  [
    right + Math.cos(nw) * 0.06,
    Math.sin(nw) * 0.06,
    right + Math.cos(se) * halfWidth,
    Math.sin(se) * halfWidth,
  ],
  [
    right + Math.cos(sw) * 0.06,
    Math.sin(sw) * 0.06,
    right + Math.cos(ne) * halfWidth,
    Math.sin(ne) * halfWidth,
  ],
  [
    left + Math.cos(ne) * 0.06,
    Math.sin(ne) * 0.06,
    left + Math.cos(sw) * halfWidth,
    Math.sin(sw) * halfWidth,
  ],
  [
    left + Math.cos(se) * 0.06,
    Math.sin(se) * 0.06,
    left + Math.cos(nw) * halfWidth,
    Math.sin(nw) * halfWidth,
  ],
]

const EYE = [
  [-0.18, -0.18, 0.18, 0.18],
  [0.18, -0.18, -0.18, 0.18],

  [
    step + halfWidth,
    -halfWidth - edge,
    step - 0.08,
    -0.08 + edge + lineWidth * 0.5,
  ],
  [
    step - halfWidth,
    -halfWidth - edge,
    step + 0.08,
    -0.08 + edge + lineWidth * 0.5,
  ],
  [step + halfWidth, halfWidth + edge, step - 0.045, -0.045 + edge],
  [step - halfWidth, halfWidth + edge, step + 0.045, -0.045 + edge],
]

const BOOKS = [
  [-0.08, -0.08, 0.08, 0.08],
  [0.08, -0.08, -0.08, 0.08],

  [step + halfWidth, halfWidth, step - 0.045, -0.045],
  [step - halfWidth, halfWidth, step + 0.045, -0.045],
  [step + halfWidth, halfWidth + edge, step - 0.08, -0.08 + edge],
  [step - halfWidth, halfWidth + edge, step + 0.08, -0.08 + edge],
]

const figures = [IT, FIRST, EYE, EIGHT, BOOKS, IT, IT]

const lines = figures[0]

const state = {
  step: 0,
}

const uniforms: Record<string, () => any> = {
  time: () => performance.now() / 1000,
  resolution: () => [canvasRef.value?.clientWidth, canvasRef.value?.clientHeight],
  lineWidth: () => lineWidth,
  step: () => step,
  state: () => state.step,
  mouse: () => mouse.value,
  scale: () => 1 / Math.min(window.devicePixelRatio, 2),
}

let interpolated = lerpMatrix(figures[0], figures[1], 0)

for (let i = 0; i < 10; i++) {
  uniforms[`line${i}`] = () => [0, 0, 0, 0]
}

lines.forEach((line, i) => {
  uniforms[`line${i}`] = () => interpolated[i]
})

uniforms.linesCount = () => lines.length

for (let i = 0; i < 10; i++) {
  uniforms[`dot${i}`] = () => [0, 0, 0, 0]
}

dots.forEach((line, i) => {
  uniforms[`dot${i}`] = () => line
})

uniforms.dotsCount = () => dots.length

onMounted(() => {
  if (canvasRef.value) {
    const dpr = Math.min(window.devicePixelRatio, 2)

    canvasRef.value.width = canvasRef.value.clientWidth * dpr
    canvasRef.value.height = canvasRef.value.clientHeight * dpr

    const handleResize = () => {
      if (!canvasRef.value)
        return

      const dpr = Math.min(window.devicePixelRatio, 2)

      canvasRef.value.width = canvasRef.value.clientWidth * dpr
      canvasRef.value.height = canvasRef.value.clientHeight * dpr
    }

    window.addEventListener('resize', handleResize)

    const handleMouseMove = (e: MouseEvent) => {
      if (!canvasRef.value)
        return

      const rect = canvasRef.value.getBoundingClientRect()
      const x = e.clientX / window.innerWidth * 2 - 1
      const y = -((e.clientY - rect.top) / rect.height * 2 - 1)
      mouse.value = [x, y]
    }

    window.addEventListener('mousemove', handleMouseMove)

    const regl = createRegl({
      canvas: canvasRef.value,
    })

    const drawFullscreenQuad = regl({
      frag: `precision highp float;

      uniform float lineWidth;
      uniform float step;
      uniform float state;
      uniform float scale;
      uniform vec2 mouse;
      
      ${lines.map((line, i) => {
        return `
          uniform vec4 line${i};
        `
      }).join(`
      `)}
      uniform float linesCount;
      ${dots.map((line, i) => {
        return `
          uniform vec2 dot${i};
        `
      }).join(`
      `)}

      uniform float dotsCount;
      uniform float time;
      uniform vec2 resolution;

      float sdOrientedBox( in vec2 p, in vec2 a, in vec2 b, float th )
      {
          float l = length(b-a);
          vec2  d = (b-a)/l;
          vec2  q = (p-(a+b)*0.5);
                q = mat2(d.x,-d.y,d.y,d.x)*q;
                q = abs(q)-vec2(l,th)*0.5;
          return length(max(q,0.0)) + min(max(q.x,q.y),0.0);    
      }

      vec2 getPosition(in vec2 uv, in vec2 p, in float step) {  
        float stepX = step;
        float stepY = step;
        
        return mod(uv + vec2(p.x, p.y), step * 2.) - vec2(stepX, stepY);
      }

      float addLine(in vec2 uv, in vec4 p) {
        vec2 center = (p.xy + p.zw) * 0.5;
        
        return sdOrientedBox(
          getPosition(uv, center, step), 
          p.xy - center, 
          p.zw - center, 
          lineWidth
        );
      }

      float addDot(in vec2 uv, in vec2 p) {
        float halfWidth = lineWidth * 0.37;
        
        return sdOrientedBox(
          getPosition(uv, p, step), 
          vec2(halfWidth, -halfWidth), 
          vec2(-halfWidth, halfWidth),
          lineWidth
        );
      }

      float ripple(vec4 p, float rad, float freq, float amp, float time) {
        return 
          sin(rad) * amp
          + sin(p.x * freq) * amp 
          + sin(p.y * freq) * amp
          + sin((p.z * 0.01) * freq) * amp
          + sin(time * 0.01) * amp
          ;
      }

      void main() {
        vec2 uv = gl_FragCoord.xy / resolution.xy*2.-1.;
        float ratio = resolution.x / resolution.y;
        uv.x *= ratio;
        uv *= scale;
        vec2 ruv = uv;

        float invscale = 1. - scale;

        vec2 mouse_scaled = mouse * vec2(1. * ratio, 1.) + vec2(1. * ratio, 1.0) * invscale;
    
        float anim = clamp(sin(time * 0.5) + 1., 0., 0.999);
        float dist = clamp(pow(distance(ruv, mouse_scaled) * 2.0, 5.5), 0., 1.);

        vec2 rippled_uv = uv + ripple(vec4(uv, uv), 0., pow(anim, 50.) * 15. + sin(time * 0.5), 0.05, time);

        uv = mix(ruv, rippled_uv, dist);
      
        // lines
        float l = 1.;
        // dots
        float d = 1.;
        float i = 0.;
        float halfWidth = lineWidth * 0.37;
        vec3 fgColor = vec3(1.0, 0.4, 0.8);
        vec3 bgColor = vec3(1., 1., 1.);

        bgColor = mix(bgColor, vec3(0.95), dist);
        
        ${lines.map((line, i) => {
          return `
            l = min(addLine(uv, line${i}), l);
            i = ${i}.;
            if (l > 0.001 && i > 0.0) {
              fgColor = vec3(0.4, 0.4, 0.8);
            }
            
            if (i > 0.0) {
              fgColor = mix(bgColor, fgColor, smoothstep(-0.007, -0.002, l));
            }
          `
        }).join(`
        `)}
        
        ${dots.map((line, i) => {
          return `
            d = min(addDot(uv, dot${i}), d);
            if (d < 0.001) {
              fgColor = fgColor;
            }
          `
        }).join(`
        `)}
        
        if (state > 5. || state < 1.) {
          d = 1.;
        }
        
        d = smoothstep(0.0, 0.002, min(l, d));

        vec3 color = mix(fgColor, bgColor, d);
        
        gl_FragColor = vec4(color, 1.);
      }
      `,
      vert: `
        precision mediump float;
        attribute vec2 position;
        varying vec2 vUv;

        void main() {
          vUv = 0.5 * (position + 1.0);
          gl_Position = vec4(position, 0, 1);
        }
      `,
      attributes: {
        position: [
          [-1, -1],
          [1, -1],
          [1, 1],
          [-1, -1],
          [1, 1],
          [-1, 1],
        ],
      },
      uniforms,
      count: 6,
    })

    const renderController = regl.frame(() => {
      drawFullscreenQuad()

      const time = performance.now() / 1000
      const step = time % (figures.length - 1)

      state.step = step

      const from = figures[Math.floor(step)]
      const to = figures[Math.floor(step) + 1]

      interpolated = lerpMatrix(from, to, (time % 1) ** 16)
    })

    onUnmounted(() => {
      renderController.cancel()
      window.removeEventListener('mousemove', handleMouseMove)
      window.removeEventListener('resize', handleResize)
    })
  }
})
</script>

<template>
  <canvas
    id="canvas"
    ref="canvasRef"
    class="absolute inset-0 w-full h-full"
  />
</template>
