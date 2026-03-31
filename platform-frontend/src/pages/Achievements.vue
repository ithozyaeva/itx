<script setup lang="ts">
import type { AchievementCategory, AchievementsResponse, UserAchievement } from '@/models/achievement'
import { Typography } from 'itx-ui-kit'
import {
  Award,
  BookOpen,
  Briefcase,
  CalendarCheck,
  CheckCircle,
  ClipboardList,
  Crown,
  FileText,
  Flame,
  Footprints,
  Gem,
  GraduationCap,
  HardHat,
  History,
  ListChecks,
  Loader2,
  Medal,
  MessageSquare,
  MessagesSquare,
  Mic,
  Package,
  Presentation,
  Share2,
  ShoppingCart,
  Star,
  Swords,
  Target,
  Trophy,
  UserCheck,
  UserPlus,
  Users,
  Zap,
} from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { achievementsService } from '@/services/achievements'
import { handleError } from '@/services/errorService'

const data = ref<AchievementsResponse | null>(null)
const isLoading = ref(true)
const activeCategory = ref<AchievementCategory | 'all'>('all')
const selectedAchievement = ref<UserAchievement | null>(null)
const showDetailDialog = ref(false)

function openDetail(achievement: UserAchievement) {
  selectedAchievement.value = achievement
  showDetailDialog.value = true
}

const iconMap: Record<string, any> = {
  'footprints': Footprints,
  'flame': Flame,
  'calendar-check': CalendarCheck,
  'medal': Medal,
  'mic': Mic,
  'presentation': Presentation,
  'star': Star,
  'trophy': Trophy,
  'crown': Crown,
  'gem': Gem,
  'message-square': MessageSquare,
  'messages-square': MessagesSquare,
  'book-open': BookOpen,
  'share-2': Share2,
  'users': Users,
  'user-plus': UserPlus,
  'user-check': UserCheck,
  'zap': Zap,
  'file-text': FileText,
  'package': Package,
  'shopping-cart': ShoppingCart,
  'clipboard-list': ClipboardList,
  'list-checks': ListChecks,
  'check-circle': CheckCircle,
  'hard-hat': HardHat,
  'briefcase': Briefcase,
  'target': Target,
  'swords': Swords,
  'history': History,
  'graduation-cap': GraduationCap,
}

const categories: { key: AchievementCategory | 'all', label: string }[] = [
  { key: 'all', label: 'Все' },
  { key: 'events', label: 'События' },
  { key: 'points', label: 'Баллы' },
  { key: 'social', label: 'Социальные' },
  { key: 'activity', label: 'Активность' },
]

const filteredItems = computed(() => {
  if (!data.value)
    return []
  if (activeCategory.value === 'all')
    return data.value.items
  return data.value.items.filter(a => a.category === activeCategory.value)
})

