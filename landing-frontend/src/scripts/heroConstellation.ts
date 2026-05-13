// Декоративная canvas-анимация в hero. Изначально была Vue-компонентом
// (см. landing-frontend/src/components/HeroConstellation.vue) — портирована
// в vanilla, чтобы не тащить Vue runtime ради 200 строк canvas-кода.

const SEED = 0x1F4A9
const NODE_COUNT = 38
const SIZE = 520
const RINGS = [0.32, 0.58, 0.86]

interface Node {
  baseX: number
  baseY: number
  ax: number
  ay: number
  freqX: number
  freqY: number
  ph: number
  size: number
  hub: boolean
  pulse: number
  posX: number
  posY: number
}

function mulberry32(seed: number) {
  let t = seed >>> 0
  return () => {
    t = (t + 0x6D2B79F5) >>> 0
    let r = t
    r = Math.imul(r ^ (r >>> 15), r | 1)
    r ^= r + Math.imul(r ^ (r >>> 7), r | 61)
    return ((r ^ (r >>> 14)) >>> 0) / 4294967296
  }
}

function buildNodes(): Node[] {
  const rnd = mulberry32(SEED)
  const nodes: Node[] = []

  nodes.push({
    baseX: 0,
    baseY: 0,
    ax: 0,
    ay: 0,
    freqX: 0,
    freqY: 0,
    ph: 0,
    size: 4.5,
    hub: true,
    pulse: 0,
    posX: 0,
    posY: 0,
  })

  const ringSum = RINGS.reduce((a, b) => a + b)
  for (const ring of RINGS) {
    const count = Math.round(NODE_COUNT * ring / ringSum)
    for (let i = 0; i < count; i++) {
      const baseAngle = (i / count) * Math.PI * 2 + rnd() * 0.4 - 0.2
      const baseR = ring + (rnd() * 0.07 - 0.035)
      const isHub = rnd() < 0.18
      const bx = Math.cos(baseAngle) * baseR
      const by = Math.sin(baseAngle) * baseR
      nodes.push({
        baseX: bx,
        baseY: by,
        ax: 0.012 + rnd() * 0.018,
        ay: 0.012 + rnd() * 0.018,
        freqX: 0.25 + rnd() * 0.35,
        freqY: 0.25 + rnd() * 0.35,
        ph: rnd() * Math.PI * 2,
        size: isHub ? 2.8 + rnd() * 1.2 : 1.2 + rnd() * 1.0,
        hub: isHub,
        pulse: 0,
        posX: bx,
        posY: by,
      })
    }
  }
  return nodes
}

interface Palette {
  raw: string
  fg: string
  hub: string
  line: string
  faint: string
  veryFaint: string
}

function readPalette(): Palette {
  const fallback = '151 60% 54%'
  const accent = getComputedStyle(document.documentElement).getPropertyValue('--accent').trim() || fallback
  return {
    raw: accent,
    fg: `hsl(${accent})`,
    hub: '#ffb547',
    line: `hsl(${accent} / 0.32)`,
    faint: `hsl(${accent} / 0.18)`,
    veryFaint: `hsl(${accent} / 0.10)`,
  }
}