async function fetchAchievements() {
  isLoading.value = true
  try {
    data.value = await achievementsService.getMyAchievements()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchAchievements()
})
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <div class="flex items-center justify-between mb-6">
      <Typography
        variant="h2"
        as="h1"
      >
        Достижения
      </Typography>
      <div
        v-if="data"
        class="flex items-center gap-2 text-sm"
      >
        <Award class="h-5 w-5 text-yellow-500" />
        <span class="font-bold">{{ data.unlockedCount }}</span>
        <span class="text-muted-foreground">/ {{ data.totalCount }}</span>
      </div>
    </div>

    <div
      v-if="isLoading"
      class="flex justify-center py-12"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <template v-else-if="data">
      <div class="flex gap-2 mb-6 flex-wrap">
        <button
          v-for="cat in categories"
          :key="cat.key"
          class="px-3 py-1.5 rounded-full text-sm font-medium transition-colors"
          :class="activeCategory === cat.key
            ? 'bg-primary text-primary-foreground'
            : 'bg-card border border-border text-muted-foreground hover:text-foreground'"
          @click="activeCategory = cat.key"
        >
          {{ cat.label }}
        </button>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="achievement in filteredItems"
          :key="achievement.id"
          class="rounded-2xl p-4 transition-colors border cursor-pointer hover:shadow-md"
          :class="achievement.unlocked
            ? 'bg-green-500/5 border-green-500/30'
            : 'bg-card border-border opacity-60'"
          @click="openDetail(achievement)"
        >
          <div class="flex items-start gap-3">
            <div
              class="flex items-center justify-center w-10 h-10 rounded-full shrink-0"
              :class="achievement.unlocked ? 'bg-green-500/20' : 'bg-primary/10'"
            >
              <component
                :is="iconMap[achievement.icon] || Award"
                class="h-5 w-5"
                :class="achievement.unlocked ? 'text-green-500' : 'text-primary'"
              />
            </div>
            <div class="flex-1 min-w-0">
              <div class="font-medium text-sm">
                {{ achievement.title }}
              </div>
              <div class="text-xs text-muted-foreground mt-0.5">
                {{ achievement.description }}
              </div>
            </div>
          </div>

          <div class="mt-3">
            <div class="flex justify-between text-xs text-muted-foreground mb-1">
              <span>{{ Math.min(achievement.progress, achievement.threshold) }} / {{ achievement.threshold }}</span>
              <span v-if="achievement.unlocked" class="text-green-500 font-medium">Получено</span>
            </div>
            <div class="w-full h-1.5 rounded-full bg-muted overflow-hidden">
              <div
                class="h-full rounded-full transition-all"
                :class="achievement.unlocked ? 'bg-green-500' : 'bg-primary'"
                :style="{ width: `${achievement.threshold > 0 ? Math.min(100, (achievement.progress / achievement.threshold) * 100) : 0}%` }"
              />
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Achievement Detail Dialog -->
    <Dialog
      v-model:open="showDetailDialog"
    >
      <DialogContent
        v-if="selectedAchievement"
        class="sm:max-w-md"
      >
        <DialogHeader>
          <div class="flex justify-center mb-4">
            <div
              class="flex items-center justify-center w-16 h-16 rounded-full"
              :class="selectedAchievement.unlocked ? 'bg-green-500/20' : 'bg-primary/10'"
            >
              <component
                :is="iconMap[selectedAchievement.icon] || Award"
                class="h-8 w-8"
                :class="selectedAchievement.unlocked ? 'text-green-500' : 'text-primary'"
              />
            </div>
          </div>
          <DialogTitle class="text-center">
            {{ selectedAchievement.title }}
          </DialogTitle>
          <DialogDescription class="text-center">
            {{ selectedAchievement.description }}
          </DialogDescription>
        </DialogHeader>

        <div class="space-y-4 pt-2">
          <!-- Unlock status -->
          <div class="flex items-center justify-center gap-2">
            <span
              v-if="selectedAchievement.unlocked"
              class="inline-flex items-center gap-1.5 text-sm font-medium text-green-500"
            >
              <CheckCircle class="h-4 w-4" />
              Получено
            </span>
            <span
              v-else
              class="text-sm text-muted-foreground"
            >
              Ещё не получено
            </span>
          </div>

          <!-- Progress bar -->
          <div>
            <div class="flex justify-between text-sm text-muted-foreground mb-1.5">
              <span>Прогресс</span>
              <span class="font-medium">{{ Math.min(selectedAchievement.progress, selectedAchievement.threshold) }} / {{ selectedAchievement.threshold }}</span>
            </div>
            <div class="w-full h-2 rounded-full bg-muted overflow-hidden">
              <div
                class="h-full rounded-full transition-all"
                :class="selectedAchievement.unlocked ? 'bg-green-500' : 'bg-primary'"
                :style="{ width: `${selectedAchievement.threshold > 0 ? Math.min(100, (selectedAchievement.progress / selectedAchievement.threshold) * 100) : 0}%` }"
              />
            </div>
            <p class="text-xs text-muted-foreground text-right mt-1">
              {{ selectedAchievement.threshold > 0 ? Math.min(100, Math.round((selectedAchievement.progress / selectedAchievement.threshold) * 100)) : 0 }}%
            </p>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  </div>
</template>