export function initHeroConstellation(canvas: HTMLCanvasElement) {
  const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  const dpr = Math.min(window.devicePixelRatio || 1, 2)
  const cssSize = SIZE
  const px = cssSize * dpr
  canvas.width = px
  canvas.height = px

  const ctx = canvas.getContext('2d')
  if (!ctx)
    return
  ctx.scale(dpr, dpr)

  const nodes = buildNodes()
  const palette = readPalette()
  const cx = cssSize / 2
  const cy = cssSize / 2
  const R = cssSize * 0.48

  let mouseX = 0
  let mouseY = 0
  let targetMX = 0
  let targetMY = 0

  function onMove(e: MouseEvent) {
    const rect = canvas.getBoundingClientRect()
    const nx = (e.clientX - rect.left - rect.width / 2) / (rect.width / 2)
    const ny = (e.clientY - rect.top - rect.height / 2) / (rect.height / 2)
    targetMX = Math.max(-1, Math.min(1, nx)) * 8
    targetMY = Math.max(-1, Math.min(1, ny)) * 8
  }
  window.addEventListener('mousemove', onMove, { passive: true })

  let raf = 0
  const start = performance.now()
  let lastPulse = -2
  let pulseT = -1

  function compute(t: number) {
    mouseX += (targetMX - mouseX) * 0.06
    mouseY += (targetMY - mouseY) * 0.06

    for (const n of nodes) {
      if (reduceMotion) {
        n.posX = n.baseX
        n.posY = n.baseY
      }
      else {
        n.posX = n.baseX + Math.cos(t * n.freqX + n.ph) * n.ax
        n.posY = n.baseY + Math.sin(t * n.freqY + n.ph * 1.3) * n.ay
      }
    }

    if (!reduceMotion) {
      if (pulseT < 0 && t - lastPulse > 5) {
        pulseT = 0
        lastPulse = t
      }
      if (pulseT >= 0) {
        pulseT += 0.0055
        if (pulseT > 1.15)
          pulseT = -1
      }
    }

    for (const n of nodes) {
      if (n.pulse > 0)
        n.pulse = Math.max(0, n.pulse - 0.018)
    }

    if (pulseT >= 0) {
      for (const n of nodes) {
        const r = Math.hypot(n.posX, n.posY)
        if (Math.abs(r - pulseT) < 0.04)
          n.pulse = 1
      }
    }
  }

  function toScreen(nx: number, ny: number): [number, number] {
    return [cx + nx * R + mouseX, cy + ny * R + mouseY]
  }

  function drawChrome() {
    ctx!.strokeStyle = palette.line
    ctx!.lineWidth = 1
    ctx!.beginPath()
    ctx!.arc(cx, cy, R * 0.96, 0, Math.PI * 2)
    ctx!.stroke()

    ctx!.strokeStyle = palette.veryFaint
    ctx!.beginPath()
    ctx!.arc(cx, cy, R * 0.62, 0, Math.PI * 2)
    ctx!.stroke()

    ctx!.strokeStyle = palette.faint
    for (let i = 0; i < 24; i++) {
      const a = (i / 24) * Math.PI * 2
      const long = i % 6 === 0
      const r1 = R * 0.96
      const r2 = R * (long ? 0.92 : 0.94)
      ctx!.beginPath()
      ctx!.moveTo(cx + Math.cos(a) * r1, cy + Math.sin(a) * r1)
      ctx!.lineTo(cx + Math.cos(a) * r2, cy + Math.sin(a) * r2)
      ctx!.stroke()
    }

    ctx!.strokeStyle = palette.veryFaint
    ctx!.beginPath()
    ctx!.moveTo(cx - R * 0.96, cy)
    ctx!.lineTo(cx + R * 0.96, cy)
    ctx!.moveTo(cx, cy - R * 0.96)
    ctx!.lineTo(cx, cy + R * 0.96)
    ctx!.stroke()
  }

  function drawScanPulse() {
    if (pulseT < 0)
      return
    const r = pulseT * R
    const alpha = Math.max(0, 1 - pulseT) * 0.45
    ctx!.strokeStyle = `rgba(255, 181, 71, ${alpha})`
    ctx!.lineWidth = 1.2
    ctx!.beginPath()
    ctx!.arc(cx + mouseX, cy + mouseY, r, 0, Math.PI * 2)
    ctx!.stroke()
  }

  function drawEdges() {
    const k = 2
    for (let i = 0; i < nodes.length; i++) {
      const a = nodes[i]
      const dists: { j: number, d: number }[] = []
      for (let j = 0; j < nodes.length; j++) {
        if (i === j)
          continue
        const dx = nodes[j].posX - a.posX
        const dy = nodes[j].posY - a.posY
        dists.push({ j, d: dx * dx + dy * dy })
      }
      dists.sort((p, q) => p.d - q.d)
      for (let n = 0; n < k && n < dists.length; n++) {
        const b = nodes[dists[n].j]
        const d = Math.sqrt(dists[n].d)
        const alpha = Math.max(0, 0.45 - d * 0.55)
        if (alpha < 0.03)
          continue
        const [x1, y1] = toScreen(a.posX, a.posY)
        const [x2, y2] = toScreen(b.posX, b.posY)
        const ap = a.pulse * 0.7
        ctx!.strokeStyle = ap > 0
          ? `rgba(255, 181, 71, ${Math.min(0.7, alpha + ap)})`
          : `hsl(${palette.raw} / ${alpha})`
        ctx!.lineWidth = 0.7
        ctx!.beginPath()
        ctx!.moveTo(x1, y1)
        ctx!.lineTo(x2, y2)
        ctx!.stroke()
      }
    }
  }

  function drawNodes() {
    for (const n of nodes) {
      const [x, y] = toScreen(n.posX, n.posY)
      const pulseGlow = n.pulse
      const radius = n.size + pulseGlow * 2.2

      if (n.hub || pulseGlow > 0) {
        const grad = ctx!.createRadialGradient(x, y, 0, x, y, radius * 5)
        if (pulseGlow > 0) {
          grad.addColorStop(0, `rgba(255, 181, 71, ${0.55 + pulseGlow * 0.35})`)
          grad.addColorStop(1, 'rgba(255, 181, 71, 0)')
        }
        else {
          grad.addColorStop(0, `hsl(${palette.raw} / 0.45)`)
          grad.addColorStop(1, `hsl(${palette.raw} / 0)`)
        }
        ctx!.fillStyle = grad
        ctx!.beginPath()
        ctx!.arc(x, y, radius * 5, 0, Math.PI * 2)
        ctx!.fill()
      }

      ctx!.fillStyle = pulseGlow > 0.3 ? palette.hub : palette.fg
      ctx!.beginPath()
      ctx!.arc(x, y, radius, 0, Math.PI * 2)
      ctx!.fill()
    }
  }

  function frame() {
    const t = (performance.now() - start) / 1000
    compute(t)

    ctx!.clearRect(0, 0, cssSize, cssSize)
    ctx!.save()
    ctx!.beginPath()
    ctx!.arc(cx, cy, R * 0.97, 0, Math.PI * 2)
    ctx!.clip()

    drawChrome()
    drawScanPulse()
    drawEdges()
    drawNodes()

    ctx!.restore()

    if (!reduceMotion)
      raf = requestAnimationFrame(frame)
  }

  if (reduceMotion) {
    frame()
  }
  else {
    raf = requestAnimationFrame(frame)
  }

  return () => {
    cancelAnimationFrame(raf)
    window.removeEventListener('mousemove', onMove)
  }
}

export const HERO_CANVAS_SIZE = SIZE
